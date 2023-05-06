// Copyright © 2023 Kube logging authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package syslogng_agent

import (
	"crypto/sha256"
	"fmt"
	"strconv"

	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/merge"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	util "github.com/cisco-open/operator-tools/pkg/utils"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/kube-logging/logging-operator/pkg/resources/nodeagent"
	"github.com/kube-logging/logging-operator/pkg/resources/templates"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

type SyslogNGAgentReconciler struct {
	agent           v1beta1.SyslogNGAgent
	agentReconciler *nodeagent.GenericAgentReconciler
	dataProvider    nodeagent.AgentDataProvider
	log             logr.Logger
}

func NewSyslogNGAgentReconciler(
	reconciler *nodeagent.GenericAgentReconciler,
	dataProvider nodeagent.AgentDataProvider,
	log logr.Logger,
	agent v1beta1.SyslogNGAgent,
) *SyslogNGAgentReconciler {
	return &SyslogNGAgentReconciler{
		agentReconciler: reconciler,
		agent:           agent,
		dataProvider:    dataProvider,
		log:             log,
	}
}

func (s *SyslogNGAgentReconciler) Reconcile() (*reconcile.Result, error) {
	metricsEnabled := s.agent.Spec.Metrics != nil
	prometheusAnnotationsEnabled := metricsEnabled && s.agent.Spec.Metrics.PrometheusAnnotations
	if spec, err := nodeagent.NodeAgentSyslogNGDefaults(metricsEnabled, prometheusAnnotationsEnabled); err != nil {
		s.agent.Spec = *spec
		return nil, errors.Wrap(err, "applying syslogNG defaults")
	} else {
		err = merge.Merge(&s.agent.Spec, spec)
		if err != nil {
			return nil, err
		}
	}

	input := syslogNGConfig{
		TargetHost: s.dataProvider.TargetHost(),
	}
	secret, err := nodeagent.GenerateConfigSecret(input, syslogNGConfigTemplate, s.dataProvider.ConfigFileName())

	resourceName := s.dataProvider.QualifiedName("")

	resourceBuilders := []reconciler.ResourceBuilder{
		s.configSecretResource(&secret),
		s.daemonSetWithConfigRollout(&secret, resourceName),
	}

	result, err := s.agentReconciler.Reconcile(resourceBuilders)
	return &result, err
}

func (s *SyslogNGAgentReconciler) daemonSetWithConfigRollout(secret *corev1.Secret, serviceAccountName string) func() (runtime.Object, reconciler.DesiredState, error) {
	return func() (runtime.Object, reconciler.DesiredState, error) {
		return s.daemonSet(secret, serviceAccountName)
	}
}

func (s *SyslogNGAgentReconciler) configSecretResource(secret *corev1.Secret) reconciler.ResourceBuilder {
	return func() (runtime.Object, reconciler.DesiredState, error) {
		err := s.agentReconciler.ChildObjectMeta(secret, s.dataProvider.QualifiedName("config"))
		return secret, reconciler.StatePresent, err
	}
}

func (n *SyslogNGAgentReconciler) daemonSet(configSecret *corev1.Secret, serviceAccountName string) (runtime.Object, reconciler.DesiredState, error) {
	var containerPorts []corev1.ContainerPort
	podSecurityContext := corev1.PodSecurityContext{}
	containerSecurityContext := corev1.SecurityContext{}
	var desired *appsv1.DaemonSet
	meta := metav1.ObjectMeta{}
	var containerName string

	if n.agent.Spec.Metrics != nil && n.agent.Spec.Metrics.Port != 0 {
		containerPorts = append(containerPorts, corev1.ContainerPort{
			Name:          "monitor",
			ContainerPort: n.agent.Spec.Metrics.Port,
			Protocol:      corev1.ProtocolTCP,
		})
	}
	podSecurityContext = corev1.PodSecurityContext{
		FSGroup:      n.agent.Spec.Security.PodSecurityContext.FSGroup,
		RunAsNonRoot: n.agent.Spec.Security.PodSecurityContext.RunAsNonRoot,
		RunAsUser:    n.agent.Spec.Security.PodSecurityContext.RunAsUser,
		RunAsGroup:   n.agent.Spec.Security.PodSecurityContext.RunAsGroup,
	}
	containerSecurityContext = corev1.SecurityContext{
		RunAsUser:                n.agent.Spec.Security.SecurityContext.RunAsUser,
		RunAsNonRoot:             n.agent.Spec.Security.SecurityContext.RunAsNonRoot,
		ReadOnlyRootFilesystem:   n.agent.Spec.Security.SecurityContext.ReadOnlyRootFilesystem,
		AllowPrivilegeEscalation: n.agent.Spec.Security.SecurityContext.AllowPrivilegeEscalation,
		Privileged:               n.agent.Spec.Security.SecurityContext.Privileged,
		SELinuxOptions:           n.agent.Spec.Security.SecurityContext.SELinuxOptions,
	}

	if err := n.agentReconciler.ChildObjectMeta(desired, n.dataProvider.QualifiedName("")); err != nil {
		return nil, nil, errors.WrapIf(err, "generating object meta for syslog-ng-agent daemonset")
	}

	containerName = n.dataProvider.GetConstants().ContainerName

	desired = n.prepareDaemonSet(configSecret, serviceAccountName, meta, podSecurityContext, containerName, containerPorts, containerSecurityContext)

	n.agent.Spec.BufferStorageVolume.WithDefaultHostPath(
		fmt.Sprintf(v1beta1.HostPath, n.dataProvider.GetConstants().LoggingName, n.dataProvider.GetConstants().VolumeName))

	// TODO take care of persistfile
	if err := n.agent.Spec.BufferStorageVolume.ApplyVolumeForPodSpec(
		n.dataProvider.GetConstants().VolumeName,
		n.dataProvider.GetConstants().ContainerName,
		n.dataProvider.GetConstants().StoragePath,
		&desired.Spec.Template.Spec); err != nil {
		return desired, reconciler.StatePresent, err
	}
	if mergeErr := merge.Merge(desired, n.agent.Spec.DaemonSetOverrides); mergeErr != nil {
		return desired, reconciler.StatePresent, errors.WrapIf(mergeErr, "unable to merge overrides to base object")
	}

	return desired, reconciler.StatePresent, nil
}

func (n *SyslogNGAgentReconciler) prepareDaemonSet(
	configSecret *corev1.Secret,
	serviceAccountName string,
	meta metav1.ObjectMeta,
	podSecurityContext corev1.PodSecurityContext,
	containerName string, containerPorts []corev1.ContainerPort,
	containerSecurityContext corev1.SecurityContext,
) *appsv1.DaemonSet {

	podMeta := metav1.ObjectMeta{
		Labels:      n.dataProvider.ResourceLabels(),
		Annotations: n.dataProvider.ResourceAnnotations(),
	}

	h := sha256.New()
	for _, d := range configSecret.Data {
		_, _ = h.Write(d)
	}
	configHash := fmt.Sprintf("%x", h.Sum(nil))
	podMeta = templates.Annotate(podMeta, "checksum", configHash)

	desired := &appsv1.DaemonSet{
		ObjectMeta: meta,
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: n.dataProvider.ResourceLabels()},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: podMeta,
				Spec: corev1.PodSpec{
					ServiceAccountName: serviceAccountName,
					Volumes:            n.generateVolume(configSecret),
					SecurityContext:    &podSecurityContext,
					Containers: []corev1.Container{
						{
							Name:            containerName,
							Ports:           containerPorts,
							VolumeMounts:    n.generateVolumeMounts(),
							SecurityContext: &containerSecurityContext,
						},
					},
				},
			},
		},
	}
	return desired
}

func (n *SyslogNGAgentReconciler) generateVolumeMounts() (v []corev1.VolumeMount) {

	v = []corev1.VolumeMount{
		{
			Name:      "containerspath",
			MountPath: n.agent.Spec.ContainersPath,
		},
		{
			Name:      "varlogspath",
			MountPath: n.agent.Spec.VarLogsPath,
		},
	}

	for vCount, vMnt := range n.agent.Spec.ExtraVolumeMounts {
		v = append(v, corev1.VolumeMount{
			Name:      "extravolumemount" + strconv.Itoa(vCount),
			ReadOnly:  util.PointerToBool(vMnt.ReadOnly),
			MountPath: vMnt.Destination,
		})
	}

	if n.agent.Spec.CustomConfigSecret == "" {
		v = append(v, corev1.VolumeMount{
			Name:      "config",
			MountPath: n.dataProvider.GetConstants().ConfigPath,
		})
	} else {
		// TODO
	}

	if n.agent.Spec.TLS.Enabled {
		tlsRelatedVolume := []corev1.VolumeMount{
			{
				Name:      "syslog-ng-tls",
				MountPath: "/syslog-ng/tls/",
			},
		}
		v = append(v, tlsRelatedVolume...)
	}

	return
}

func (n *SyslogNGAgentReconciler) generateVolume(configSecret *corev1.Secret) (v []corev1.Volume) {
	v = []corev1.Volume{
		{
			Name: "containerspath",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: n.agent.Spec.ContainersPath,
				},
			},
		},
		{
			Name: "varlogspath",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: n.agent.Spec.VarLogsPath,
				},
			},
		},
	}

	for vCount, vMnt := range n.agent.Spec.ExtraVolumeMounts {
		v = append(v, corev1.Volume{
			Name: "extravolumemount" + strconv.Itoa(vCount),
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: vMnt.Source,
				},
			}})
	}

	if n.agent.Spec.CustomConfigSecret == "" {
		volume := corev1.Volume{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: configSecret.Name,
					Items: []corev1.KeyToPath{
						{
							// TODO replace with constants from the syslogng_agent package
							Key:  "BaseConfigNameSyslogNG",
							Path: "BaseConfigNameSyslogNG",
						},
					},
				},
			},
		}

		v = append(v, volume)
	} else {
		v = append(v, corev1.Volume{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: n.agent.Spec.CustomConfigSecret,
				},
			},
		})
	}
	if n.agent.Spec.TLS.Enabled {
		tlsRelatedVolume := corev1.Volume{
			Name: "syslog-ng-tls",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: n.agent.Spec.TLS.SecretName,
				},
			},
		}
		v = append(v, tlsRelatedVolume)
	}
	return
}
