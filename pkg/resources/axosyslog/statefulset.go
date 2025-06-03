// Copyright © 2025 Kube logging authors
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

package axosyslog

import (
	"fmt"
	"path/filepath"

	"github.com/cisco-open/operator-tools/pkg/reconciler"
	"github.com/cisco-open/operator-tools/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

const (
	LabelAppName               = "app.kubernetes.io/name"
	LabelAppComponent          = "app.kubernetes.io/component"
	commonAxoSyslogObjectValue = "axosyslog"
)

func StatefulSet(object any) (runtime.Object, reconciler.DesiredState, error) {
	axoSyslog, ok := object.(*v1beta1.AxoSyslog)
	if !ok {
		return nil, reconciler.StateAbsent, fmt.Errorf("expected *v1beta1.AxoSyslog, got %T", axoSyslog)
	}

	statefulset := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      commonAxoSyslogObjectValue,
			Namespace: axoSyslog.Namespace,
			Labels: map[string]string{
				LabelAppName:      commonAxoSyslogObjectValue,
				LabelAppComponent: commonAxoSyslogObjectValue,
			},
		},
		Spec: appsv1.StatefulSetSpec{
			PodManagementPolicy: appsv1.OrderedReadyPodManagement,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					LabelAppName:      commonAxoSyslogObjectValue,
					LabelAppComponent: commonAxoSyslogObjectValue,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						LabelAppName:      commonAxoSyslogObjectValue,
						LabelAppComponent: commonAxoSyslogObjectValue,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            commonAxoSyslogObjectValue,
							Image:           axoSyslog.Spec.Image.RepositoryWithTag(),
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								{
									Name:          "otlp-grpc",
									ContainerPort: 4317,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							Args: []string{
								"--cfgfile=/etc/axosyslog/config/axosyslog.conf",
								"--no-caps",
								"-Fe",
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/etc/axosyslog/config",
								},
							},
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
							Env: []corev1.EnvVar{{Name: "BUFFER_PATH", Value: "/buffers"}},
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									Exec: &corev1.ExecAction{
										Command: []string{"/usr/sbin/syslog-ng-ctl", "query", "get", "global.sdata_updates.processed"},
									},
								},
								InitialDelaySeconds: 30,
								TimeoutSeconds:      0,
								PeriodSeconds:       10,
								SuccessThreshold:    0,
								FailureThreshold:    3,
							},
						},
						{
							Name:            "config-reloader",
							Image:           axoSyslog.Spec.ConfigReloadImage.RepositoryWithTag(),
							ImagePullPolicy: corev1.PullIfNotPresent,
							Args: []string{
								"-cfgjson",
								generateConfigReloaderConfig(),
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/etc/axosyslog/config",
								},
							},
						},
						// TODO: Add syslog-ng-metrics sidecar
					},
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: axoSyslogConfigName,
									},
								},
							},
						},
					},
					SecurityContext: &corev1.PodSecurityContext{
						FSGroup: utils.IntPointer64(101),
					},
				},
			},
			ServiceName: fmt.Sprintf("%s-headless", commonAxoSyslogObjectValue),
		},
	}
	// TODO: merge with sts overrides

	return statefulset, reconciler.StatePresent, nil
}

func generateConfigReloaderConfig() string {
	return fmt.Sprintf(`
	{
		"events": {
		  "onFileCreate": {
			"%s" : [
			  {
				"exec": {
				  "key": "info",
				  "command": "echo $(date) config secret changed!"
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
	`, filepath.Join("/etc/syslog-ng/config", "..data"), "/tmp/syslog-ng/syslog-ng.ctl")
}
