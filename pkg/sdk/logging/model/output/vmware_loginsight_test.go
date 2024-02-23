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

package output_test

import (
	"testing"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/require"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/render"
)

func TestVMwareLogInsight(t *testing.T) {
	// language=yaml
	CONFIG := []byte(`---
agent_id: test_agent_id
authentication: basic
buffer:
  disabled: false
config_param:
  source: log_source
flatten_hashes: true
flatten_hashes_separator: _
http_conn_debug: true
http_method: post
host: 127.0.0.1
log_text_keys:
- log
- message
- msg
max_batch_size: 400000
path: api/v1/events/ingest
port: 443
raise_on_error: true
rate_limit_msec: 10
request_retries: 10
request_timeout: 10
ssl_verify: false
scheme: https
serializer: json
shorten_keys:
  kubernetes_: k8s_
  namespace: ns
  labels_: ""
  _name: ""
  _hash: ""
  container_: ""`)

	expected := `
  <match **>
	@type vmware_loginsight
	@id test
	agent_id test_agent_id
	authentication basic
    config_param {"source":"log_source"}
	flatten_hashes true
	flatten_hashes_separator _
	host 127.0.0.1
	http_conn_debug true
	http_method post
	log_text_keys ["log","message","msg"]
    max_batch_size 400000
    path api/v1/events/ingest
    port 443
    raise_on_error true
    rate_limit_msec 10
	request_retries 10
	request_timeout 10
	scheme https
	serializer json
	shorten_keys {"_hash":"","_name":"","container_":"","kubernetes_":"k8s_","labels_":"","namespace":"ns"}
	ssl_verify false
	<buffer tag,time>
      @type file
      chunk_limit_size 8MB
      path /buffers/test.*.buffer
      retry_forever true
      timekey 10m
      timekey_wait 1m
	</buffer>
  </match>
`
	es := &output.VMwareLogInsightOutput{}
	require.NoError(t, yaml.Unmarshal(CONFIG, es))
	test := render.NewOutputPluginTest(t, es)
	test.DiffResult(expected)
}
