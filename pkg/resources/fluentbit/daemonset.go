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
	"strings"

	"github.com/banzaicloud/logging-operator/pkg/resources/templates"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/api/v1beta1"
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

	labels := util.MergeLabels(r.Logging.Spec.FluentbitSpec.Labels, r.getFluentBitLabels())
	meta := r.FluentbitObjectMeta(fluentbitDaemonSetName)
	meta.Annotations = util.MergeLabels(meta.Annotations, r.Logging.Spec.FluentbitSpec.DaemonSetAnnotations)
	podMeta := metav1.ObjectMeta{
		Labels:      labels,
		Annotations: r.Logging.Spec.FluentbitSpec.Annotations,
	}

	if r.configs != nil {
		for key, config := range r.configs {
			h := sha256.New()
			_, _ = h.Write(config)
			podMeta = templates.Annotate(podMeta, fmt.Sprintf("checksum/%s", key), fmt.Sprintf("%x", h.Sum(nil)))
		}
	}

	containers := []corev1.Container{
		*r.fluentbitContainer(),
	}
	if c := r.bufferMetricsSidecarContainer(); c != nil {
		containers = append(containers, *c)
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
					NodeSelector:       r.Logging.Spec.FluentbitSpec.NodeSelector,
					Affinity:           r.Logging.Spec.FluentbitSpec.Affinity,
					PriorityClassName:  r.Logging.Spec.FluentbitSpec.PodPriorityClassName,
					SecurityContext: &corev1.PodSecurityContext{
						FSGroup:      r.Logging.Spec.FluentbitSpec.Security.PodSecurityContext.FSGroup,
						RunAsNonRoot: r.Logging.Spec.FluentbitSpec.Security.PodSecurityContext.RunAsNonRoot,
						RunAsUser:    r.Logging.Spec.FluentbitSpec.Security.PodSecurityContext.RunAsUser,
						RunAsGroup:   r.Logging.Spec.FluentbitSpec.Security.PodSecurityContext.RunAsGroup,
					},
					ImagePullSecrets: r.Logging.Spec.FluentbitSpec.Image.ImagePullSecrets,
					DNSPolicy:        r.Logging.Spec.FluentbitSpec.DNSPolicy,
					DNSConfig:        r.Logging.Spec.FluentbitSpec.DNSConfig,
					HostNetwork:      r.Logging.Spec.FluentbitSpec.HostNetwork,

					Containers: containers,
				},
			},
			UpdateStrategy: r.Logging.Spec.FluentbitSpec.UpdateStrategy,
		},
	}

	r.Logging.Spec.FluentbitSpec.PositionDB.WithDefaultHostPath(
		fmt.Sprintf(v1beta1.HostPath, r.Logging.Name, TailPositionVolume))
	r.Logging.Spec.FluentbitSpec.BufferStorageVolume.WithDefaultHostPath(
		fmt.Sprintf(v1beta1.HostPath, r.Logging.Name, BufferStorageVolume))

	if err := r.Logging.Spec.FluentbitSpec.PositionDB.ApplyVolumeForPodSpec(TailPositionVolume, containerName, "/tail-db", &desired.Spec.Template.Spec); err != nil {
		return desired, reconciler.StatePresent, err
	}
	if err := r.Logging.Spec.FluentbitSpec.BufferStorageVolume.ApplyVolumeForPodSpec(BufferStorageVolume, containerName, r.Logging.Spec.FluentbitSpec.BufferStorage.StoragePath, &desired.Spec.Template.Spec); err != nil {
		return desired, reconciler.StatePresent, err
	}

	return desired, reconciler.StatePresent, nil
}

func (r *Reconciler) fluentbitContainer() *corev1.Container {
	return &corev1.Container{
		Name:            containerName,
		Image:           r.Logging.Spec.FluentbitSpec.Image.RepositoryWithTag(),
		ImagePullPolicy: corev1.PullPolicy(r.Logging.Spec.FluentbitSpec.Image.PullPolicy),
		Ports:           r.generatePortsMetrics(),
		Resources:       r.Logging.Spec.FluentbitSpec.Resources,
		VolumeMounts:    r.generateVolumeMounts(),
		SecurityContext: &corev1.SecurityContext{
			RunAsUser:                r.Logging.Spec.FluentbitSpec.Security.SecurityContext.RunAsUser,
			RunAsNonRoot:             r.Logging.Spec.FluentbitSpec.Security.SecurityContext.RunAsNonRoot,
			ReadOnlyRootFilesystem:   r.Logging.Spec.FluentbitSpec.Security.SecurityContext.ReadOnlyRootFilesystem,
			AllowPrivilegeEscalation: r.Logging.Spec.FluentbitSpec.Security.SecurityContext.AllowPrivilegeEscalation,
			Privileged:               r.Logging.Spec.FluentbitSpec.Security.SecurityContext.Privileged,
			SELinuxOptions:           r.Logging.Spec.FluentbitSpec.Security.SecurityContext.SELinuxOptions,
		},
		Env:            r.Logging.Spec.FluentbitSpec.EnvVars,
		LivenessProbe:  r.Logging.Spec.FluentbitSpec.LivenessProbe,
		ReadinessProbe: r.Logging.Spec.FluentbitSpec.ReadinessProbe,
	}
}

func (r *Reconciler) generatePortsMetrics() (containerPorts []corev1.ContainerPort) {
	if r.Logging.Spec.FluentbitSpec.Metrics != nil && r.Logging.Spec.FluentbitSpec.Metrics.Port != 0 {
		containerPorts = append(containerPorts, corev1.ContainerPort{
			Name:          "monitor",
			ContainerPort: r.Logging.Spec.FluentbitSpec.Metrics.Port,
			Protocol:      corev1.ProtocolTCP,
		})
	}
	return
}

func (r *Reconciler) generateVolumeMounts() (v []corev1.VolumeMount) {
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

	for vCount, vMnt := range r.Logging.Spec.FluentbitSpec.ExtraVolumeMounts {
		v = append(v, corev1.VolumeMount{
			Name:      "extravolumemount" + strconv.Itoa(vCount),
			ReadOnly:  *vMnt.ReadOnly,
			MountPath: vMnt.Destination,
		})
	}

	if r.Logging.Spec.FluentbitSpec.CustomConfigSecret == "" {
		v = append(v, corev1.VolumeMount{
			Name:      "config",
			MountPath: "/fluent-bit/etc/fluent-bit.conf",
			SubPath:   BaseConfigName,
		})
		if r.Logging.Spec.FluentbitSpec.EnableUpstream {
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

	if *r.Logging.Spec.FluentbitSpec.TLS.Enabled {
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
		volume := corev1.Volume{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: r.Logging.QualifiedName(fluentBitSecretConfigName),
					Items: []corev1.KeyToPath{
						{
							Key:  BaseConfigName,
							Path: BaseConfigName,
						},
					},
				},
			},
		}
		if r.Logging.Spec.FluentbitSpec.EnableUpstream {
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
					SecretName: r.Logging.Spec.FluentbitSpec.CustomConfigSecret,
				},
			},
		})
	}
	if *r.Logging.Spec.FluentbitSpec.TLS.Enabled {
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

func (r *Reconciler) generatePortsBufferVolumeMetrics() []corev1.ContainerPort {
	port := int32(defaultBufferVolumeMetricsPort)
	if r.Logging.Spec.FluentbitSpec.Metrics != nil && r.Logging.Spec.FluentbitSpec.BufferVolumeMetrics.Port != 0 {
		port = r.Logging.Spec.FluentbitSpec.BufferVolumeMetrics.Port
	}
	return []corev1.ContainerPort{
		{
			Name:          "buffer-metrics",
			ContainerPort: port,
			Protocol:      corev1.ProtocolTCP,
		},
	}
}

func (r *Reconciler) bufferMetricsSidecarContainer() *corev1.Container {
	if r.Logging.Spec.FluentbitSpec.BufferVolumeMetrics != nil {
		port := int32(defaultBufferVolumeMetricsPort)
		if r.Logging.Spec.FluentbitSpec.BufferVolumeMetrics.Port != 0 {
			port = r.Logging.Spec.FluentbitSpec.BufferVolumeMetrics.Port
		}
		portParam := fmt.Sprintf("--web.listen-address=:%d", port)
		args := []string{portParam}
		if len(r.Logging.Spec.FluentbitSpec.BufferVolumeArgs) != 0 {
			args = append(args, r.Logging.Spec.FluentbitSpec.BufferVolumeArgs...)
		} else {
			args = append(args, "--collector.disable-defaults", "--collector.filesystem", "--collector.textfile", "--collector.textfile.directory=/prometheus/node_exporter/textfile_collector/")
		}

		nodeExporterCmd := fmt.Sprintf("nodeexporter -> ./bin/node_exporter %v", strings.Join(args, " "))
		bufferSizeCmd := "buffersize -> /prometheus/buffer-size.sh"

		return &corev1.Container{
			Name:            "buffer-metrics-sidecar",
			Image:           r.Logging.Spec.FluentbitSpec.BufferVolumeImage.RepositoryWithTag(),
			ImagePullPolicy: corev1.PullPolicy(r.Logging.Spec.FluentbitSpec.BufferVolumeImage.PullPolicy),
			Args: []string{
				"--exec", nodeExporterCmd,
				"--exec", bufferSizeCmd,
			},
			Env: []corev1.EnvVar{
				{
					Name:  "BUFFER_PATH",
					Value: r.Logging.Spec.FluentbitSpec.BufferStorage.StoragePath,
				},
			},
			Ports: r.generatePortsBufferVolumeMetrics(),
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      BufferStorageVolume,
					MountPath: r.Logging.Spec.FluentbitSpec.BufferStorage.StoragePath,
				},
			},
		}
	}
	return nil
}
