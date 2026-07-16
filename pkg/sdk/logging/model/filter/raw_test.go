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

	"github.com/cisco-open/operator-tools/pkg/secret"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/filter"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/render"
)

func TestConfigureCustomFilterViaRawConfiguration(t *testing.T) {
	CONFIG := []byte(`
config: |
  @type my_filter
  <my_section>
    foo bar
    tags ["web", "api", "db"]
  </my_section>
`)

	expected := `
<filter **>
  @type my_filter
  @id test
  <my_section>
    foo bar
    tags ["web", "api", "db"]
  </my_section>
</filter>
`
	parser := &filter.Raw{}
	require.NoError(t, yaml.Unmarshal(CONFIG, parser))
	test := render.NewOutputPluginTest(t, parser)
	test.DiffResult(expected)
}

func TestConfigureUnexposedFilterViaRawConfiguration(t *testing.T) {
	CONFIG := []byte(`
config: |
  @type ua_parser
  flatten true
  key_name ua_string
`)

	expected := `
<filter **>
  @type ua_parser
  @id test
  flatten true
  key_name ua_string
</filter>
`
	parser := &filter.Raw{}
	require.NoError(t, yaml.Unmarshal(CONFIG, parser))
	test := render.NewOutputPluginTest(t, parser)
	test.DiffResult(expected)
}

type mockSecretLoader struct{}

func (m mockSecretLoader) Load(secret *secret.Secret) (string, error) {
	return "", nil
}

func TestRawConfigurationMissingType(t *testing.T) {
	CONFIG := []byte(`
config: |
  <my_section>
    foo bar
    tags ["web", "api", "db"]
  </my_section>
`)

	parser := &filter.Raw{}
	require.NoError(t, yaml.Unmarshal(CONFIG, parser))

	_, err := parser.ToDirective(mockSecretLoader{}, "test")
	require.Error(t, err)
	require.Equal(t, "raw filter config must specify @type", err.Error())
}

func TestRawConfigurationUnclosedSection(t *testing.T) {
	CONFIG := []byte(`
config: |
  @type my_filter
  <my_section>
    foo bar
    tags ["web", "api", "db"]
`)

	parser := &filter.Raw{}
	require.NoError(t, yaml.Unmarshal(CONFIG, parser))

	_, err := parser.ToDirective(mockSecretLoader{}, "test")
	require.Error(t, err)
	require.Equal(t, "unexpected end of raw config: missing closing tag </my_section>", err.Error())
}
