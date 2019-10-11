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
	"strconv"

	corev1 "k8s.io/api/core/v1"
)

// +kubebuilder:object:generate=true

// FluentbitSpec defines the desired state of Fluentbit
type FluentbitSpec struct {
	Annotations map[string]string           `json:"annotations,omitempty"`
	Image       ImageSpec                   `json:"image,omitempty"`
	TLS         FluentbitTLS                `json:"tls,omitempty"`
	TargetHost  string                      `json:"targetHost,omitempty"`
	TargetPort  int32                       `json:"targetPort,omitempty"`
	Resources   corev1.ResourceRequirements `json:"resources,omitempty"`
	Parser      string                      `json:"parser,omitempty"`
	Tolerations []corev1.Toleration         `json:"tolerations,omitempty"`
	Metrics     *Metrics                    `json:"metrics,omitempty"`
}

// +kubebuilder:object:generate=true

// FluentbitTLS defines the TLS configs
type FluentbitTLS struct {
	Enabled    bool   `json:"enabled"`
	SecretName string `json:"secretName"`
	SharedKey  string `json:"sharedKey,omitempty"`
}

// GetPrometheusPortFromAnnotation gets the port value from annotation
func (spec FluentbitSpec) GetPrometheusPortFromAnnotation() int32 {
	var err error
	var port int64
	if spec.Annotations != nil {
		port, err = strconv.ParseInt(spec.Annotations["prometheus.io/port"], 10, 32)
		if err != nil {
			panic(err)
		}
	}
	return int32(port)
}
