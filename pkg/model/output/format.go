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

// +kubebuilder:object:generate=true
type Format struct {
	// Output line formatting: out_file,json,ltsv,csv,msgpack,hash,single_value (default: json)
	// +kubebuilder:validation:Enum=out_file;json;ltsv;csv;msgpack;hash;single_value
	Type string `json:"type,omitempty"`
}

func (f *Format) ToDirective(secretLoader secret.SecretLoader) (types.Directive, error) {
	metadata := types.PluginMeta{
		Directive: "format",
	}
	if f.Type != "" {
		metadata.Type = f.Type
	} else {
		metadata.Type = "json"
	}
	f.Type = ""
	return types.NewFlatDirective(metadata, f, secretLoader)
}
