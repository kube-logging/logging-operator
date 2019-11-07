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

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
)

// ImageSpec struct hold information about image specification
type ImageSpec struct {
	Repository string `json:"repository"`
	Tag        string `json:"tag"`
	PullPolicy string `json:"pullPolicy"`
}

// Metrics defines the service monitor endpoints
type Metrics struct {
	Interval              string `json:"interval,omitempty"`
	Timeout               string `json:"timeout,omitempty"`
	Port                  int32  `json:"port,omitempty"`
	Path                  string `json:"path,omitempty"`
	ServiceMonitor        bool   `json:"serviceMonitor,omitempty"`
	PrometheusAnnotations bool   `json:"prometheusAnnotations,omitempty"`
}

type KubernetesStorage struct {
	HostPath *corev1.HostPathVolumeSource `json:"host_path,omitempty"`
	EmptyDir *corev1.EmptyDirVolumeSource `json:"emptyDir,omitempty"`
}

// GetVolume returns a default emptydir volume if none configured
func (storage *KubernetesStorage) GetVolume(name string) corev1.Volume {
	volume := corev1.Volume{
		Name: name,
	}
	if storage != nil {
		if storage.HostPath != nil {
			volume.VolumeSource = corev1.VolumeSource{
				HostPath: storage.HostPath,
			}
			return volume
		} else if storage.EmptyDir != nil {
			volume.VolumeSource = corev1.VolumeSource{
				EmptyDir: storage.EmptyDir,
			}
			return volume
		}
	}
	// return a default emptydir volume if none configured
	volume.VolumeSource = corev1.VolumeSource{
		EmptyDir: &corev1.EmptyDirVolumeSource{},
	}
	return volume
}

// Security defines Fluentd, Fluentbit deployment security properties
type Security struct {
	ServiceAccount               string                     `json:"serviceAccount,omitempty"`
	RoleBasedAccessControlCreate *bool                      `json:"roleBasedAccessControlCreate,omitempty"`
	PodSecurityPolicyCreate      bool                       `json:"podSecurityPolicyCreate,omitempty"`
	SecurityContext              *corev1.SecurityContext    `json:"securityContext,omitempty"`
	PodSecurityContext           *corev1.PodSecurityContext `json:"podSecurityContext,omitempty"`
}
