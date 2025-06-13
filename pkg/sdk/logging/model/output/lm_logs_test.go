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

package output_test

import (
	"testing"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/render"
	"sigs.k8s.io/yaml"

	"github.com/stretchr/testify/require"
)

func TestLMLogsOutputConfig(t *testing.T) {
	CONFIG := []byte(`
company_name: mycompany
company_domain: logicmonitor.com
resource_mapping: '{"kubernetes.host": "system.hostname"}'
access_id:
  value: "my_access_id"
access_key:
  value: "my_access_key"
flush_interval: 30s
debug: true
include_metadata: false
device_less_logs: false
buffer:
  chunk_limit_records: 1000
  flush_interval: 5s
  retry_max_times: 3
`)

	expected := `
  <match **>
  @type lm
  @id test
  access_id my_access_id
  access_key my_access_key
  company_domain logicmonitor.com
  company_name mycompany
  debug true
  device_less_logs false
  flush_interval 30s
  include_metadata false
  resource_mapping {"kubernetes.host": "system.hostname"}
  <buffer tag,time>
	@type file
	chunk_limit_records 1000
	flush_interval 5s
	path /buffers/test.*.buffer
	retry_forever true
	retry_max_times 3
	timekey 10m
	timekey_wait 1m
  </buffer>
</match>
`
	s := &output.LMLogsOutputConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, s))
	test := render.NewOutputPluginTest(t, s)
	test.DiffResult(expected)
}

func TestLMLogsOutputConfigWithBearerToken(t *testing.T) {
	CONFIG := []byte(`
company_name: mycompany
resource_mapping: '{"event_key": "lm_property"}'
bearer_token:
  value: "my_bearer_token"
`)

	expected := `
  <match **>
  @type lm
  @id test
  bearer_token my_bearer_token
  company_domain logicmonitor.com
  company_name mycompany
  flush_interval 60s
  resource_mapping {"event_key": "lm_property"}
  <buffer tag,time>
	@type file
	path /buffers/test.*.buffer
	retry_forever true
	timekey 10m
	timekey_wait 1m
  </buffer>
</match>
`
	s := &output.LMLogsOutputConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, s))
	test := render.NewOutputPluginTest(t, s)
	test.DiffResult(expected)
}

func TestLMLogsOutputConfigWithHTTPProxy(t *testing.T) {
	CONFIG := []byte(`
company_name: mycompany
resource_mapping: '{"kubernetes.namespace": "system.groups"}'
access_id:
  value: "test_id"
access_key:
  value: "test_key"
http_proxy: "http://user:pass@proxy.server:8080"
force_encoding: "UTF-8"
`)

	expected := `
  <match **>
  @type lm
  @id test
  access_id test_id
  access_key test_key
  company_domain logicmonitor.com
  company_name mycompany
  flush_interval 60s
  force_encoding UTF-8
  http_proxy http://user:pass@proxy.server:8080
  resource_mapping {"kubernetes.namespace": "system.groups"}
  <buffer tag,time>
	@type file
	path /buffers/test.*.buffer
	retry_forever true
	timekey 10m
	timekey_wait 1m
  </buffer>
</match>
`
	s := &output.LMLogsOutputConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, s))
	test := render.NewOutputPluginTest(t, s)
	test.DiffResult(expected)
}

func TestLMLogsOutputConfigWithoutResourceMapping(t *testing.T) {
	CONFIG := []byte(`
company_name: mycompany
access_id:
  value: "my_access_id"
access_key:
  value: "my_access_key"
`)

	expected := `
  <match **>
  @type lm
  @id test
  access_id my_access_id
  access_key my_access_key
  company_domain logicmonitor.com
  company_name mycompany
  flush_interval 60s
  <buffer tag,time>
	@type file
	path /buffers/test.*.buffer
	retry_forever true
	timekey 10m
	timekey_wait 1m
  </buffer>
</match>
`
	s := &output.LMLogsOutputConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, s))
	test := render.NewOutputPluginTest(t, s)
	test.DiffResult(expected)
}
