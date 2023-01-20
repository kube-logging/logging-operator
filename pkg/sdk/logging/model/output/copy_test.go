// Copyright © 2020 Banzai Cloud
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

	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/output"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/render"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	CONFIG := []byte(`
stores:
- nullout: {}
- prometheus:
    metrics:
    - name: fluentd_output_messages_total
      desc: The total number of output messages.
      type: counter
    labels:
      source: $.source
`)

	expected := `
  <match **>
    @type copy
    @id test
    <store>
      @type null
    </store>
    <store>
      @type prometheus
      <metric>
        desc The total number of output messages.
        name fluentd_output_messages_total
        type counter
      </metric>
      <labels>
        source $.source
      </labels>
    </store>
  </match>
`

	copy := &output.CopyOutputConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, copy))
	test := render.NewOutputPluginTest(t, copy)
	test.DiffResult(expected)
}
