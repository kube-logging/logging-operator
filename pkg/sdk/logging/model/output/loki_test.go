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

func TestLoki(t *testing.T) {
	CONFIG := []byte(`
url: http://loki:3100
configure_kubernetes_labels: true
labels:
  name: $.name
extra_labels:
  testing: "testing"
buffer:
  timekey: 1m
  timekey_wait: 30s
  timekey_use_utc: true
`)
	expected := `
  <match **>
    @type loki
    @id test
    extra_labels {"testing":"testing"}
    extract_kubernetes_labels true
	include_thread_label true
    line_format json
    remove_keys ["kubernetes"]
    url http://loki:3100
    <label>
      container $.kubernetes.container_name
      container_id $.kubernetes.docker_id
      host $.kubernetes.host
      name $.name
      namespace $.kubernetes.namespace_name
      pod $.kubernetes.pod_name
      pod_id $.kubernetes.pod_id
    </label>
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
	loki := &output.LokiOutput{}
	require.NoError(t, yaml.Unmarshal(CONFIG, loki))
	test := render.NewOutputPluginTest(t, loki)
	test.DiffResult(expected)
}
