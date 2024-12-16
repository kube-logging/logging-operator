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

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/render"
)

func TestKafka(t *testing.T) {
	CONFIG := []byte(`
brokers: kafka-headless.kafka.svc.cluster.local:29092
default_topic: topic
sasl_over_ssl: false
ssl_verify_hostname: false
format:
  type: json
buffer:
  timekey: 1m
  timekey_wait: 30s
  timekey_use_utc: true
`)
	expected := `
  <match **>
    @type kafka2
    @id test
    brokers kafka-headless.kafka.svc.cluster.local:29092
    default_topic topic
    sasl_over_ssl false
    ssl_verify_hostname false
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
	kafka := &output.KafkaOutputConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, kafka))
	test := render.NewOutputPluginTest(t, kafka)
	test.DiffResult(expected)
}

func TestRdkafka(t *testing.T) {
	CONFIG := []byte(`
brokers: kafka-headless.kafka.svc.cluster.local:29092
default_topic: topic
use_rdkafka: true
ssl_verify_hostname: false
rdkafka_options:
  sasl.mechanisms: PLAIN
  sasl.username: user
  security.protocol: SASL_SSL
  ssl.ca.location: /etc/ssl/certs/ca-certificates.crt
  ssl.certificate.location: /etc/ssl/certs/tls.crt
  ssl.key.location: /etc/ssl/certs/tls.key
  ssl.key.password: password 
format:
  type: json
buffer:
  timekey: 1m
  timekey_wait: 30s
  timekey_use_utc: true
`)
	expected := `
  <match **>
    @type rdkafka2
    @id test
    brokers kafka-headless.kafka.svc.cluster.local:29092
    default_topic topic
    ssl_verify_hostname false
    <buffer tag,time>
      @type file
      path /buffers/test.*.buffer
      retry_forever true
      timekey 1m
      timekey_use_utc true
      timekey_wait 30s
    </buffer>
    <rdkafka_options>
      sasl.mechanisms PLAIN
      sasl.username user
      security.protocol SASL_SSL
      ssl.ca.location /etc/ssl/certs/ca-certificates.crt
      ssl.certificate.location /etc/ssl/certs/tls.crt
      ssl.key.location /etc/ssl/certs/tls.key
      ssl.key.password password
    </rdkafka_options>
    <format>
      @type json
    </format>
  </match>
`
	kafka := &output.KafkaOutputConfig{}
	require.NoError(t, yaml.Unmarshal(CONFIG, kafka))
	test := render.NewOutputPluginTest(t, kafka)
	test.DiffResult(expected)
}
