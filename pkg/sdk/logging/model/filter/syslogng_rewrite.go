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

// +name:"Syslog-NG rewrite"
// +weight:"200"
type _hugoRewrite interface{} //nolint:deadcode,unused

// +kubebuilder:object:generate=true
// +docName:"[Syslog-NG Rewrite Filter](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/77)"
// The syslog-ng rewrite filter can be used to replace message parts.
type _docRewrite interface{} //nolint:deadcode,unused

// +name:"Syslog-NG Rewrite"
// +url:"https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/77"
// +version:"more info"
// +description:"Rewrite parts of the message"
// +status:"GA"
type _metaRewrite interface{} //nolint:deadcode,unused

// +kubebuilder:object:generate=true
type RewriteConfig struct {
	Rename     *RenameConfig     `json:"rename,omitempty"`
	Set        *SetConfig        `json:"set,omitempty"`
	Substitute *SubstituteConfig `json:"subst,omitempty"`
	Unset      *UnsetConfig      `json:"unset,omitempty"`
}

func (c RewriteConfig) RenderAsSyslogNGConfig(ctx syslogng.Context) error {
	return syslogng.AllOf(
		syslogng.String("rewrite { "),
		syslogng.AllOf(
			syslogng.RenderIf(c.Rename != nil, c.Rename),
			syslogng.RenderIf(c.Set != nil, c.Set),
			syslogng.RenderIf(c.Substitute != nil, c.Substitute),
			syslogng.RenderIf(c.Unset != nil, c.Unset),
		),
		syslogng.String(" };"),
	).RenderAsSyslogNGConfig(ctx)
}

// +kubebuilder:object:generate=true
// https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/78#TOPIC-1829213
type RenameConfig struct {
	OldFieldName string     `json:"oldName"`
	NewFieldName string     `json:"newName"`
	Condition    *Condition `json:"condition,omitempty"`
}

func (c RenameConfig) RenderAsSyslogNGConfig(ctx syslogng.Context) error {
	return syslogng.AllOf(
		syslogng.Printf("rename(%q %q", c.OldFieldName, c.NewFieldName),
		syslogng.RenderIf(c.Condition != nil, syslogng.AllOf(
			syslogng.String(" "),
			c.Condition,
		)),
		syslogng.String(")"),
	).RenderAsSyslogNGConfig(ctx)
}

// +kubebuilder:object:generate=true
// https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/77#TOPIC-1829207
type SetConfig struct {
	FieldName string     `json:"field"` // NOTE: this is specified as `value(<field name>)` in the syslog-ng config
	Value     string     `json:"value"`
	Condition *Condition `json:"condition,omitempty"`
}

func (c SetConfig) RenderAsSyslogNGConfig(ctx syslogng.Context) error {
	return syslogng.AllOf(
		syslogng.Printf("set(%q value(%q)", c.Value, c.FieldName),
		syslogng.RenderIf(c.Condition != nil, syslogng.AllOf(
			syslogng.String(" "),
			c.Condition,
		)),
		syslogng.String(")"),
	).RenderAsSyslogNGConfig(ctx)
}

// +kubebuilder:object:generate=true
// https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/77#TOPIC-1829206
type SubstituteConfig struct {
	Pattern     string       `json:"pattern"`
	Replacement string       `json:"replace"`
	FieldName   string       `json:"field"`
	Flags       PatternFlags `json:"flags,omitempty"` // https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/81#TOPIC-1829224
	Type        string       `json:"type,omitempty"`  // https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/81#TOPIC-1829223
	Condition   *Condition   `json:"condition,omitempty"`
}

func (c SubstituteConfig) RenderAsSyslogNGConfig(ctx syslogng.Context) error {
	return syslogng.AllOf(
		syslogng.Printf("subst(%q %q value(%q)", c.Pattern, c.Replacement, c.FieldName),
		syslogng.RenderIf(len(c.Flags) > 0, c.Flags),
		syslogng.RenderIf(c.Type != "", syslogng.Printf("type(%q)", c.Type)),
		syslogng.RenderIf(c.Condition != nil, syslogng.AllOf(
			syslogng.String(" "),
			c.Condition,
		)),
		syslogng.String(")"),
	).RenderAsSyslogNGConfig(ctx)
}

// +kubebuilder:object:generate=true
// https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/78#TOPIC-1829212
type UnsetConfig struct {
	FieldName string     `json:"field"` // NOTE: this is specified as `value(<field name>)` in the syslog-ng config
	Condition *Condition `json:"condition,omitempty"`
}

func (c UnsetConfig) RenderAsSyslogNGConfig(ctx syslogng.Context) error {
	return syslogng.AllOf(
		syslogng.Printf("unset(value(%q)", c.FieldName),
		syslogng.RenderIf(c.Condition != nil, syslogng.AllOf(
			syslogng.String(" "),
			c.Condition,
		)),
		syslogng.String(")"),
	).RenderAsSyslogNGConfig(ctx)
}

// +kubebuilder:object:generate=true
type Condition MatchExpr

func (c Condition) RenderAsSyslogNGConfig(ctx syslogng.Context) error {
	return syslogng.AllOf(
		syslogng.String("condition("),
		MatchExpr(c),
		syslogng.String(")"),
	).RenderAsSyslogNGConfig(ctx)
}
