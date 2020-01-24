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
	"crypto/sha256"
	"fmt"
	"strconv"

	"github.com/banzaicloud/logging-operator/pkg/resources/templates"
	"github.com/banzaicloud/logging-operator/pkg/sdk/api/v1beta1"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	util "github.com/banzaicloud/operator-tools/pkg/utils"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	TailPositionVolume  = "positiondb"
	BufferStorageVolume = "buffers"
)

func (r *Reconciler) daemonSet() (runtime.Object, reconciler.DesiredState, error) {
	var containerPorts []corev1.ContainerPort
	if r.Logging.Spec.FluentbitSpec.Metrics != nil && r.Logging.Spec.FluentbitSpec.Metrics.Port != 0 {
		containerPorts = append(containerPorts, corev1.ContainerPort{
			Name:          "monitor",
			ContainerPort: r.Logging.Spec.FluentbitSpec.Metrics.Port,
			Protocol:      corev1.ProtocolTCP,
		})
	}

	labels := util.MergeLabels(r.Logging.Spec.FluentbitSpec.Labels, r.getFluentBitLabels())
	meta := r.FluentbitObjectMeta(fluentbitDaemonSetName)
	podMeta := metav1.ObjectMeta{
		Labels:      labels,
		Annotations: r.Logging.Spec.FluentbitSpec.Annotations,
	}

	if r.desiredConfig != "" {
		h := sha256.New()
		_, _ = h.Write([]byte(r.desiredConfig))
		templates.Annotate(podMeta, "checksum/config", fmt.Sprintf("%x", h.Sum(nil)))
	}

	desired := &appsv1.DaemonSet{
		ObjectMeta: meta,
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: util.MergeLabels(r.Logging.Spec.FluentbitSpec.Labels, r.getFluentBitLabels())},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: podMeta,
				Spec: corev1.PodSpec{
					ServiceAccountName: r.getServiceAccount(),
					Volumes:            r.generateVolume(),
					Tolerations:        r.Logging.Spec.FluentbitSpec.Tolerations,
					PriorityClassName:  r.Logging.Spec.FluentbitSpec.PodPriorityClassName,
					SecurityContext: &corev1.PodSecurityContext{
						FSGroup:      r.Logging.Spec.FluentbitSpec.Security.PodSecurityContext.FSGroup,
						RunAsNonRoot: r.Logging.Spec.FluentbitSpec.Security.PodSecurityContext.RunAsNonRoot,
						RunAsUser:    r.Logging.Spec.FluentbitSpec.Security.PodSecurityContext.RunAsUser,
						RunAsGroup:   r.Logging.Spec.FluentbitSpec.Security.PodSecurityContext.RunAsGroup,
					},
					Containers: []corev1.Container{
						{
							Name:            "fluent-bit",
							Image:           r.Logging.Spec.FluentbitSpec.Image.Repository + ":" + r.Logging.Spec.FluentbitSpec.Image.Tag,
							ImagePullPolicy: corev1.PullPolicy(r.Logging.Spec.FluentbitSpec.Image.PullPolicy),
							Ports:           containerPorts,
							Resources:       r.Logging.Spec.FluentbitSpec.Resources,
							VolumeMounts:    r.generateVolumeMounts(),
							SecurityContext: &corev1.SecurityContext{
								RunAsUser:                r.Logging.Spec.FluentbitSpec.Security.SecurityContext.RunAsUser,
								RunAsNonRoot:             r.Logging.Spec.FluentbitSpec.Security.SecurityContext.RunAsNonRoot,
								ReadOnlyRootFilesystem:   r.Logging.Spec.FluentbitSpec.Security.SecurityContext.ReadOnlyRootFilesystem,
								AllowPrivilegeEscalation: r.Logging.Spec.FluentbitSpec.Security.SecurityContext.AllowPrivilegeEscalation,
								Privileged:               r.Logging.Spec.FluentbitSpec.Security.SecurityContext.Privileged,
							},
							LivenessProbe:  r.Logging.Spec.FluentbitSpec.LivenessProbe,
							ReadinessProbe: r.Logging.Spec.FluentbitSpec.ReadinessProbe,
						},
					},
				},
			},
		},
	}

	return desired, reconciler.StatePresent, nil
}

func (r *Reconciler) generateVolumeMounts() (v []corev1.VolumeMount) {
	v = []corev1.VolumeMount{
		{
			Name:      "varlibcontainers",
			ReadOnly:  true,
			MountPath: "/var/lib/docker/containers",
		},
		{
			Name:      TailPositionVolume,
			MountPath: "/tail-db",
		},
		{
			Name:      BufferStorageVolume,
			MountPath: r.Logging.Spec.FluentbitSpec.BufferStorage.StoragePath,
		},
		{
			Name:      "varlogs",
			ReadOnly:  true,
			MountPath: "/var/log/",
		},
	}

	for vCount, vMnt := range r.Logging.Spec.FluentbitSpec.ExtraVolumeMounts {
		v = append(v, corev1.VolumeMount{
			Name:      "extravolumemount" + strconv.Itoa(vCount),
			ReadOnly:  vMnt.ReadOnly,
			MountPath: vMnt.Destination,
		})
	}

	if r.Logging.Spec.FluentbitSpec.CustomConfigSecret == "" {
		v = append(v, corev1.VolumeMount{
			Name:      "config",
			MountPath: "/fluent-bit/etc/fluent-bit.conf",
			SubPath:   "fluent-bit.conf",
		})
	} else {
		v = append(v, corev1.VolumeMount{
			Name:      "config",
			MountPath: "/fluent-bit/etc/",
		})
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
					Path: r.Logging.Spec.FluentbitSpec.MountPath,
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
	}

	for vCount, vMnt := range r.Logging.Spec.FluentbitSpec.ExtraVolumeMounts {
		v = append(v, corev1.Volume{
			Name: "extravolumemount" + strconv.Itoa(vCount),
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: vMnt.Source,
				},
			}})
	}

	if r.Logging.Spec.FluentbitSpec.CustomConfigSecret == "" {
		v = append(v, corev1.Volume{
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
		})
	} else {
		v = append(v, corev1.Volume{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: r.Logging.Spec.FluentbitSpec.CustomConfigSecret,
				},
			},
		})
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
	r.Logging.Spec.FluentbitSpec.PositionDB.WithDefaultHostPath(
		fmt.Sprintf(v1beta1.HostPath, r.Logging.Name, TailPositionVolume))
	r.Logging.Spec.FluentbitSpec.BufferStorageVolume.WithDefaultHostPath(
		fmt.Sprintf(v1beta1.HostPath, r.Logging.Name, BufferStorageVolume))

	v = append(v, r.Logging.Spec.FluentbitSpec.PositionDB.GetVolume(TailPositionVolume))
	v = append(v, r.Logging.Spec.FluentbitSpec.BufferStorageVolume.GetVolume(BufferStorageVolume))
	return
}
