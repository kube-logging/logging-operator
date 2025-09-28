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

package output

import (
	"github.com/cisco-open/operator-tools/pkg/secret"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/types"
)

// +name:"RabbitMQ"
// +weight:"200"
type _hugoRabbitMQ interface{} //nolint:deadcode,unused

// +docName:"RabbitMQ plugin for Fluentd"
/*
Sends logs to RabbitMQ Queues. For details, see [https://github.com/nttcom/fluent-plugin-rabbitmq](https://github.com/nttcom/fluent-plugin-rabbitmq).

## Example output configurations

```yaml
spec:
  rabbitmq:
    host: rabbitmq-master.prod.svc.cluster.local
		user: test-user
		pass: test-pass
		port: 5672
    buffer:
      tags: "[]"
      flush_interval: 10s
```
*/
type _docRabbitMQ interface{} //nolint:deadcode,unused

// +name:"RabbitMQ"
// +url:"https://github.com/nttcom/fluent-plugin-rabbitmq"
// +version:"0.1.5"
// +description:"Sends logs to RabbitMQ Queues."
// +status:"GA"
type _metaRabbitMQ interface{} //nolint:deadcode,unused

// +kubebuilder:object:generate=true
// +docName:"Output Config"
type RabbitMQOutputConfig struct {
	// Host
	Host string `json:"host,omitempty"`
	// Hosts
	Hosts []string `json:"hosts,omitempty"`
	// Port
	Port int `json:"port,omitempty"`
	// Username
	Username *secret.Secret `json:"user,omitempty"`
	// Password
	Password *secret.Secret `json:"pass,omitempty"`
	// VHost
	VHost string `json:"vhost,omitempty"`
	// Connection Timeout in seconds
	ConnectionTimeoutInSeconds int `json:"connection_timeout,omitempty"`
	// Continuation Timeout in seconds
	ContinuationTimeoutInSeconds int `json:"continuation_timeout,omitempty"`
	// Automatic network failure recovery
	AutomaticallyRecover bool `json:"automatically_recover,omitempty"`
	// Network Recovery Interval in seconds
	NetworkRecoveryIntervalInSeconds int `json:"network_recovery_interval,omitempty"`
	// Recovery Attempts
	RecoveryAttempts int `json:"recovery_attempts,omitempty"`
	// Auth Mechanism
	AuthMechanism string `json:"auth_mechanism,omitempty"`
	// Heartbeat Timeout in seconds
	HeartbeatIntervalInSeconds int `json:"heartbeat,omitempty"`
	// Maximum permissible size of a frame
	FrameMax int `json:"frame_max,omitempty"`

	// Enable TLS or not
	TLS bool `json:"tls,omitempty"`
	// Path to TLS certificate file
	TLSCert string `json:"tls_cert,omitempty"`
	// Path to TLS key file
	TLSKey string `json:"tls_key,omitempty"`
	// Path to TLS CA certificates files
	TLSCACertificates []string `json:"tls_ca_certificates,omitempty"`
	// Verify Peer or not
	VerifyPeer bool `json:"verify_peer,omitempty"`

	// Name of the exchange
	Exchange string `json:"exchange"`
	// Type of the exchange
	ExchangeType string `json:"exchange_type"`
	// Exchange durability
	ExchangeDurable bool `json:"exchange_durable,omitempty"`
	// Weather to declare exchange or not
	ExchangeNoDeclare string `json:"exchange_no_declare,omitempty"`

	// Messages are persistent to disk
	Persistent bool `json:"persistent,omitempty"`
	// Routing key to route messages
	RoutingKey string `json:"routing_key,omitempty"`
	// Id to specify message_id
	IdKey string `json:"id_key,omitempty"`
	// Time of record is used as timestamp in AMQP message
	Timestamp bool `json:"timestamp,omitempty"`
	// Message content type
	ContentType string `json:"content_type,omitempty"`
	// Message content encoding
	ContentEncoding string `json:"content_encoding,omitempty"`
	// Message message time-to-live in seconds
	ExpirationInSeconds int `json:"expiration,omitempty"`
	// Message type
	MessageType string `json:"message_type,omitempty"`
	// Message priority
	Priority int `json:"priority,omitempty"`
	// Application Id
	AppId string `json:"app_id,omitempty"`

	// +docLink:"Format,../format/"
	Format *Format `json:"format,omitempty"`
	// +docLink:"Buffer,../buffer/"
	Buffer *Buffer `json:"buffer,omitempty"`
}

func (c *RabbitMQOutputConfig) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	const pluginType = "rabbitMQ"
	rabbitMQ := &types.OutputPlugin{
		PluginMeta: types.PluginMeta{
			Type:      pluginType,
			Directive: "match",
			Tag:       "**",
			Id:        id,
		},
	}
	if params, err := types.NewStructToStringMapper(secretLoader).StringsMap(c); err != nil {
		return nil, err
	} else {
		rabbitMQ.Params = params
	}

	if c.Buffer == nil {
		c.Buffer = &Buffer{}
	}

	if buffer, err := c.Buffer.ToDirective(secretLoader, id); err != nil {
		return nil, err
	} else {
		rabbitMQ.SubDirectives = append(rabbitMQ.SubDirectives, buffer)
	}

	if c.Format != nil {
		if format, err := c.Format.ToDirective(secretLoader, ""); err != nil {
			return nil, err
		} else {
			rabbitMQ.SubDirectives = append(rabbitMQ.SubDirectives, format)
		}
	}

	return rabbitMQ, nil
}
