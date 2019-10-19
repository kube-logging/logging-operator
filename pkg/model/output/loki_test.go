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

	"github.com/banzaicloud/logging-operator/pkg/model/output"
	"github.com/banzaicloud/logging-operator/pkg/model/render"
	"github.com/ghodss/yaml"
)

func TestLoki(t *testing.T) {
	CONFIG := []byte(`
url: http://loki:3100
configure_kubernetes_labels: true
buffer:
  timekey: 1m
  timekey_wait: 30s
  timekey_use_utc: true
`)
	expected := `
  <match **>
    @type kubernetes_loki
    @id test-kubernetes_loki
    extract_kubernetes_labels true
    line_format json
    url http://loki:3100
    <label>
      container ${record.dig("kubernetes", "container_name")}
      container_id ${record.dig("kubernetes", "docker_id")}
      host ${record.dig("kubernetes", "host")}
      namespace ${record.dig("kubernetes", "namespace_name")}
      pod ${record.dig("kubernetes", "pod_name")}
      pod_id ${record.dig("kubernetes", "pod_id")}
    </label>
    <buffer tag,time>
      @type file
      path /buffers/default.*.buffer
      retry_forever true
      timekey 1m
      timekey_use_utc true
      timekey_wait 30s
    </buffer>
  </match>
`
	loki := &output.LokiOutput{}
	yaml.Unmarshal(CONFIG, loki)
	test := render.NewOutputPluginTest(t, loki)
	test.DiffResult(expected)
}
