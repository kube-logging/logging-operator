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

// +name:"Syslog-NG parser"
// +weight:"200"
type _hugoParser interface{} //nolint:deadcode,unused

// +kubebuilder:object:generate=true
// +docName:"[Syslog-NG Parser Filter](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.36/administration-guide/90)"
type _docParser interface{} //nolint:deadcode,unused

// +name:"Syslog-NG Parser"
// +url:"https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.36/administration-guide/90"
// +version:"more info"
// +description:"Parse data from records"
// +status:"GA"
type _metaParser interface{} //nolint:deadcode,unused

// +kubebuilder:object:generate=true
// +docName:"[Parser](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.36/administration-guide/82#TOPIC-1768819)"
type ParserConfig struct {
	Regexp *RegexpParser `json:"regexp,omitempty" syslog-ng:"parser-drv,name=regexp-parser"`
}

// +kubebuilder:object:generate=true
// +docName:"[Regexp parser](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.36/administration-guide/90)"
type RegexpParser struct {
	// The regular expression patterns that you want to find a match. regexp-parser() supports multiple patterns, and stops the processing at the first successful match.
	Patterns []string `json:"patterns"`
	// Insert a prefix before the name part of the parsed name-value pairs to help further processing.
	Prefix string `json:"prefix,omitempty"`
	// Specify a template of the record fields to match against.
	Template string `json:"template,omitempty"`
	// Pattern flags
	Flags []string `json:"flags,omitempty"` // https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/81#TOPIC-1829224
}
