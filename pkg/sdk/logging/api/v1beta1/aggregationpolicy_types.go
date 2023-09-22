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

type AggregationPolicySpec struct {
	// LoggingRef identifies the logging that this policy applies to
	LoggingRef string `json:"loggingRef"`

	// Agent is the name of the specific agent that this policy should be applied to
	// Leave it empty if it should apply to all agents.
	Agent string `json:"agent,omitempty"`

	// WatchNamespaceTargets refers to the list of logging resources specified by a label selector to forward logs to
	// Filtering of namespaces will happen based on the watchNamespaces and watchNamespaceSelector fields of the target logging resource
	WatchNamespaceTargets metav1.LabelSelector `json:"watchNamespaceTargets"`
}

type AggregationPolicyStatus struct {
	// Enumerate all loggings with all the destination namespaces expanded
	Tenants []Tenant `json:"tenants,omitempty"`

	// Enumerate problems that prohibits this policy to take effect
	Problems []string `json:"problems,omitempty"`
}

type Tenant struct {
	Name       string   `json:"name"`
	Namespaces []string `json:"namespaces,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=aggregationpolicies,scope=Cluster,shortName=agp,categories=logging-all
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Problems",type="integer",JSONPath=".status.problemsCount",description="Number of problems"
// +kubebuilder:storageversion

// AggregationPolicy defines
type AggregationPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AggregationPolicySpec   `json:"spec,omitempty"`
	Status AggregationPolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type AggregationPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AggregationPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AggregationPolicy{}, &AggregationPolicyList{})
}
