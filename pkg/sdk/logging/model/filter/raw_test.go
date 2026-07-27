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
	"bytes"
	"strings"
	"testing"

	"github.com/cisco-open/operator-tools/pkg/secret"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/filter"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/render"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/types"
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

func TestRawConfigurationRejectsExcessiveNesting(t *testing.T) {
	depth := 64
	parser := &filter.Raw{
		Config: "@type my_filter\n" +
			strings.Repeat("<s>\n", depth) +
			strings.Repeat("</s>\n", depth),
	}

	_, err := parser.ToDirective(mockSecretLoader{}, "test")
	require.Error(t, err)
	require.Equal(t, "raw config nesting is too deep: exceeds maximum of 32 levels", err.Error())
}

func TestRawConfigurationAllowsRealisticNesting(t *testing.T) {
	depth := 32
	parser := &filter.Raw{
		Config: "@type my_filter\n" +
			strings.Repeat("<s>\n", depth) +
			strings.Repeat("</s>\n", depth),
	}

	_, err := parser.ToDirective(mockSecretLoader{}, "test")
	require.NoError(t, err)
}

func TestRawConfigurationRejectsOversizedConfig(t *testing.T) {
	parser := &filter.Raw{
		Config: "@type my_filter\n" + strings.Repeat("key value\n", 8*1024),
	}

	_, err := parser.ToDirective(mockSecretLoader{}, "test")
	require.Error(t, err)
	require.Equal(t, "raw config is too large: 81936 bytes exceeds maximum of 65536", err.Error())
}

// Rendered size must stay bounded: indentation grows with nesting depth, so an
// unbounded config could otherwise amplify a small Flow into gigabytes of config.
func TestRawConfigurationRenderedSizeIsBounded(t *testing.T) {
	depth := 32
	body := strings.Repeat("k v\n", 8*1024)
	parser := &filter.Raw{
		Config: "@type my_filter\n" +
			strings.Repeat("<s>\n", depth) + body + strings.Repeat("</s>\n", depth),
	}

	directive, err := parser.ToDirective(mockSecretLoader{}, "test")
	require.NoError(t, err)

	var buf bytes.Buffer
	renderer := render.FluentRender{Out: &buf, Indent: 2}
	require.NoError(t, renderer.RenderDirectives([]types.Directive{directive}, 0))

	require.Less(t, buf.Len(), 4*1024*1024, "rendered raw filter must stay bounded")
}

func TestConfigureRawFilterWithNestedFilterSection(t *testing.T) {
	CONFIG := []byte(`
config: |
  @type my_filter
  <filter mytag>
    key value
  </filter>
`)

	expected := `
<filter **>
  @type my_filter
  @id test
  <filter mytag>
    key value
  </filter>
</filter>
`
	parser := &filter.Raw{}
	require.NoError(t, yaml.Unmarshal(CONFIG, parser))
	test := render.NewOutputPluginTest(t, parser)
	test.DiffResult(expected)
}

func TestRawConfigurationUnclosedNestedFilterSection(t *testing.T) {
	CONFIG := []byte(`
config: |
  @type my_filter
  <filter>
    key value
`)

	parser := &filter.Raw{}
	require.NoError(t, yaml.Unmarshal(CONFIG, parser))

	_, err := parser.ToDirective(mockSecretLoader{}, "test")
	require.Error(t, err)
	require.Equal(t, "unexpected end of raw config: missing closing tag </filter>", err.Error())
}

func TestRawConfigurationStrayTopLevelClosingTag(t *testing.T) {
	CONFIG := []byte(`
config: |
  @type my_filter
  key1 val1
  </filter>
  key2 val2
`)

	parser := &filter.Raw{}
	require.NoError(t, yaml.Unmarshal(CONFIG, parser))

	_, err := parser.ToDirective(mockSecretLoader{}, "test")
	require.Error(t, err)
	require.Equal(t, "unexpected closing tag in raw config: </filter>", err.Error())
}

func TestRawConfigurationMismatchedClosingTag(t *testing.T) {
	CONFIG := []byte(`
config: |
  @type my_filter
  <my_section>
    foo bar
  </other_section>
`)

	parser := &filter.Raw{}
	require.NoError(t, yaml.Unmarshal(CONFIG, parser))

	_, err := parser.ToDirective(mockSecretLoader{}, "test")
	require.Error(t, err)
	require.Equal(t, "unexpected closing tag in raw config: </other_section>", err.Error())
}

func TestRawConfigurationWrappedInFilterTags(t *testing.T) {
	CONFIG := []byte(`
config: |
  <filter **>
    @type my_filter
    key1 val1
  </filter>
`)

	parser := &filter.Raw{}
	require.NoError(t, yaml.Unmarshal(CONFIG, parser))

	_, err := parser.ToDirective(mockSecretLoader{}, "test")
	require.Error(t, err)
	require.Equal(t, "raw filter config must not include the enclosing <filter> tags, provide only the filter body", err.Error())
}
