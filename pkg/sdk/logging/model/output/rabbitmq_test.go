// Copyright © 2020 Banzai Cloud
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

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/render"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

func TestRabbitMQ(t *testing.T) {
	CONFIG := []byte(`
host: rabbitmq-master.prod.svc.cluster.local
port: 5672
user:
  value: test-user
pass:
  value: test-pass
exchange: test-exchange
exchange_type: fanout
format:
  type: json
buffer:
  timekey: 1m
  timekey_wait: 30s
  timekey_use_utc: true
`)

	expected := `
  <match **>
    @type rabbitmq
    @id test
    exchange test-exchange
    exchange_type fanout
    host rabbitmq-master.prod.svc.cluster.local
    pass test-pass
    port 5672
    user test-user
    <buffer tag,time>
      @type file
      path /buffers/test.*.buffer
      retry_forever true
      timekey 1m
      timekey_use_utc true
      timekey_wait 30s
    </buffer>
    <format>
      @type json
    </format>
  </match>
`

	rabbitmq := &output.RabbitMQOutputConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, rabbitmq))
	test := render.NewOutputPluginTest(t, rabbitmq)
	test.DiffResult(expected)
}
