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

package filter_test

import (
	"testing"

	"github.com/ghodss/yaml"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/filter"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/render"
	"github.com/stretchr/testify/require"
)

func TestPrometheus(t *testing.T) {
	CONFIG := []byte(`
metrics:
- type: counter
  name: message_foo_counter
  desc: "The total number of foo in message."
  key: foo
  labels:
    foo: bar
labels:
  tag: ${tag}
  host: ${hostname}
  namespace: $.kubernetes.namespace
`)
	expected := `
<filter **>
  @type prometheus
  @id test
	<metric>
	  desc The total number of foo in message.
	  key foo
	  name message_foo_counter
	  type counter
	  <labels>
		foo bar
	  </labels>
	</metric>
	<labels>
      host ${hostname}
      namespace $.kubernetes.namespace
      tag ${tag}
	</labels>
</filter>
`
	parser := &filter.PrometheusConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, parser))
	test := render.NewOutputPluginTest(t, parser)
	test.DiffResult(expected)
}
