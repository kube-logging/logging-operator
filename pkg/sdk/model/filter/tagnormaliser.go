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

// +docName:"Fluentd Plugin to re-tag based on log metadata"
//More info at https://github.com/banzaicloud/fluent-plugin-tag-normaliser
//
//Available kubernetes metadata
//
//| Parameter | Description | Example |
//|-----------|-------------|---------|
//| ${pod_name} | Pod name | understood-butterfly-logging-demo-7dcdcfdcd7-h7p9n |
//| ${container_name} | Container name inside the Pod | logging-demo |
//| ${namespace_name} | Namespace name | default |
//| ${pod_id} | Kubernetes UUID for Pod | 1f50d309-45a6-11e9-b795-025000000001  |
//| ${labels} | Kubernetes Pod labels. This is a nested map. You can access nested attributes via `.`  | {"app":"logging-demo", "pod-template-hash":"7dcdcfdcd7" }  |
//| ${host} | Node hostname the Pod runs on | docker-desktop |
//| ${docker_id} | Docker UUID of the container | 3a38148aa37aa3... |
type _docTagNormaliser interface{}

// +name:"Tag Normaliser"
// +url:"https://github.com/banzaicloud/fluent-plugin-tag-normaliser"
// +version:"more info"
// +description:"Re-tag based on log metadata"
// +status:"GA"
type _metaTagNormaliser interface{}

// +docName:"Tag Normaliser parameters"
type TagNormaliser struct {
	// Re-Tag log messages info at [github](https://github.com/banzaicloud/fluent-plugin-tag-normaliser)
	Format string `json:"format,omitempty" plugin:"default:${namespace_name}.${pod_name}.${container_name}"`
}

func (t *TagNormaliser) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	pluginType := "tag_normaliser"
	return types.NewFlatDirective(types.PluginMeta{
		Type:      pluginType,
		Directive: "match",
		Tag:       "kubernetes.**",
		Id:        id + "_" + pluginType,
	}, t, secretLoader)
}
