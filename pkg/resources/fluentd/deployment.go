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

package fluentd

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
func (r *Reconciler) deployment() runtime.Object {
	deploymentName := "fluentd"
	if r.Fluentd.Labels["release"] != "" {
		deploymentName = r.Fluentd.Labels["release"] + "-fluentd"
	}

	return &appsv1.Deployment{
		ObjectMeta: templates.FluentdObjectMeta(deploymentName, util.MergeLabels(r.Fluentd.Labels, labelSelector), r.Fluentd),
		Spec: appsv1.DeploymentSpec{
			Replicas: util.IntPointer(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: labelSelector,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: util.MergeLabels(r.Fluentd.Labels, labelSelector),
					// TODO Move annotations to configuration
					Annotations: map[string]string{
						"prometheus.io/scrape": "true",
						"prometheus.io/path":   "/metrics",
						"prometheus.io/port":   "25000",
					},
				},
				Spec: corev1.PodSpec{
					Volumes: generateVolume(r.Fluentd),
					InitContainers: []corev1.Container{
						{
							Name:    "volume-mount-hack",
							Image:   "busybox",
							Command: []string{"sh", "-c", "chmod -R 777 /buffers"},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "buffer",
									MountPath: "/buffers",
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "fluentd",
							Image: "banzaicloud/fluentd:v1.1.4",
							Ports: []corev1.ContainerPort{
								{
									Name:          "monitor",
									ContainerPort: 25000,
									Protocol:      "TCP",
								},
								{
									Name:          "fluent-input",
									ContainerPort: 24240,
									Protocol:      "TCP",
								},
							},

							VolumeMounts: generateVolumeMounts(r.Fluentd),
						},
						*newConfigMapReloader(),
					},
					//ServiceAccountName: "",
				},
			},
		},
	}
}

func newConfigMapReloader() *corev1.Container {
	return &corev1.Container{
		Name:  "config-reloader",
		Image: "jimmidyson/configmap-reload:v0.2.2",
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

func generateVolumeMounts(fluentd *loggingv1alpha1.Fluentd) (v []corev1.VolumeMount) {
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
			Name:      "buffer",
			MountPath: "/buffers",
		},
	}
	if fluentd.Spec.TLS.Enabled {
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

func generateVolume(fluentd *loggingv1alpha1.Fluentd) (v []corev1.Volume) {
	v = []corev1.Volume{
		{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "fluentd-config",
					},
				},
			},
		},
		{
			Name: "app-config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "fluentd-app-config",
					},
				},
			},
		},
		{
			Name: "buffer",
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: "fluentd-buffer",
					ReadOnly:  false,
				},
			},
		},
	}
	if fluentd.Spec.TLS.Enabled {
		tlsRelatedVolume := corev1.Volume{
			Name: "fluentd-tls",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: fluentd.Spec.TLS.SecretName,
				},
			},
		}
		v = append(v, tlsRelatedVolume)
	}
	return
}
