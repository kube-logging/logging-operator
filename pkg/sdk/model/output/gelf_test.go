// Copyright © 2021 Banzai Cloud
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

func TestGELFOutputConfig(t *testing.T) {
	CONFIG := []byte(`
host: gelf-host
port: 12201
buffer:
  flush_thread_count: 8
  flush_interval: 5s
  chunk_limit_size: 8M
  queue_limit_length: 512
  retry_max_interval: 30
  retry_forever: true
`)
	expected := `
  <match **>
    @type gelf
    @id test
    host gelf-host
    port 12201
    <buffer tag,time>
      @type file
      chunk_limit_size 8M
      flush_interval 5s
      flush_thread_count 8
      path /buffers/test.*.buffer
      queue_limit_length 512
      retry_forever true
      retry_max_interval 30
      timekey 10m
      timekey_wait 10m
    </buffer>
  </match>
`
	s := &output.GELFOutputConfig{}
	yaml.Unmarshal(CONFIG, s)
	test := render.NewOutputPluginTest(t, s)
	test.DiffResult(expected)
}
