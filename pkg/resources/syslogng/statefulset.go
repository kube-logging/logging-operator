// Copyright © 2022 Banzai Cloud
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

package syslogng

import (
	"fmt"
	"strings"

	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	util "github.com/banzaicloud/operator-tools/pkg/utils"
	"github.com/spf13/cast"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *Reconciler) statefulset() (runtime.Object, reconciler.DesiredState, error) {
	spec := r.statefulsetSpec()

	r.Logging.Spec.SyslogNGSpec.BufferStorageVolume.WithDefaultHostPath(
		fmt.Sprintf(v1beta1.HostPath, r.Logging.Name, r.Logging.QualifiedName(v1beta1.DefaultSyslogNGBufferStorageVolumeName)),
	)
	if !r.Logging.Spec.SyslogNGSpec.DisablePvc {
		err := r.Logging.Spec.SyslogNGSpec.BufferStorageVolume.ApplyPVCForStatefulSet(containerName, bufferPath, spec, func(name string) metav1.ObjectMeta {
			return r.SyslogNGObjectMeta(name, ComponentSyslogNG)
		})
		if err != nil {
			return nil, reconciler.StatePresent, err
		}
	} else {
		err := r.Logging.Spec.SyslogNGSpec.BufferStorageVolume.ApplyVolumeForPodSpec(r.Logging.QualifiedName(v1beta1.DefaultSyslogNGBufferStorageVolumeName), containerName, bufferPath, &spec.Template.Spec)
		if err != nil {
			return nil, reconciler.StatePresent, err
		}
	}
	for _, n := range r.Logging.Spec.SyslogNGSpec.ExtraVolumes {
		if err := n.ApplyVolumeForPodSpec(&spec.Template.Spec); err != nil {
			return nil, reconciler.StatePresent, err
		}
	}

	desired := &appsv1.StatefulSet{
		ObjectMeta: r.SyslogNGObjectMeta(StatefulSetName, ComponentSyslogNG),
		Spec:       *spec,
	}

	desired.Annotations = util.MergeLabels(desired.Annotations, r.Logging.Spec.SyslogNGSpec.StatefulSetAnnotations)

	return desired, reconciler.StatePresent, nil
}

func (r *Reconciler) statefulsetSpec() *appsv1.StatefulSetSpec {
	var initContainers []corev1.Container
	if c := r.volumeMountHackContainer(); c != nil {
		initContainers = append(initContainers, *c)
	}

	containers := []corev1.Container{
		syslogNGContainer(r.Logging.Spec.SyslogNGSpec),
	}
	if c := r.bufferMetricsSidecarContainer(); c != nil {
		containers = append(containers, *c)
	}

	if c := r.syslogNGMetricsSidecarContainer(); c != nil {
		containers = append(containers, *c)
	}

	sts := &appsv1.StatefulSetSpec{
		PodManagementPolicy: appsv1.PodManagementPolicyType(r.Logging.Spec.SyslogNGSpec.Scaling.PodManagementPolicy),
		Selector: &metav1.LabelSelector{
			MatchLabels: r.Logging.GetSyslogNGLabels(ComponentSyslogNG),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: r.generatePodMeta(),
			Spec: corev1.PodSpec{
				Volumes:                   r.generateVolume(),
				ServiceAccountName:        r.getServiceAccount(),
				InitContainers:            initContainers,
				ImagePullSecrets:          r.Logging.Spec.SyslogNGSpec.Image.ImagePullSecrets,
				Containers:                containers,
				NodeSelector:              r.Logging.Spec.SyslogNGSpec.NodeSelector,
				Tolerations:               r.Logging.Spec.SyslogNGSpec.Tolerations,
				Affinity:                  r.Logging.Spec.SyslogNGSpec.Affinity,
				TopologySpreadConstraints: r.Logging.Spec.SyslogNGSpec.TopologySpreadConstraints,
				PriorityClassName:         r.Logging.Spec.SyslogNGSpec.PodPriorityClassName,
				DNSPolicy:                 r.Logging.Spec.SyslogNGSpec.DNSPolicy,
				DNSConfig:                 r.Logging.Spec.SyslogNGSpec.DNSConfig,
				SecurityContext: &corev1.PodSecurityContext{
					RunAsNonRoot: r.Logging.Spec.SyslogNGSpec.Security.PodSecurityContext.RunAsNonRoot,
					FSGroup:      r.Logging.Spec.SyslogNGSpec.Security.PodSecurityContext.FSGroup,
					RunAsUser:    r.Logging.Spec.SyslogNGSpec.Security.PodSecurityContext.RunAsUser,
					RunAsGroup:   r.Logging.Spec.SyslogNGSpec.Security.PodSecurityContext.RunAsGroup},
			},
		},
		ServiceName: r.Logging.QualifiedName(ServiceName + "-headless"),
	}

	if r.Logging.Spec.SyslogNGSpec.Scaling.Replicas > 0 {
		sts.Replicas = util.IntPointer(cast.ToInt32(r.Logging.Spec.SyslogNGSpec.Scaling.Replicas))
	}

	return sts
}

func syslogNGContainer(spec *v1beta1.SyslogNGSpec) corev1.Container {
	envVars := append(spec.EnvVars,
		corev1.EnvVar{Name: "BUFFER_PATH", Value: bufferPath},
	)

	container := corev1.Container{
		Name:            "syslog-ng",
		Image:           spec.Image.RepositoryWithTag(),
		ImagePullPolicy: corev1.PullPolicy(spec.Image.PullPolicy),
		Ports:           generatePorts(spec),
		VolumeMounts:    generateVolumeMounts(spec),
		Resources:       spec.Resources,
		SecurityContext: &corev1.SecurityContext{
			RunAsUser:                spec.Security.SecurityContext.RunAsUser,
			RunAsGroup:               spec.Security.SecurityContext.RunAsGroup,
			ReadOnlyRootFilesystem:   spec.Security.SecurityContext.ReadOnlyRootFilesystem,
			AllowPrivilegeEscalation: spec.Security.SecurityContext.AllowPrivilegeEscalation,
			Privileged:               spec.Security.SecurityContext.Privileged,
			RunAsNonRoot:             spec.Security.SecurityContext.RunAsNonRoot,
			SELinuxOptions:           spec.Security.SecurityContext.SELinuxOptions,
		},
		Env:            envVars,
		LivenessProbe:  spec.LivenessProbe,
		ReadinessProbe: generateReadinessCheck(spec),
	}

	if spec.SyslogNGOutLogrotate != nil && spec.SyslogNGOutLogrotate.Enabled {
		container.Args = []string{
			"--control=/tmp/syslog-ng/syslog-ng.ctl",
		}
	}

	return container
}

func (r *Reconciler) generatePodMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Labels: r.Logging.GetSyslogNGLabels(ComponentSyslogNG),
	}
	if r.Logging.Spec.SyslogNGSpec.Annotations != nil {
		meta.Annotations = r.Logging.Spec.SyslogNGSpec.Annotations
	}
	return meta
}

func generatePortsBufferVolumeMetrics(spec *v1beta1.SyslogNGSpec) []corev1.ContainerPort {
	port := int32(defaultBufferVolumeMetricsPort)
	if spec.Metrics != nil && spec.BufferVolumeMetrics.Port != 0 {
		port = spec.BufferVolumeMetrics.Port
	}
	return []corev1.ContainerPort{
		{
			Name:          "buffer-metrics",
			ContainerPort: port,
			Protocol:      "TCP",
		},
	}
}

func generatePorts(spec *v1beta1.SyslogNGSpec) []corev1.ContainerPort {
	ports := []corev1.ContainerPort{}
	if spec.PortTCP != nil {
		ports = append(ports, corev1.ContainerPort{
			Name:          "syslog-ng-tcp",
			ContainerPort: *spec.PortTCP,
			Protocol:      "TCP",
		})
	}
	if spec.PortUDP != nil {
		ports = append(ports, corev1.ContainerPort{
			Name:          "syslog-ng-udp",
			ContainerPort: *spec.PortUDP,
			Protocol:      "UDP",
		})
	}
	return ports
}

func generateVolumeMounts(spec *v1beta1.SyslogNGSpec) []corev1.VolumeMount {
	res := []corev1.VolumeMount{
		{
			Name:      "config",
			MountPath: "/etc/syslog-ng/",
		},
	}
	if spec != nil && spec.TLS.Enabled {
		res = append(res, corev1.VolumeMount{
			Name:      "syslog-ng-tls",
			MountPath: "/syslog-ng/tls/",
		})
	}
	if spec != nil && spec.Metrics != nil {
		res = append(res, corev1.VolumeMount{
			Name:      "syslog-ng-socket",
			MountPath: "/tmp/syslog-ng",
		})
	}

	return res
}

func (r *Reconciler) generateVolume() (v []corev1.Volume) {
	v = []corev1.Volume{
		{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: r.Logging.QualifiedName(SecretConfigName),
				},
			},
		},
		{
			Name: "app-config",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: r.Logging.QualifiedName(AppSecretConfigName),
				},
			},
		},
		{
			Name: "output-secret",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: r.Logging.QualifiedName(OutputSecretName),
				},
			},
		},
	}
	if r.Logging.Spec.SyslogNGSpec.TLS.Enabled {
		tlsRelatedVolume := corev1.Volume{
			Name: "syslog-ng-tls",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: r.Logging.Spec.SyslogNGSpec.TLS.SecretName,
				},
			},
		}
		v = append(v, tlsRelatedVolume)
	}
	if r.Logging.Spec.SyslogNGSpec.Metrics != nil {
		socketVolume := corev1.Volume{
			Name: "syslog-ng-socket",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{
					Medium:    corev1.StorageMediumDefault,
					SizeLimit: nil,
				},
			},
		}
		v = append(v, socketVolume)
	}

	return
}

func (r *Reconciler) volumeMountHackContainer() *corev1.Container {
	if r.Logging.Spec.SyslogNGSpec.VolumeMountChmod {
		return &corev1.Container{
			Name:            "volume-mount-hack",
			Image:           r.Logging.Spec.SyslogNGSpec.VolumeModImage.RepositoryWithTag(),
			ImagePullPolicy: corev1.PullPolicy(r.Logging.Spec.SyslogNGSpec.VolumeModImage.PullPolicy),
			Command:         []string{"sh", "-c", "chmod -R 777 " + bufferPath},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      r.Logging.QualifiedName(v1beta1.DefaultSyslogNGBufferStorageVolumeName),
					MountPath: bufferPath,
				},
			},
		}
	}
	return nil
}

func (r *Reconciler) syslogNGMetricsSidecarContainer() *corev1.Container {
	if r.Logging.Spec.SyslogNGSpec.Metrics != nil {
		return &corev1.Container{
			Name:            "exporter",
			ImagePullPolicy: corev1.PullPolicy(r.Logging.Spec.SyslogNGSpec.PrometheusExporterImage.PullPolicy),
			Image:           r.Logging.Spec.SyslogNGSpec.PrometheusExporterImage.RepositoryWithTag(),
			Resources:       r.Logging.Spec.SyslogNGSpec.PrometheusExporterResources,
			Ports: []corev1.ContainerPort{
				{
					Name:          "exporter",
					ContainerPort: r.Logging.Spec.SyslogNGSpec.Metrics.Port,
					Protocol:      "TCP",
				},
			},
			Args: []string{
				"--socket.path=/tmp/syslog-ng.ctl",
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "syslog-ng-socket",
					MountPath: "/tmp",
				},
			},
		}

	}
	return nil
}

func (r *Reconciler) bufferMetricsSidecarContainer() *corev1.Container {
	if r.Logging.Spec.SyslogNGSpec.BufferVolumeMetrics != nil {
		port := int32(defaultBufferVolumeMetricsPort)
		if r.Logging.Spec.SyslogNGSpec.BufferVolumeMetrics.Port != 0 {
			port = r.Logging.Spec.SyslogNGSpec.BufferVolumeMetrics.Port
		}
		portParam := fmt.Sprintf("--web.listen-address=:%d", port)
		args := []string{portParam}
		if len(r.Logging.Spec.SyslogNGSpec.BufferVolumeArgs) != 0 {
			args = append(args, r.Logging.Spec.SyslogNGSpec.BufferVolumeArgs...)
		} else {
			args = append(args, "--collector.disable-defaults", "--collector.filesystem")
		}
		customRunner := fmt.Sprintf("./bin/node_exporter %v", strings.Join(args, " "))
		return &corev1.Container{
			Name:            "buffer-metrics-sidecar",
			Image:           r.Logging.Spec.SyslogNGSpec.BufferVolumeImage.RepositoryWithTag(),
			ImagePullPolicy: corev1.PullPolicy(r.Logging.Spec.SyslogNGSpec.BufferVolumeImage.PullPolicy),
			Args:            []string{"--startup", customRunner},
			Ports:           generatePortsBufferVolumeMetrics(r.Logging.Spec.SyslogNGSpec),
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      r.Logging.QualifiedName(v1beta1.DefaultSyslogNGBufferStorageVolumeName),
					MountPath: bufferPath,
				},
			},
		}
	}
	return nil
}

func generateReadinessCheck(spec *v1beta1.SyslogNGSpec) *corev1.Probe {
	if spec.ReadinessProbe != nil {
		return spec.ReadinessProbe
	}

	if spec.ReadinessDefaultCheck.BufferFreeSpace || spec.ReadinessDefaultCheck.BufferFileNumber {
		check := []string{"/bin/sh", "-c"}
		bash := []string{}
		if spec.ReadinessDefaultCheck.BufferFreeSpace {
			bash = append(bash,
				fmt.Sprintf("FREESPACE_THRESHOLD=%d", spec.ReadinessDefaultCheck.BufferFreeSpaceThreshold),
				"FREESPACE_CURRENT=$(df -h $BUFFER_PATH  | grep / | awk '{ print $5}' | sed 's/%//g')",
				"if [ \"$FREESPACE_CURRENT\" -gt \"$FREESPACE_THRESHOLD\" ] ; then exit 1; fi",
			)
		}
		if spec.ReadinessDefaultCheck.BufferFileNumber {
			bash = append(bash,
				fmt.Sprintf("MAX_FILE_NUMBER=%d", spec.ReadinessDefaultCheck.BufferFileNumberMax),
				"FILE_NUMBER_CURRENT=$(find $BUFFER_PATH -type f -name *.buffer | wc -l)",
				"if [ \"$FILE_NUMBER_CURRENT\" -gt \"$MAX_FILE_NUMBER\" ] ; then exit 1; fi",
			)
		}
		return &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				Exec: &corev1.ExecAction{
					Command: append(check, strings.Join(bash, "\n")),
				},
			},
			InitialDelaySeconds: spec.ReadinessDefaultCheck.InitialDelaySeconds,
			TimeoutSeconds:      spec.ReadinessDefaultCheck.TimeoutSeconds,
			PeriodSeconds:       spec.ReadinessDefaultCheck.PeriodSeconds,
			SuccessThreshold:    spec.ReadinessDefaultCheck.SuccessThreshold,
			FailureThreshold:    spec.ReadinessDefaultCheck.FailureThreshold,
		}
	}
	return nil
}
