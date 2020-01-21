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
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/filter"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FlowSpec is the Kubernetes spec for Flows
type FlowSpec struct {
	Selectors  map[string]string `json:"selectors"`
	Filters    []Filter          `json:"filters,omitempty"`
	LoggingRef string            `json:"loggingRef,omitempty"`
	OutputRefs []string          `json:"outputRefs"`
}

// Filter definition for FlowSpec
type Filter struct {
	StdOut            *filter.StdOutFilterConfig `json:"stdout,omitempty"`
	Parser            *filter.ParserConfig       `json:"parser,omitempty"`
	TagNormaliser     *filter.TagNormaliser      `json:"tag_normaliser,omitempty"`
	Dedot             *filter.DedotFilterConfig  `json:"dedot,omitempty"`
	RecordTransformer *filter.RecordTransformer  `json:"record_transformer,omitempty"`
	GeoIP             *filter.GeoIP              `json:"geoip,omitempty"`
	Concat            *filter.Concat             `json:"concat,omitempty"`
	DetectExceptions  *filter.DetectExceptions   `json:"detectExceptions,omitempty"`
	Grep              *filter.GrepConfig         `json:"grep,omitempty"`
	Prometheus        *filter.PrometheusConfig   `json:"prometheus,omitempty"`
}

// FlowStatus defines the observed state of Flow
type FlowStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Logging",type=string,JSONPath=`.spec.loggingRef`

// Flow Kubernetes object
type Flow struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FlowSpec   `json:"spec,omitempty"`
	Status FlowStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FlowList contains a list of Flow
type FlowList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Flow `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Flow{}, &FlowList{})
}
