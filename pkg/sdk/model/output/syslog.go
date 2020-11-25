// Copyright © 2020 Banzai Cloud
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
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/types"
	"github.com/banzaicloud/operator-tools/pkg/secret"
)

// +name:"Syslog"
// +weight:"200"
type _hugoSyslog interface{}

// +kubebuilder:object:generate=true
// +docName:"[Syslog Output](https://github.com/cloudfoundry/fluent-plugin-syslog_rfc5424)"
// Fluentd output plugin for remote syslog with RFC5424 headers logs.
type _docSyslog interface{}

// +name:"Syslog"
// +url:"https://github.com/cloudfoundry/fluent-plugin-syslog_rfc5424"
// +version:"0.9.0.rc.5"
// +description:"Output plugin writes events to syslog"
// +status:"GA"
type _metaSyslog interface{}

// +kubebuilder:object:generate=true
type SyslogOutputConfig struct {
	// Destination host address
	Host string `json:"host"`
	// Destination host port (default: "514")
	Port int `json:"port,omitempty"`
	// Transport Protocol (default: "tls")
	Transport string `json:"transport,omitempty"`
	// skip ssl validation (default: false)
	Insecure *bool `json:"insecure,omitempty"`
	// file path to ca to trust
	TrustedCaPath *secret.Secret `json:"trusted_ca_path,omitempty"`
	// +docLink:"Format,../format/"
	Format *FormatRfc5424 `json:"format,omitempty"`
	// +docLink:"Buffer,../buffer/"
	Buffer *Buffer `json:"buffer,omitempty"`
}

// #### Example `File` output configurations
// ```
//apiVersion: logging.banzaicloud.io/v1beta1
//kind: Output
//metadata:
//  name: demo-output
//spec:
//  syslog:
//    host: SYSLOG-HOST
//    port: 123
//    format:
//      app_name_field: example.custom_field_1
//      proc_id_field: example.custom_field_2
//    buffer:
//      timekey: 1m
//      timekey_wait: 10s
//      timekey_use_utc: true
// ```
//
// #### Fluentd Config Result
// ```
//  <match **>
//	@type syslog_rfc5424
//	@id test_syslog
//	host SYSLOG-HOST
//	port 123
//  <format>
//    @type syslog_rfc5424
//    app_name_field example.custom_field_1
//    proc_id_field example.custom_field_2
//  </format>
//	<buffer tag,time>
//	  @type file
//	  path /buffers/test_file.*.buffer
//	  retry_forever true
//	  timekey 1m
//	  timekey_use_utc true
//	  timekey_wait 30s
//	</buffer>
//  </match>
// ```
type _expSyslog interface{}

func (s *SyslogOutputConfig) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	const pluginType = "syslog_rfc5424"
	syslog := &types.OutputPlugin{
		PluginMeta: types.PluginMeta{
			Type:      pluginType,
			Directive: "match",
			Tag:       "**",
			Id:        id,
		},
	}
	if params, err := types.NewStructToStringMapper(secretLoader).StringsMap(s); err != nil {
		return nil, err
	} else {
		syslog.Params = params
	}
	if s.Buffer != nil {
		if buffer, err := s.Buffer.ToDirective(secretLoader, id); err != nil {
			return nil, err
		} else {
			syslog.SubDirectives = append(syslog.SubDirectives, buffer)
		}
	}
	if s.Format != nil {
		if format, err := s.Format.ToDirective(secretLoader, ""); err != nil {
			return nil, err
		} else {
			syslog.SubDirectives = append(syslog.SubDirectives, format)
		}
	}
	return syslog, nil
}
