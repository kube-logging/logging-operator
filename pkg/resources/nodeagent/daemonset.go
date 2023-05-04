// Copyright © 2021 Cisco Systems, Inc. and/or its affiliates
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

package nodeagent

import (
	"crypto/sha256"
	"fmt"
	"strconv"

	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/merge"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	util "github.com/cisco-open/operator-tools/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kube-logging/logging-operator/pkg/resources/templates"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

const (
	TailPositionVolume  = "positiondb"
	BufferStorageVolume = "buffers"
)

func (n *nodeAgentInstance) daemonSet() (runtime.Object, reconciler.DesiredState, error) {
	var containerPorts []corev1.ContainerPort
	podSecurityContext := corev1.PodSecurityContext{}
	containerSecurityContext := corev1.SecurityContext{}
	var desired *appsv1.DaemonSet
	meta := metav1.ObjectMeta{}
	var containerName string

	if n.nodeAgent.FluentbitSpec != nil {
		if n.nodeAgent.FluentbitSpec.Metrics != nil && n.nodeAgent.FluentbitSpec.Metrics.Port != 0 {
			containerPorts = append(containerPorts, corev1.ContainerPort{
				Name:          "monitor",
				ContainerPort: n.nodeAgent.FluentbitSpec.Metrics.Port,
				Protocol:      corev1.ProtocolTCP,
			})
		}
		podSecurityContext = corev1.PodSecurityContext{
			FSGroup:      n.nodeAgent.FluentbitSpec.Security.PodSecurityContext.FSGroup,
			RunAsNonRoot: n.nodeAgent.FluentbitSpec.Security.PodSecurityContext.RunAsNonRoot,
			RunAsUser:    n.nodeAgent.FluentbitSpec.Security.PodSecurityContext.RunAsUser,
			RunAsGroup:   n.nodeAgent.FluentbitSpec.Security.PodSecurityContext.RunAsGroup,
		}
		containerSecurityContext = corev1.SecurityContext{
			RunAsUser:                n.nodeAgent.FluentbitSpec.Security.SecurityContext.RunAsUser,
			RunAsNonRoot:             n.nodeAgent.FluentbitSpec.Security.SecurityContext.RunAsNonRoot,
			ReadOnlyRootFilesystem:   n.nodeAgent.FluentbitSpec.Security.SecurityContext.ReadOnlyRootFilesystem,
			AllowPrivilegeEscalation: n.nodeAgent.FluentbitSpec.Security.SecurityContext.AllowPrivilegeEscalation,
			Privileged:               n.nodeAgent.FluentbitSpec.Security.SecurityContext.Privileged,
			SELinuxOptions:           n.nodeAgent.FluentbitSpec.Security.SecurityContext.SELinuxOptions,
		}
		meta = n.NodeAgentObjectMeta(DaemonSetNameFluentbit)
		containerName = containerNameFluentbit

		desired = n.prepareDaemonSet(meta, podSecurityContext, containerName, containerPorts, containerSecurityContext)

		n.nodeAgent.FluentbitSpec.PositionDB.WithDefaultHostPath(
			fmt.Sprintf(v1beta1.HostPath, n.logging.Name, TailPositionVolume))
		n.nodeAgent.FluentbitSpec.BufferStorageVolume.WithDefaultHostPath(
			fmt.Sprintf(v1beta1.HostPath, n.logging.Name, BufferStorageVolume))
		if err := n.nodeAgent.FluentbitSpec.PositionDB.ApplyVolumeForPodSpec(TailPositionVolume, containerNameFluentbit, "/tail-db", &desired.Spec.Template.Spec); err != nil {
			return desired, reconciler.StatePresent, err
		}
		if err := n.nodeAgent.FluentbitSpec.BufferStorageVolume.ApplyVolumeForPodSpec(BufferStorageVolume, containerNameFluentbit, n.nodeAgent.FluentbitSpec.BufferStorage.StoragePath, &desired.Spec.Template.Spec); err != nil {
			return desired, reconciler.StatePresent, err
		}
		if mergeErr := merge.Merge(desired, n.nodeAgent.FluentbitSpec.DaemonSetOverrides); mergeErr != nil {
			return desired, reconciler.StatePresent, errors.WrapIf(mergeErr, "unable to merge overrides to base object")
		}
		return desired, reconciler.StatePresent, nil
	}

	if n.nodeAgent.SyslogNGSpec != nil {
		if n.nodeAgent.SyslogNGSpec.Metrics != nil && n.nodeAgent.SyslogNGSpec.Metrics.Port != 0 {
			containerPorts = append(containerPorts, corev1.ContainerPort{
				Name:          "monitor",
				ContainerPort: n.nodeAgent.SyslogNGSpec.Metrics.Port,
				Protocol:      corev1.ProtocolTCP,
			})
		}
		podSecurityContext = corev1.PodSecurityContext{
			FSGroup:      n.nodeAgent.SyslogNGSpec.Security.PodSecurityContext.FSGroup,
			RunAsNonRoot: n.nodeAgent.SyslogNGSpec.Security.PodSecurityContext.RunAsNonRoot,
			RunAsUser:    n.nodeAgent.SyslogNGSpec.Security.PodSecurityContext.RunAsUser,
			RunAsGroup:   n.nodeAgent.SyslogNGSpec.Security.PodSecurityContext.RunAsGroup,
		}
		containerSecurityContext = corev1.SecurityContext{
			RunAsUser:                n.nodeAgent.SyslogNGSpec.Security.SecurityContext.RunAsUser,
			RunAsNonRoot:             n.nodeAgent.SyslogNGSpec.Security.SecurityContext.RunAsNonRoot,
			ReadOnlyRootFilesystem:   n.nodeAgent.SyslogNGSpec.Security.SecurityContext.ReadOnlyRootFilesystem,
			AllowPrivilegeEscalation: n.nodeAgent.SyslogNGSpec.Security.SecurityContext.AllowPrivilegeEscalation,
			Privileged:               n.nodeAgent.SyslogNGSpec.Security.SecurityContext.Privileged,
			SELinuxOptions:           n.nodeAgent.SyslogNGSpec.Security.SecurityContext.SELinuxOptions,
		}
		meta = n.NodeAgentObjectMeta(DaemonSetNameSyslogNG)
		containerName = containerNameSyslogNG

		desired = n.prepareDaemonSet(meta, podSecurityContext, containerName, containerPorts, containerSecurityContext)

		n.nodeAgent.SyslogNGSpec.BufferStorageVolume.WithDefaultHostPath(
			fmt.Sprintf(v1beta1.HostPath, n.logging.Name, BufferStorageVolume))

		// TODO take care of persistfile
		if err := n.nodeAgent.SyslogNGSpec.BufferStorageVolume.ApplyVolumeForPodSpec(BufferStorageVolume, containerNameSyslogNG, n.nodeAgent.SyslogNGSpec.BufferStorage.StoragePath, &desired.Spec.Template.Spec); err != nil {
			return desired, reconciler.StatePresent, err
		}
		if mergeErr := merge.Merge(desired, n.nodeAgent.SyslogNGSpec.DaemonSetOverrides); mergeErr != nil {
			return desired, reconciler.StatePresent, errors.WrapIf(mergeErr, "unable to merge overrides to base object")
		}
	}
	return desired, reconciler.StatePresent, nil
}

func (n *nodeAgentInstance) prepareDaemonSet(meta metav1.ObjectMeta, podSecurityContext corev1.PodSecurityContext, containerName string, containerPorts []corev1.ContainerPort, containerSecurityContext corev1.SecurityContext) *appsv1.DaemonSet {
	podMeta := metav1.ObjectMeta{
		Labels:      n.getNodeAgentLabels(),
		Annotations: n.nodeAgent.Metadata.Annotations,
	}

	if n.configs != nil {
		for key, config := range n.configs {
			h := sha256.New()
			_, _ = h.Write(config)
			podMeta = templates.Annotate(podMeta, fmt.Sprintf("checksum/%s", key), fmt.Sprintf("%x", h.Sum(nil)))
		}
	}

	desired := &appsv1.DaemonSet{
		ObjectMeta: meta,
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: n.getNodeAgentLabels()},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: podMeta,
				Spec: corev1.PodSpec{
					ServiceAccountName: n.getServiceAccount(),
					Volumes:            n.generateVolume(),
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

func (n *nodeAgentInstance) generateVolumeMounts() (v []corev1.VolumeMount) {
	if n.nodeAgent.FluentbitSpec != nil {
		v = []corev1.VolumeMount{
			{
				Name:      "containerspath",
				MountPath: n.nodeAgent.FluentbitSpec.ContainersPath,
			},
			{
				Name:      "varlogspath",
				MountPath: n.nodeAgent.FluentbitSpec.VarLogsPath,
			},
		}

		for vCount, vMnt := range n.nodeAgent.FluentbitSpec.ExtraVolumeMounts {
			v = append(v, corev1.VolumeMount{
				Name:      "extravolumemount" + strconv.Itoa(vCount),
				ReadOnly:  util.PointerToBool(vMnt.ReadOnly),
				MountPath: vMnt.Destination,
			})
		}

		if n.nodeAgent.FluentbitSpec.CustomConfigSecret == "" {
			v = append(v, corev1.VolumeMount{
				Name:      "config",
				MountPath: "/fluent-bit/conf_operator",
			})
			if util.PointerToBool(n.nodeAgent.FluentbitSpec.EnableUpstream) {
				v = append(v, corev1.VolumeMount{
					Name:      "config",
					MountPath: "/fluent-bit/conf_upstream",
				})
			}
		} else {
			v = append(v, corev1.VolumeMount{
				Name:      "config",
				MountPath: "/fluent-bit/etc/",
			})
		}

		if n.nodeAgent.FluentbitSpec != nil && n.nodeAgent.FluentbitSpec.TLS != nil && util.PointerToBool(n.nodeAgent.FluentbitSpec.TLS.Enabled) {
			tlsRelatedVolume := []corev1.VolumeMount{
				{
					Name:      "fluent-bit-tls",
					MountPath: "/fluent-bit/tls/",
				},
			}
			v = append(v, tlsRelatedVolume...)
		}

	} else if n.nodeAgent.SyslogNGSpec != nil {
		v = []corev1.VolumeMount{
			{
				Name:      "containerspath",
				MountPath: n.nodeAgent.SyslogNGSpec.ContainersPath,
			},
			{
				Name:      "varlogspath",
				MountPath: n.nodeAgent.SyslogNGSpec.VarLogsPath,
			},
		}

		for vCount, vMnt := range n.nodeAgent.SyslogNGSpec.ExtraVolumeMounts {
			v = append(v, corev1.VolumeMount{
				Name:      "extravolumemount" + strconv.Itoa(vCount),
				ReadOnly:  util.PointerToBool(vMnt.ReadOnly),
				MountPath: vMnt.Destination,
			})
		}

		if n.nodeAgent.SyslogNGSpec.CustomConfigSecret == "" {
			v = append(v, corev1.VolumeMount{
				Name:      "config",
				MountPath: "/etc/syslog-ng/config/syslog-ng.conf",
			})
		} else {
			// TODO
		}

		if n.nodeAgent.SyslogNGSpec != nil && n.nodeAgent.SyslogNGSpec.TLS.Enabled {
			tlsRelatedVolume := []corev1.VolumeMount{
				{
					Name:      "syslog-ng-tls",
					MountPath: "/syslog-ng/tls/",
				},
			}
			v = append(v, tlsRelatedVolume...)
		}
	}

	return
}

func (n *nodeAgentInstance) generateVolume() (v []corev1.Volume) {
	if n.nodeAgent.FluentbitSpec != nil {
		v = []corev1.Volume{
			{
				Name: "containerspath",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: n.nodeAgent.FluentbitSpec.ContainersPath,
					},
				},
			},
			{
				Name: "varlogspath",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: n.nodeAgent.FluentbitSpec.VarLogsPath,
					},
				},
			},
		}

		for vCount, vMnt := range n.nodeAgent.FluentbitSpec.ExtraVolumeMounts {
			v = append(v, corev1.Volume{
				Name: "extravolumemount" + strconv.Itoa(vCount),
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: vMnt.Source,
					},
				}})
		}

		if n.nodeAgent.FluentbitSpec.CustomConfigSecret == "" {
			volume := corev1.Volume{
				Name: "config",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: n.QualifiedName(SecretConfigNameFluentbit),
						Items: []corev1.KeyToPath{
							{
								Key:  BaseConfigNameFluentbit,
								Path: BaseConfigNameFluentbit,
							},
						},
					},
				},
			}
			if util.PointerToBool(n.nodeAgent.FluentbitSpec.EnableUpstream) {
				volume.VolumeSource.Secret.Items = append(volume.VolumeSource.Secret.Items, corev1.KeyToPath{
					Key:  UpstreamConfigNameFluentbit,
					Path: UpstreamConfigNameFluentbit,
				})
			}
			v = append(v, volume)
		} else {
			v = append(v, corev1.Volume{
				Name: "config",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: n.nodeAgent.FluentbitSpec.CustomConfigSecret,
					},
				},
			})
		}
		if n.nodeAgent.FluentbitSpec.TLS != nil && util.PointerToBool(n.nodeAgent.FluentbitSpec.TLS.Enabled) {
			tlsRelatedVolume := corev1.Volume{
				Name: "fluent-bit-tls",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: n.nodeAgent.FluentbitSpec.TLS.SecretName,
					},
				},
			}
			v = append(v, tlsRelatedVolume)
		}

	} else if n.nodeAgent.SyslogNGSpec != nil {
		v = []corev1.Volume{
			{
				Name: "containerspath",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: n.nodeAgent.SyslogNGSpec.ContainersPath,
					},
				},
			},
			{
				Name: "varlogspath",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: n.nodeAgent.SyslogNGSpec.VarLogsPath,
					},
				},
			},
		}

		for vCount, vMnt := range n.nodeAgent.SyslogNGSpec.ExtraVolumeMounts {
			v = append(v, corev1.Volume{
				Name: "extravolumemount" + strconv.Itoa(vCount),
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: vMnt.Source,
					},
				}})
		}

		if n.nodeAgent.SyslogNGSpec.CustomConfigSecret == "" {
			volume := corev1.Volume{
				Name: "config",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: n.QualifiedName(secretConfigNameSyslogNG),
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
						SecretName: n.nodeAgent.SyslogNGSpec.CustomConfigSecret,
					},
				},
			})
		}
		if n.nodeAgent.SyslogNGSpec.TLS.Enabled {
			tlsRelatedVolume := corev1.Volume{
				Name: "syslog-ng-tls",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: n.nodeAgent.SyslogNGSpec.TLS.SecretName,
					},
				},
			}
			v = append(v, tlsRelatedVolume)
		}
	}

	return
}
