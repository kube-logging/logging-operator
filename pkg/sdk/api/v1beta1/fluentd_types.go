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
	"github.com/banzaicloud/operator-tools/pkg/volume"
	corev1 "k8s.io/api/core/v1"
)

// +kubebuilder:object:generate=true

// FluentdSpec defines the desired state of Fluentd
type FluentdSpec struct {
	Annotations map[string]string `json:"annotations,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	TLS         FluentdTLS        `json:"tls,omitempty"`
	Image       ImageSpec         `json:"image,omitempty"`
	DisablePvc  bool              `json:"disablePvc,omitempty"`
	// BufferStorageVolume is by default configured as PVC using FluentdPvcSpec
	// +docLink:"volume.KubernetesVolume,https://github.com/banzaicloud/operator-tools/tree/master/docs/types"
	BufferStorageVolume  volume.KubernetesVolume     `json:"bufferStorageVolume,omitempty"`
	VolumeMountChmod     bool                        `json:"volumeMountChmod,omitempty"`
	VolumeModImage       ImageSpec                   `json:"volumeModImage,omitempty"`
	ConfigReloaderImage  ImageSpec                   `json:"configReloaderImage,omitempty"`
	Resources            corev1.ResourceRequirements `json:"resources,omitempty"`
	LivenessProbe        *corev1.Probe               `json:"livenessProbe,omitempty"`
	LivenessDefaultCheck bool                        `json:"livenessDefaultCheck,omitempty"`
	ReadinessProbe       *corev1.Probe               `json:"readinessProbe,omitempty"`
	Port                 int32                       `json:"port,omitempty"`
	Tolerations          []corev1.Toleration         `json:"tolerations,omitempty"`
	NodeSelector         map[string]string           `json:"nodeSelector,omitempty"`
	Metrics              *Metrics                    `json:"metrics,omitempty"`
	Security             *Security                   `json:"security,omitempty"`
	Scaling              *FluentdScaling             `json:"scaling,omitempty"`
	// +kubebuilder:validation:enum=fatal,error,warn,info,debug,trace
	LogLevel             string `json:"logLevel,omitempty"`
	PodPriorityClassName string `json:"podPriorityClassName,omitempty"`
	// +kubebuilder:validation:enum=stdout,null
	FluentLogDestination string `json:"fluentLogDestination,omitempty"`
	// FluentOutLogrotate sends fluent's stdout to file and rotates it
	FluentOutLogrotate *FluentOutLogrotate `json:"fluentOutLogrotate,omitempty"`
}

// +kubebuilder:object:generate=true

type FluentOutLogrotate struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path,omitempty"`
	Age     string `json:"age,omitempty"`
	Size    string `json:"size,omitempty"`
}

// +kubebuilder:object:generate=true

// FluentdScaling enables configuring the scaling behaviour of the fluentd statefulset
type FluentdScaling struct {
	Replicas int `json:"replicas"`
}

// +kubebuilder:object:generate=true

// FluentdTLS defines the TLS configs
type FluentdTLS struct {
	Enabled    bool   `json:"enabled"`
	SecretName string `json:"secretName"`
	SharedKey  string `json:"sharedKey,omitempty"`
}
