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

// +name:"Stdout"
// +url:"https://docs.fluentd.org/filter/stdout"
// +version:"more info"
// +description:"Prints events to stdout"
// +status:"GA"
type _metaStdOut interface{}

// +kubebuilder:object:generate=true

type StdOutFilterConfig struct {
}

func NewStdOutFilterConfig() *StdOutFilterConfig {
	return &StdOutFilterConfig{}
}

func (c *StdOutFilterConfig) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	pluginType := "stdout"
	return types.NewFlatDirective(types.PluginMeta{
		Type:      pluginType,
		Directive: "filter",
		Tag:       "**",
		Id:        id + "_" + pluginType,
	}, c, secretLoader)
}
