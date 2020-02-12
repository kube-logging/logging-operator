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

func TestWebsocket(t *testing.T) {
	CONFIG := []byte(`
host: 192.168.1.1
port: 8080
use_msgpack: true
add_time: false
add_tag: false
buffered_messages: 100
token:
  value: SomeToken
`)
	expected := `
  <match **>
    @type websocket
		@id test_websocket
		add_tag false
		add_time false
		buffered_messages 100
		host 192.168.1.1
		port 8080
		token SomeToken
		use_msgpack true
  </match>
`
	s := &output.WebsocketOutput{}
	yaml.Unmarshal(CONFIG, s)
	test := render.NewOutputPluginTest(t, s)
	test.DiffResult(expected)
}
