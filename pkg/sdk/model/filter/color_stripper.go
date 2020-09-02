// Copyright © 2020 Banzai Cloud
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
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/types"
	"github.com/banzaicloud/operator-tools/pkg/secret"
)

// +name:"Color Stripper"
// +weight:"200"
type _hugoColorStripper interface{}

// +kubebuilder:object:generate=true
// +docName:"[Color Stripper](https://github.com/mattheworiordan/fluent-plugin-color-stripper)"
// Fluentd Filter plugin to strip ANSI color from log lines.
type _docColorStripper interface{}

// +name:"Color Stripper"
// +url:"https://github.com/mattheworiordan/fluent-plugin-color-stripper"
// +version:"1.0.0"
// +description:"Strips ANSI color from log lines"
// +status:"GA"
type _metaColorStripper interface{}

// +kubebuilder:object:generate=true
type ColorStripper struct {
	// A comma-delimited list of keys to strip colors from
	StripFields string `json:"strip_fields,omitempty"`
	// Tag to apply to output
	Tag string `json:"tag,omitempty" plugin:"default:formatted"`
}

// #### Example `Color Stripper` filter configurations
// ```yaml
//apiVersion: logging.banzaicloud.io/v1beta1
//kind: Flow
//metadata:
//  name: demo-flow
//spec:
//  filters:
//    - colorStripper: {}
//  selectors: {}
//  outputRefs:
//    - demo-output
// ```
//
// #### Fluentd Config Result
// ```yaml
//<filter **>
//  @type color_stripper
//  @id test_color_stripper
//  tag formatted
//</filter>
// ```
type _expColorStripper interface{}

func (c *ColorStripper) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	pluginType := "color_stripper"
	return types.NewFlatDirective(types.PluginMeta{
		Type:      pluginType,
		Directive: "filter",
		Tag:       "**",
		Id:        id,
	}, c, secretLoader)
}
