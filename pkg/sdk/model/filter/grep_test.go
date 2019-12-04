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

	"github.com/banzaicloud/logging-operator/pkg/sdk/model/filter"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/render"
	"github.com/ghodss/yaml"
)

func TestGrepRegexp(t *testing.T) {
	CONFIG := []byte(`
regexp:
  - key: elso
    pattern: /^5\d\d$/
`)
	expected := `
<filter **>
@type grep
@id test_grep
<regexp>
  key elso
  pattern /^5\d\d$/
</regexp>
</filter>
`
	grep := &filter.GrepConfig{}
	yaml.Unmarshal(CONFIG, grep)
	test := render.NewOutputPluginTest(t, grep)
	test.DiffResult(expected)
}

func TestGrepExclude(t *testing.T) {
	CONFIG := []byte(`
exclude:
  - key: elso
    pattern: /^5\d\d$/
`)
	expected := `
<filter **>
@type grep
@id test_grep
<exclude>
  key elso
  pattern /^5\d\d$/
</exclude>
</filter>
`
	grep := &filter.GrepConfig{}
	yaml.Unmarshal(CONFIG, grep)
	test := render.NewOutputPluginTest(t, grep)
	test.DiffResult(expected)
}

func TestGrepOr(t *testing.T) {
	CONFIG := []byte(`
or:
  - regexp:
    - key: elso
      pattern: /^5\d\d$/
    - key: masodik
      pattern: /\.css$/
`)
	expected := `
<filter **>
@type grep
@id test_grep
<or>
  <regexp>
	key elso
	pattern /^5\d\d$/
  </regexp>
  <regexp>
	key masodik
	pattern /\.css$/
  </regexp>
</or>
</filter>
`
	grep := &filter.GrepConfig{}
	yaml.Unmarshal(CONFIG, grep)
	test := render.NewOutputPluginTest(t, grep)
	test.DiffResult(expected)
}

func TestGrepAnd(t *testing.T) {
	CONFIG := []byte(`
and:
  - exclude:
    - key: elso
      pattern: /^5\d\d$/
    - key: masodik
      pattern: /\.css$/
`)
	expected := `
<filter **>
@type grep
@id test_grep
<and>
  <exclude>
	key elso
	pattern /^5\d\d$/
  </exclude>
  <exclude>
	key masodik
	pattern /\.css$/
  </exclude>
</and>
</filter>
`
	grep := &filter.GrepConfig{}
	yaml.Unmarshal(CONFIG, grep)
	test := render.NewOutputPluginTest(t, grep)
	test.DiffResult(expected)
}

func TestGrepMulti(t *testing.T) {
	CONFIG := []byte(`
exclude:
  - key: elso
    pattern: /^5\d\d$/
  - key: masodik
    pattern: /\.css$/
and:
  - regexp:
    - key: elso
      pattern: /^5\d\d$/
    - key: masodik
      pattern: /\.css$/
`)
	expected := `
<filter **>
@type grep
@id test_grep
<exclude>
  key elso
  pattern /^5\d\d$/
</exclude>
<exclude>
  key masodik
  pattern /\.css$/
</exclude>
<and>
  <regexp>
	key elso
	pattern /^5\d\d$/
  </regexp>
  <regexp>
	key masodik
	pattern /\.css$/
  </regexp>
</and>
</filter>
`
	grep := &filter.GrepConfig{}
	yaml.Unmarshal(CONFIG, grep)
	test := render.NewOutputPluginTest(t, grep)
	test.DiffResult(expected)
}
