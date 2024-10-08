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

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/render"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

func TestRelabel(t *testing.T) {
	CONFIG := []byte(`
label: '@new-label'
`)

	expected := `
  <match **>
    @type relabel
    @id test
	@label @new-label
  </match>
`

	relabel := &output.RelabelOutputConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, relabel))
	test := render.NewOutputPluginTest(t, relabel)
	test.DiffResult(expected)
}
