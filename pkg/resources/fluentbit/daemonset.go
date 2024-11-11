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

	"github.com/cisco-open/operator-tools/pkg/reconciler"
	util "github.com/cisco-open/operator-tools/pkg/utils"

	"github.com/kube-logging/logging-operator/pkg/resources/templates"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"

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

	labels := util.MergeLabels(r.fluentbitSpec.Labels, r.getFluentBitLabels())
	meta := r.FluentbitObjectMeta(fluentbitDaemonSetName)
	meta.Annotations = util.MergeLabels(meta.Annotations, r.fluentbitSpec.DaemonSetAnnotations)
	podMeta := metav1.ObjectMeta{
		Labels:      labels,
		Annotations: r.fluentbitSpec.Annotations,
	}
	imagePullSecrets := r.fluentbitSpec.Image.ImagePullSecrets

	if r.fluentbitSpec.ConfigHotReload == nil && r.configs != nil {
		for key, config := range r.configs {
			h := sha256.New()
			_, _ = h.Write(config)
			podMeta = templates.Annotate(podMeta, fmt.Sprintf("checksum/%s", key), fmt.Sprintf("%x", h.Sum(nil)))
		}
	}

	containers := []corev1.Container{
		*r.fluentbitContainer(),
	}
	if r.fluentbitSpec.ConfigHotReload != nil {
		containers = append(containers, newConfigMapReloader(r.fluentbitSpec))
		imagePullSecrets = append(imagePullSecrets, r.fluentbitSpec.ConfigHotReload.Image.ImagePullSecrets...)
	}
	if c := r.bufferMetricsSidecarContainer(); c != nil {
		containers = append(containers, *c)
	}

	desired := &appsv1.DaemonSet{
		ObjectMeta: meta,
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: util.MergeLabels(r.fluentbitSpec.Labels, r.getFluentBitLabels())},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: podMeta,
				Spec: corev1.PodSpec{
					ServiceAccountName: r.getServiceAccount(),
					Volumes:            r.generateVolume(),
					Tolerations:        r.fluentbitSpec.Tolerations,
					NodeSelector:       r.fluentbitSpec.NodeSelector,
					Affinity:           r.fluentbitSpec.Affinity,
					PriorityClassName:  r.fluentbitSpec.PodPriorityClassName,
					SecurityContext: &corev1.PodSecurityContext{
						FSGroup:        r.fluentbitSpec.Security.PodSecurityContext.FSGroup,
						RunAsNonRoot:   r.fluentbitSpec.Security.PodSecurityContext.RunAsNonRoot,
						RunAsUser:      r.fluentbitSpec.Security.PodSecurityContext.RunAsUser,
						RunAsGroup:     r.fluentbitSpec.Security.PodSecurityContext.RunAsGroup,
						SeccompProfile: r.fluentbitSpec.Security.SecurityContext.SeccompProfile,
					},
					ImagePullSecrets: imagePullSecrets,
					DNSPolicy:        r.fluentbitSpec.DNSPolicy,
					DNSConfig:        r.fluentbitSpec.DNSConfig,
					HostNetwork:      r.fluentbitSpec.HostNetwork,

					Containers: containers,
				},
			},
			UpdateStrategy: r.fluentbitSpec.UpdateStrategy,
		},
	}

	r.fluentbitSpec.PositionDB.WithDefaultHostPath(
		fmt.Sprintf(v1beta1.HostPath, r.nameProvider.Name(), TailPositionVolume))
	r.fluentbitSpec.BufferStorageVolume.WithDefaultHostPath(
		fmt.Sprintf(v1beta1.HostPath, r.nameProvider.Name(), BufferStorageVolume))

	if err := r.fluentbitSpec.PositionDB.ApplyVolumeForPodSpec(TailPositionVolume, containerName, "/tail-db", &desired.Spec.Template.Spec); err != nil {
		return desired, reconciler.StatePresent, err
	}
	if err := r.fluentbitSpec.BufferStorageVolume.ApplyVolumeForPodSpec(BufferStorageVolume, containerName, r.fluentbitSpec.BufferStorage.StoragePath, &desired.Spec.Template.Spec); err != nil {
		return desired, reconciler.StatePresent, err
	}

	return desired, reconciler.StatePresent, nil
}

func (r *Reconciler) fluentbitContainer() *corev1.Container {
	args := []string{
		StockBinPath, "-c", fmt.Sprintf("%s/%s", OperatorConfigPath, BaseConfigName),
	}
	if r.fluentbitSpec.ConfigHotReload != nil {
		args = append(args, "--enable-hot-reload")
	}
	return &corev1.Container{
		Name:            containerName,
		Image:           r.fluentbitSpec.Image.RepositoryWithTag(),
		ImagePullPolicy: corev1.PullPolicy(r.fluentbitSpec.Image.PullPolicy),
		Ports:           r.generatePortsMetrics(),
		Resources:       r.fluentbitSpec.Resources,
		VolumeMounts:    r.generateVolumeMounts(),
		SecurityContext: r.fluentbitSpec.Security.SecurityContext,
		Command:         args,
		Env:             r.fluentbitSpec.EnvVars,
		LivenessProbe:   r.fluentbitSpec.LivenessProbe,
		ReadinessProbe:  r.fluentbitSpec.ReadinessProbe,
	}
}

func (r *Reconciler) generatePortsMetrics() (containerPorts []corev1.ContainerPort) {
	if r.fluentbitSpec.Metrics != nil && r.fluentbitSpec.Metrics.Port != 0 {
		containerPorts = append(containerPorts, corev1.ContainerPort{
			Name:          "monitor",
			ContainerPort: r.fluentbitSpec.Metrics.Port,
			Protocol:      corev1.ProtocolTCP,
		})
	}
	return
}

func newConfigMapReloader(spec *v1beta1.FluentbitSpec) corev1.Container {
	var args []string
	vm := []corev1.VolumeMount{
		{
			Name:      "config",
			MountPath: OperatorConfigPath,
		},
	}

	args = append(args,
		fmt.Sprintf("--volume-dir=%s", OperatorConfigPath),
		fmt.Sprintf("--webhook-url=http://127.0.0.1:%d/api/v2/reload", spec.Metrics.Port),
	)

	c := corev1.Container{
		Name:            "config-reloader",
		ImagePullPolicy: corev1.PullPolicy(spec.ConfigHotReload.Image.PullPolicy),
		Image:           spec.ConfigHotReload.Image.RepositoryWithTag(),
		Resources:       spec.ConfigHotReload.Resources,
		Args:            args,
		VolumeMounts:    vm,
		SecurityContext: spec.Security.SecurityContext,
	}

	return c
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
		{
			Name:      "config",
			MountPath: OperatorConfigPath,
		},
	}

	for vCount, vMnt := range r.fluentbitSpec.ExtraVolumeMounts {
		v = append(v, corev1.VolumeMount{
			Name:      "extravolumemount" + strconv.Itoa(vCount),
			ReadOnly:  *vMnt.ReadOnly,
			MountPath: vMnt.Destination,
		})
	}

	if *r.fluentbitSpec.TLS.Enabled {
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
					Path: r.fluentbitSpec.MountPath,
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

	for vCount, vMnt := range r.fluentbitSpec.ExtraVolumeMounts {
		v = append(v, corev1.Volume{
			Name: "extravolumemount" + strconv.Itoa(vCount),
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: vMnt.Source,
				},
			}})
	}

	if r.fluentbitSpec.CustomConfigSecret == "" {
		volume := corev1.Volume{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: r.nameProvider.ComponentName(fluentBitSecretConfigName),
				},
			},
		}
		v = append(v, volume)
	} else {
		v = append(v, corev1.Volume{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: r.fluentbitSpec.CustomConfigSecret,
				},
			},
		})
	}
	if *r.fluentbitSpec.TLS.Enabled {
		tlsRelatedVolume := corev1.Volume{
			Name: "fluent-bit-tls",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: r.fluentbitSpec.TLS.SecretName,
				},
			},
		}
		v = append(v, tlsRelatedVolume)
	}
	return
}

func (r *Reconciler) generatePortsBufferVolumeMetrics() []corev1.ContainerPort {
	port := int32(defaultBufferVolumeMetricsPort)
	if r.fluentbitSpec.Metrics != nil && r.fluentbitSpec.BufferVolumeMetrics.Port != 0 {
		port = r.fluentbitSpec.BufferVolumeMetrics.Port
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
	if r.fluentbitSpec.BufferVolumeMetrics != nil {
		port := int32(defaultBufferVolumeMetricsPort)
		if r.fluentbitSpec.BufferVolumeMetrics.Port != 0 {
			port = r.fluentbitSpec.BufferVolumeMetrics.Port
		}
		portParam := fmt.Sprintf("--web.listen-address=:%d", port)
		args := []string{portParam}
		if len(r.fluentbitSpec.BufferVolumeArgs) != 0 {
			args = append(args, r.fluentbitSpec.BufferVolumeArgs...)
		} else {
			args = append(args, "--collector.disable-defaults", "--collector.filesystem", "--collector.textfile", "--collector.textfile.directory=/prometheus/node_exporter/textfile_collector/")
		}

		nodeExporterCmd := fmt.Sprintf("nodeexporter -> ./bin/node_exporter %v", strings.Join(args, " "))
		bufferSizeCmd := "buffersize -> /prometheus/buffer-size.sh"

		return &corev1.Container{
			Name:            "buffer-metrics-sidecar",
			Image:           r.fluentbitSpec.BufferVolumeImage.RepositoryWithTag(),
			ImagePullPolicy: corev1.PullPolicy(r.fluentbitSpec.BufferVolumeImage.PullPolicy),
			Args: []string{
				"--exec", nodeExporterCmd,
				"--exec", bufferSizeCmd,
			},
			Env: []corev1.EnvVar{
				{
					Name:  "BUFFER_PATH",
					Value: r.fluentbitSpec.BufferStorage.StoragePath,
				},
			},
			Ports: r.generatePortsBufferVolumeMetrics(),
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      BufferStorageVolume,
					MountPath: r.fluentbitSpec.BufferStorage.StoragePath,
				},
			},
			Resources:       r.fluentbitSpec.BufferVolumeResources,
			SecurityContext: r.fluentbitSpec.Security.SecurityContext,
			LivenessProbe:   r.fluentbitSpec.BufferVolumeLivenessProbe,
		}
	}
	return nil
}
