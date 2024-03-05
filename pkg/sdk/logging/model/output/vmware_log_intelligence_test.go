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

	"github.com/ghodss/yaml"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/render"

	"github.com/stretchr/testify/require"
)

func TestVMwareLogIntelligenceOutputConfig(t *testing.T) {
	CONFIG := []byte(`
endpoint_url: https://data.upgrade.symphony-dev.com/le-mans/v1/streams/ingestion-pipeline-stream
verify_ssl: true
http_compress: false
headers: 
  content_type: "application/json"
  authorization: 
   value: "Bearer 12345"
  structure: simple
buffer:
  chunk_limit_records: 300
  flush_interval: 3s
  retry_max_times: 3
`)

	expected := `
  <match **>
  @type vmware_log_intelligence
  @id test
  endpoint_url https://data.upgrade.symphony-dev.com/le-mans/v1/streams/ingestion-pipeline-stream
  http_compress false
  verify_ssl true
  <headers>
    Authorization Bearer 12345
    Content-Type application/json
    structure simple
  </headers>
  <buffer tag,time>
	@type file
	chunk_limit_records 300
	chunk_limit_size 8MB
	flush_interval 3s
	path /buffers/test.*.buffer
	retry_forever true
	retry_max_times 3
	timekey 10m
	timekey_wait 1m
  </buffer>
</match>
`
	s := &output.VMwareLogIntelligenceOutputConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, s))
	test := render.NewOutputPluginTest(t, s)
	test.DiffResult(expected)
}

func TestVMwareLogIntelligenceOutputConfigWithDefaultHeaderValues(t *testing.T) {
	CONFIG := []byte(`
endpoint_url: https://data.upgrade.symphony-dev.com/le-mans/v1/streams/ingestion-pipeline-stream
verify_ssl: true
http_compress: false
headers:
  authorization:
    value: "Bearer 12345"
buffer:
  chunk_limit_records: 300
  flush_interval: 3s
  retry_max_times: 3
`)

	expected := `
  <match **>
  @type vmware_log_intelligence
  @id test
  endpoint_url https://data.upgrade.symphony-dev.com/le-mans/v1/streams/ingestion-pipeline-stream
  http_compress false
  verify_ssl true
  <headers>
    Authorization Bearer 12345
    Content-Type application/json
    structure simple
  </headers>
  <buffer tag,time>
	@type file
	chunk_limit_records 300
	chunk_limit_size 8MB
	flush_interval 3s
	path /buffers/test.*.buffer
	retry_forever true
	retry_max_times 3
	timekey 10m
	timekey_wait 1m
  </buffer>
</match>
`
	s := &output.VMwareLogIntelligenceOutputConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, s))
	test := render.NewOutputPluginTest(t, s)
	test.DiffResult(expected)
}

func TestVMwareLogIntelligenceOutputConfigWithFormat(t *testing.T) {
	CONFIG := []byte(`
endpoint_url: https://data.upgrade.symphony-dev.com/le-mans/v1/streams/ingestion-pipeline-stream
verify_ssl: true
headers:
  content_type: "application/json"
  authorization:
    value: "Bearer 12345"
  structure: simple
buffer:
  chunk_limit_records: 300
  flush_interval: 3s
  retry_max_times: 3
format:
  type: json
`)

	expected := `
  <match **>
  @type vmware_log_intelligence
  @id test
  endpoint_url https://data.upgrade.symphony-dev.com/le-mans/v1/streams/ingestion-pipeline-stream
  verify_ssl true
  <headers>
    Authorization Bearer 12345
    Content-Type application/json
    structure simple
  </headers>
  <buffer tag,time>
	@type file
	chunk_limit_records 300
	chunk_limit_size 8MB
	flush_interval 3s
	path /buffers/test.*.buffer
	retry_forever true
	retry_max_times 3
	timekey 10m
	timekey_wait 1m
  </buffer>
  <format>
    @type json
  </format>
</match>
`
	s := &output.VMwareLogIntelligenceOutputConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, s))
	test := render.NewOutputPluginTest(t, s)
	test.DiffResult(expected)
}
