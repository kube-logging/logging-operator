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
)

func TestLogDNAOutput(t *testing.T) {
	CONFIG := []byte(`
api_key: xxxxxxxxxxxxxxxxxxxxxxxxxxy
hostname: logging-operator
app: myapp
tags: web,dev
request_timeout: 30000 ms
ingester_domain: https://logs.logdna.com
ingester_endpoint: /logs/ingest
buffer:
  timekey: 1m
  timekey_wait: 30s
  timekey_use_utc: true
`)
	expected := `
  <match **>
	@type logdna
	@id test_logdna
	api_key xxxxxxxxxxxxxxxxxxxxxxxxxxy
	app myapp
	hostname logging-operator
	ingester_domain https://logs.logdna.com
	ingester_endpoint /logs/ingest
	request_timeout 30000 ms
	tags web,dev
	<buffer tag,time>
	  @type file
	  chunk_limit_size 8MB
	  path /buffers/test_logdna.*.buffer
	  retry_forever true
	  timekey 1m
	  timekey_use_utc true
	  timekey_wait 30s
	</buffer>
  </match>
`
	logdna := &output.LogDNAOutput{}
	err := yaml.Unmarshal(CONFIG, logdna)
	if err != nil {
		t.Fatalf(err.Error())
	}
	test := render.NewOutputPluginTest(t, logdna)
	test.DiffResult(expected)
}
