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

// +docName:"File output plugin for Fluentd"
//This plugin has been designed to output logs or metrics to File.
//More info at https://docs.fluentd.org/output/file
//
// #### Example output configurations
// ```
// spec:
//  file:
//    path: /tmp/logs/${tag}/%Y/%m/%d.%H.%M
//    buffer:
//      timekey: 1m
//      timekey_wait: 10s
//      timekey_use_utc: true
// ```
type _docFile interface{}

// +name:"File"
// +url:"https://docs.fluentd.org/output/file"
// +version:"more info"
// +description:"Output plugin writes events to files"
// +status:"GA"
type _metaFile interface{}

// +kubebuilder:object:generate=true
type FileOutputConfig struct {
	// The Path of the file. The actual path is path + time + ".log" by default.
	Path string `json:"path"`
	// +docLink:"Buffer,./buffer.md"
	Buffer *Buffer `json:"buffer,omitempty"`
}

func (c *FileOutputConfig) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	pluginType := "file"
	pluginID := id + "_" + pluginType
	file := &types.OutputPlugin{
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
		file.Params = params
	}
	if c.Buffer != nil {
		if buffer, err := c.Buffer.ToDirective(secretLoader, pluginID); err != nil {
			return nil, err
		} else {
			file.SubDirectives = append(file.SubDirectives, buffer)
		}
	}
	return file, nil
}
