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

func TestDetectExceptions(t *testing.T) {
	CONFIG := []byte(`
multiline_flush_interval: 0.1
languages:
  - java
  - python
`)
	expected := `
<match kubernetes.**>
  @type detect_exceptions
  @id test
  languages ["java","python"]
  multiline_flush_interval 0.1
  remove_tag_prefix kubernetes
</match>
`
	ed := &filter.DetectExceptions{}
	require.NoError(t, yaml.Unmarshal(CONFIG, ed))
	test := render.NewOutputPluginTest(t, ed)
	test.DiffResult(expected)
}
