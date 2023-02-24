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

func TestParser(t *testing.T) {
	CONFIG := []byte(`
remove_key_name_field: true
reserve_data: true
emit_invalid_record_to_error: false
parse:
  type: nginx
`)
	expected := `
<filter **>
  @type parser
  @id test
  emit_invalid_record_to_error false
  key_name message
  remove_key_name_field true
  reserve_data true
  <parse>
    @type nginx
  </parse>
</filter>
`
	parser := &filter.ParserConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, parser))
	test := render.NewOutputPluginTest(t, parser)
	test.DiffResult(expected)
}

func TestParserMultiParser(t *testing.T) {
	CONFIG := []byte(`
remove_key_name_field: true
reserve_data: true
parse:
  type: multi_format
  patterns:
  - format: nginx
  - format: regexp
    expression: /asdsada/
  - format: none
`)
	expected := `
<filter **>
  @type parser
  @id test
  key_name message
  remove_key_name_field true
  reserve_data true
  <parse>
    @type multi_format
    <pattern>
      format nginx
    </pattern>
    <pattern>
      expression /asdsada/
      format regexp
    </pattern>
    <pattern>
      format none
    </pattern>
  </parse>
</filter>
`
	parser := &filter.ParserConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, parser))
	test := render.NewOutputPluginTest(t, parser)
	test.DiffResult(expected)
}

func TestParserMultiLineParser(t *testing.T) {
	CONFIG := []byte(`
remove_key_name_field: true
reserve_data: true
parse:
  type: multiline
  format_firstline: '/^Started/'
  multiline:
  - '/Started (?<method>[^ ]+) "(?<path>[^"]+)" for (?<host>[^ ]+) at (?<time>[^ ]+ [^ ]+ [^ ]+)\n/'
  - '/Processing by (?<controller>[^\u0023]+)\u0023(?<controller_method>[^ ]+) as (?<format>[^ ]+?)\n/'
  - '/(  Parameters: (?<parameters>[^ ]+)\n)?/'
  - '/  Rendered (?<template>[^ ]+) within (?<layout>.+) \([\d\.]+ms\)\n/'
  - '/Completed (?<code>[^ ]+) [^ ]+ in (?<runtime>[\d\.]+)ms \(Views: (?<view_runtime>[\d\.]+)ms \| ActiveRecord: (?<ar_runtime>[\d\.]+)ms\)/'
`)

	expected := `
<filter **>
  @type parser
  @id test
  key_name message
  remove_key_name_field true
  reserve_data true
  <parse>
    @type multiline
    format1 /Started (?<method>[^ ]+) "(?<path>[^"]+)" for (?<host>[^ ]+) at (?<time>[^ ]+ [^ ]+ [^ ]+)\n/
    format2 /Processing by (?<controller>[^\u0023]+)\u0023(?<controller_method>[^ ]+) as (?<format>[^ ]+?)\n/
    format3 /(  Parameters: (?<parameters>[^ ]+)\n)?/
    format4 /  Rendered (?<template>[^ ]+) within (?<layout>.+) \([\d\.]+ms\)\n/
    format5 /Completed (?<code>[^ ]+) [^ ]+ in (?<runtime>[\d\.]+)ms \(Views: (?<view_runtime>[\d\.]+)ms \| ActiveRecord: (?<ar_runtime>[\d\.]+)ms\)/
    format_firstline /^Started/
  </parse>
</filter>
`
	parser := &filter.ParserConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, parser))
	test := render.NewOutputPluginTest(t, parser)
	test.DiffResult(expected)
}

func TestParserGrokSingleParser(t *testing.T) {
	CONFIG := []byte(`
remove_key_name_field: true
reserve_data: true
parse:
  type: grok
  grok_failure_key: grokFailure
  grok_pattern: "%{GREEDYDATA:grokMessage}"
`)

	expected := `
<filter **>
  @type parser
  @id test
  key_name message
  remove_key_name_field true
  reserve_data true
  <parse>
    @type grok
    grok_failure_key grokFailure
    grok_pattern %{GREEDYDATA:grokMessage}
  </parse>
</filter>
`
	parser := &filter.ParserConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, parser))
	test := render.NewOutputPluginTest(t, parser)
	test.DiffResult(expected)
}

func TestParserGrokMultiParser(t *testing.T) {
	CONFIG := []byte(`
remove_key_name_field: true
reserve_data: true
parse:
  type: grok
  grok_failure_key: grokFailure
  grok_patterns:
    - pattern: "%{GREEDYDATA:firstMessage}"
    - pattern: "%{GREEDYDATA:secondMessage}"
`)

	expected := `
<filter **>
  @type parser
  @id test
  key_name message
  remove_key_name_field true
  reserve_data true
  <parse>
    @type grok
    grok_failure_key grokFailure
    <grok>
      pattern %{GREEDYDATA:firstMessage}
      time_key time
    </grok>
    <grok>
      pattern %{GREEDYDATA:secondMessage}
      time_key time
    </grok>    
  </parse>
</filter>
`
	parser := &filter.ParserConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, parser))
	test := render.NewOutputPluginTest(t, parser)
	test.DiffResult(expected)
}
