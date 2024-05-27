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

	"github.com/ghodss/yaml"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/render"
	"github.com/stretchr/testify/require"
)

func TestForward(t *testing.T) {
	CONFIG := []byte(`
  servers: 
    - host: 192.168.1.3
      name: myserver1
      port: 24224
      weight: 60
    - host: 192.168.1.4
      name: myserver2
      port: 24223
      weight: 40
  buffer:
    timekey: 1m
    timekey_wait: 30s
    timekey_use_utc: true
  keepalive: true
  keepalive_timeout: 20
  time_as_integer: true
  send_timeout: 60
`)
	expected := `
  <match **>
    @type forward
    @id test
    keepalive true
    keepalive_timeout 20
    send_timeout 60
    time_as_integer true
    <buffer tag,time>
      @type file
      path /buffers/test.*.buffer
      retry_forever true
      timekey 1m
      timekey_use_utc true
      timekey_wait 30s
    </buffer>
    <server>
      host 192.168.1.3
      name myserver1
      port 24224
      weight 60
    </server>
    <server>
      host 192.168.1.4
      name myserver2
      port 24223
      weight 40
    </server>
  </match>
`
	g := &output.ForwardOutput{}
	require.NoError(t, yaml.Unmarshal(CONFIG, g))
	test := render.NewOutputPluginTest(t, g)
	test.DiffResult(expected)
}
