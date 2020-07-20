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

package output_test

import (
	"testing"

	"github.com/banzaicloud/logging-operator/pkg/sdk/model/output"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/render"
	"github.com/ghodss/yaml"
)

func TestLogDNA(t *testing.T) {
	CONFIG := []byte(`
api_key: xxxxxxxxxxxxxxxxxxxxxxxxxxy
hostname: logging-operator

app: myapp
`)
	expected := `
  <match **>
	@type logdna
	@id test_logdna
	api_key xxxxxxxxxxxxxxxxxxxxxxxxxxy
	app myapp
    hostname logging-operator
  </match>
`
	s := &output.LogDNAOutput{}
	yaml.Unmarshal(CONFIG, s)
	test := render.NewOutputPluginTest(t, s)
	test.DiffResult(expected)
}
