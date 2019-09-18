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
	"github.com/banzaicloud/logging-operator/api/v1alpha2"
	"github.com/banzaicloud/logging-operator/pkg/resources/templates"
	"github.com/banzaicloud/logging-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *Reconciler) statefulset() runtime.Object {
	spec := *r.statefulsetSpec()
	if !r.Logging.Spec.FluentdSpec.DisablePvc {
		spec.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{
			{
				ObjectMeta: templates.FluentdObjectMeta(
					r.Logging.QualifiedName(bufferVolumeName), util.MergeLabels(r.Logging.Labels, labelSelector), r.Logging),
				Spec: r.Logging.Spec.FluentdSpec.FluentdPvcSpec,
			},
		}
	}
	return &appsv1.StatefulSet{
		ObjectMeta: templates.FluentdObjectMeta(
			r.Logging.QualifiedName(StatefulSetName), util.MergeLabels(r.Logging.Labels, labelSelector), r.Logging),
		Spec: spec,
	}
}

func (r *Reconciler) statefulsetSpec() *appsv1.StatefulSetSpec {
	return &appsv1.StatefulSetSpec{
		Replicas: util.IntPointer(1),
		Selector: &metav1.LabelSelector{
			MatchLabels: labelSelector,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: r.generatePodMeta(),
			Spec: corev1.PodSpec{
				Volumes: r.generateVolume(),
				InitContainers: []corev1.Container{
					{
						Name:            "volume-mount-hack",
						Image:           r.Logging.Spec.FluentdSpec.VolumeModImage.Repository + ":" + r.Logging.Spec.FluentdSpec.VolumeModImage.Tag,
						ImagePullPolicy: corev1.PullPolicy(r.Logging.Spec.FluentdSpec.VolumeModImage.PullPolicy),
						Command:         []string{"sh", "-c", "chmod -R 777 /buffers"},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      r.Logging.QualifiedName(bufferVolumeName),
								MountPath: "/buffers",
							},
						},
					},
				},
				Containers: []corev1.Container{
					*r.fluentContainer(),
					*newConfigMapReloader(r.Logging.Spec.FluentdSpec.ConfigReloaderImage),
				},
				NodeSelector: r.Logging.Spec.FluentdSpec.NodeSelector,
				Tolerations:  r.Logging.Spec.FluentdSpec.Tolerations,
			},
		},
	}
}

func (r *Reconciler) fluentContainer() *corev1.Container {
	return &corev1.Container{
		Name:            "fluentd",
		Image:           r.Logging.Spec.FluentdSpec.Image.Repository + ":" + r.Logging.Spec.FluentdSpec.Image.Tag,
		ImagePullPolicy: corev1.PullPolicy(r.Logging.Spec.FluentdSpec.Image.PullPolicy),
		Ports:           generatePorts(r.Logging.Spec.FluentdSpec),
		VolumeMounts:    r.generateVolumeMounts(),
		Resources:       r.Logging.Spec.FluentdSpec.Resources,
	}
}

func (r *Reconciler) generatePodMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Labels: util.MergeLabels(r.Logging.Labels, labelSelector),
	}
	if r.Logging.Spec.FluentdSpec.Annotations != nil {
		meta.Annotations = r.Logging.Spec.FluentdSpec.Annotations
	}
	return meta
}

func newConfigMapReloader(spec v1alpha2.ImageSpec) *corev1.Container {
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

func generatePorts(spec *v1alpha2.FluentdSpec) []corev1.ContainerPort {
	ports := []corev1.ContainerPort{
		{
			Name:          "fluent-input",
			ContainerPort: spec.Port,
			Protocol:      "TCP",
		},
	}
	if spec.GetPrometheusPortFromAnnotation() != 0 {
		ports = append(ports, corev1.ContainerPort{
			Name:          "monitor",
			ContainerPort: spec.GetPrometheusPortFromAnnotation(),
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
			Name:      r.Logging.QualifiedName(bufferVolumeName),
			MountPath: "/buffers",
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
	}
	if !r.Logging.Spec.FluentdSpec.DisablePvc {
		bufferVolume := corev1.Volume{
			Name: r.Logging.QualifiedName(bufferVolumeName),
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: r.Logging.QualifiedName(bufferVolumeName),
					ReadOnly:  false,
				},
			},
		}
		v = append(v, bufferVolume)
	} else {
		bufferVolume := corev1.Volume{
			Name: r.Logging.QualifiedName(bufferVolumeName),
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		}
		v = append(v, bufferVolume)
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
