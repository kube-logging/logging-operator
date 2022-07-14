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

// +kubebuilder:object:generate=true
type MatchExpr struct {
	// +docLink:"And Directive,#And-Directive"
	And []MatchExpr `json:"and,omitempty"`
	// +docLink:"Regexp Directive,#Regexp-Directive"
	Regexp *RegexpMatchExpr `json:"regexp,omitempty"`
	// +docLink:"Exclude Directive,#Exclude-Directive"
	Not []MatchExpr `json:"not,omitempty"`
	// +docLink:"Or Directive,#Or-Directive"
	Or []MatchExpr `json:"or,omitempty"`
}

// +kubebuilder:object:generate=true
// +docName:"[Regexp Directive](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/68#TOPIC-1829171) {#Regexp-Directive}"
// Specify filtering rule.
type RegexpMatchExpr struct {
	// Pattern expression to evaluate
	Pattern string `json:"pattern"`
	// Specify a template of the record fields to match against.
	Template string `json:"template"`
	// Specify a field name of the record to match against the value of.
	Value string `json:"value"`
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
