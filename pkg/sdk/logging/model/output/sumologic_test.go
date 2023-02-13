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

func TestSumologic(t *testing.T) {
	CONFIG := []byte(`
data_type: metrics
metric_data_format: carbon2
log_format: json
source_category: prod/someapp/logs
source_name: AppA
compress: true
buffer:
  type: file
  timekey_wait: 5s
  timekey: 30s
`)
	expected := `
  <match **>
    @type sumologic
    @id test
    compress true
    data_type metrics
    log_format json
    metric_data_format carbon2
    source_category prod/someapp/logs
    source_name AppA
    <buffer tag,time>
      @type file
	  chunk_limit_size 8MB
      path /buffers/test.*.buffer
      retry_forever true
      timekey 30s
      timekey_wait 5s
    </buffer>
  </match>
`
	s := &output.SumologicOutput{}
	require.NoError(t, yaml.Unmarshal(CONFIG, s))
	test := render.NewOutputPluginTest(t, s)
	test.DiffResult(expected)
}
