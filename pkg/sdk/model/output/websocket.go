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

// +name:"Websocket"
// +url:"https://github.com/banzaicloud/fluent-plugin-websocket"
// +version:"0.1.8"
// +description:"Websocket server to send logs to connected clients"
// +status:"GA"
type _metaWebsocket interface{}

// +docName:"Websocket output plugin for Fluentd"
// This plugin works as websocket server which can output JSON string or MessagePack binary.
//
// More info and examples at https://github.com/banzaicloud/fluent-plugin-websocket
type _docWebsocket interface{}

// +kubebuilder:object:generate=true
// +docName:"Output Config"
type WebsocketOutput struct {
	// WebSocket server IP address. (default: 0.0.0.0 (ANY))
	Host string `json:"host,omitempty"`
	// WebSocket server port. (default: 8080)
	Port int `json:"port,omitempty"`
	// Send MessagePack format binary. Otherwise, you send JSON format text. (default: false)
	UseMsgpack *bool `json:"use_msgpack,omitempty"`
	// Add timestamp to the data. (default: false)
	AddTime *bool `json:"add_time,omitempty"`
	// Add fluentd tag to the data. (default: true)
	AddTag *bool `json:"add_tag,omitempty"`
	// The number of messages to be buffered. The new connection receives them. (default: 0)
	BufferedMessages int `json:"buffered_messages,omitempty"`
	// Authentication token. Passed as get param. If set to nil, authentication is disabled. (default: nil)
	Token *secret.Secret `json:"token,omitempty"`
}

func (w *WebsocketOutput) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	pluginType := "websocket"
	pluginID := id + "_" + pluginType
	return types.NewFlatDirective(types.PluginMeta{
		Type:      pluginType,
		Directive: "match",
		Tag:       "**",
		Id:        pluginID,
	}, w, secretLoader)
}
