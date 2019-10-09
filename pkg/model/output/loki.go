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
	"github.com/banzaicloud/logging-operator/pkg/model/secret"
	"github.com/banzaicloud/logging-operator/pkg/model/types"
)

// +docName:"Loki output plugin "
//Fluentd output plugin to ship logs to a Loki server.
//More info at https://github.com/banzaicloud/fluent-plugin-kubernetes-loki
type _docLoki interface{}

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
	// Set of labels to include with every Loki stream.(default: nil)
	ExtraLabels bool `json:"extra_labels,omitempty"`
	// +docLink:"Buffer,./buffer.md"
	Buffer *Buffer `json:"buffer,omitempty"`
}

func (l *LokiOutput) ToDirective(secretLoader secret.SecretLoader) (types.Directive, error) {
	loki := &types.OutputPlugin{
		PluginMeta: types.PluginMeta{
			Type:      "kubernetes_loki",
			Directive: "match",
			Tag:       "**",
		},
	}
	if params, err := types.NewStructToStringMapper(secretLoader).StringsMap(l); err != nil {
		return nil, err
	} else {
		loki.Params = params
	}
	if l.Buffer != nil {
		if buffer, err := l.Buffer.ToDirective(secretLoader); err != nil {
			return nil, err
		} else {
			loki.SubDirectives = append(loki.SubDirectives, buffer)
		}
	}
	return loki, nil
}
