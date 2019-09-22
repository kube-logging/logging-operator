package output

import (
	"github.com/banzaicloud/logging-operator/pkg/model/secret"
	"github.com/banzaicloud/logging-operator/pkg/model/types"
)

// +kubebuilder:object:generate=true

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

func (a *AzureStorage) ToDirective(secretLoader secret.SecretLoader) (types.Directive, error) {
	azure := &types.OutputPlugin{
		PluginMeta: types.PluginMeta{
			Type:      "azurestorage",
			Directive: "match",
			Tag:       "**",
		},
	}
	if params, err := types.NewStructToStringMapper(secretLoader).StringsMap(a); err != nil {
		return nil, err
	} else {
		azure.Params = params
	}
	if a.Buffer != nil {
		if buffer, err := a.Buffer.ToDirective(secretLoader); err != nil {
			return nil, err
		} else {
			azure.SubDirectives = append(azure.SubDirectives, buffer)
		}
	}
	return azure, nil
}
