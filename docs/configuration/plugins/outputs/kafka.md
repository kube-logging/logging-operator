---
title: Kafka
weight: 200
generated_file: true
---

# Kafka output plugin for Fluentd
## Overview


For details, see [https://github.com/fluent/fluent-plugin-kafka](https://github.com/fluent/fluent-plugin-kafka).

For an example deployment, see [Transport Nginx Access Logs into Kafka with Logging Operator](../../../../quickstarts/kafka-nginx/).

## Example output configurations

```yaml
spec:
  kafka:
    brokers: kafka-headless.kafka.svc.cluster.local:29092
    default_topic: topic
    sasl_over_ssl: false
    format:
      type: json
    buffer:
      tags: topic
      timekey: 1m
      timekey_wait: 30s
      timekey_use_utc: true
```


## Configuration
## Kafka

Send your logs to Kafka

### ack_timeout (int, optional) {#kafka-ack_timeout}

How long the producer waits for acks. The unit is seconds

Default: nil => Uses default of ruby-kafka library

### brokers (string, required) {#kafka-brokers}

The list of all seed brokers, with their host and port information. 


### buffer (*Buffer, optional) {#kafka-buffer}

[Buffer](../buffer/) 


### client_id (string, optional) {#kafka-client_id}

Client ID

Default: "kafka"

### compression_codec (string, optional) {#kafka-compression_codec}

The codec the producer uses to compress messages . The available options are gzip and snappy.

Default: nil

### default_message_key (string, optional) {#kafka-default_message_key}

The name of default message key .

Default: nil

### default_partition_key (string, optional) {#kafka-default_partition_key}

The name of default partition key .

Default: nil

### default_topic (string, optional) {#kafka-default_topic}

The name of default topic .

Default: nil

### discard_kafka_delivery_failed (bool, optional) {#kafka-discard_kafka_delivery_failed}

Discard the record where Kafka DeliveryFailed occurred

Default: false

### exclude_partion_key (bool, optional) {#kafka-exclude_partion_key}

Exclude Partition key

Default: false

### exclude_topic_key (bool, optional) {#kafka-exclude_topic_key}

Exclude Topic key

Default: false

### format (*Format, required) {#kafka-format}

[Format](../format/) 


### get_kafka_client_log (bool, optional) {#kafka-get_kafka_client_log}

Get Kafka Client log

Default: false

### headers (map[string]string, optional) {#kafka-headers}

Headers

Default: {}

### headers_from_record (map[string]string, optional) {#kafka-headers_from_record}

Headers from Record

Default: {}

### idempotent (bool, optional) {#kafka-idempotent}

Idempotent

Default: false

### kafka_agg_max_bytes (int, optional) {#kafka-kafka_agg_max_bytes}

Maximum value of total message size to be included in one batch transmission. .

Default: 4096

### kafka_agg_max_messages (int, optional) {#kafka-kafka_agg_max_messages}

Maximum number of messages to include in one batch transmission. .

Default: nil

### keytab (*secret.Secret, optional) {#kafka-keytab}


### max_send_retries (int, optional) {#kafka-max_send_retries}

Number of times to retry sending of messages to a leader

Default: 1

### message_key_key (string, optional) {#kafka-message_key_key}

Message Key

Default: "message_key"

### partition_key (string, optional) {#kafka-partition_key}

Partition

Default: "partition"

### partition_key_key (string, optional) {#kafka-partition_key_key}

Partition Key

Default: "partition_key"

### password (*secret.Secret, optional) {#kafka-password}

Password when using PLAIN/SCRAM SASL authentication 


### principal (string, optional) {#kafka-principal}


### required_acks (int, optional) {#kafka-required_acks}

The number of acks required per request .

Default: -1

### ssl_ca_cert (*secret.Secret, optional) {#kafka-ssl_ca_cert}

CA certificate 


### ssl_ca_certs_from_system (*bool, optional) {#kafka-ssl_ca_certs_from_system}

System's CA cert store

Default: false

### ssl_client_cert (*secret.Secret, optional) {#kafka-ssl_client_cert}

Client certificate 


### ssl_client_cert_chain (*secret.Secret, optional) {#kafka-ssl_client_cert_chain}

Client certificate chain 


### ssl_client_cert_key (*secret.Secret, optional) {#kafka-ssl_client_cert_key}

Client certificate key 


### ssl_verify_hostname (*bool, optional) {#kafka-ssl_verify_hostname}

Verify certificate hostname 


### sasl_over_ssl (bool, required) {#kafka-sasl_over_ssl}

SASL over SSL

Default: true

### scram_mechanism (string, optional) {#kafka-scram_mechanism}

If set, use SCRAM authentication with specified mechanism. When unset, default to PLAIN authentication 


### slow_flush_log_threshold (string, optional) {#kafka-slow_flush_log_threshold}

The threshold for chunk flush performance check. Parameter type is float, not time, default: 20.0 (seconds) If chunk flush takes longer time than this threshold, Fluentd logs a warning message and increases the  `fluentd_output_status_slow_flush_count` metric. 


### topic_key (string, optional) {#kafka-topic_key}

Topic Key

Default: "topic"

### use_default_for_unknown_topic (bool, optional) {#kafka-use_default_for_unknown_topic}

Use default for unknown topics

Default: false

### username (*secret.Secret, optional) {#kafka-username}

Username when using PLAIN/SCRAM SASL authentication 



