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

package fluentbit

import (
	"github.com/banzaicloud/logging-operator/pkg/k8sutil"
	"github.com/banzaicloud/logging-operator/pkg/resources/templates"
	"github.com/banzaicloud/logging-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// TODO in case of rbac add created serviceAccount name
func (r *Reconciler) daemonSet() (runtime.Object, k8sutil.DesiredState) {

	var containerPorts []corev1.ContainerPort

	if r.Logging.Spec.FluentbitSpec.Metrics != nil && r.Logging.Spec.FluentbitSpec.Metrics.Port != 0 {
		containerPorts = append(containerPorts, corev1.ContainerPort{
			Name:          "monitor",
			ContainerPort: r.Logging.Spec.FluentbitSpec.Metrics.Port,
			Protocol:      corev1.ProtocolTCP,
		})
	}

	labels := util.MergeLabels(r.Logging.Labels, r.getFluentBitLabels())

	return &appsv1.DaemonSet{
		ObjectMeta: templates.FluentbitObjectMeta(
			r.Logging.QualifiedName(fluentbitDaemonSetName), labels, r.Logging),
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: util.MergeLabels(r.Logging.Labels, r.getFluentBitLabels())},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: r.Logging.Spec.FluentbitSpec.Annotations,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: r.Logging.QualifiedName(serviceAccountName),
					Volumes:            r.generateVolume(),
					Tolerations:        r.Logging.Spec.FluentbitSpec.Tolerations,
					Containers: []corev1.Container{
						{
							Name:            "fluent-bit",
							Image:           r.Logging.Spec.FluentbitSpec.Image.Repository + ":" + r.Logging.Spec.FluentbitSpec.Image.Tag,
							ImagePullPolicy: corev1.PullPolicy(r.Logging.Spec.FluentbitSpec.Image.PullPolicy),
							Ports:           containerPorts,
							Resources:       r.Logging.Spec.FluentbitSpec.Resources,
							VolumeMounts:    r.generateVolumeMounts(),
						},
					},
				},
			},
		},
	}, k8sutil.StatePresent
}

func (r *Reconciler) generateVolumeMounts() (v []corev1.VolumeMount) {
	v = []corev1.VolumeMount{
		{
			Name:      "varlibcontainers",
			ReadOnly:  true,
			MountPath: "/var/lib/docker/containers",
		},
		{
			Name:      "config",
			MountPath: "/fluent-bit/etc/fluent-bit.conf",
			SubPath:   "fluent-bit.conf",
		},
		{
			Name:      "positions",
			MountPath: "/tail-db",
		},
		{
			Name:      "varlogs",
			ReadOnly:  true,
			MountPath: "/var/log/",
		},
	}
	if r.Logging.Spec.FluentbitSpec.TLS.Enabled {
		tlsRelatedVolume := []corev1.VolumeMount{
			{
				Name:      "fluent-bit-tls",
				MountPath: "/fluent-bit/tls/",
			},
		}
		v = append(v, tlsRelatedVolume...)
	}
	return
}

func (r *Reconciler) generateVolume() (v []corev1.Volume) {
	v = []corev1.Volume{
		{
			Name: "varlibcontainers",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/var/lib/docker/containers",
				},
			},
		},
		{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: r.Logging.QualifiedName(fluentBitSecretConfigName),
					Items: []corev1.KeyToPath{
						{
							Key:  "fluent-bit.conf",
							Path: "fluent-bit.conf",
						},
					},
				},
			},
		},
		{
			Name: "varlogs",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/var/log",
				},
			},
		},
		{
			Name: "positions",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
	}
	if r.Logging.Spec.FluentbitSpec.TLS.Enabled {
		tlsRelatedVolume := corev1.Volume{
			Name: "fluent-bit-tls",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: r.Logging.Spec.FluentbitSpec.TLS.SecretName,
				},
			},
		}
		v = append(v, tlsRelatedVolume)
	}
	return
}
