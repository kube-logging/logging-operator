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

func TestAwsElasticsearch(t *testing.T) {
	CONFIG := []byte(`
logstash_format: true
include_tag_key: true
tag_key: "@log_name"
flush_interval: 1s
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
	@type aws-elasticsearch-service
	@id test
	exception_backup true
	fail_on_detecting_es_version_retry_exceed true
	fail_on_putting_template_retry_exceed true
	flush_interval 1s
	include_tag_key true
	logstash_format true
	reload_connections true
	ssl_verify true
	tag_key @log_name
	utc_index true
	verify_es_version_at_startup true
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
	awsEs := &output.AwsElasticsearchOutputConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, awsEs))
	test := render.NewOutputPluginTest(t, awsEs)
	test.DiffResult(expected)
}
