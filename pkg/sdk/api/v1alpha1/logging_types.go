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

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=loggings,scope=Cluster,categories=logging-all

// Logging is the Schema for the loggings API
type Logging struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   v1beta1.LoggingSpec   `json:"spec,omitempty"`
	Status v1beta1.LoggingStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// LoggingList contains a list of Logging
type LoggingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Logging `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Logging{}, &LoggingList{})
}
