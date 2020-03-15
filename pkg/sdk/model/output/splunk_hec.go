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
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/types"
	"github.com/banzaicloud/operator-tools/pkg/secret"
)

// +docName:"Splunk via Hec output plugin for Fluentd"
//More info at https://github.com/splunk/fluent-plugin-splunk-hec
//
// #### Example output configurations
// ```
// spec:
//   SplunkHec:
//     host: splunk.default.svc.cluster.local
//     port: 8088
//     protocol: http
// ```
type _docSplunkHec interface{}

// +name:"Splunk Hec"
// +url:"https://github.com/splunk/fluent-plugin-splunk-hec/releases/tag/1.2.1
// +version:"1.2.1"
// +description:"Fluent Plugin Splunk Hec Release 1.2.1"
// +status:"GA"
type _metaSplunkHec interface{}

// +kubebuilder:object:generate=true
// +docName:"SplunkHecOutput"
// SplunkHecOutput sends your logs to Splunk via Hec
type SplunkHecOutput struct {
	// You can specify SplunkHec host by this parameter.
	HecHost string `json:"hec_host"`
	// The port number for the Hec token or the Hec load balancer. (default:8088)
	HecPort int `json:"hec_port,omitempty"`
	// This is the protocol to use for calling the Hec API. Available values are: http, https. (default:https)
	Protocol string `json:"protocol,omitempty"`
	// Identifier for the Hec token.
	// +docLink:"Secret,./secret.md"
	HecToken *secret.Secret `json:"hec_token"`
	// When data_type is set to "metric", the ingest API will treat every key-value pair in the input event as a metric name-value pair. Set metrics_from_event to false to disable this behavior and use metric_name_key and metric_value_key to define metrics. (Default:true)
	MetricsFromEvent bool `json:"metrics_from_event,omitempty"`
	// Field name that contains the metric name. This parameter only works in conjunction with the metrics_from_event paramter. When this prameter is set, the metrics_from_event parameter is automatically set to false.
	MetricsNameKey string `json:"metrics_name_key,omitempty"`
	// Field name that contains the metric value, this parameter is required when metric_name_key is configured.
	MetricsValueKey string `json:"metrics_value_key,omitempty"`
}

func (c *SplunkHecOutput) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	pluginType := "splunk_hec"
	pluginID := id + "_" + pluginType
	SplunkHec := &types.OutputPlugin{
		PluginMeta: types.PluginMeta{
			Type:      pluginType,
			Directive: "match",
			Tag:       "**",
			Id:        pluginID,
		},
	}
	if params, err := types.NewStructToStringMapper(secretLoader).StringsMap(c); err != nil {
		return nil, err
	} else {
		SplunkHec.Params = params
	}
	return SplunkHec, nil
}
