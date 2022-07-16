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

package filter

import (
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/render/syslogng"
)

// +name:"Syslog-NG match"
// +weight:"200"
type _hugoMatch interface{} //nolint:deadcode,unused

// +kubebuilder:object:generate=true
// +docName:"[Syslog-NG Match Filter](https://docs.fluentd.org/filter/grep)"
// The match filter can be used to selectively keep records
type _docMatch interface{} //nolint:deadcode,unused

// +name:"Syslog-NG Match"
// +url:"https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/65#TOPIC-1829159"
// +version:"more info"
// +description:"Selectively keep records"
// +status:"GA"
type _metaMatch interface{} //nolint:deadcode,unused

// +kubebuilder:object:generate=true
type MatchConfig MatchExpr

func (c MatchConfig) RenderAsSyslogNGConfig(ctx syslogng.Context) error {
	return syslogng.AllOf(
		syslogng.String("filter { "),
		MatchExpr(c),
		syslogng.String(" };"),
	).RenderAsSyslogNGConfig(ctx)
}

// +kubebuilder:object:generate=true
type MatchExpr struct {
	// +docLink:"And Directive,#And-Directive"
	And AndExpr `json:"and,omitempty"`
	// +docLink:"Regexp Directive,#Regexp-Directive"
	Regexp *RegexpMatchExpr `json:"regexp,omitempty"`
	// +docLink:"Exclude Directive,#Exclude-Directive"
	Not *NotExpr `json:"not,omitempty"`
	// +docLink:"Or Directive,#Or-Directive"
	Or OrExpr `json:"or,omitempty"`
}

func (e MatchExpr) RenderAsSyslogNGConfig(ctx syslogng.Context) error {
	return syslogng.AllOf(
		syslogng.RenderIf(len(e.And) > 0, e.And),
		syslogng.RenderIf(e.Regexp != nil, e.Regexp),
		syslogng.RenderIf(e.Not != nil, e.Not),
		syslogng.RenderIf(len(e.Or) > 0, e.Or),
	).RenderAsSyslogNGConfig(ctx)
}

// +kubebuilder:object:generate=true
// +docName:"[Regexp Directive](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/68#TOPIC-1829171) {#Regexp-Directive}"
// Specify filtering rule.
type RegexpMatchExpr struct {
	// Pattern expression to evaluate
	Pattern string `json:"pattern"`
	// Specify a template of the record fields to match against.
	Template string `json:"template,omitempty"`
	// Specify a field name of the record to match against the value of.
	Value string `json:"value,omitempty"`
	// Pattern flags
	Flags PatternFlags `json:"flags,omitempty"` // https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/81#TOPIC-1829224
	// Pattern type
	Type string `json:"type,omitempty"` // https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/81#TOPIC-1829223
}

func (e RegexpMatchExpr) RenderAsSyslogNGConfig(ctx syslogng.Context) error {
	return syslogng.AllOf(
		syslogng.Printf("match(%q", e.Pattern),
		syslogng.RenderIf(e.Template != "", syslogng.Printf(" template(%q)", e.Template)),
		syslogng.RenderIf(e.Value != "", syslogng.Printf(" value(%q)", e.Value)),
		syslogng.RenderIf(len(e.Flags) > 0, syslogng.AllOf(syslogng.String(" "), e.Flags)),
		syslogng.RenderIf(e.Type != "", syslogng.Printf(" type(%q)", e.Type)),
		syslogng.String(")"),
	).RenderAsSyslogNGConfig(ctx)
}

// #### Example `Regexp` filter configurations
// ```yaml
//apiVersion: logging.banzaicloud.io/v1beta1
//kind: Flow
//metadata:
//  name: demo-flow
//spec:
//  filters:
//    - match:
//        regexp:
//        - value: first
//          pattern: ^5\d\d$
//  match: {}
//  localOutputRefs:
//    - demo-output
// ```
//
// #### Syslog-NG Config Result
// ```
//    filter {
//        match("^5\d\d$" value("first"));
//    }
// ```
type _expRegexpMatch interface{} //nolint:deadcode,unused

// +kubebuilder:object:generate=true
type AndExpr []MatchExpr

func (e AndExpr) RenderAsSyslogNGConfig(ctx syslogng.Context) error {
	const operator = " and "
	return syslogng.AllOf(
		syslogng.String("("),
		// TODO: fugly
		syslogng.ConfigRendererFunc(func(ctx syslogng.Context) error {
			for i, e := range e {
				if err := syslogng.AllOf(
					syslogng.RenderIf(i > 0, syslogng.String(operator)),
					syslogng.NoIndent(e),
				).RenderAsSyslogNGConfig(ctx); err != nil {
					return err
				}
			}
			return nil
		}),
		syslogng.String(")"),
	).RenderAsSyslogNGConfig(ctx)
}

// +kubebuilder:object:generate=true
type NotExpr MatchExpr

func (e *NotExpr) RenderAsSyslogNGConfig(ctx syslogng.Context) error {
	if e != nil {
		return syslogng.AllOf(
			syslogng.String("(not "),
			(*MatchExpr)(e),
			syslogng.String(")"),
		).RenderAsSyslogNGConfig(ctx)
	}
	return nil
}

// +kubebuilder:object:generate=true
type OrExpr []MatchExpr

func (e OrExpr) RenderAsSyslogNGConfig(ctx syslogng.Context) error {
	const operator = " or "
	return syslogng.AllOf(
		syslogng.String("("),
		// TODO: fugly
		syslogng.ConfigRendererFunc(func(ctx syslogng.Context) error {
			for i, e := range e {
				if err := syslogng.AllOf(
					syslogng.RenderIf(i > 0, syslogng.String(operator)),
					syslogng.NoIndent(e),
				).RenderAsSyslogNGConfig(ctx); err != nil {
					return err
				}
			}
			return nil
		}),
		syslogng.String(")"),
	).RenderAsSyslogNGConfig(ctx)
}

// +kubebuilder:object:generate=true
type PatternFlags []string

func (fs PatternFlags) RenderAsSyslogNGConfig(ctx syslogng.Context) error {
	return syslogng.AllOf(
		syslogng.String("flags("),
		// TODO: fugly
		syslogng.ConfigRendererFunc(func(ctx syslogng.Context) error {
			for i, f := range fs {
				if err := syslogng.AllOf(
					syslogng.RenderIf(i > 0, syslogng.String(" ")),
					syslogng.Printf("%q", string(f)),
				).RenderAsSyslogNGConfig(ctx); err != nil {
					return err
				}
			}
			return nil
		}),
		syslogng.String(")"),
	).RenderAsSyslogNGConfig(ctx)
}
