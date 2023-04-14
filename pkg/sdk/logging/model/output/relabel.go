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

package output

import (
	"github.com/cisco-open/operator-tools/pkg/secret"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/types"
)

// +name:"Relabel"
// +weight:"200"
type _hugoRelabel interface{} //nolint:deadcode,unused

// +name:"Relabel"
// +url:"https://docs.fluentd.org/output/relabel"
// +version:"more info"
// +description:"Relabel output plugin re-labels events."
// +status:"GA"
type _metaRelabel interface{} //nolint:deadcode,unused

// +kubebuilder:object:generate=true
// +docName:"Output Config"
type RelabelOutputConfig struct {
	// Specifies new label for events
	Label string `json:"label"`
}

func (c *RelabelOutputConfig) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	const pluginType = "relabel"
	relabel := &types.OutputPlugin{
		PluginMeta: types.PluginMeta{
			Type:      pluginType,
			Directive: "match",
			Tag:       "**",
			Id:        id,
			Label:     c.Label,
		},
	}
	if params, err := types.NewStructToStringMapper(secretLoader).StringsMap(c); err != nil {
		return nil, err
	} else {
		relabel.Params = params
		delete(relabel.Params, "label")
	}
	return relabel, nil
}
