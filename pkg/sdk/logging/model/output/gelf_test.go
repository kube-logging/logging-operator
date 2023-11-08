// Copyright © 2021 Cisco Systems, Inc. and/or its affiliates
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

	"github.com/ghodss/yaml"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/render"
	"github.com/stretchr/testify/require"
)

func TestGELFOutputConfig(t *testing.T) {
	CONFIG := []byte(`
host: gelf-host
port: 12201
`)
	expected := `
  <match **>
    @type gelf
    @id test
    host gelf-host
    port 12201
  </match>
`
	s := &output.GELFOutputConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, s))
	test := render.NewOutputPluginTest(t, s)
	test.DiffResult(expected)
}
