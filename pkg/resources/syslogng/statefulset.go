// Copyright © 2022 Banzai Cloud
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

package syslogng

import (
	"fmt"
	"strings"

	"emperror.dev/errors"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/banzaicloud/operator-tools/pkg/merge"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	util "github.com/banzaicloud/operator-tools/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

//func statefulSetDefaults(userDefined **v1beta1.SyslogNGSpec) (*v1beta1.SyslogNGSpec, error) {
//
//}
func (r *Reconciler) statefulset() (runtime.Object, reconciler.DesiredState, error) {

	containers := []corev1.Container{
		syslogNGContainer(r.Logging.Spec.SyslogNGSpec),
		configReloadContainer(r.Logging.Spec.SyslogNGSpec),
	}
	if c := r.bufferMetricsSidecarContainer(); c != nil {
		containers = append(containers, *c)
	}

	if c := r.syslogNGMetricsSidecarContainer(); c != nil {
		containers = append(containers, *c)
	}

	desired := &appsv1.StatefulSet{
		ObjectMeta: r.Logging.SyslogNGObjectMeta(StatefulSetName, ComponentSyslogNG),
		Spec: appsv1.StatefulSetSpec{
			PodManagementPolicy: appsv1.OrderedReadyPodManagement,
			Selector: &metav1.LabelSelector{
				MatchLabels: r.Logging.GetSyslogNGLabels(ComponentSyslogNG),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: r.Logging.GetSyslogNGLabels(ComponentSyslogNG),
				},
				Spec: corev1.PodSpec{
					Containers: containers,
					Volumes:    r.generateVolume(),
					SecurityContext: &corev1.PodSecurityContext{
						FSGroup: util.IntPointer64(101),
					},
				},
			},
			ServiceName: r.Logging.QualifiedName(ServiceName + "-headless"),
		},
	}
	if !r.Logging.Spec.SyslogNGSpec.SkipRoleBasedAccessControlCreate {
		desired.Spec.Template.Spec.ServiceAccountName = r.Logging.QualifiedName(defaultServiceAccountName)
	}
	err := merge.Merge(desired, r.Logging.Spec.SyslogNGSpec.StatefulSetOverrides)
	if err != nil {
		return desired, reconciler.StatePresent, errors.WrapIf(err, "unable to merge overrides to base object")
	}

	return desired, reconciler.StatePresent, nil
}

func syslogNGContainer(spec *v1beta1.SyslogNGSpec) corev1.Container {

	container := corev1.Container{
		Name:            containerName,
		Image:           v1beta1.RepositoryWithTag(imageRepository, imageTag),
		ImagePullPolicy: corev1.PullIfNotPresent,
		Ports: []corev1.ContainerPort{{
			Name:          "syslog-ng-tcp",
			ContainerPort: 601,
			Protocol:      corev1.ProtocolTCP,
		}},
		VolumeMounts: generateVolumeMounts(spec),
		Resources: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse("400M"),
				corev1.ResourceCPU:    resource.MustParse("1000m"),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse("100M"),
				corev1.ResourceCPU:    resource.MustParse("500m"),
			},
		},
		Env: []corev1.EnvVar{{Name: "BUFFER_PATH", Value: bufferPath}},
		LivenessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				Exec: &corev1.ExecAction{
					Command: []string{`syslog-ng-ctl --control=/tmp/syslog-ng/syslog-ng.ctl query get "destination.file.processed"`}},
			},
			InitialDelaySeconds: 30,
			TimeoutSeconds:      0,
			PeriodSeconds:       10,
			SuccessThreshold:    0,
			FailureThreshold:    3,
		},
		ReadinessProbe: generateReadinessCheck(spec),
	}

	return container
}

func generatePortsBufferVolumeMetrics(spec *v1beta1.SyslogNGSpec) []corev1.ContainerPort {
	port := int32(defaultBufferVolumeMetricsPort)
	if spec.Metrics != nil && spec.BufferVolumeMetrics.Port != 0 {
		port = spec.BufferVolumeMetrics.Port
	}
	return []corev1.ContainerPort{
		{
			Name:          "buffer-metrics",
			ContainerPort: port,
			Protocol:      corev1.ProtocolTCP,
		},
	}
}

func generateVolumeMounts(spec *v1beta1.SyslogNGSpec) []corev1.VolumeMount {
	res := []corev1.VolumeMount{
		{
			Name:      configVolumeName,
			MountPath: configDir,
		},
	}
	if spec != nil && spec.TLS.Enabled {
		res = append(res, corev1.VolumeMount{
			Name:      tlsVolumeName,
			MountPath: "/syslog-ng/tls/",
		})
	}
	if spec != nil {
		res = append(res, corev1.VolumeMount{
			Name:      socketVolumeName,
			MountPath: "/tmp/syslog-ng",
		})
	}

	return res
}

func (r *Reconciler) generateVolume() (v []corev1.Volume) {
	v = []corev1.Volume{
		{
			Name: configVolumeName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: r.Logging.QualifiedName(SecretConfigName),
				},
			},
		},
	}
	if r.Logging.Spec.SyslogNGSpec.TLS.Enabled {
		tlsRelatedVolume := corev1.Volume{
			Name: tlsVolumeName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: r.Logging.Spec.SyslogNGSpec.TLS.SecretName,
				},
			},
		}
		v = append(v, tlsRelatedVolume)
	}
	socketVolume := corev1.Volume{
		Name: socketVolumeName,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{
				Medium:    corev1.StorageMediumDefault,
				SizeLimit: nil,
			},
		},
	}
	v = append(v, socketVolume)

	return
}

func (r *Reconciler) syslogNGMetricsSidecarContainer() *corev1.Container {
	if r.Logging.Spec.SyslogNGSpec.Metrics != nil {
		return &corev1.Container{
			Name:            "exporter",
			ImagePullPolicy: corev1.PullIfNotPresent,
			Image:           v1beta1.RepositoryWithTag(prometheusExporterImageRepository, prometheusExporterImageTag),
			Ports: []corev1.ContainerPort{
				{
					Name:          metricsPortName,
					ContainerPort: metricsPortNumber,
					Protocol:      corev1.ProtocolTCP,
				},
			},
			Args: []string{
				"--socket.path=/tmp/syslog-ng.ctl",
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      socketVolumeName,
					MountPath: "/tmp",
				},
			},
		}

	}
	return nil
}

func (r *Reconciler) bufferMetricsSidecarContainer() *corev1.Container {
	if r.Logging.Spec.SyslogNGSpec.BufferVolumeMetrics != nil {
		port := int32(defaultBufferVolumeMetricsPort)
		if r.Logging.Spec.SyslogNGSpec.BufferVolumeMetrics.Port != 0 {
			port = r.Logging.Spec.SyslogNGSpec.BufferVolumeMetrics.Port
		}
		portParam := fmt.Sprintf("--web.listen-address=:%d", port)
		args := []string{portParam, "--collector.disable-defaults", "--collector.filesystem"}

		customRunner := fmt.Sprintf("./bin/node_exporter %v", strings.Join(args, " "))
		return &corev1.Container{
			Name:            "buffer-metrics-sidecar",
			Image:           v1beta1.RepositoryWithTag(bufferVolumeImageRepository, bufferVolumeImageTag),
			ImagePullPolicy: corev1.PullIfNotPresent,
			Args:            []string{"--startup", customRunner},
			Ports:           generatePortsBufferVolumeMetrics(r.Logging.Spec.SyslogNGSpec),
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      r.Logging.QualifiedName(bufferStorageVolumeName),
					MountPath: bufferPath,
				},
			},
		}
	}
	return nil
}

func generateReadinessCheck(spec *v1beta1.SyslogNGSpec) *corev1.Probe {

	if spec.ReadinessDefaultCheck.BufferFreeSpace || spec.ReadinessDefaultCheck.BufferFileNumber {
		check := []string{"/bin/sh", "-c"}
		bash := []string{}
		if spec.ReadinessDefaultCheck.BufferFreeSpace {
			if spec.ReadinessDefaultCheck.BufferFreeSpaceThreshold == 0 {
				spec.ReadinessDefaultCheck.BufferFreeSpaceThreshold = 90
			}
			bash = append(bash,
				fmt.Sprintf("FREESPACE_THRESHOLD=%d", spec.ReadinessDefaultCheck.BufferFreeSpaceThreshold),
				"FREESPACE_CURRENT=$(df -h $BUFFER_PATH  | grep / | awk '{ print $5}' | sed 's/%//g')",
				"if [ \"$FREESPACE_CURRENT\" -gt \"$FREESPACE_THRESHOLD\" ] ; then exit 1; fi",
			)
		}
		if spec.ReadinessDefaultCheck.BufferFileNumber {
			if spec.ReadinessDefaultCheck.BufferFileNumberMax == 0 {
				spec.ReadinessDefaultCheck.BufferFileNumberMax = 5000
			}

			bash = append(bash,
				fmt.Sprintf("MAX_FILE_NUMBER=%d", spec.ReadinessDefaultCheck.BufferFileNumberMax),
				"FILE_NUMBER_CURRENT=$(find $BUFFER_PATH -type f -name *.buffer | wc -l)",
				"if [ \"$FILE_NUMBER_CURRENT\" -gt \"$MAX_FILE_NUMBER\" ] ; then exit 1; fi",
			)
		}
		if spec.ReadinessDefaultCheck.InitialDelaySeconds == 0 {
			spec.ReadinessDefaultCheck.InitialDelaySeconds = 5
		}
		if spec.ReadinessDefaultCheck.TimeoutSeconds == 0 {
			spec.ReadinessDefaultCheck.TimeoutSeconds = 3
		}
		if spec.ReadinessDefaultCheck.PeriodSeconds == 0 {
			spec.ReadinessDefaultCheck.PeriodSeconds = 30
		}
		if spec.ReadinessDefaultCheck.SuccessThreshold == 0 {
			spec.ReadinessDefaultCheck.SuccessThreshold = 3
		}
		if spec.ReadinessDefaultCheck.FailureThreshold == 0 {
			spec.ReadinessDefaultCheck.FailureThreshold = 1
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

func configReloadContainer(spec *v1beta1.SyslogNGSpec) corev1.Container {
	//TODO: ADD TLS reload watch
	container := corev1.Container{
		Name:            "config-reloader",
		Image:           v1beta1.RepositoryWithTag(configReloaderImageRepository, configReloaderImageTag),
		ImagePullPolicy: corev1.PullIfNotPresent,
		Args: []string{
			"-cfgjson",
			generateConfigReloaderConfig(configDir),
		},
		VolumeMounts: generateVolumeMounts(spec),
	}

	return container
}

func generateConfigReloaderConfig(configDir string) string {
	return fmt.Sprintf(`
	{
		"events": {
		  "onFileCreate": {
			"%s..data" : [
			  {
				"exec": {
				  "key": "info",
				  "command": "echo config secret changed!"
				}  
			  },
			  {
				"exec": {
					"key": "reload",
					"command": "echo RELOAD | socat - UNIX-CONNECT:%s"
				}
			  }  
			]
		  }
		}
	  }	  
	`, configDir, socketPath)
}
