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
	Annotations         map[string]string                `json:"annotations,omitempty"`
	TLS                 FluentdTLS                       `json:"tls,omitempty"`
	Image               ImageSpec                        `json:"image,omitempty"`
	FluentdPvcSpec      corev1.PersistentVolumeClaimSpec `json:"fluentdPvcSpec,omitempty"`
	DisablePvc          bool                             `json:"disablePvc,omitempty"`
	VolumeModImage      ImageSpec                        `json:"volumeModImage,omitempty"`
	ConfigReloaderImage ImageSpec                        `json:"configReloaderImage,omitempty"`
	Resources           corev1.ResourceRequirements      `json:"resources,omitempty"`
	Port                int32                            `json:"port,omitempty"`
	Tolerations         []corev1.Toleration              `json:"tolerations,omitempty"`
	NodeSelector        map[string]string                `json:"nodeSelector,omitempty"`
	Metrics             *Metrics                         `json:"metrics,omitempty"`
}

// +kubebuilder:object:generate=true

// FluentdTLS defines the TLS configs
type FluentdTLS struct {
	Enabled    bool   `json:"enabled"`
	SecretName string `json:"secretName"`
	SharedKey  string `json:"sharedKey,omitempty"`
}

// GetPrometheusPortFromAnnotation gets the port value from annotation
//func (spec FluentdSpec) GetPrometheusPortFromAnnotation() int32 {
//	var err error
//	var port int64
//	if spec.Annotations != nil {
//		port, err = strconv.ParseInt(spec.Annotations["prometheus.io/port"], 10, 32)
//		if err != nil {
//			return 0
//		}
//	}
//	return int32(port)
//}
//
//func (spec FluentdSpec) GetPrometheusPathFromAnnotation() string {
//	var path string
//	if spec.Annotations != nil {
//		path = spec.Annotations["prometheus.io/path"]
//
//	}
//	return path
//}
