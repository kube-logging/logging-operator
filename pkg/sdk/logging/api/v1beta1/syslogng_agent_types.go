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

// +name:"SyslogNGFlowSpec"
// +weight:"200"
type _hugoSyslogNGAgentSpec interface{} //nolint:deadcode,unused

// +name:"SyslogNGFlowSpec"
// +version:"v1beta1"
// +description:"SyslogNGFlowSpec is the Kubernetes spec for SyslogNGFlows"
type _metaSyslogNGAgentSpec interface{} //nolint:deadcode,unused

// SyslogNGAgentSpec is the Kubernetes spec for SyslogNGAgent
type SyslogNGAgentSpec struct {
	Disabled bool `json:"disabled,omitempty"`
}

// SyslogNGAgentStatus is the Kubernetes status for SyslogNGAgent
type SyslogNGAgentStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:categories=logging-all
// +kubebuilder:subresource:status
// +kubebuilder:storageversion

// SyslogNGAgent Kubernetes object
type SyslogNGAgent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeAgentSyslogNG   `json:"spec,omitempty"`
	Status SyslogNGAgentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SyslogNGAgentList contains a list of SyslogNGAgent
type SyslogNGAgentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SyslogNGAgent `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SyslogNGAgent{}, &SyslogNGAgentList{})
}
