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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true

// DefaultClusterFlow is the Schema for the defaultclusterflows API
type DefaultClusterFlow struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DefaultClusterFlowSpec `json:"spec,omitempty"`
	Status FlowStatus             `json:"status,omitempty"`
}

// DefaultClusterFlowSpec is the Kubernetes spec for a cluster wide default Flow
type DefaultClusterFlowSpec struct {
	Filters    []Filter `json:"filters,omitempty"`
	LoggingRef string   `json:"loggingRef,omitempty"`
	OutputRefs []string `json:"outputRefs"`
}

// +kubebuilder:object:root=true

// DefaultClusterFlowList contains a list of DefaultClusterFlow
type DefaultClusterFlowList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DefaultClusterFlow `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DefaultClusterFlow{}, &DefaultClusterFlowList{})
}
