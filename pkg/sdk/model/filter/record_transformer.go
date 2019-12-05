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

// +name:"Record Transformer"
// +url:"https://docs.fluentd.org/filter/record_transformer"
// +version:"more info"
// +description:"Mutates/transforms incoming event streams."
// +status:"GA"
type _metaRecordTransformer interface{}

// +kubebuilder:object:generate=true
type RecordTransformer struct {
	// A comma-delimited list of keys to delete
	RemoveKeys string `json:"remove_keys,omitempty"`
	// A comma-delimited list of keys to keep.
	KeepKeys string `json:"keep_keys,omitempty"`
	// Create new Hash to transform incoming data (default: false)
	RenewRecord bool `json:"renew_record,omitempty"`
	// Specify field name of the record to overwrite the time of events. Its value must be unix time.
	RenewTimeKey string `json:"renew_time_key,omitempty"`
	// When set to true, the full Ruby syntax is enabled in the ${...} expression. (default: false)
	EnableRuby bool `json:"enable_ruby,omitempty"`
	// Use original value type. (default: true)
	AutoTypecast bool `json:"auto_typecast,omitempty"`
	// Add records docs at: https://docs.fluentd.org/filter/record_transformer
	// Records are represented as maps: `key: value`
	Records []Record `json:"records,omitempty"`
}

// Parameters inside record directives are considered to be new key-value pairs
type Record map[string]string

func (r *Record) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	recordSet := types.PluginMeta{
		Directive: "record",
	}
	directive := &types.GenericDirective{
		PluginMeta: recordSet,
		Params:     *r,
	}
	return directive, nil
}

func (r *RecordTransformer) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	pluginType := "record_transformer"
	recordTransformer := &types.GenericDirective{
		PluginMeta: types.PluginMeta{
			Type:      pluginType,
			Directive: "filter",
			Tag:       "**",
			Id:        id + "_" + pluginType,
		},
	}
	if params, err := types.NewStructToStringMapper(secretLoader).StringsMap(r); err != nil {
		return nil, err
	} else {
		recordTransformer.Params = params
	}
	if len(r.Records) > 0 {
		for _, record := range r.Records {
			if meta, err := record.ToDirective(secretLoader, ""); err != nil {
				return nil, err
			} else {
				recordTransformer.SubDirectives = append(recordTransformer.SubDirectives, meta)
			}
		}
	}
	return recordTransformer, nil
}
