/*
 * Copyright © 2019 Banzai Cloud
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

import (
	"strconv"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FluentbitSpec defines the desired state of Fluentbit
// +k8s:openapi-gen=true
type FluentbitSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Namespace   string                      `json:"namespace"`
	Annotations map[string]string           `json:"annotations"`
	Image       ImageSpec                   `json:"image"`
	TLS         FluentbitTLS                `json:"tls"`
	Resources   corev1.ResourceRequirements `json:"resources,omitempty"`
	Tolerations []corev1.Toleration         `json:"tolerations,omitempty"`
}

// FluentbitTLS defines the TLS configs
type FluentbitTLS struct {
	Enabled    bool   `json:"enabled"`
	SecretName string `json:"secretName"`
	SecretType string `json:"secretType,omitempty"`
	SharedKey  string `json:"sharedKey"`
}

// FluentbitStatus defines the observed state of Fluentbit
// +k8s:openapi-gen=true
type FluentbitStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Fluentbit is the Schema for the fluentbits API
// +k8s:openapi-gen=true
type Fluentbit struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FluentbitSpec   `json:"spec,omitempty"`
	Status FluentbitStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FluentbitList contains a list of Fluentbit
type FluentbitList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Fluentbit `json:"items"`
}

// GetPrometheusPortFromAnnotation gets the port value from annotation
func (spec FluentbitSpec) GetPrometheusPortFromAnnotation() int32 {
	port, err := strconv.ParseInt(spec.Annotations["prometheus.io/port"], 10, 32)
	if err != nil {
		panic(err)
	}
	return int32(port)
}

func init() {
	SchemeBuilder.Register(&Fluentbit{}, &FluentbitList{})
}
