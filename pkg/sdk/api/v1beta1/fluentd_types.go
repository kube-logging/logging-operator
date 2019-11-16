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

// +kubebuilder:object:generate=true

// FluentdSpec defines the desired state of Fluentd
type FluentdSpec struct {
	Annotations map[string]string `json:"annotations,omitempty"`
	TLS         FluentdTLS        `json:"tls,omitempty"`
	Image       ImageSpec         `json:"image,omitempty"`
	// Deprecated, use BufferStorageVolume to configure PVC explicitly
	FluentdPvcSpec corev1.PersistentVolumeClaimSpec `json:"fluentdPvcSpec,omitempty"`
	DisablePvc     bool                             `json:"disablePvc,omitempty"`
	// BufferStorageVolume is the alternative volume to use if PVC is disabled
	BufferStorageVolume KubernetesStorage           `json:"bufferStorageVolume,omitempty"`
	VolumeMountChmod    bool                        `json:"volumeMountChmod,omitempty"`
	VolumeModImage      ImageSpec                   `json:"volumeModImage,omitempty"`
	ConfigReloaderImage ImageSpec                   `json:"configReloaderImage,omitempty"`
	Resources           corev1.ResourceRequirements `json:"resources,omitempty"`
	Port                int32                       `json:"port,omitempty"`
	Tolerations         []corev1.Toleration         `json:"tolerations,omitempty"`
	NodeSelector        map[string]string           `json:"nodeSelector,omitempty"`
	Metrics             *Metrics                    `json:"metrics,omitempty"`
	Security            *Security                   `json:"security,omitempty"`
	Scaling             *FluentdScaling             `json:"scaling,omitempty"`
	// +kubebuilder:validation:enum=fatal,error,warn,info,debug,trace
	LogLevel string `json:"logLevel,omitempty"`
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
