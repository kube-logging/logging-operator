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
	"github.com/banzaicloud/logging-operator/pkg/model/secret"
	"github.com/banzaicloud/logging-operator/pkg/model/types"
)

// +docName:"Fluentd concat filter"
// Fluentd Filter plugin to concatenate multiline log separated in multiple events.
// More information at https://github.com/fluent-plugins-nursery/fluent-plugin-concat
//
// #### Example record configurations
// ```
// spec:
//  filters:
//    - concat:
//        key: message
//        partial_key: "partial_message"
//        separator: ""
// ```
type _docConcat interface{}

// +kubebuilder:object:generate=true
type Concat struct {
	//The key for part of multiline log.
	Key string `json:"key"`
	//The separator of lines. (default: "\n")
	Separator string `json:"separator,omitempty"`
	//The number of lines. This is exclusive with multiline_start_regex.
	NLines int `json:"n_lines,omitempty"`
	//The regexp to match beginning of multiline. This is exclusive with n_lines.
	MultilineStartRegexp string `json:"multiline_start_regexp,omitempty"`
	//The regexp to match ending of multiline. This is exclusive with n_lines.
	MultilineEndRegexp string `json:"multiline_end_regexp,omitempty"`
	//The regexp to match continuous lines. This is exclusive with n_lines.
	ContinuousLineRegexp string `json:"continuous_line_regexp,omitempty"`
	//The key to determine which stream an event belongs to.
	StreamIdentityKey string `json:"stream_identity_key,omitempty"`
	//The number of seconds after which the last received event log will be flushed. If specified 0, wait for next line forever.
	FlushInterval int `json:"flush_interval,omitempty"`
	//The label name to handle events caused by timeout.
	TimeoutLabel string `json:"timeout_label,omitempty"`
	//Use timestamp of first record when buffer is flushed. (default: False)
	UseFirstTimestamp bool `json:"use_first_timestamp,omitempty"`
	//The field name that is the reference to concatenate records
	PartialKey string `json:"partial_key,omitempty"`
	//The value stored in the field specified by partial_key that represent partial log
	PartialValue string `json:"partial_value,omitempty"`
	//If true, keep partial_key in concatenated records (default:False)
	KeepPartialKey bool `json:"keep_partial_key,omitempty"`
	//Use partial metadata to concatenate multiple records
	UsePartialMetadata string `json:"use_partial_metadata,omitempty"`
	//If true, keep partial metadata
	KeepPartialMetadata string `json:"keep_partial_metadata,omitempty"`
}

func (c *Concat) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	pluginType := "concat"
	return types.NewFlatDirective(types.PluginMeta{
		Type:      pluginType,
		Directive: "filter",
		Tag:       "**",
		Id:        id + "_" + pluginType,
	}, c, secretLoader)
}
