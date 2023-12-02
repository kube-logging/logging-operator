// Copyright © 2022 Cisco Systems, Inc. and/or its affiliates
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

func TestOpenSearch(t *testing.T) {
	CONFIG := []byte(`
host: opensearch-cluster.default.svc.cluster.local
port: 9200
scheme: https
ssl_version: TLSv1_2
ssl_verify: false
buffer:
  timekey: 1m
  timekey_wait: 30s
  timekey_use_utc: true
`)
	expected := `
  <match **>
  @type opensearch
  @id test
  catch_transport_exception_on_retry true
  emit_error_label_event true
  exception_backup true
  fail_on_detecting_os_version_retry_exceed true
  fail_on_putting_template_retry_exceed true
  host opensearch-cluster.default.svc.cluster.local
  http_backend_excon_nonblock true
  port 9200
  reload_connections true
  scheme https
  ssl_verify false
  ssl_version TLSv1_2
  utc_index true
  verify_os_version_at_startup true
  <buffer tag,time>
    @type file
  chunk_limit_size 8MB
    path /buffers/test.*.buffer
    retry_forever true
    timekey 1m
    timekey_use_utc true
    timekey_wait 30s
  </buffer>
  </match>
`
	es := &output.OpenSearchOutput{}
	require.NoError(t, yaml.Unmarshal(CONFIG, es))
	test := render.NewOutputPluginTest(t, es)
	test.DiffResult(expected)
}

func TestAwsOpenSearch(t *testing.T) {
	CONFIG := []byte(`
host: opensearch-cluster.default.svc.cluster.local
port: 9200
scheme: https
ssl_version: TLSv1_2
ssl_verify: false
endpoint:
  url: https://CLUSTER_ENDPOINT_URL
  region: eu-west-1
  access_key_id:
    value: aws-key
  secret_access_key:
    value: aws_secret
buffer:
  timekey: 1m
  timekey_wait: 30s
  timekey_use_utc: true
`)
	expected := `
  <match **>
  @type opensearch
  @id test
  catch_transport_exception_on_retry true
  emit_error_label_event true
  exception_backup true
  fail_on_detecting_os_version_retry_exceed true
  fail_on_putting_template_retry_exceed true
  host opensearch-cluster.default.svc.cluster.local
  http_backend_excon_nonblock true
  port 9200
  reload_connections true
  scheme https
  ssl_verify false
  ssl_version TLSv1_2
  utc_index true
  verify_os_version_at_startup true
  <endpoint>
    access_key_id aws-key
    region eu-west-1
    secret_access_key aws_secret
    url https://CLUSTER_ENDPOINT_URL
  </endpoint>
  <buffer tag,time>
    @type file
  chunk_limit_size 8MB
    path /buffers/test.*.buffer
    retry_forever true
    timekey 1m
    timekey_use_utc true
    timekey_wait 30s
  </buffer>
  </match>
`
	es := &output.OpenSearchOutput{}
	require.NoError(t, yaml.Unmarshal(CONFIG, es))
	test := render.NewOutputPluginTest(t, es)
	test.DiffResult(expected)
}

func TestOpenSearchDataStream(t *testing.T) {
	CONFIG := []byte(`
host: opensearch-cluster.default.svc.cluster.local
port: 9200
scheme: https
ssl_version: TLSv1_2
ssl_verify: false
data_stream_enable: true
data_stream_name: test-ds
buffer:
  timekey: 1m
  timekey_wait: 30s
  timekey_use_utc: true
`)
	expected := `
  <match **>
    @type opensearch_data_stream
    @id test
    catch_transport_exception_on_retry true
    data_stream_name test-ds
    emit_error_label_event true
    exception_backup true
    fail_on_detecting_os_version_retry_exceed true
    fail_on_putting_template_retry_exceed true
    host opensearch-cluster.default.svc.cluster.local
    http_backend_excon_nonblock true
    port 9200
    reload_connections true
    scheme https
    ssl_verify false
    ssl_version TLSv1_2
    utc_index true
    verify_os_version_at_startup true
    <buffer tag,time>
      @type file
      chunk_limit_size 8MB
      path /buffers/test.*.buffer
      retry_forever true
      timekey 1m
      timekey_use_utc true
      timekey_wait 30s
    </buffer>
  </match>
`
	es := &output.OpenSearchOutput{}
	require.NoError(t, yaml.Unmarshal(CONFIG, es))
	test := render.NewOutputPluginTest(t, es)
	test.DiffResult(expected)
}
