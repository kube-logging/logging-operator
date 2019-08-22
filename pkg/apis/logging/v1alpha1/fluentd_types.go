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

// FluentdSpec defines the desired state of Fluentd
// +k8s:openapi-gen=true
type FluentdSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Namespace           string                           `json:"namespace"`
	Annotations         map[string]string                `json:"annotations"`
	TLS                 FluentdTLS                       `json:"tls"`
	Image               ImageSpec                        `json:"image"`
	FluentdPvcSpec      corev1.PersistentVolumeClaimSpec `json:"fluentdPvcSpec"`
	VolumeModImage      ImageSpec                        `json:"volumeModImage"`
	ConfigReloaderImage ImageSpec                        `json:"configReloaderImage"`
	Resources           corev1.ResourceRequirements      `json:"resources,omitempty"`
	ServiceType         corev1.ServiceType               `json:"serviceType,omitempty"`
	Tolerations         []corev1.Toleration              `json:"tolerations,omitempty"`
}

// FluentdTLS defines the TLS configs
type FluentdTLS struct {
	Enabled    bool   `json:"enabled"`
	SecretName string `json:"secretName"`
	SharedKey  string `json:"sharedKey"`
}

// FluentdStatus defines the observed state of Fluentd
// +k8s:openapi-gen=true
type FluentdStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Fluentd is the Schema for the fluentds API
// +k8s:openapi-gen=true
type Fluentd struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FluentdSpec   `json:"spec,omitempty"`
	Status FluentdStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FluentdList contains a list of Fluentd
type FluentdList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Fluentd `json:"items"`
}

// GetPrometheusPortFromAnnotation gets the port value from annotation
func (spec FluentdSpec) GetPrometheusPortFromAnnotation() int32 {
	port, err := strconv.ParseInt(spec.Annotations["prometheus.io/port"], 10, 32)
	if err != nil {
		panic(err)
	}
	return int32(port)
}

// GetServiceType gets the service type if set or ClusterIP as the default
func (spec FluentdSpec) GetServiceType() corev1.ServiceType {
	if spec.ServiceType == "" {
		return corev1.ServiceTypeClusterIP
	}
	return spec.ServiceType
}

func init() {
	SchemeBuilder.Register(&Fluentd{}, &FluentdList{})
}
