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
	Host                             string         `json:"host,omitempty"`
	Hosts                            []string       `json:"hosts,omitempty"`
	Port                             int            `json:"port,omitempty"`
	User                             *secret.Secret `json:"user,omitempty"`
	Pass                             *secret.Secret `json:"pass,omitempty"`
	VHost                            string         `json:"vhost,omitempty"`
	ConnectionTimeoutInSeconds       int            `json:"connection_timeout,omitempty"`
	NetworkRecoveryIntervalInSeconds int            `json:"network_recovery_interval,omitempty"`
	ContinuationTimeoutInSeconds     int            `json:"continuation_timeout,omitempty"`
	RecoveryAttempts                 int            `json:"recovery_attempts,omitempty"`
	AutomaticallyRecover             bool           `json:"automatically_recover,omitempty"`
	HeartbeatIntervalInSeconds       int            `json:"heartbeat,omitempty"`
	FrameMax                         int            `json:"frame_max,omitempty"`

	TLS               bool     `json:"tls,omitempty"`
	TLSCert           string   `json:"tls_cert,omitempty"`
	TLSKey            string   `json:"tls_key,omitempty"`
	TLSCACertificates []string `json:"tls_ca_certificates,omitempty"`
	VerifyPeer        bool     `json:"verify_peer,omitempty"`

	Exchange          string `json:"exchange"`
	ExchangeType      string `json:"exchange_type"`
	ExchangeDurable   string `json:"exchange_durable,omitempty"`
	ExchangeNoDeclare string `json:"exchange_no_declare,omitempty"`

	RoutingKey      string `json:"routing_key,omitempty"`
	IdKey           string `json:"id_key,omitempty"`
	Timestamp       string `json:"timestamp,omitempty"`
	ContentType     string `json:"content_type,omitempty"`
	ContentEncoding string `json:"content_encoding,omitempty"`
	Expiration      int    `json:"expiration,omitempty"`
	MessageType     string `json:"message_type,omitempty"`
	Priority        int    `json:"priority,omitempty"`
	AppId           int    `json:"app_id,omitempty"`
	Persistent      bool   `json:"persistent,omitempty"`

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
