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

func (r *Reconciler) deployment() runtime.Object {
	deploymentName := "fluentd"
	if r.Fluentd.Labels["release"] != "" {
		deploymentName = r.Fluentd.Labels["release"] + "-fluentd"
	}

	deployment := appsv1.Deployment{
		ObjectMeta: templates.FluentdObjectMeta(deploymentName, util.MergeLabels(r.Fluentd.Labels, labelSelector), r.Fluentd),
		Spec: appsv1.DeploymentSpec{
			Replicas: util.IntPointer(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: labelSelector,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      util.MergeLabels(r.Fluentd.Labels, labelSelector),
					Annotations: r.Fluentd.Spec.Annotations,
				},
				Spec: corev1.PodSpec{
					Volumes: generateVolume(r.Fluentd),
					InitContainers: []corev1.Container{
						{
							Name:            "volume-mount-hack",
							Image:           r.Fluentd.Spec.VolumeModImage.Repository + ":" + r.Fluentd.Spec.VolumeModImage.Tag,
							ImagePullPolicy: corev1.PullPolicy(r.Fluentd.Spec.VolumeModImage.PullPolicy),
							Command:         []string{"sh", "-c", "chmod -R 777 /buffers"},
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
							Name:            "fluentd",
							Image:           r.Fluentd.Spec.Image.Repository + ":" + r.Fluentd.Spec.Image.Tag,
							ImagePullPolicy: corev1.PullPolicy(r.Fluentd.Spec.Image.PullPolicy),
							Ports: []corev1.ContainerPort{
								{
									Name:          "monitor",
									ContainerPort: r.Fluentd.Spec.GetPrometheusPortFromAnnotation(),
									Protocol:      "TCP",
								},
								{
									Name:          "fluent-input",
									ContainerPort: 24240,
									Protocol:      "TCP",
								},
							},

							VolumeMounts: generateVolumeMounts(r.Fluentd),
							Resources:    r.Fluentd.Spec.Resources,
						},
						*newConfigMapReloader(r.Fluentd.Spec.ConfigReloaderImage),
					},
					Tolerations: r.Fluentd.Spec.Tolerations,
				},
			},
		},
	}
	if r.Fluentd.Spec.DeploymentStrategy != "" || r.Fluentd.Spec.DeploymentStrategy != appsv1.RollingUpdateDeploymentStrategyType {
		deployment.Spec.Strategy = appsv1.DeploymentStrategy{
			Type: r.Fluentd.Spec.DeploymentStrategy,
		}
	}
	return &deployment
}

func newConfigMapReloader(spec loggingv1alpha1.ImageSpec) *corev1.Container {
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
