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
	"fmt"
	"reflect"
	"strings"

	"github.com/cisco-open/operator-tools/pkg/reconciler"
	util "github.com/cisco-open/operator-tools/pkg/utils"
	"github.com/spf13/cast"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

func (r *Reconciler) statefulset() (runtime.Object, reconciler.DesiredState, error) {
	spec := r.statefulsetSpec()

	r.fluentdSpec.BufferStorageVolume.WithDefaultHostPath(
		fmt.Sprintf(v1beta1.HostPath, r.Logging.Name, r.Logging.QualifiedName(v1beta1.DefaultFluentdBufferStorageVolumeName)),
	)
	if !r.fluentdSpec.DisablePvc {
		err := r.fluentdSpec.BufferStorageVolume.ApplyPVCForStatefulSet(containerName, bufferPath, spec, func(name string) metav1.ObjectMeta {
			return r.FluentdObjectMeta(name, ComponentFluentd)
		})
		if err != nil {
			return nil, reconciler.StatePresent, err
		}
	} else {
		err := r.fluentdSpec.BufferStorageVolume.ApplyVolumeForPodSpec(r.Logging.QualifiedName(v1beta1.DefaultFluentdBufferStorageVolumeName), containerName, bufferPath, &spec.Template.Spec)
		if err != nil {
			return nil, reconciler.StatePresent, err
		}
	}
	for _, n := range r.fluentdSpec.ExtraVolumes {
		if n.Volume == nil || n.Volume.PersistentVolumeClaim == nil {
			continue
		}
		if isPersistentVolumeClaimSpecEmpty(n.Volume.PersistentVolumeClaim.PersistentVolumeClaimSpec) {
			if err := n.ApplyVolumeForPodSpec(&spec.Template.Spec); err != nil {
				return nil, reconciler.StatePresent, err
			}
		} else {
			if err := n.Volume.ApplyPVCForStatefulSet(n.ContainerName, n.Path, spec, func(name string) metav1.ObjectMeta {
				return r.FluentdObjectMeta(name, ComponentFluentd)
			}); err != nil {
				return nil, reconciler.StatePresent, err
			}
		}
	}

	desired := &appsv1.StatefulSet{
		ObjectMeta: r.FluentdObjectMeta(StatefulSetName, ComponentFluentd),
		Spec:       *spec,
	}

	desired.Annotations = util.MergeLabels(desired.Annotations, r.fluentdSpec.StatefulSetAnnotations)

	return desired, reconciler.StatePresent, nil
}

func isPersistentVolumeClaimSpecEmpty(pvcSpec corev1.PersistentVolumeClaimSpec) bool {
	empty := corev1.PersistentVolumeClaimSpec{}
	return reflect.DeepEqual(pvcSpec, empty)
}

func (r *Reconciler) statefulsetSpec() *appsv1.StatefulSetSpec {
	var initContainers []corev1.Container

	if c := r.tmpDirHackContainer(); c != nil {
		initContainers = append(initContainers, *c)
	}
	if c := r.volumeMountHackContainer(); c != nil {
		initContainers = append(initContainers, *c)
	}
	if i := generateInitContainer(*r.fluentdSpec); i != nil {
		initContainers = append(initContainers, *i)
	}

	containers := []corev1.Container{
		fluentContainer(r.fluentdSpec),
		*newConfigMapReloader(r.fluentdSpec),
	}
	if c := r.bufferMetricsSidecarContainer(); c != nil {
		containers = append(containers, *c)
	}
	if len(r.fluentdSpec.SidecarContainers) != 0 {
		containers = append(containers, r.fluentdSpec.SidecarContainers...)
	}

	sts := &appsv1.StatefulSetSpec{
		PodManagementPolicy: appsv1.PodManagementPolicyType(r.fluentdSpec.Scaling.PodManagementPolicy),
		Selector: &metav1.LabelSelector{
			MatchLabels: r.Logging.GetFluentdLabels(ComponentFluentd, *r.fluentdSpec),
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: r.generatePodMeta(),
			Spec: corev1.PodSpec{
				Volumes:                   r.generateVolume(),
				ServiceAccountName:        r.getServiceAccount(),
				InitContainers:            initContainers,
				ImagePullSecrets:          r.fluentdSpec.Image.ImagePullSecrets,
				Containers:                containers,
				NodeSelector:              r.fluentdSpec.NodeSelector,
				Tolerations:               r.fluentdSpec.Tolerations,
				Affinity:                  r.fluentdSpec.Affinity,
				TopologySpreadConstraints: r.fluentdSpec.TopologySpreadConstraints,
				PriorityClassName:         r.fluentdSpec.PodPriorityClassName,
				DNSPolicy:                 r.fluentdSpec.DNSPolicy,
				DNSConfig:                 r.fluentdSpec.DNSConfig,
				SecurityContext:           r.fluentdSpec.Security.PodSecurityContext,
			},
		},
		ServiceName: r.Logging.QualifiedName(ServiceName + "-headless"),
	}

	if r.fluentdSpec.Scaling.Replicas > 0 {
		sts.Replicas = util.IntPointer(cast.ToInt32(r.fluentdSpec.Scaling.Replicas))
	}

	return sts
}

func fluentContainer(spec *v1beta1.FluentdSpec) corev1.Container {
	envVars := append(spec.EnvVars,
		corev1.EnvVar{Name: "BUFFER_PATH", Value: bufferPath},
	)

	container := corev1.Container{
		Name:            "fluentd",
		Image:           spec.Image.RepositoryWithTag(),
		ImagePullPolicy: corev1.PullPolicy(spec.Image.PullPolicy),
		Ports:           generatePorts(spec),
		VolumeMounts:    generateVolumeMounts(spec),
		Resources:       spec.Resources,
		SecurityContext: spec.Security.SecurityContext,
		Env:             envVars,
		LivenessProbe:   spec.LivenessProbe,
		ReadinessProbe:  generateReadinessCheck(spec),
	}

	if spec.FluentOutLogrotate != nil && spec.FluentOutLogrotate.Enabled {
		container.Args = []string{
			"fluentd",
			"-o", spec.FluentOutLogrotate.Path,
			"--log-rotate-age", spec.FluentOutLogrotate.Age,
			"--log-rotate-size", spec.FluentOutLogrotate.Size,
		}
	}

	container.Args = append(container.Args, spec.ExtraArgs...)

	return container
}

func (r *Reconciler) generatePodMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Labels: r.Logging.GetFluentdLabels(ComponentFluentd, *r.fluentdSpec),
	}
	if r.fluentdSpec.Annotations != nil {
		meta.Annotations = r.fluentdSpec.Annotations
	}
	return meta
}

func newConfigMapReloader(spec *v1beta1.FluentdSpec) *corev1.Container {
	var args []string
	vm := []corev1.VolumeMount{
		{
			Name:      "app-config",
			MountPath: "/fluentd/app-config",
		},
	}

	if spec.CompressConfigFile {
		args = append(args,
			"--volume-dir-archive=/tmp/archive",
			"--dir-for-unarchive=/fluentd/app-config",
			"--webhook-url=http://127.0.0.1:24444/api/config.reload",
		)
		vm = append(vm, corev1.VolumeMount{
			Name:      "app-config-compress",
			MountPath: "tmp/archive",
		})
	} else {
		args = append(args,
			"--volume-dir=/fluentd/etc",
			"--volume-dir=/fluentd/app-config",
			"--webhook-url=http://127.0.0.1:24444/api/config.reload",
		)
		vm = append(vm, corev1.VolumeMount{
			Name:      "config",
			MountPath: "/fluentd/etc",
		})
	}

	c := &corev1.Container{
		Name:            "config-reloader",
		ImagePullPolicy: corev1.PullPolicy(spec.ConfigReloaderImage.PullPolicy),
		Image:           spec.ConfigReloaderImage.RepositoryWithTag(),
		Resources:       spec.ConfigReloaderResources,
		Args:            args,
		VolumeMounts:    vm,
		SecurityContext: spec.Security.SecurityContext,
	}

	return c
}

func generatePortsBufferVolumeMetrics(spec *v1beta1.FluentdSpec) []corev1.ContainerPort {
	port := int32(defaultBufferVolumeMetricsPort)
	if spec.Metrics != nil && spec.BufferVolumeMetrics.Port != 0 {
		port = spec.BufferVolumeMetrics.Port
	}
	return []corev1.ContainerPort{
		{
			Name:          "buffer-metrics",
			ContainerPort: port,
			Protocol:      "TCP",
		},
	}
}

func generatePorts(spec *v1beta1.FluentdSpec) []corev1.ContainerPort {
	ports := []corev1.ContainerPort{
		{
			Name:          "fluent-input",
			ContainerPort: spec.Port,
			Protocol:      "TCP",
		},
	}
	if spec.Metrics != nil && spec.Metrics.Port != 0 {
		ports = append(ports, corev1.ContainerPort{
			Name:          "monitor",
			ContainerPort: spec.Metrics.Port,
			Protocol:      "TCP",
		})
	}
	return ports
}

func generateVolumeMounts(spec *v1beta1.FluentdSpec) []corev1.VolumeMount {
	res := []corev1.VolumeMount{
		{
			Name:      "config",
			MountPath: "/fluentd/etc/",
		},
		{
			Name:      "app-config",
			MountPath: "/fluentd/app-config",
		},
		{
			Name:      "output-secret",
			MountPath: OutputSecretPath,
		},
	}
	if spec != nil && spec.TLS.Enabled {
		res = append(res, corev1.VolumeMount{
			Name:      "fluentd-tls",
			MountPath: "/fluentd/tls/",
		})
	}
	if isFluentdReadOnlyRootFilesystem(spec) {
		res = append(res, corev1.VolumeMount{
			Name:      "tmp",
			SubPath:   "fluentd",
			MountPath: "/tmp",
		})
	}
	return res
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
			Name: "output-secret",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: r.Logging.QualifiedName(OutputSecretName),
				},
			},
		},
	}

	if isFluentdReadOnlyRootFilesystem(r.fluentdSpec) {
		v = append(v, corev1.Volume{
			Name:         "tmp",
			VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
		})
	}

	if r.fluentdSpec.CompressConfigFile {
		v = append(v, corev1.Volume{
			Name: "app-config",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		})
		v = append(v, corev1.Volume{
			Name: "app-config-compress",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: r.Logging.QualifiedName(AppSecretConfigName),
				},
			},
		})
	} else {
		v = append(v, corev1.Volume{
			Name: "app-config",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: r.Logging.QualifiedName(AppSecretConfigName),
				},
			},
		})
	}

	if r.fluentdSpec.TLS.Enabled {
		tlsRelatedVolume := corev1.Volume{
			Name: "fluentd-tls",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: r.fluentdSpec.TLS.SecretName,
				},
			},
		}
		v = append(v, tlsRelatedVolume)
	}
	return
}

func (r *Reconciler) tmpDirHackContainer() *corev1.Container {
	if isFluentdReadOnlyRootFilesystem(r.fluentdSpec) {
		return &corev1.Container{
			Command:         []string{"sh", "-c", "mkdir -p /mnt/tmp/fluentd/; chmod +t /mnt/tmp/fluentd"},
			Image:           r.fluentdSpec.Image.RepositoryWithTag(),
			ImagePullPolicy: corev1.PullPolicy(r.fluentdSpec.Image.PullPolicy),
			Name:            "tmp-dir-hack",
			Resources:       r.fluentdSpec.Resources,
			SecurityContext: r.fluentdSpec.Security.SecurityContext,
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "tmp",
					MountPath: "/mnt/tmp"},
			},
		}
	}
	return nil
}

func (r *Reconciler) volumeMountHackContainer() *corev1.Container {
	if r.fluentdSpec.VolumeMountChmod {
		return &corev1.Container{
			Name:            "volume-mount-hack",
			Image:           r.fluentdSpec.VolumeModImage.RepositoryWithTag(),
			ImagePullPolicy: corev1.PullPolicy(r.fluentdSpec.VolumeModImage.PullPolicy),
			Command:         []string{"sh", "-c", "chmod -R 777 " + bufferPath},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      r.bufferVolumeName(),
					MountPath: bufferPath,
				},
			},
		}
	}
	return nil
}

func (r *Reconciler) bufferMetricsSidecarContainer() *corev1.Container {
	if r.fluentdSpec.BufferVolumeMetrics != nil {
		port := int32(defaultBufferVolumeMetricsPort)
		if r.fluentdSpec.BufferVolumeMetrics.Port != 0 {
			port = r.fluentdSpec.BufferVolumeMetrics.Port
		}
		portParam := fmt.Sprintf("--web.listen-address=:%d", port)
		args := []string{portParam}
		if len(r.fluentdSpec.BufferVolumeArgs) != 0 {
			args = append(args, r.fluentdSpec.BufferVolumeArgs...)
		} else {
			args = append(args, "--collector.disable-defaults", "--collector.filesystem", "--collector.textfile", "--collector.textfile.directory=/prometheus/node_exporter/textfile_collector/")
		}

		nodeExporterCmd := fmt.Sprintf("nodeexporter -> ./bin/node_exporter %v", strings.Join(args, " "))
		bufferSizeCmd := "buffersize -> /prometheus/buffer-size.sh"

		return &corev1.Container{
			Name:            "buffer-metrics-sidecar",
			Image:           r.fluentdSpec.BufferVolumeImage.RepositoryWithTag(),
			ImagePullPolicy: corev1.PullPolicy(r.fluentdSpec.BufferVolumeImage.PullPolicy),
			Args: []string{
				"--exec", nodeExporterCmd,
				"--exec", bufferSizeCmd,
			},
			Env: []corev1.EnvVar{
				{
					Name:  "BUFFER_PATH",
					Value: bufferPath,
				},
			},
			Ports: generatePortsBufferVolumeMetrics(r.fluentdSpec),
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      r.bufferVolumeName(),
					MountPath: bufferPath,
				},
			},
			Resources:       r.fluentdSpec.BufferVolumeResources,
			SecurityContext: r.fluentdSpec.Security.SecurityContext,
		}
	}
	return nil
}

func (r *Reconciler) bufferVolumeName() string {
	volumeName := r.Logging.QualifiedName(v1beta1.DefaultFluentdBufferStorageVolumeName)
	if r.fluentdSpec.BufferStorageVolume.PersistentVolumeClaim != nil {
		volumeName = r.Logging.QualifiedName(r.fluentdSpec.BufferStorageVolume.PersistentVolumeClaim.PersistentVolumeSource.ClaimName)
	}
	return volumeName
}

func generateReadinessCheck(spec *v1beta1.FluentdSpec) *corev1.Probe {
	if spec.ReadinessProbe != nil {
		return spec.ReadinessProbe
	}

	if spec.ReadinessDefaultCheck.BufferFreeSpace || spec.ReadinessDefaultCheck.BufferFileNumber {
		check := []string{"/bin/sh", "-c"}
		bash := []string{}
		if spec.ReadinessDefaultCheck.BufferFreeSpace {
			bash = append(bash,
				fmt.Sprintf("FREESPACE_THRESHOLD=%d", spec.ReadinessDefaultCheck.BufferFreeSpaceThreshold),
				"FREESPACE_CURRENT=$(df -h $BUFFER_PATH  | grep / | awk '{ print $5}' | sed 's/%//g')",
				"if [ \"$FREESPACE_CURRENT\" -gt \"$FREESPACE_THRESHOLD\" ] ; then exit 1; fi",
			)
		}
		if spec.ReadinessDefaultCheck.BufferFileNumber {
			bash = append(bash,
				fmt.Sprintf("MAX_FILE_NUMBER=%d", spec.ReadinessDefaultCheck.BufferFileNumberMax),
				"FILE_NUMBER_CURRENT=$(find $BUFFER_PATH -type f -name *.buffer | wc -l)",
				"if [ \"$FILE_NUMBER_CURRENT\" -gt \"$MAX_FILE_NUMBER\" ] ; then exit 1; fi",
			)
		}
		return &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				Exec: &corev1.ExecAction{
					Command: append(check, strings.Join(bash, "\n")),
				},
			},
			InitialDelaySeconds: spec.ReadinessDefaultCheck.InitialDelaySeconds,
			TimeoutSeconds:      spec.ReadinessDefaultCheck.TimeoutSeconds,
			PeriodSeconds:       spec.ReadinessDefaultCheck.PeriodSeconds,
			SuccessThreshold:    spec.ReadinessDefaultCheck.SuccessThreshold,
			FailureThreshold:    spec.ReadinessDefaultCheck.FailureThreshold,
		}
	}
	return nil
}

func generateInitContainer(spec v1beta1.FluentdSpec) *corev1.Container {
	if spec.CompressConfigFile {
		return &corev1.Container{
			Name:            "init-config-reloader",
			Image:           spec.ConfigReloaderImage.RepositoryWithTag(),
			ImagePullPolicy: corev1.PullPolicy(spec.Image.PullPolicy),
			Resources:       spec.ConfigReloaderResources,
			Args: []string{
				"--init-mode=true",
				"--volume-dir-archive=/tmp/archive",
				"--dir-for-unarchive=/fluentd/app-config",
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "app-config",
					MountPath: "/fluentd/app-config",
				},
				{
					Name:      "app-config-compress",
					MountPath: "/tmp/archive",
				},
			},
			SecurityContext: spec.Security.SecurityContext,
		}
	}
	return nil
}

func isFluentdReadOnlyRootFilesystem(spec *v1beta1.FluentdSpec) bool {
	if spec.Security.SecurityContext.ReadOnlyRootFilesystem != nil {
		return *spec.Security.SecurityContext.ReadOnlyRootFilesystem
	}

	return false
}
