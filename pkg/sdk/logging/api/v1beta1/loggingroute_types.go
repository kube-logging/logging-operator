// Copyright © 2023 Kube logging authors
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

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type LoggingRouteSpec struct {
	// Source identifies the logging that this policy applies to
	Source string `json:"source"`

	// Targets refers to the list of logging resources specified by a label selector to forward logs to
	// Filtering of namespaces will happen based on the watchNamespaces and watchNamespaceSelector fields of the target logging resource
	Targets metav1.LabelSelector `json:"targets"`
}

type LoggingRouteStatus struct {
	// Enumerate all loggings with all the destination namespaces expanded
	Tenants []Tenant `json:"tenants,omitempty"`

	// Enumerate problems that prohibits this route to take effect and populate the tenants field
	Problems []string `json:"problems,omitempty"`

	// Notices highlights non-blocker issues the user should pay attention to
	Notices []string `json:"notices,omitempty"`
}

type Tenant struct {
	Name       string   `json:"name"`
	Namespaces []string `json:"namespaces,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=loggingroutes,scope=Cluster,shortName=lr,categories=logging-all
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Problems",type="integer",JSONPath=".status.problemsCount",description="Number of problems"
// +kubebuilder:storageversion

// LoggingRoute defines
type LoggingRoute struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LoggingRouteSpec   `json:"spec,omitempty"`
	Status LoggingRouteStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type LoggingRouteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LoggingRoute `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LoggingRoute{}, &LoggingRouteList{})
}
