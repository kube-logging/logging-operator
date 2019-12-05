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

// +docName:"Exception Detector"
//This filter plugin consumes a log stream of JSON objects which contain single-line log messages. If a consecutive sequence of log messages form an exception stack trace, they forwarded as a single, combined JSON object. Otherwise, the input log data is forwarded as is.
//More info at https://github.com/GoogleCloudPlatform/fluent-plugin-detect-exceptions
//
// > Note: As Tag management is not supported yet, this Plugin is **mutually exclusive** with [Tag normaliser](./tagnormaliser.md)
//
// #### Example output configurations
// ```
//filters:
//  spec:
//    detectExceptions:
//      languages: java, python
//      multiline_flush_interval: 0.1
// ```
type _docExceptionDetector interface{}

// +name:"Exception Detector"
// +url:"https://github.com/GoogleCloudPlatform/fluent-plugin-detect-exceptions"
// +version:"more info"
// +description:"Exception Detector"
// +status:"GA"
type _metaDDetectExceptions interface{}

// +kubebuilder:object:generate=true
// +docName:"Output Config"
type DetectExceptions struct {
	// The field which contains the raw message text in the input JSON data. (default: "")
	Message string `json:"message,omitempty"`
	// The prefix to be removed from the input tag when outputting a record. (default: "")
	RemoveTagPrefix string `json:"remove_tag_prefix,omitempty"`
	// The interval of flushing the buffer for multiline format. (default: nil)
	MultilineFlushInterval string `json:"multiline_flush_interval,omitempty"`
	// Programming languages for which to detect exceptions. (default: [])
	Languages []string `json:"languages,omitempty"`
	// Maximum number of lines to flush (0 means no limit) (default: 1000)
	MaxLines int `json:"max_lines,omitempty"`
	// Maximum number of bytes to flush (0 means no limit) (default: 0)
	MaxBytes int `json:"max_bytes,omitempty"`
	// Separate log streams by this field in the input JSON data. (default: "")
	Stream string `json:"stream,omitempty"`
}

func (d *DetectExceptions) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	pluginType := "detect_exceptions"
	pluginID := id + "_" + pluginType
	detector := &types.OutputPlugin{
		PluginMeta: types.PluginMeta{
			Type:      pluginType,
			Directive: "match",
			Tag:       "kubernetes.**",
			Id:        pluginID,
		},
	}
	detect := d.DeepCopy()
	detect.RemoveTagPrefix = "kubernetes"
	if params, err := types.NewStructToStringMapper(secretLoader).StringsMap(detect); err != nil {
		return nil, err
	} else {
		detector.Params = params
	}
	return detector, nil
}
