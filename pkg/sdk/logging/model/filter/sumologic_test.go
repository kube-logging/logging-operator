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

package filter_test

import (
	"testing"

	"github.com/ghodss/yaml"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/filter"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/render"
	"github.com/stretchr/testify/require"
)

func TestSumoLogic(t *testing.T) {
	CONFIG := []byte(`
source_name: "elso"
source_host: "masodik"
`)
	expected := `
<filter **>
  @type kubernetes_sumologic
  @id test
  source_host masodik
  source_name elso
</filter>
`
	parser := &filter.SumoLogic{}
	require.NoError(t, yaml.Unmarshal(CONFIG, parser))
	test := render.NewOutputPluginTest(t, parser)
	test.DiffResult(expected)
}
