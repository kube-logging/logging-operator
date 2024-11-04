// Copyright © 2024 Kube logging authors
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

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/render"
)

func TestNull(t *testing.T) {
	CONFIG := []byte(`
never_flush: false
`)

	expected := `
  <match **>
    @type null
	@id test
	never_flush false
  </match>
`

	null := &output.NullOutputConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, null))
	test := render.NewOutputPluginTest(t, null)
	test.DiffResult(expected)
}
