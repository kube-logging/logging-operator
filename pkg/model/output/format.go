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
