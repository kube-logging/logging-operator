/*
 * Copyright © 2019 Banzai Cloud
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package fluentbit

import (
	loggingv1alpha1 "github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1"
	"github.com/banzaicloud/logging-operator/pkg/resources/templates"
	"github.com/banzaicloud/logging-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// TODO in case of rbac add created serviceAccount name
func (r *Reconciler) daemonSet() runtime.Object {

	var containerPorts []corev1.ContainerPort

	if _, ok := r.Fluentbit.Spec.Annotations["prometheus.io/port"]; ok {
		containerPorts = append(containerPorts, corev1.ContainerPort{
			Name:          "monitor",
			ContainerPort: r.Fluentbit.Spec.GetPrometheusPortFromAnnotation(),
			Protocol:      corev1.ProtocolTCP,
		})
	}

	return &appsv1.DaemonSet{
		ObjectMeta: templates.FluentbitObjectMeta(fluentbitDeaemonSetName, util.MergeLabels(r.Fluentbit.Labels, labelSelector), r.Fluentbit),
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: util.MergeLabels(r.Fluentbit.Labels, labelSelector)},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      util.MergeLabels(r.Fluentbit.Labels, labelSelector),
					Annotations: r.Fluentbit.Spec.Annotations,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: serviceAccountName,
					Volumes:            generateVolume(r.Fluentbit),
					Containers: []corev1.Container{
						{
							Name:            "fluent-bit",
							Image:           r.Fluentbit.Spec.Image.Repository + ":" + r.Fluentbit.Spec.Image.Tag,
							ImagePullPolicy: corev1.PullPolicy(r.Fluentbit.Spec.Image.PullPolicy),
							Ports:           containerPorts,
							Resources:       r.Fluentbit.Spec.Resources,
							VolumeMounts:    generateVolumeMounts(r.Fluentbit),
						},
					},
					Tolerations: r.Fluentbit.Spec.Tolerations,
				},
			},
		},
	}
}

func generateVolumeMounts(fluentbit *loggingv1alpha1.Fluentbit) (v []corev1.VolumeMount) {
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
	if fluentbit.Spec.TLS.Enabled {
		tlsRelatedVolume := []corev1.VolumeMount{
			{
				Name:      "fluent-tls",
				MountPath: "/fluent-bit/tls",
			},
		}
		v = append(v, tlsRelatedVolume...)
	}
	return
}

func generateVolume(fluentbit *loggingv1alpha1.Fluentbit) (v []corev1.Volume) {
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
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "fluent-bit-config",
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
	if fluentbit.Spec.TLS.Enabled {
		tlsRelatedVolume := corev1.Volume{
			Name: "fluent-tls",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: fluentbit.Spec.TLS.SecretName,
				},
			},
		}
		v = append(v, tlsRelatedVolume)
	}
	return
}
