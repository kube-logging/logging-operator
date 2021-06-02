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

package v1alpha1

import (
	"github.com/banzaicloud/logging-operator/pkg/sdk/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +name:"ClusterFlow"
// +weight:"200"
type _hugoClusterFlow interface{}

// +name:"ClusterFlow"
// +version:"v1alpha1"
// +description:"ClusterFlow is the Schema for the clusterflows API"
type _metaClusterFlow interface{}

// +kubebuilder:object:root=true
// +kubebuilder:resource:categories=logging-all
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Active",type="boolean",JSONPath=".status.active",description="Is the flow active?"
// +kubebuilder:printcolumn:name="Problems",type="integer",JSONPath=".status.problemsCount",description="Number of problems"

// ClusterFlow is the Schema for the clusterflows API
type ClusterFlow struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Name of the logging cluster to be attached
	Spec   v1beta1.ClusterFlowSpec `json:"spec,omitempty"`
	Status v1beta1.FlowStatus      `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterFlowList contains a list of ClusterFlow
type ClusterFlowList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterFlow `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterFlow{}, &ClusterFlowList{})
}
