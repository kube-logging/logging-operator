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
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/render"
	"github.com/stretchr/testify/require"
)

func TestSplunkHEC(t *testing.T) {
	CONFIG := []byte(`
hec_host: splunk.default.svc.cluster.local
hec_port: 8088
insecure_ssl: true
index: foo
fields:
  dummy: " "
host: foo
protocol: http
metrics_from_event: true
metric_name_key: foo
metric_value_key: foo
coerce_to_utf8: true
non_utf8_replacement_string: bar
sourcetype: foo
source: foo
buffer:
  timekey: 1m
  timekey_wait: 30s
  timekey_use_utc: true
`)
	expected := `
	<match **>
		@type splunk_hec
		@id test
		coerce_to_utf8 true
		hec_host splunk.default.svc.cluster.local
		hec_port 8088
		host foo
		index foo
		insecure_ssl true
		metric_name_key foo
		metric_value_key foo
		metrics_from_event true
		non_utf8_replacement_string bar
		protocol http
		source foo
		sourcetype foo
		<buffer tag,time>
		  @type file
		  chunk_limit_size 8MB
		  path /buffers/test.*.buffer
		  retry_forever true
		  timekey 1m
		  timekey_use_utc true
		  timekey_wait 30s
		</buffer>
		<fields>
		  dummy
		</fields>
	</match>
`
	es := &output.SplunkHecOutput{}
	require.NoError(t, yaml.Unmarshal(CONFIG, es))
	test := render.NewOutputPluginTest(t, es)
	test.DiffResult(expected)
}
