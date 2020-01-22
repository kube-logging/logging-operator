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
	spec := *r.statefulsetSpec()
	if !r.Logging.Spec.FluentdSpec.DisablePvc {
		spec.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{
			{
				ObjectMeta: r.FluentdObjectMeta(r.Logging.Spec.FluentdSpec.BufferStorageVolume.PersistentVolumeClaim.PersistentVolumeSource.ClaimName),
				Spec:       r.Logging.Spec.FluentdSpec.BufferStorageVolume.PersistentVolumeClaim.PersistentVolumeClaimSpec,
				Status: corev1.PersistentVolumeClaimStatus{
					Phase: corev1.ClaimPending,
				},
			},
		}
	}
	desired := &appsv1.StatefulSet{
		ObjectMeta: r.FluentdObjectMeta(StatefulSetName),
		Spec:       spec,
	}

	return desired, reconciler.StatePresent, nil
}

func (r *Reconciler) statefulsetSpec() *appsv1.StatefulSetSpec {
	initContainers := make([]corev1.Container, 0)

	if r.Logging.Spec.FluentdSpec.VolumeMountChmod {
		initContainers = append(initContainers, corev1.Container{
			Name:            "volume-mount-hack",
			Image:           r.Logging.Spec.FluentdSpec.VolumeModImage.Repository + ":" + r.Logging.Spec.FluentdSpec.VolumeModImage.Tag,
			ImagePullPolicy: corev1.PullPolicy(r.Logging.Spec.FluentdSpec.VolumeModImage.PullPolicy),
			Command:         []string{"sh", "-c", "chmod -R 777", BufferPath},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      r.Logging.QualifiedName(bufferVolumeName),
					MountPath: BufferPath,
				},
			},
		})
	}

	return &appsv1.StatefulSetSpec{
		Replicas: util.IntPointer(cast.ToInt32(r.Logging.Spec.FluentdSpec.Scaling.Replicas)),
		Selector: &metav1.LabelSelector{
			MatchLabels: r.getFluentdLabels(),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: r.generatePodMeta(),
			Spec: corev1.PodSpec{
				Volumes:            r.generateVolume(),
				ServiceAccountName: r.getServiceAccount(),
				InitContainers:     initContainers,
				Containers: []corev1.Container{
					*r.fluentContainer(),
					*newConfigMapReloader(r.Logging.Spec.FluentdSpec.ConfigReloaderImage),
				},
				NodeSelector:      r.Logging.Spec.FluentdSpec.NodeSelector,
				Tolerations:       r.Logging.Spec.FluentdSpec.Tolerations,
				PriorityClassName: r.Logging.Spec.FluentdSpec.PodPriorityClassName,
				SecurityContext: &corev1.PodSecurityContext{
					RunAsNonRoot: r.Logging.Spec.FluentdSpec.Security.PodSecurityContext.RunAsNonRoot,
					FSGroup:      r.Logging.Spec.FluentdSpec.Security.PodSecurityContext.FSGroup,
					RunAsUser:    r.Logging.Spec.FluentdSpec.Security.PodSecurityContext.RunAsUser,
					RunAsGroup:   r.Logging.Spec.FluentdSpec.Security.PodSecurityContext.RunAsGroup},
			},
		},
	}
}

func (r *Reconciler) fluentContainer() *corev1.Container {
	container := &corev1.Container{
		Name:            "fluentd",
		Image:           r.Logging.Spec.FluentdSpec.Image.Repository + ":" + r.Logging.Spec.FluentdSpec.Image.Tag,
		ImagePullPolicy: corev1.PullPolicy(r.Logging.Spec.FluentdSpec.Image.PullPolicy),
		Ports:           generatePorts(r.Logging.Spec.FluentdSpec),
		VolumeMounts:    r.generateVolumeMounts(),
		Resources:       r.Logging.Spec.FluentdSpec.Resources,
		Env: []corev1.EnvVar{
			{
				Name:  "BUFFER_PATH",
				Value: BufferPath,
			},
		},
		SecurityContext: &corev1.SecurityContext{
			RunAsUser:                r.Logging.Spec.FluentdSpec.Security.SecurityContext.RunAsUser,
			RunAsGroup:               r.Logging.Spec.FluentdSpec.Security.SecurityContext.RunAsGroup,
			ReadOnlyRootFilesystem:   r.Logging.Spec.FluentdSpec.Security.SecurityContext.ReadOnlyRootFilesystem,
			AllowPrivilegeEscalation: r.Logging.Spec.FluentdSpec.Security.SecurityContext.AllowPrivilegeEscalation,
			Privileged:               r.Logging.Spec.FluentdSpec.Security.SecurityContext.Privileged,
			RunAsNonRoot:             r.Logging.Spec.FluentdSpec.Security.SecurityContext.RunAsNonRoot,
		},
		LivenessProbe:  r.Logging.Spec.FluentdSpec.LivenessProbe,
		ReadinessProbe: r.Logging.Spec.FluentdSpec.ReadinessProbe,
	}

	if r.Logging.Spec.FluentdSpec.FluentOutLogrotate != nil && r.Logging.Spec.FluentdSpec.FluentOutLogrotate.Enabled {
		container.Args = []string{
			"fluentd",
			"-o", r.Logging.Spec.FluentdSpec.FluentOutLogrotate.Path,
			"--log-rotate-age", r.Logging.Spec.FluentdSpec.FluentOutLogrotate.Age,
			"--log-rotate-size", r.Logging.Spec.FluentdSpec.FluentOutLogrotate.Size,
		}
	}

	return container
}

func (r *Reconciler) generatePodMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Labels: r.getFluentdLabels(),
	}
	if r.Logging.Spec.FluentdSpec.Annotations != nil {
		meta.Annotations = r.Logging.Spec.FluentdSpec.Annotations
	}
	return meta
}

func generateLoggingRefLabels(loggingRef string) map[string]string {
	return map[string]string{"app.kubernetes.io/managed-by": loggingRef}
}

func newConfigMapReloader(spec v1beta1.ImageSpec) *corev1.Container {
	return &corev1.Container{
		Name:            "config-reloader",
		ImagePullPolicy: corev1.PullPolicy(spec.PullPolicy),
		Image:           spec.Repository + ":" + spec.Tag,
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

func (r *Reconciler) generateVolumeMounts() (v []corev1.VolumeMount) {
	v = []corev1.VolumeMount{
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
		{
			Name:      r.Logging.QualifiedName(bufferVolumeName),
			MountPath: BufferPath,
		},
	}
	if r.Logging.Spec.FluentdSpec.TLS.Enabled {
		tlsRelatedVolume := []corev1.VolumeMount{
			{
				Name:      "fluentd-tls",
				MountPath: "/fluentd/tls/",
			},
		}
		v = append(v, tlsRelatedVolume...)
	}
	return
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
	v = append(v,
		r.Logging.Spec.FluentdSpec.BufferStorageVolume.GetVolume(
			r.Logging.QualifiedName(bufferVolumeName),
			fmt.Sprintf(v1beta1.HostPath, r.Logging.Name, r.Logging.QualifiedName(bufferVolumeName)),
		))
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
