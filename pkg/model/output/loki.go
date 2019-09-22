package output

import (
	"github.com/banzaicloud/logging-operator/pkg/model/secret"
	"github.com/banzaicloud/logging-operator/pkg/model/types"
)

// +kubebuilder:object:generate=true
// +docName:"Loki"
// Fluentd output plugin to ship logs to a Loki server.
type LokiOutput struct {
	// The url of the Loki server to send logs to. (default:https://logs-us-west1.grafana.net)
	Url string `json:"url,omitempty"`
	// Specify a username if the Loki server requires authentication.
	// +docLink:"Secret,./secret.md"
	Username *secret.Secret `json:"username,omitempty"`
	// Specify password if the Loki server requires authentication.
	// +docLink:"Secret,./secret.md"
	Password *secret.Secret `json:"password,omitempty"`
	// Loki is a multi-tenant log storage platform and all requests sent must include a tenant.
	Tenant string `json:"tenant,omitempty"`
	// Set of labels to include with every Loki stream.(default: nil)
	ExtraLabels bool `json:"extra_labels,omitempty"`
	// +docLink:"Buffer,./buffer.md"
	Buffer *Buffer `json:"buffer,omitempty"`
}

func (l *LokiOutput) ToDirective(secretLoader secret.SecretLoader) (types.Directive, error) {
	loki := &types.OutputPlugin{
		PluginMeta: types.PluginMeta{
			Type:      "kubernetes_loki",
			Directive: "match",
			Tag:       "**",
		},
	}
	if params, err := types.NewStructToStringMapper(secretLoader).StringsMap(l); err != nil {
		return nil, err
	} else {
		loki.Params = params
	}
	if l.Buffer != nil {
		if buffer, err := l.Buffer.ToDirective(secretLoader); err != nil {
			return nil, err
		} else {
			loki.SubDirectives = append(loki.SubDirectives, buffer)
		}
	}
	return loki, nil
}
