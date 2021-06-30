// Copyright © 2019 Banzai Cloud
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

package fluentd

import (
	"fmt"

	"github.com/banzaicloud/logging-operator/pkg/sdk/api/v1beta1"
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

	r.Logging.Spec.FluentdSpec.BufferStorageVolume.WithDefaultHostPath(
		fmt.Sprintf(v1beta1.HostPath, r.Logging.Name, r.Logging.QualifiedName(v1beta1.DefaultFluentdBufferStorageVolumeName)),
	)
	if !r.Logging.Spec.FluentdSpec.DisablePvc {
		err := r.Logging.Spec.FluentdSpec.BufferStorageVolume.ApplyPVCForStatefulSet(containerName, bufferPath, spec, func(name string) metav1.ObjectMeta {
			return r.FluentdObjectMeta(name, ComponentFluentd)
		})
		if err != nil {
			return nil, reconciler.StatePresent, err
		}
	} else {
		err := r.Logging.Spec.FluentdSpec.BufferStorageVolume.ApplyVolumeForPodSpec(v1beta1.DefaultFluentdBufferStorageVolumeName, containerName, bufferPath, &spec.Template.Spec)
		if err != nil {
			return nil, reconciler.StatePresent, err
		}
	}

	desired := &appsv1.StatefulSet{
		ObjectMeta: r.FluentdObjectMeta(StatefulSetName, ComponentFluentd),
		Spec:       *spec,
	}

	return desired, reconciler.StatePresent, nil
}

func (r *Reconciler) statefulsetSpec() *appsv1.StatefulSetSpec {
	var initContainers []corev1.Container
	if c := r.volumeMountHackContainer(); c != nil {
		initContainers = append(initContainers, *c)
	}

	containers := []corev1.Container{
		fluentContainer(r.Logging.Spec.FluentdSpec),
		*newConfigMapReloader(r.Logging.Spec.FluentdSpec),
	}
	if c := r.bufferMetricsSidecarContainer(); c != nil {
		containers = append(containers, *c)
	}

	return &appsv1.StatefulSetSpec{
		Replicas:            util.IntPointer(cast.ToInt32(r.Logging.Spec.FluentdSpec.Scaling.Replicas)),
		PodManagementPolicy: appsv1.PodManagementPolicyType(r.Logging.Spec.FluentdSpec.Scaling.PodManagementPolicy),
		Selector: &metav1.LabelSelector{
			MatchLabels: r.getFluentdLabels(ComponentFluentd),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: r.generatePodMeta(),
			Spec: corev1.PodSpec{
				Volumes:            r.generateVolume(),
				ServiceAccountName: r.getServiceAccount(),
				InitContainers:     initContainers,
				ImagePullSecrets:   r.Logging.Spec.FluentdSpec.Image.ImagePullSecrets,
				Containers:         containers,
				NodeSelector:       r.Logging.Spec.FluentdSpec.NodeSelector,
				Tolerations:        r.Logging.Spec.FluentdSpec.Tolerations,
				Affinity:           r.Logging.Spec.FluentdSpec.Affinity,
				PriorityClassName:  r.Logging.Spec.FluentdSpec.PodPriorityClassName,
				SecurityContext: &corev1.PodSecurityContext{
					RunAsNonRoot: r.Logging.Spec.FluentdSpec.Security.PodSecurityContext.RunAsNonRoot,
					FSGroup:      r.Logging.Spec.FluentdSpec.Security.PodSecurityContext.FSGroup,
					RunAsUser:    r.Logging.Spec.FluentdSpec.Security.PodSecurityContext.RunAsUser,
					RunAsGroup:   r.Logging.Spec.FluentdSpec.Security.PodSecurityContext.RunAsGroup},
			},
		},
		ServiceName: r.Logging.QualifiedName(ServiceName + "-headless"),
	}
}

func fluentContainer(spec *v1beta1.FluentdSpec) corev1.Container {
	envVars := append(spec.EnvVars,
		corev1.EnvVar{Name: "BUFFER_PATH", Value: bufferPath},
	)

	container := corev1.Container{
		Name:            "fluentd",
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
		ReadinessProbe: spec.ReadinessProbe,
	}

	if spec.FluentOutLogrotate != nil && spec.FluentOutLogrotate.Enabled {
		container.Args = []string{
			"fluentd",
			"-o", spec.FluentOutLogrotate.Path,
			"--log-rotate-age", spec.FluentOutLogrotate.Age,
			"--log-rotate-size", spec.FluentOutLogrotate.Size,
		}
	}

	return container
}

func (r *Reconciler) generatePodMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Labels: r.getFluentdLabels(ComponentFluentd),
	}
	if r.Logging.Spec.FluentdSpec.Annotations != nil {
		meta.Annotations = r.Logging.Spec.FluentdSpec.Annotations
	}
	return meta
}

func generateLoggingRefLabels(loggingRef string) map[string]string {
	return map[string]string{"app.kubernetes.io/managed-by": loggingRef}
}

func newConfigMapReloader(spec *v1beta1.FluentdSpec) *corev1.Container {
	return &corev1.Container{
		Name:            "config-reloader",
		ImagePullPolicy: corev1.PullPolicy(spec.ConfigReloaderImage.PullPolicy),
		Image:           spec.ConfigReloaderImage.RepositoryWithTag(),
		Resources:       spec.ConfigReloaderResources,
		Args: []string{
			"-volume-dir=/fluentd/etc",
			"-volume-dir=/fluentd/app-config/",
			"-webhook-url=http://127.0.0.1:24444/api/config.reload",
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "config",
				MountPath: "/fluentd/etc",
			},
			{
				Name:      "app-config",
				MountPath: "/fluentd/app-config/",
			},
		},
	}
}

func generatePortsBufferVolumeMetrics(spec *v1beta1.FluentdSpec) []corev1.ContainerPort {
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

func generatePorts(spec *v1beta1.FluentdSpec) []corev1.ContainerPort {
	ports := []corev1.ContainerPort{
		{
			Name:          "fluent-input",
			ContainerPort: spec.Port,
			Protocol:      "TCP",
		},
	}
	if spec.Metrics != nil && spec.Metrics.Port != 0 {
		ports = append(ports, corev1.ContainerPort{
			Name:          "monitor",
			ContainerPort: spec.Metrics.Port,
			Protocol:      "TCP",
		})
	}
	return ports
}

func generateVolumeMounts(spec *v1beta1.FluentdSpec) []corev1.VolumeMount {
	res := []corev1.VolumeMount{
		{
			Name:      "config",
			MountPath: "/fluentd/etc/",
		},
		{
			Name:      "app-config",
			MountPath: "/fluentd/app-config/",
		},
		{
			Name:      "output-secret",
			MountPath: OutputSecretPath,
		},
	}
	if spec != nil && spec.TLS.Enabled {
		res = append(res, corev1.VolumeMount{
			Name:      "fluentd-tls",
			MountPath: "/fluentd/tls/",
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
	if r.Logging.Spec.FluentdSpec.TLS.Enabled {
		tlsRelatedVolume := corev1.Volume{
			Name: "fluentd-tls",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: r.Logging.Spec.FluentdSpec.TLS.SecretName,
				},
			},
		}
		v = append(v, tlsRelatedVolume)
	}
	return
}

func (r *Reconciler) volumeMountHackContainer() *corev1.Container {
	if r.Logging.Spec.FluentdSpec.VolumeMountChmod {
		return &corev1.Container{
			Name:            "volume-mount-hack",
			Image:           r.Logging.Spec.FluentdSpec.VolumeModImage.RepositoryWithTag(),
			ImagePullPolicy: corev1.PullPolicy(r.Logging.Spec.FluentdSpec.VolumeModImage.PullPolicy),
			Command:         []string{"sh", "-c", "chmod -R 777 " + bufferPath},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      r.Logging.QualifiedName(v1beta1.DefaultFluentdBufferStorageVolumeName),
					MountPath: bufferPath,
				},
			},
		}
	}
	return nil
}

func (r *Reconciler) bufferMetricsSidecarContainer() *corev1.Container {
	if r.Logging.Spec.FluentdSpec.BufferVolumeMetrics != nil {
		port := int32(defaultBufferVolumeMetricsPort)
		if r.Logging.Spec.FluentdSpec.BufferVolumeMetrics.Port != 0 {
			port = r.Logging.Spec.FluentdSpec.BufferVolumeMetrics.Port
		}
		portParam := fmt.Sprintf("--web.listen-address=:%d", port)
		args := []string{portParam}
		if len(r.Logging.Spec.FluentdSpec.BufferVolumeArgs) != 0 {
			args = append(args, r.Logging.Spec.FluentdSpec.BufferVolumeArgs...)
		} else {
			args = append(args, "--collector.disable-defaults", "--collector.filesystem")
		}
		return &corev1.Container{
			Name:            "buffer-metrics-sidecar",
			Image:           r.Logging.Spec.FluentdSpec.BufferVolumeImage.RepositoryWithTag(),
			ImagePullPolicy: corev1.PullPolicy(r.Logging.Spec.FluentdSpec.BufferVolumeImage.PullPolicy),
			Args:            args,
			Ports:           generatePortsBufferVolumeMetrics(r.Logging.Spec.FluentdSpec),
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      r.Logging.QualifiedName(v1beta1.DefaultFluentdBufferStorageVolumeName),
					MountPath: bufferPath,
				},
			},
		}
	}
	return nil
}
