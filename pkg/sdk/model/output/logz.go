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
	"fmt"

	"github.com/banzaicloud/logging-operator/pkg/sdk/model/types"
	"github.com/banzaicloud/operator-tools/pkg/secret"
)

// +docName:"LogZ output plugin for Fluentd"
//More info at https://github.com/logzio/fluent-plugin-logzio
//>Example Deployment: [Save all logs to LogZ](../../../docs/example-logz.md)
//
// #### Example output configurations
// ```
// spec:
//   logz:
//     url: https://listener.logz.io
//     port: 8071
//     token: 12345
//     output:
//       include_time: true
//       include_tags: false
//       tags_fieldname: @log_name
//     buffer:
//       type: file
//       flush_thread_count: 4
//       flush_interval: 3s
//       chunk_limit_size: 16m
//       queue_limit_length: 4096
//     http:
//       idle_timeout: 10
// ```
type _docLogZ interface{}

// +name:"LogZ"
// +url:"https://github.com/uken/fluent-plugin-elasticsearch/releases/tag/v3.7.0"
// +version:"3.7.0"
// +description:"Send your logs to LogZ"
// +status:"GA"
type _metaLogZ interface{}

// +kubebuilder:object:generate=true
// +docName:"Logzio"
// LogZ Send your logs to LogZ.io
type LogZOutput struct {
	// +docLink:"Shared Credentials,#Shared-Credentials"
	Endpoint *Endpoint `json:"endpoint"`
	// Token (LogZ API Token).
	// +docLink:"Secret,./secret.md"
	//Token             *secret.Secret `json:"token,omitempty"`
	OutputIncludeTime bool `json:"output_include_time,omitempty"`
	OutputIncludeTags bool `json:"output_include_tags,omitempty"`
	HTTPIdleTimeout   int  `json:"http_idle_timeout,omitempty"`

	// +docLink:"Buffer,./buffer.md"
	Buffer *Buffer `json:"buffer,omitempty"`
}

// Endpoint thing
// +kubebuilder:object:generate=true
// +docName:"Endpoint"
type Endpoint struct {
	URL  string `json:"url,omitempty" plugin:"default:"https://listener.logz.io"`
	Port int    `json:"port,omitempty" plugin:"default:8071"`
	// LogZ API Token.
	// +docLink:"Secret,./secret.md"
	Token *secret.Secret `json:"token,omitempty"`
}

// ToDirective converts Endpoint struct to fluentd configuration.
func (e *Endpoint) ToDirective(secretLoader secret.SecretLoader) (types.Directive, error) {
	metadata := types.PluginMeta{}

	endpoint := e.DeepCopy()
	return types.NewFlatDirective(metadata, endpoint, secretLoader)
}

func (e *LogZOutput) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	pluginType := "logzio_buffered"
	pluginID := id + "_" + pluginType
	logz := &types.OutputPlugin{
		PluginMeta: types.PluginMeta{
			Type:      pluginType,
			Directive: "match",
			Tag:       "**",
			Id:        pluginID,
		},
	}
	if params, err := types.NewStructToStringMapper(secretLoader).StringsMap(e); err != nil {
		return nil, err
	} else {
		logz.Params = params
	}

	// Construct endpoint_url parameter.
	url := "https://listener.logz.io"
	port := 8071
	connectionString := ""
	if e.Endpoint != nil {
		if e.Endpoint.Port != 0 {
			port = e.Endpoint.Port
		}
		if e.Endpoint.URL != "" {
			url = e.Endpoint.URL
		}
		connectionString = fmt.Sprintf("%s:%d?type=my_type", url, port)

		// decrypt token secrent
		endpoint, err := e.Endpoint.ToDirective(secretLoader)
		if err != nil {
			return nil, err
		}
		if token, ok := endpoint.GetParams()["token"]; ok {
			connectionString = fmt.Sprintf("%s&token=%s", connectionString, token)
		}
	}

	//n := map[string]string{"endpoint_url": connectionString}
	// add endpoint_url to parameters
	logz.Params.Merge(map[string]string{"endpoint_url": connectionString})

	// logz.Params = params
	if e.Buffer != nil {
		if buffer, err := e.Buffer.ToDirective(secretLoader, pluginID); err != nil {
			return nil, err
		} else {
			logz.SubDirectives = append(logz.SubDirectives, buffer)
		}
	}
	return logz, nil
}
