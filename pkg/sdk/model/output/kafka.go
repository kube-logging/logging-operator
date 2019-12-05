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

package output

import (
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/secret"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/types"
)

// +docName:"Kafka output plugin for Fluentd"
//  More info at https://github.com/fluent/fluent-plugin-kafka
//>Example Deployment: [Transport Nginx Access Logs into Kafka with Logging Operator](../../../docs/example-kafka-nginx.md)
//
// #### Example output configurations
// ```
// spec:
//   kafka:
//     brokers: kafka-headless.kafka.svc.cluster.local:29092
//     default_topic: topic
//     sasl_over_ssl: false
//     format:
//       type: json
//     buffer:
//       tags: topic
//       timekey: 1m
//       timekey_wait: 30s
//       timekey_use_utc: true
// ```
type _docKafka interface{}

// +name:"Kafka"
// +url:"https://github.com/fluent/fluent-plugin-kafka/releases/tag/v0.12.1"
// +version:"0.12.1"
// +description:"Send your logs to Kafka"
// +status:"GA"
type _metaKafka interface{}

// +kubebuilder:object:generate=true
// +docName:"Kafka"
// Send your logs to Kafka
type KafkaOutputConfig struct {

	// The list of all seed brokers, with their host and port information.
	Brokers string `json:"brokers"`
	// Topic Key (default: "topic")
	TopicKey string `json:"topic_key,omitempty"`
	// Partition (default: "partition")
	PartitionKey string `json:"partition_key,omitempty"`
	// Partition Key (default: "partition_key")
	PartitionKeyKey string `json:"partition_key_key,omitempty"`
	// Message Key (default: "message_key")
	MessageKeyKey string `json:"message_key_key,omitempty"`
	// The name of default topic (default: nil).
	DefaultTopic string `json:"default_topic,omitempty"`
	// The name of default partition key (default: nil).
	DefaultPartitionKey string `json:"default_partition_key,omitempty"`
	// The name of default message key (default: nil).
	DefaultMessageKey string `json:"default_message_key,omitempty"`
	// Exclude Topic key (default: false)
	ExcludeTopicKey bool `json:"exclude_topic_key,omitempty"`
	// Exclude Partition key (default: false)
	ExcludePartitionKey bool `json:"exclude_partion_key,omitempty"`
	// Get Kafka Client log (default: false)
	GetKafkaClientLog bool `json:"get_kafka_client_log,omitempty"`
	// Headers (default: {})
	Headers map[string]string `json:"headers,omitempty"`
	// Headers from Record (default: {})
	HeadersFromRecord map[string]string `json:"headers_from_record,omitempty"`
	// Use default for unknown topics (default: false)
	UseDefaultForUnknownTopic bool `json:"use_default_for_unknown_topic,omitempty"`
	// Idempotent (default: false)
	Idempotent bool `json:"idempotent,omitempty"`
	// SASL over SSL (default: true)
	// +kubebuilder:validation:Optional
	SaslOverSSL bool `json:"sasl_over_ssl"`
	// Number of times to retry sending of messages to a leader (default: 1)
	MaxSendRetries int `json:"max_send_retries,omitempty"`
	// The number of acks required per request (default: -1).
	RequiredAcks int `json:"required_acks,omitempty"`
	// How long the producer waits for acks. The unit is seconds (default: nil => Uses default of ruby-kafka library)
	AckTimeout int `json:"ack_timeout,omitempty"`
	// The codec the producer uses to compress messages (default: nil). The available options are gzip and snappy.
	CompressionCodec string `json:"compression_codec,omitempty"`
	// +docLink:"Format,./format.md"
	Format *Format `json:"format"`
	// +docLink:"Buffer,./buffer.md"
	Buffer *Buffer `json:"buffer,omitempty"`
}

func (e *KafkaOutputConfig) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	pluginType := "kafka2"
	pluginID := id + "_" + pluginType
	kafka := &types.OutputPlugin{
		PluginMeta: types.PluginMeta{
			Type:      pluginType,
			Directive: "match",
			Tag:       "**",
			Id:        pluginID,
		},
	}
	if params, err := types.NewStructToStringMapper(secretLoader).StringsMap(e); err != nil {
		return nil, err
	} else {
		kafka.Params = params
	}
	if e.Buffer != nil {
		if buffer, err := e.Buffer.ToDirective(secretLoader, pluginID); err != nil {
			return nil, err
		} else {
			kafka.SubDirectives = append(kafka.SubDirectives, buffer)
		}
	}
	if e.Format != nil {
		if format, err := e.Format.ToDirective(secretLoader, ""); err != nil {
			return nil, err
		} else {
			kafka.SubDirectives = append(kafka.SubDirectives, format)
		}
	}
	return kafka, nil
}
