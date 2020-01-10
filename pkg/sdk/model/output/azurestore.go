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

// +docName:"Azure Storage output plugin for Fluentd"
//Azure Storage output plugin buffers logs in local file and upload them to Azure Storage periodically.
//More info at https://github.com/htgc/fluent-plugin-azurestorage
type _docAzure interface{}

// +name:"Azure Storage"
// +url:"https://github.com/htgc/fluent-plugin-azurestorage/releases/tag/v0.1.0"
// +version:"0.1.0"
// +description:"Store logs in Azure Storage"
// +status:"GA"
type _metaAzure interface{}

// +kubebuilder:object:generate=true
// +docName:"Output Config"
type AzureStorage struct {
	// Path prefix of the files on Azure
	Path string `json:"path,omitempty"`
	// Your azure storage account
	// +docLink:"Secret,./secret.md"
	AzureStorageAccount *secret.Secret `json:"azure_storage_account"`
	// Your azure storage access key
	// +docLink:"Secret,./secret.md"
	AzureStorageAccessKey *secret.Secret `json:"azure_storage_access_key"`
	// Your azure storage container
	AzureContainer string `json:"azure_container"`
	// Azure storage type currently only "blob" supported (default: blob)
	AzureStorageType string `json:"azure_storage_type,omitempty"`
	// Object key format (default: %{path}%{time_slice}_%{index}.%{file_extension})
	AzureObjectKeyFormat string `json:"azure_object_key_format,omitempty"`
	// Store as: gzip, json, text, lzo, lzma2 (default: gzip)
	StoreAs string `json:"store_as,omitempty"`
	// Automatically create container if not exists(default: true)
	AutoCreateContainer bool `json:"auto_create_container,omitempty"`
	// Compat format type: out_file, json, ltsv (default: out_file)
	Format string `json:"format,omitempty" plugin:"default:json"`
	// +docLink:"Buffer,./buffer.md"
	Buffer *Buffer `json:"buffer,omitempty"`
}

func (a *AzureStorage) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	pluginType := "azurestorage"
	pluginID := id + "_" + pluginType
	azure := &types.OutputPlugin{
		PluginMeta: types.PluginMeta{
			Type:      pluginType,
			Directive: "match",
			Tag:       "**",
			Id:        pluginID,
		},
	}
	if params, err := types.NewStructToStringMapper(secretLoader).StringsMap(a); err != nil {
		return nil, err
	} else {
		azure.Params = params
	}
	if a.Buffer != nil {
		if buffer, err := a.Buffer.ToDirective(secretLoader, pluginID); err != nil {
			return nil, err
		} else {
			azure.SubDirectives = append(azure.SubDirectives, buffer)
		}
	}
	return azure, nil
}
