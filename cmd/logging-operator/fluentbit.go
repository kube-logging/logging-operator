package main

import (
	corev1 "k8s.io/api/core/v1"
extensionv1 "k8s.io/api/extensions/v1beta1"
metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
"k8s.io/apimachinery/pkg/runtime/schema"
	"github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1"
)

// What inputs we neeed?
func generateConfig() string {
	return ""
}

func newFluentBitConifg(cr *v1alpha1.LoggingOperator) *corev1.ConfigMap {
	labels := map[string]string {
		"app": "fluent-bit",
	}
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind: "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "fluent-bit",
			Namespace: cr.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "LoggingOperator",
				}),
			},
			Labels: labels,
		},
		Data: map[string]string {
			"fluent-bit.conf": generateConfig(),
		},
	}
}

// TODO the options should come from the operator configuration
func newFluentBitDaemonSet(cr *v1alpha1.LoggingOperator) *extensionv1.DaemonSet {
	labels := map[string]string {
		"app": "fluent-bit",
	}
	return &extensionv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			Kind: "DaemonSet",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "fluent-bit",
			Namespace: cr.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "LoggingOperator",
				}),
			},
			Labels: labels,
		},
		Spec: extensionv1.DaemonSetSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "fluent-bit",
					Labels: labels,
					// TODO Move annotations to configuration
					Annotations: map[string]string{
						"prometheus.io/scrape": "true",
						"prometheus.io/path": "/metrics",
						"prometheus.io/port": "24231",
					},
					OwnerReferences: []metav1.OwnerReference{
						*metav1.NewControllerRef(cr, schema.GroupVersionKind{
							Group:   v1alpha1.SchemeGroupVersion.Group,
							Version: v1alpha1.SchemeGroupVersion.Version,
							Kind:    "LoggingOperator",
						}),
					},
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "varlibcontainers",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "/var/lib/docker/containers",
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							// TODO move to configuration
							Name:    "fluent-bit",
							Image:   "fluent/fluent-bit:latest",
							Ports:      []corev1.ContainerPort{
								{
									Name:          "monitor",
									ContainerPort: 2020,
									Protocol:      "TCP",
								},
							},
							// TODO Get this from config
							Resources: corev1.ResourceRequirements{
								Limits:   nil,
								Requests: nil,
							},
							VolumeMounts:  []corev1.VolumeMount{
								{
									Name:             "varlibcontainers",
									ReadOnly:         true,
									MountPath:        "/var/lib/docker/containers",
								},
							},
						},
					},
				},
			},
		},
	}
}