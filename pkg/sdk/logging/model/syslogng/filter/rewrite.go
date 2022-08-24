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
	GroupUnset *GroupUnsetConfig `json:"group_unset,omitempty" syslog-ng:"rewrite-drv,name=groupunset"`
	Rename     *RenameConfig     `json:"rename,omitempty" syslog-ng:"rewrite-drv,name=rename"`
	Set        *SetConfig        `json:"set,omitempty" syslog-ng:"rewrite-drv,name=set"`
	Substitute *SubstituteConfig `json:"subst,omitempty" syslog-ng:"rewrite-drv,name=subst"`
	Unset      *UnsetConfig      `json:"unset,omitempty" syslog-ng:"rewrite-drv,name=unset"`
}

// +kubebuilder:object:generate=true
// https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/78#TOPIC-1829213
type RenameConfig struct {
	OldFieldName string     `json:"oldName" syslog-ng:"pos=0"`
	NewFieldName string     `json:"newName" syslog-ng:"pos=1"`
	Condition    *MatchExpr `json:"condition,omitempty" syslog-ng:"name=condition,optional"`
}

// +kubebuilder:object:generate=true
// https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/77#TOPIC-1829207
type SetConfig struct {
	FieldName string     `json:"field" syslog-ng:"name=value"` // NOTE: this is specified as `value(<field name>)` in the syslog-ng config
	Value     string     `json:"value" syslog-ng:"pos=0"`
	Condition *MatchExpr `json:"condition,omitempty" syslog-ng:"name=condition,optional"`
}

// +kubebuilder:object:generate=true
// https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/77#TOPIC-1829206
type SubstituteConfig struct {
	Pattern     string     `json:"pattern" syslog-ng:"pos=0"`
	Replacement string     `json:"replace" syslog-ng:"pos=1"`
	FieldName   string     `json:"field" syslog-ng:"name=value"`
	Flags       []string   `json:"flags,omitempty" syslog-ng:"name=flags,optional"` // https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/81#TOPIC-1829224
	Type        string     `json:"type,omitempty" syslog-ng:"name=type,optional"`   // https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/81#TOPIC-1829223
	Condition   *MatchExpr `json:"condition,omitempty" syslog-ng:"name=condition,optional"`
}

// +kubebuilder:object:generate=true
// https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/78#TOPIC-1829212
type UnsetConfig struct {
	FieldName string     `json:"field" syslog-ng:"name=value"` // NOTE: this is specified as `value(<field name>)` in the syslog-ng config
	Condition *MatchExpr `json:"condition,omitempty" syslog-ng:"name=condition,optional"`
}

// +kubebuilder:object:generate=true
// https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/78#TOPIC-1829212
type GroupUnsetConfig struct {
	Pattern   string     `json:"pattern" syslog-ng:"name=values"` // NOTE: this is specified as `value(<field name>)` in the syslog-ng config
	Condition *MatchExpr `json:"condition,omitempty" syslog-ng:"name=condition,optional"`
}
