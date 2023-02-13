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

package output_test

import (
	"testing"

	"github.com/ghodss/yaml"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/render"
	"github.com/stretchr/testify/require"
)

func TestSyslogOutputConfig(t *testing.T) {
	CONFIG := []byte(`
host: SYSLOG-HOST
port: 123
format:
  app_name_field: example.custom_field_1
  proc_id_field: example.custom_field_2
buffer:
  timekey: 1m
  timekey_wait: 30s
  timekey_use_utc: true
`)
	expected := `
  <match **>
	@type syslog_rfc5424
	@id test
	host SYSLOG-HOST
	port 123
	<buffer tag,time>
	  @type file
	  chunk_limit_size 8MB
	  path /buffers/test.*.buffer
	  retry_forever true
	  timekey 1m
	  timekey_use_utc true
	  timekey_wait 30s
	</buffer>
	<format>
	  @type syslog_rfc5424
	  app_name_field example.custom_field_1
	  proc_id_field example.custom_field_2
	</format>
  </match>
`
	f := &output.SyslogOutputConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, f))
	test := render.NewOutputPluginTest(t, f)
	test.DiffResult(expected)
}

func TestSyslogOutputConfigEmptyFormat(t *testing.T) {
	CONFIG := []byte(`
host: SYSLOG-HOST
port: 123
buffer:
  timekey: 1m
  timekey_wait: 30s
  timekey_use_utc: true
`)
	expected := `
  <match **>
	@type syslog_rfc5424
	@id test
	host SYSLOG-HOST
	port 123
	<buffer tag,time>
	  @type file
	  chunk_limit_size 8MB
	  path /buffers/test.*.buffer
	  retry_forever true
	  timekey 1m
	  timekey_use_utc true
	  timekey_wait 30s
	</buffer>
	<format>
	  @type syslog_rfc5424
	</format>
  </match>
`
	f := &output.SyslogOutputConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, f))
	test := render.NewOutputPluginTest(t, f)
	test.DiffResult(expected)
}

func TestSyslogOutputConfigForMutualTLS(t *testing.T) {
	CONFIG := []byte(`
allow_self_signed_cert: true
enable_system_cert_store: true
fqdn: Test-Fqdn
host: SYSLOG-HOST
port: 123
verify_fqdn: false
version: TLSv1_2
format:
  app_name_field: example.custom_field_1
  proc_id_field: example.custom_field_2
buffer:
  timekey: 1m
  timekey_wait: 30s
  timekey_use_utc: true
`)
	expected := `
  <match **>
	@type syslog_rfc5424
	@id test
	allow_self_signed_cert true
	enable_system_cert_store true
	fqdn Test-Fqdn
	host SYSLOG-HOST
	port 123
	verify_fqdn false
	version TLSv1_2
	<buffer tag,time>
	  @type file
	  chunk_limit_size 8MB
	  path /buffers/test.*.buffer
	  retry_forever true
	  timekey 1m
	  timekey_use_utc true
	  timekey_wait 30s
	</buffer>
	<format>
	  @type syslog_rfc5424
	  app_name_field example.custom_field_1
	  proc_id_field example.custom_field_2
	</format>
  </match>
`
	f := &output.SyslogOutputConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, f))
	test := render.NewOutputPluginTest(t, f)
	test.DiffResult(expected)
}
