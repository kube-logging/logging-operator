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

func TestRecordModifier(t *testing.T) {
	CONFIG := []byte(`
records:
- foo: "bar"
`)
	expected := `
<filter **>
  @type record_modifier
  @id test
  <record>
    foo bar
  </record>
</filter>
`
	parser := &filter.RecordModifier{}
	require.NoError(t, yaml.Unmarshal(CONFIG, parser))
	test := render.NewOutputPluginTest(t, parser)
	test.DiffResult(expected)
}

func TestRecordModifierWithAllOptions(t *testing.T) {
	CONFIG := []byte(`
prepare_value: "require 'foo'; @foo = Foo.new"
char_encoding: "utf-8:euc-jp"
remove_keys: "key1, key2"
whitelist_keys: "key1, key2"
replaces:
- key: "key1"
  expression: "/^(?<start>.+).{2}(?<end>.+)$/"
  replace: "\\k<start>ors\\k<end>"
records:
- foo: "bar"
`)
	expected := `
<filter **>
  @type record_modifier
	@id test
	char_encoding utf-8:euc-jp
	prepare_value require 'foo'; @foo = Foo.new
	remove_keys key1, key2
	whitelist_keys key1, key2
	<replace>
		expression /^(?<start>.+).{2}(?<end>.+)$/
		key key1
		replace \k<start>ors\k<end>
	</replace>
  <record>
    foo bar
  </record>
</filter>
`
	parser := &filter.RecordModifier{}
	require.NoError(t, yaml.Unmarshal(CONFIG, parser))
	test := render.NewOutputPluginTest(t, parser)
	test.DiffResult(expected)
}
