// Copyright © 2021 Banzai Cloud
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

package nodeagent

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

func (n *nodeAgentInstance) daemonSet() (runtime.Object, reconciler.DesiredState, error) {
	var containerPorts []corev1.ContainerPort
	if n.nodeAgent.FluentbitSpec.Metrics != nil && n.nodeAgent.FluentbitSpec.Metrics.Port != 0 {
		containerPorts = append(containerPorts, corev1.ContainerPort{
			Name:          "monitor",
			ContainerPort: n.nodeAgent.FluentbitSpec.Metrics.Port,
			Protocol:      corev1.ProtocolTCP,
		})
	}

	labels := util.MergeLabels(n.nodeAgent.FluentbitSpec.Labels, n.getFluentBitLabels())
	meta := n.NodeAgentObjectMeta(fluentbitDaemonSetName)
	podMeta := metav1.ObjectMeta{
		Labels:      labels,
		Annotations: n.nodeAgent.FluentbitSpec.Annotations,
	}

	if n.configs != nil {
		for key, config := range n.configs {
			h := sha256.New()
			_, _ = h.Write(config)
			templates.Annotate(podMeta, fmt.Sprintf("checksum/%s", key), fmt.Sprintf("%x", h.Sum(nil)))
		}
	}

	desired := &appsv1.DaemonSet{
		ObjectMeta: meta,
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: util.MergeLabels(n.nodeAgent.FluentbitSpec.Labels, n.getFluentBitLabels())},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: podMeta,
				Spec: corev1.PodSpec{
					ServiceAccountName: n.getServiceAccount(),
					Volumes:            n.generateVolume(),
					Tolerations:        n.nodeAgent.FluentbitSpec.Tolerations,
					NodeSelector:       n.nodeAgent.FluentbitSpec.NodeSelector,
					Affinity:           n.nodeAgent.FluentbitSpec.Affinity,
					PriorityClassName:  n.nodeAgent.FluentbitSpec.PodPriorityClassName,
					SecurityContext: &corev1.PodSecurityContext{
						FSGroup:      n.nodeAgent.FluentbitSpec.Security.PodSecurityContext.FSGroup,
						RunAsNonRoot: n.nodeAgent.FluentbitSpec.Security.PodSecurityContext.RunAsNonRoot,
						RunAsUser:    n.nodeAgent.FluentbitSpec.Security.PodSecurityContext.RunAsUser,
						RunAsGroup:   n.nodeAgent.FluentbitSpec.Security.PodSecurityContext.RunAsGroup,
					},
					ImagePullSecrets: n.nodeAgent.FluentbitSpec.Image.ImagePullSecrets,
					Containers: []corev1.Container{
						{
							Name:            containerName,
							Image:           n.nodeAgent.FluentbitSpec.Image.Repository + ":" + n.nodeAgent.FluentbitSpec.Image.Tag,
							ImagePullPolicy: corev1.PullPolicy(n.nodeAgent.FluentbitSpec.Image.PullPolicy),
							Ports:           containerPorts,
							Resources:       n.nodeAgent.FluentbitSpec.Resources,
							VolumeMounts:    n.generateVolumeMounts(),
							SecurityContext: &corev1.SecurityContext{
								RunAsUser:                n.nodeAgent.FluentbitSpec.Security.SecurityContext.RunAsUser,
								RunAsNonRoot:             n.nodeAgent.FluentbitSpec.Security.SecurityContext.RunAsNonRoot,
								ReadOnlyRootFilesystem:   n.nodeAgent.FluentbitSpec.Security.SecurityContext.ReadOnlyRootFilesystem,
								AllowPrivilegeEscalation: n.nodeAgent.FluentbitSpec.Security.SecurityContext.AllowPrivilegeEscalation,
								Privileged:               n.nodeAgent.FluentbitSpec.Security.SecurityContext.Privileged,
								SELinuxOptions:           n.nodeAgent.FluentbitSpec.Security.SecurityContext.SELinuxOptions,
							},
							LivenessProbe:  n.nodeAgent.FluentbitSpec.LivenessProbe,
							ReadinessProbe: n.nodeAgent.FluentbitSpec.ReadinessProbe,
						},
					},
				},
			},
		},
	}

	n.nodeAgent.FluentbitSpec.PositionDB.WithDefaultHostPath(
		fmt.Sprintf(v1beta1.HostPath, n.logging.Name, TailPositionVolume))
	n.nodeAgent.FluentbitSpec.BufferStorageVolume.WithDefaultHostPath(
		fmt.Sprintf(v1beta1.HostPath, n.logging.Name, BufferStorageVolume))

	if err := n.nodeAgent.FluentbitSpec.PositionDB.ApplyVolumeForPodSpec(TailPositionVolume, containerName, "/tail-db", &desired.Spec.Template.Spec); err != nil {
		return desired, reconciler.StatePresent, err
	}
	if err := n.nodeAgent.FluentbitSpec.BufferStorageVolume.ApplyVolumeForPodSpec(BufferStorageVolume, containerName, n.nodeAgent.FluentbitSpec.BufferStorage.StoragePath, &desired.Spec.Template.Spec); err != nil {
		return desired, reconciler.StatePresent, err
	}

	return desired, reconciler.StatePresent, nil
}

func (n *nodeAgentInstance) generateVolumeMounts() (v []corev1.VolumeMount) {
	v = []corev1.VolumeMount{
		{
			Name:      "varlibcontainers",
			ReadOnly:  true,
			MountPath: "/var/lib/docker/containers",
		},
		{
			Name:      "varlogs",
			ReadOnly:  true,
			MountPath: "/var/log/",
		},
	}

	for vCount, vMnt := range n.nodeAgent.FluentbitSpec.ExtraVolumeMounts {
		v = append(v, corev1.VolumeMount{
			Name:      "extravolumemount" + strconv.Itoa(vCount),
			ReadOnly:  vMnt.ReadOnly,
			MountPath: vMnt.Destination,
		})
	}

	if n.nodeAgent.FluentbitSpec.CustomConfigSecret == "" {
		v = append(v, corev1.VolumeMount{
			Name:      "config",
			MountPath: "/fluent-bit/etc/fluent-bit.conf",
			SubPath:   BaseConfigName,
		})
		if n.nodeAgent.FluentbitSpec.EnableUpstream {
			v = append(v, corev1.VolumeMount{
				Name:      "config",
				MountPath: "/fluent-bit/etc/upstream.conf",
				SubPath:   UpstreamConfigName,
			})
		}
	} else {
		v = append(v, corev1.VolumeMount{
			Name:      "config",
			MountPath: "/fluent-bit/etc/",
		})
	}

	if n.nodeAgent.FluentbitSpec.TLS.Enabled {
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

func (n *nodeAgentInstance) generateVolume() (v []corev1.Volume) {
	v = []corev1.Volume{
		{
			Name: "varlibcontainers",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: n.nodeAgent.FluentbitSpec.MountPath,
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

	for vCount, vMnt := range n.nodeAgent.FluentbitSpec.ExtraVolumeMounts {
		v = append(v, corev1.Volume{
			Name: "extravolumemount" + strconv.Itoa(vCount),
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: vMnt.Source,
				},
			}})
	}

	if n.nodeAgent.FluentbitSpec.CustomConfigSecret == "" {
		volume := corev1.Volume{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: n.logging.QualifiedName(fluentBitSecretConfigName),
					Items: []corev1.KeyToPath{
						{
							Key:  BaseConfigName,
							Path: BaseConfigName,
						},
					},
				},
			},
		}
		if n.nodeAgent.FluentbitSpec.EnableUpstream {
			volume.VolumeSource.Secret.Items = append(volume.VolumeSource.Secret.Items, corev1.KeyToPath{
				Key:  UpstreamConfigName,
				Path: UpstreamConfigName,
			})
		}
		v = append(v, volume)
	} else {
		v = append(v, corev1.Volume{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: n.nodeAgent.FluentbitSpec.CustomConfigSecret,
				},
			},
		})
	}
	if n.nodeAgent.FluentbitSpec.TLS.Enabled {
		tlsRelatedVolume := corev1.Volume{
			Name: "fluent-bit-tls",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: n.nodeAgent.FluentbitSpec.TLS.SecretName,
				},
			},
		}
		v = append(v, tlsRelatedVolume)
	}
	return
}
