package filter

import (
	"github.com/banzaicloud/logging-operator/pkg/model/secret"
	"github.com/banzaicloud/logging-operator/pkg/model/types"
)

// +docName:"Fluentd Plugin to re-tag based on log metadata"
//More info at https://github.com/banzaicloud/fluent-plugin-tag-normaliser
//
//Available kubernetes metadata
//
//| Parameter | Description | Example |
//|-----------|-------------|---------|
//| ${pod_name} | Pod name | understood-butterfly-nginx-logging-demo-7dcdcfdcd7-h7p9n |
//| ${container_name} | Container name inside the Pod | nginx-logging-demo |
//| ${namespace_name} | Namespace name | default |
//| ${pod_id} | Kubernetes UUID for Pod | 1f50d309-45a6-11e9-b795-025000000001  |
//| ${labels} | Kubernetes Pod labels. This is a nested map. You can access nested attributes via `.`  | {"app":"nginx-logging-demo", "pod-template-hash":"7dcdcfdcd7" }  |
//| ${host} | Node hostname the Pod runs on | docker-desktop |
//| ${docker_id} | Docker UUID of the container | 3a38148aa37aa3... |
type _docTagNormaliser interface{}

// +docName:"Tag Normaliser parameters"
type TagNormaliser struct {
	// Re-Tag log messages info at [github](https://github.com/banzaicloud/fluent-plugin-tag-normaliser)
	Format string `json:"format,omitempty" plugin:"default:${namespace_name}.${pod_name}.${container_name}"`
}

func (t *TagNormaliser) ToDirective(secretLoader secret.SecretLoader) (types.Directive, error) {
	return types.NewFlatDirective(types.PluginMeta{
		Type:      "tag_normaliser",
		Directive: "match",
		Tag:       "kubernetes.**",
	}, t, secretLoader)
}
