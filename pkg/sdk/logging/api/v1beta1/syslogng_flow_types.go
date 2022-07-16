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

	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/filter"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/render/syslogng"
)

// +name:"SyslogNGFlowSpec"
// +weight:"200"
type _hugoSyslogNGFlowSpec interface{} //nolint:deadcode,unused

// +name:"SyslogNGFlowSpec"
// +version:"v1beta1"
// +description:"SyslogNGFlowSpec is the Kubernetes spec for SyslogNGFlows"
type _metaSyslogNGFlowSpec interface{} //nolint:deadcode,unused

// SyslogNGFlowSpec is the Kubernetes spec for SyslogNGFlows
type SyslogNGFlowSpec struct {
	Match            *SyslogNGMatch   `json:"match,omitempty"`
	Filters          []SyslogNGFilter `json:"filters,omitempty"`
	LoggingRef       string           `json:"loggingRef,omitempty"`
	GlobalOutputRefs []string         `json:"globalOutputRefs,omitempty"`
	LocalOutputRefs  []string         `json:"localOutputRefs,omitempty"`
}

type SyslogNGMatch filter.MatchConfig

// Filter definition for SyslogNGFlowSpec
type SyslogNGFilter struct {
	Match   *filter.MatchConfig   `json:"match,omitempty"`
	Rewrite *filter.RewriteConfig `json:"rewrite,omitempty"`
}

type SyslogNGFlowStatus FlowStatus

// +kubebuilder:object:root=true
// +kubebuilder:resource:categories=logging-all
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Active",type="boolean",JSONPath=".status.active",description="Is the flow active?"
// +kubebuilder:printcolumn:name="Problems",type="integer",JSONPath=".status.problemsCount",description="Number of problems"
// +kubebuilder:storageversion

// Flow Kubernetes object
type SyslogNGFlow struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SyslogNGFlowSpec   `json:"spec,omitempty"`
	Status SyslogNGFlowStatus `json:"status,omitempty"`
}

func (f SyslogNGFlow) RenderAsSyslogNGConfig(ctx syslogng.Context) error {
	const sourceName = "the_input"
	return syslogng.AllOf(
		syslogng.Printf("log flow_%s_%s {\n", f.Namespace, f.Name),
		syslogng.Indent(syslogng.AllOf(
			syslogng.Indentation(),
			syslogng.String("source("+sourceName+");\n"),
			syslogng.Indentation(),
			filter.MatchConfig{
				Not: (*filter.NotExpr)(&filter.MatchExpr{
					Regexp: &filter.RegexpMatchExpr{
						Pattern: f.Namespace,
						Value:   ".kubernetes.namespace_name",
						Type:    "string",
					},
				}),
			},
			syslogng.String(" # filter messages from other namespaces"),
			syslogng.String("\n"),
			syslogng.RenderIf(f.Spec.Match != nil, syslogng.AllOf(
				syslogng.Indentation(),
				(*filter.MatchConfig)(f.Spec.Match),
				syslogng.String(" # flow match"),
				syslogng.String("\n"),
			)),
			syslogng.ConfigRendererFunc(func(ctx syslogng.Context) error {
				for _, filter := range f.Spec.Filters {
					if err := syslogng.AllOf(
						syslogng.Indentation(),
						syslogng.RenderIf(filter.Match != nil, filter.Match),
						syslogng.RenderIf(filter.Rewrite != nil, filter.Rewrite),
						syslogng.String("\n"),
					).RenderAsSyslogNGConfig(ctx); err != nil {
						return err
					}
				}
				for _, ref := range f.Spec.LocalOutputRefs {
					c := syslogng.AllOf(
						syslogng.Indentation(),
						syslogng.Printf("destination(output_%s_%s);\n", f.Namespace, ref),
					)
					if err := c.RenderAsSyslogNGConfig(ctx); err != nil {
						return err
					}
				}
				for _, ref := range f.Spec.GlobalOutputRefs {
					c := syslogng.AllOf(
						syslogng.Indentation(),
						syslogng.Printf("destination(clusteroutput_%s_%s);\n", ctx.ControlNamespace, ref),
					)
					if err := c.RenderAsSyslogNGConfig(ctx); err != nil {
						return err
					}
				}
				return nil
			}),
		)),
		syslogng.String("};\n"),
	).RenderAsSyslogNGConfig(ctx)
}

// +kubebuilder:object:root=true

// FlowList contains a list of Flow
type SyslogNGFlowList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SyslogNGFlow `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SyslogNGFlow{}, &SyslogNGFlowList{})
}
