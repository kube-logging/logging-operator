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

func TestParser(t *testing.T) {
	CONFIG := []byte(`
remove_key_name_field: true
reserve_data: true
parsers:
- type: nginx
`)
	expected := `
<filter **>
  @type parser
  @id test_parser
  key_name message
  remove_key_name_field true
  reserve_data true
  <parse>
    @type nginx
  </parse>
</filter>
`
	parser := &filter.ParserConfig{}
	yaml.Unmarshal(CONFIG, parser)
	test := render.NewOutputPluginTest(t, parser)
	test.DiffResult(expected)
}

func TestParserMultiParser(t *testing.T) {
	CONFIG := []byte(`
remove_key_name_field: true
reserve_data: true
parsers:
- type: multi_format
  patterns:
  - format: nginx
  - format: regexp
    expression: /asdsada/
  - format: none
`)
	expected := `
<filter **>
  @type parser
  @id test_parser
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
	yaml.Unmarshal(CONFIG, parser)
	test := render.NewOutputPluginTest(t, parser)
	test.DiffResult(expected)
}
