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

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/filter"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/render"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

func TestConcat(t *testing.T) {
	testData := []struct {
		Config   []byte
		Expected string
	}{
		{
			Config: []byte(`
partial_key: "partial_message"
n_lines: 10
`),
			Expected: `
<filter **>
  @type concat
  @id test
  key message
  n_lines 10
  partial_key partial_message
  separator "\n"
</filter>
			`,
		},
		{
			Config: []byte(`
partial_key: "partial_message"
n_lines: 10
separator: ""
`),
			Expected: `
<filter **>
  @type concat
  @id test
  key message
  n_lines 10
  partial_key partial_message
  separator
</filter>
			`,
		},
	}

	for _, d := range testData {
		parser := &filter.Concat{}
		require.NoError(t, yaml.Unmarshal(d.Config, parser))
		test := render.NewOutputPluginTest(t, parser)
		test.DiffResult(d.Expected)
	}
}
