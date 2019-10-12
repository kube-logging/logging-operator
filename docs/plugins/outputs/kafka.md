# Kafka output plugin for Fluentd
## Overview
  More info at https://github.com/fluent/fluent-plugin-kafka
>Example Deployment: [Transport Nginx Access Logs into Kafka with Logging Operator](../../../docs/example-kafka-nginx.md)

 #### Example output configurations
 ```
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
### Kafka
#### Send your logs to Kafka

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| brokers | string | Yes | - | The list of all seed brokers, with their host and port information.<br> |
| topic_key | string | No |  "topic" | Topic Key <br> |
| partition_key | string | No |  "partition" | Partition <br> |
| partition_key_key | string | No |  "partition_key" | Partition Key <br> |
| message_key_key | string | No |  "message_key" | Message Key <br> |
| default_topic | string | No |  nil | The name of default topic .<br> |
| default_partition_key | string | No |  nil | The name of default partition key .<br> |
| default_message_key | string | No |  nil | The name of default message key .<br> |
| exclude_topic_key | bool | No |  false | Exclude Topic key <br> |
| exclude_partion_key | bool | No |  false | Exclude Partition key <br> |
| get_kafka_client_log | bool | No |  false | Get Kafka Client log <br> |
| headers | map[string]string | No |  {} | Headers <br> |
| headers_from_record | map[string]string | No |  {} | Headers from Record <br> |
| use_default_for_unknown_topic | bool | No |  false | Use default for unknown topics <br> |
| idempotent | bool | No |  false | Idempotent <br> |
| sasl_over_ssl | bool | Yes |  true | SASL over SSL <br> |
| max_send_retries | int | No |  1 | Number of times to retry sending of messages to a leader <br> |
| required_acks | int | No |  -1 | The number of acks required per request .<br> |
| ack_timeout | int | No |  nil => Uses default of ruby-kafka library | How long the producer waits for acks. The unit is seconds <br> |
| compression_codec | string | No |  nil | The codec the producer uses to compress messages . The available options are gzip and snappy.<br> |
| format | *Format | Yes | - | [Format](./format.md)<br> |
| buffer | *Buffer | No | - | [Buffer](./buffer.md)<br> |
