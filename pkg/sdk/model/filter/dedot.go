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
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/secret"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/types"
)

// +name:"Dedot"
// +url:"https://github.com/lunardial/fluent-plugin-dedot_filter"
// +version:"more info"
// +description:"Concatenate multiline log separated in multiple events"
// +status:"GA"
type _metaDedot interface{}

// +kubebuilder:object:generate=true

// +docName:"Fluentd Filter plugin to de-dot field name for elasticsearch."
// More info at https://github.com/lunardial/fluent-plugin-dedot_filter
type DedotFilterConfig struct {
	// Will cause the plugin to recurse through nested structures (hashes and arrays), and remove dots in those key-names too.
	Nested bool `json:"de_dot_nested,omitempty" plugin:"default:true"`

	// Separator (default:_)
	Separator string `json:"de_dot_separator,omitempty"`
}

func NewDedotFilterConfig() *DedotFilterConfig {
	return &DedotFilterConfig{}
}

func (c *DedotFilterConfig) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	pluginType := "dedot"
	return types.NewFlatDirective(types.PluginMeta{
		Type:      pluginType,
		Directive: "filter",
		Tag:       "**",
		Id:        id + "_" + pluginType,
	}, c, secretLoader)
}
