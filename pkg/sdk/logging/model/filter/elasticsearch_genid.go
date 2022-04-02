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
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/types"
	"github.com/banzaicloud/operator-tools/pkg/secret"
)

// +kubebuilder:object:generate=true
type ElasticsearchGenId struct {
	// Separator (default:_)
	Hash_id_key string `json:"hash_id_key,omitempty"`
}

// #### Example `Dedot` filter configurations
// ```yaml
//apiVersion: logging.banzaicloud.io/v1beta1
//kind: Flow
//metadata:
//  name: demo-flow
//spec:
//  filters:
//    - dedot:
//        de_dot_separator: "-"
//        de_dot_nested: true
//  selectors: {}
//  localOutputRefs:
//    - demo-output
// ```
//
// #### Fluentd Config Result
// ```yaml
//<filter **>
//  @type dedot
//  @id test_dedot
//  de_dot_nested true
//  de_dot_separator -
//</filter>
// ```

func NewElasticsearchGenId() *ElasticsearchGenId {
	return &ElasticsearchGenId{}
}

func (c *ElasticsearchGenId) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	const pluginType = "elasticsearch_genid"
	return types.NewFlatDirective(types.PluginMeta{
		Type:      pluginType,
		Directive: "filter",
		Tag:       "**",
		Id:        id,
	}, c, secretLoader)
}
