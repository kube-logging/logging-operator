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

package output

import (
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/secret"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/types"
	"github.com/banzaicloud/logging-operator/pkg/sdk/util"
)

// +docName:"Loki output plugin "
//Fluentd output plugin to ship logs to a Loki server.
//More info at https://github.com/banzaicloud/fluent-plugin-kubernetes-loki
//>Example: [Store Nginx Access Logs in Grafana Loki with Logging Operator](../../../docs/example-loki-nginx.md)
//
// #### Example output configurations
// ```
// spec:
//   loki:
//     url: http://loki:3100
//     buffer:
//       timekey: 1m
//       timekey_wait: 30s
//       timekey_use_utc: true
// ```
type _docLoki interface{}

// +name:"Grafana Loki"
// +url:"https://github.com/grafana/loki/tree/master/fluentd/fluent-plugin-grafana-loki"
// +version:"1.2.2"
// +description:"Transfer logs to Loki"
// +status:"GA"
type _metaLoki interface{}

// +kubebuilder:object:generate=true
// +docName:"Output Config"
type LokiOutput struct {
	// The url of the Loki server to send logs to. (default:https://logs-us-west1.grafana.net)
	Url string `json:"url,omitempty"`
	// Specify a username if the Loki server requires authentication.
	// +docLink:"Secret,./secret.md"
	Username *secret.Secret `json:"username,omitempty"`
	// Specify password if the Loki server requires authentication.
	// +docLink:"Secret,./secret.md"
	Password *secret.Secret `json:"password,omitempty"`
	// Loki is a multi-tenant log storage platform and all requests sent must include a tenant.
	Tenant string `json:"tenant,omitempty"`
	// Set of labels to include with every Loki stream.
	Labels Label `json:"labels,omitempty"`
	// Set of extra labels to include with every Loki stream.
	ExtraLabels map[string]string `json:"extra_labels,omitempty"`
	// Format to use when flattening the record to a log line: json, key_value (default: key_value)
	LineFormat string `json:"line_format,omitempty" plugin:"default:json"`
	// Extract kubernetes labels as loki labels (default: false)
	ExtractKubernetesLabels bool `json:"extract_kubernetes_labels,omitempty"`
	// Comma separated list of needless record keys to remove (default: [])
	RemoveKeys []string `json:"remove_keys,omitempty"`
	// If a record only has 1 key, then just set the log line to the value and discard the key. (default: false)
	DropSingleKey bool `json:"drop_single_key,omitempty"`
	// Configure Kubernetes metadata in a Prometheus like format
	ConfigureKubernetesLabels bool `json:"configure_kubernetes_labels,omitempty"`
	// +docLink:"Buffer,./buffer.md"
	Buffer *Buffer `json:"buffer,omitempty"`
}

type Label map[string]string

func (r Label) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	recordSet := types.PluginMeta{
		Directive: "label",
	}
	directive := &types.GenericDirective{
		PluginMeta: recordSet,
		Params:     r,
	}
	return directive, nil
}

func (r Label) merge(input Label) {
	for k, v := range input {
		r[k] = v
	}
}

func (l *LokiOutput) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	pluginType := "loki"
	pluginID := id + "_" + pluginType
	loki := &types.OutputPlugin{
		PluginMeta: types.PluginMeta{
			Type:      pluginType,
			Directive: "match",
			Tag:       "**",
			Id:        pluginID,
		},
	}
	if l.ConfigureKubernetesLabels {
		if l.Labels == nil {
			l.Labels = Label{}
		}
		l.Labels.merge(Label{
			"namespace":    `$.kubernetes.namespace_name`,
			"pod":          `$.kubernetes.pod_name`,
			"container_id": `$.kubernetes.docker_id`,
			"container":    `$.kubernetes.container_name`,
			"pod_id":       `$.kubernetes.pod_id`,
			"host":         `$.kubernetes.host`,
		})

		if l.RemoveKeys != nil {
			if !util.Contains(l.RemoveKeys, "kubernetes") {
				l.RemoveKeys = append(l.RemoveKeys, "kubernetes")
			}
		} else {
			l.RemoveKeys = []string{"kubernetes"}
		}
		l.ExtractKubernetesLabels = true
		// Prevent meta configuration from marshalling
		l.ConfigureKubernetesLabels = false
	}
	if params, err := types.NewStructToStringMapper(secretLoader).StringsMap(l); err != nil {
		return nil, err
	} else {
		loki.Params = params
	}
	if l.Labels != nil {
		if meta, err := l.Labels.ToDirective(secretLoader, ""); err != nil {
			return nil, err
		} else {
			loki.SubDirectives = append(loki.SubDirectives, meta)
		}
	}
	if l.Buffer != nil {
		if buffer, err := l.Buffer.ToDirective(secretLoader, pluginID); err != nil {
			return nil, err
		} else {
			loki.SubDirectives = append(loki.SubDirectives, buffer)
		}
	}
	return loki, nil
}
