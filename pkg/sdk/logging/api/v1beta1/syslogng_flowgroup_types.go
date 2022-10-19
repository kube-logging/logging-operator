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

// SyslogNGFlowGroupSpec is the Kubernetes spec for SyslogNGFlowGroups
type SyslogNGFlowGroupSpec struct {
	SyslogNGSpec `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:categories=logging-all
// +kubebuilder:storageversion

// SyslogNGFlowGroup Kubernetes object
type SyslogNGFlowGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec SyslogNGFlowGroupSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=sfgl

// SyslogNGFlowGroupList contains a list of SyslogNGFlowFlow
type SyslogNGFlowGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FlowGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SyslogNGFlowGroup{}, &SyslogNGFlowGroupList{})
}
