// Copyright © 2024 Kube logging authors
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
	"github.com/cisco-open/operator-tools/pkg/secret"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/types"
)

// +name:"LogicMonitor Logs"
// +weight:"200"
type _hugoLMLogs interface{} //nolint:deadcode,unused

// +docName:"LogicMonitor Logs output plugin for Fluentd"
/*
LogicMonitor Logs output plugin for Fluentd

Sends log records to LogicMonitor Logs via the LM API.

For details, see [https://github.com/logicmonitor/lm-logs-fluentd](https://github.com/logicmonitor/lm-logs-fluentd).

## Example output configurations

```yaml
spec:
  lmLogs:
    company_name: mycompany
    access_id:
      valueFrom:
        secretKeyRef:
          name: lm-credentials
          key: access_id
    access_key:
      valueFrom:
        secretKeyRef:
          name: lm-credentials
          key: access_key
    resource_mapping: '{"kubernetes.host": "system.hostname"}'
    flush_interval: 60s
    debug: false
```
*/
type _docLMLogs interface{} //nolint:deadcode,unused

// +name:"LogicMonitorLogs"
// +url:"https://github.com/logicmonitor/lm-logs-fluentd/releases/tag/v.1.2.5"
// +version:"v1.2.5"
// +description:"Send your logs to LogicMonitor Logs"
// +status:"GA"
type _metaLMLogs interface{} //nolint:deadcode,unused

// +kubebuilder:object:generate=true
// +docName:"LogicMonitorLogs"
type LMLogsOutputConfig struct {
	// LogicMonitor account name
	CompanyName string `json:"company_name" plugin:"required"`
	// LogicMonitor account domain. For eg. for url test.logicmonitor.com, company_domain is logicmonitor.com (default: logicmonitor.com)
	CompanyDomain string `json:"company_domain,omitempty" plugin:"default:logicmonitor.com"`
	// The mapping that defines the source of the log event to the LM resource. In this case, the <event_key> in the incoming event is mapped to the value of <lm_property>
	ResourceMapping string `json:"resource_mapping"`
	// LM API Token access ID
	// +docLink:"Secret,../secret/"
	AccessID *secret.Secret `json:"access_id,omitempty"`
	// LM API Token access key
	// +docLink:"Secret,../secret/"
	AccessKey *secret.Secret `json:"access_key,omitempty"`
	// LM API Bearer Token. Either specify access_id and access_key both or bearer_token. If all specified, LMv1 token(access_id and access_key) will be used for authentication with LogicMonitor
	// +docLink:"Secret,../secret/"
	BearerToken *secret.Secret `json:"bearer_token,omitempty"`
	// Defines the time in seconds to wait before sending batches of logs to LogicMonitor (default: 60s)
	FlushInterval string `json:"flush_interval,omitempty" plugin:"default:60s"`
	// When true, logs more information to the fluentd console
	Debug *bool `json:"debug,omitempty"`
	// Specify charset when logs contains invalid utf-8 characters
	ForceEncoding string `json:"force_encoding,omitempty"`
	// When true, appends additional metadata to the log (default: false)
	IncludeMetadata *bool `json:"include_metadata,omitempty"`
	// When true, do not map log with any resource. record must have service when true (default: false)
	DeviceLessLogs *bool `json:"device_less_logs,omitempty"`
	// http proxy string eg. http://user:pass@proxy.server:port
	HTTPProxy string `json:"http_proxy,omitempty"`
	// +docLink:"Buffer,../buffer/"
	Buffer *Buffer `json:"buffer,omitempty"`
	// +docLink:"Format,../format/"
	Format *Format `json:"format,omitempty"`
}

func (l *LMLogsOutputConfig) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	const pluginType = "lm"
	lmLogs := &types.OutputPlugin{
		PluginMeta: types.PluginMeta{
			Type:      pluginType,
			Directive: "match",
			Tag:       "**",
			Id:        id,
		},
	}

	if params, err := types.NewStructToStringMapper(secretLoader).StringsMap(l); err != nil {
		return nil, err
	} else {
		lmLogs.Params = params
	}

	if l.Buffer == nil {
		l.Buffer = &Buffer{}
	}
	if buffer, err := l.Buffer.ToDirective(secretLoader, id); err != nil {
		return nil, err
	} else {
		lmLogs.SubDirectives = append(lmLogs.SubDirectives, buffer)
	}

	if l.Format != nil {
		if format, err := l.Format.ToDirective(secretLoader, ""); err != nil {
			return nil, err
		} else {
			lmLogs.SubDirectives = append(lmLogs.SubDirectives, format)
		}
	}
	return lmLogs, nil
}
