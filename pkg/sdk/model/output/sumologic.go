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
)

// +docName:"SumoLogic output plugin for Fluentd"
//This plugin has been designed to output logs or metrics to SumoLogic via a HTTP collector endpoint
//More info at https://github.com/SumoLogic/fluentd-output-sumologic
type _docSumoLogic interface{}

// +name:"SumoLogic"
// +url:"https://github.com/SumoLogic/fluentd-output-sumologic/releases/tag/1.6.1"
// +version:"0.6.1"
// +description:"Send your logs to Sumologic"
// +status:"GA"
type _metaSumologic interface{}

// +kubebuilder:object:generate=true
// +docName:"Output Config"
type SumologicOutput struct {
	// The type of data that will be sent to Sumo Logic, either logs or metrics (default: logs)
	DataType string `json:"data_type,omitempty"`
	// SumoLogic HTTP Collector URL
	Endpoint *secret.Secret `json:"endpoint"`
	// Verify ssl certificate. (default: true)
	VerifySsl bool `json:"verify_ssl,omitempty"`
	// The format of metrics you will be sending, either graphite or carbon2 or prometheus (default: graphite)
	MetricDataFormat string `json:"metric_data_format,omitempty"`
	// Format to post logs into Sumo. (default: json)
	LogFormat string `json:"log_format,omitempty"`
	// Used to specify the key when merging json or sending logs in text format (default: message)
	LogKey string `json:"log_key,omitempty"`
	// Set _sourceCategory metadata field within SumoLogic (default: nil)
	SourceCategory string `json:"source_category,omitempty"`
	// Set _sourceName metadata field within SumoLogic - overrides source_name_key (default is nil)
	SourceName string `json:"source_name"`
	// Set as source::path_key's value so that the source_name can be extracted from Fluentd's buffer (default: source_name)
	SourceNameKey string `json:"source_name_key,omitempty"`
	// Set _sourceHost metadata field within SumoLogic (default: nil)
	SourceHost string `json:"source_host,omitempty"`
	// Set timeout seconds to wait until connection is opened. (default: 60)
	OpenTimeout int `json:"open_timeout,omitempty"`
	// Add timestamp (or timestamp_key) field to logs before sending to sumologic (default: true)
	AddTimestamp bool `json:"add_timestamp,omitempty"`
	// Field name when add_timestamp is on (default: timestamp)
	TimestampKey string `json:"timestamp_key,omitempty"`
	// Add the uri of the proxy environment if present.
	ProxyUri string `json:"proxy_uri,omitempty"`
	// Option to disable cookies on the HTTP Client. (default: false)
	DisableCookies bool `json:"disable_cookies,omitempty"`
}

func (s *SumologicOutput) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	pluginType := "sumologic"
	pluginID := id + "_" + pluginType
	return types.NewFlatDirective(types.PluginMeta{
		Type:      pluginType,
		Directive: "match",
		Tag:       "**",
		Id:        pluginID,
	}, s, secretLoader)
}
