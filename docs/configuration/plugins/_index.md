---
title: Supported Plugins
generated_file: true
---

# Supported Plugins

For more information please click on the plugin name
<center>

| Name | Profile | Description | Status |Version |
|:---|---|:---|:---:|---:|
| **[Security](common/security/)** | common |  |  | []() |
| **[Transport](common/transport/)** | common |  |  | []() |
| **[Concat](filters/concat/)** | filters | Fluentd Filter plugin to concatenate multiline log separated in multiple events. | GA | [2.5.0](https://github.com/fluent-plugins-nursery/fluent-plugin-concat) |
| **[Dedot](filters/dedot/)** | filters | Concatenate multiline log separated in multiple events | GA | [1.0.0](https://github.com/lunardial/fluent-plugin-dedot_filter) |
| **[Exception Detector](filters/detect_exceptions/)** | filters | Exception Detector | GA | [0.0.13](https://github.com/GoogleCloudPlatform/fluent-plugin-detect-exceptions) |
| **[Enhance K8s Metadata](filters/enhance_k8s/)** | filters | Fluentd output plugin to add extra Kubernetes metadata to the events. | GA | [2.0.0](https://github.com/SumoLogic/sumologic-kubernetes-collection/tree/main/fluent-plugin-enhance-k8s-metadata) |
| **[Geo IP](filters/geoip/)** | filters | Fluentd GeoIP filter | GA | [1.3.2](https://github.com/y-ken/fluent-plugin-geoip) |
| **[Grep](filters/grep/)** | filters | Grep events by the values | GA | [more info](https://docs.fluentd.org/filter/grep) |
| **[Kubernetes Events Timestamp](filters/kube_events_timestamp/)** | filters | Fluentd Filter plugin to select particular timestamp into an additional field | GA | [0.1.4](https://github.com/banzaicloud/fluentd-filter-kube-events-timestamp) |
| **[Parser](filters/parser/)** | filters | Parses a string field in event records and mutates its event record with the parsed result. | GA | [more info](https://docs.fluentd.org/filter/parser) |
| **[Prometheus](filters/prometheus/)** | filters | Prometheus Filter Plugin to count Incoming Records | GA | [2.0.2](https://github.com/fluent/fluent-plugin-prometheus#prometheus-outputfilter-plugin) |
| **[Record Modifier](filters/record_modifier/)** | filters | Modify each event record. | GA | [2.1.0](https://github.com/repeatedly/fluent-plugin-record-modifier) |
| **[Record Transformer](filters/record_transformer/)** | filters | Mutates/transforms incoming event streams. | GA | [more info](https://docs.fluentd.org/filter/record_transformer) |
| **[Stdout](filters/stdout/)** | filters | Prints events to stdout | GA | [more info](https://docs.fluentd.org/filter/stdout) |
| **[SumoLogic](filters/sumologic/)** | filters | Sumo Logic collection solution for Kubernetes | GA | [2.3.1](https://github.com/SumoLogic/sumologic-kubernetes-collection) |
| **[Tag Normaliser](filters/tagnormaliser/)** | filters | Re-tag based on log metadata | GA | [0.1.1](https://github.com/banzaicloud/fluent-plugin-tag-normaliser) |
| **[Throttle](filters/throttle/)** | filters | A sentry plugin to throttle logs. Logs are grouped by a configurable key. When a group exceeds a configuration rate, logs are dropped for this group. | GA | [0.0.5](https://github.com/rubrikinc/fluent-plugin-throttle) |
| **[Amazon Elasticsearch](outputs/aws_elasticsearch/)** | outputs | Fluent plugin for Amazon Elasticsearch | Testing | [2.4.1](https://github.com/atomita/fluent-plugin-aws-elasticsearch-service) |
| **[Azure Storage](outputs/azurestore/)** | outputs | Store logs in Azure Storage | GA | [0.2.1](https://github.com/microsoft/fluent-plugin-azure-storage-append-blob) |
| **[Buffer](outputs/buffer/)** | outputs | Fluentd event buffer | GA | [mode info](https://docs.fluentd.org/configuration/buffer-section) |
| **[Amazon CloudWatch](outputs/cloudwatch/)** | outputs | Send your logs to AWS CloudWatch | GA | [0.14.0](https://github.com/fluent-plugins-nursery/fluent-plugin-cloudwatch-logs/releases/tag/v0.14.0) |
| **[Datadog](outputs/datadog/)** | outputs | Send your logs to Datadog | Testing | [0.13.0](https://github.com/DataDog/fluent-plugin-datadog/releases/tag/v0.13.0) |
| **[Elasticsearch](outputs/elasticsearch/)** | outputs | Send your logs to Elasticsearch | GA | [5.0.5](https://github.com/uken/fluent-plugin-elasticsearch/releases/tag/v5.0.5) |
| **[File](outputs/file/)** | outputs | Output plugin writes events to files | GA | [more info](https://docs.fluentd.org/output/file) |
| **[Format](outputs/format/)** | outputs | Specify how to format output record. | GA | [more info](https://docs.fluentd.org/configuration/format-section) |
| **[Format rfc5424](outputs/format_rfc5424/)** | outputs | Specify how to format output record. | GA | [more info](https://github.com/cloudfoundry/fluent-plugin-syslog_rfc5424#format-section) |
| **[Forward](outputs/forward/)** | outputs | Forwards events to other fluentd nodes. | GA | [more info](https://docs.fluentd.org/output/forward) |
| **[Google Cloud Storage](outputs/gcs/)** | outputs | Store logs in Google Cloud Storage | GA | [0.4.0](https://github.com/banzaicloud/fluent-plugin-gcs) |
| **[Gelf](outputs/gelf/)** | outputs | Output plugin writes events to GELF | Testing | [1.0.8](https://github.com/hotschedules/fluent-plugin-gelf-hs) |
| **[Http](outputs/http/)** | outputs | Sends logs to HTTP/HTTPS endpoints. | GA | [more info](https://docs.fluentd.org/output/http) |
| **[Kafka](outputs/kafka/)** | outputs | Send your logs to Kafka | GA | [0.17.0](https://github.com/fluent/fluent-plugin-kafka/releases/tag/v0.17.0) |
| **[Amazon Kinesis Firehose](outputs/kinesis_firehose/)** | outputs | Fluent plugin for Amazon Kinesis | Testing | [3.4.0](https://github.com/awslabs/aws-fluent-plugin-kinesis/releases/tag/v3.4.0) |
| **[Amazon Kinesis Stream](outputs/kinesis_stream/)** | outputs | Fluent plugin for Amazon Kinesis | GA | [3.4.0](https://github.com/awslabs/aws-fluent-plugin-kinesis/releases/tag/v3.4.0) |
| **[LogDNA](outputs/logdna/)** | outputs | Send your logs to LogDNA | GA | [0.4.0](https://github.com/logdna/fluent-plugin-logdna/releases/tag/v0.4.0) |
| **[LogZ](outputs/logz/)** | outputs | Store logs in LogZ.io | GA | [0.0.21](https://github.com/logzio/fluent-plugin-logzio/releases/tag/v0.0.21) |
| **[Grafana Loki](outputs/loki/)** | outputs | Transfer logs to Loki | GA | [1.2.16](https://github.com/grafana/loki/tree/master/fluentd/fluent-plugin-grafana-loki) |
| **[NewRelic Logs](outputs/newrelic/)** | outputs | Send logs to New Relic Logs | GA | [1.2.1](https://github.com/newrelic/newrelic-fluentd-output) |
| **[Alibaba Cloud Storage](outputs/oss/)** | outputs | Store logs the Alibaba Cloud Object Storage Service | GA | [0.0.2](https://github.com/aliyun/fluent-plugin-oss) |
| **[Redis](outputs/redis/)** | outputs | Sends logs to Redis endpoints. | GA | [0.3.5](https://github.com/fluent-plugins-nursery/fluent-plugin-redis) |
| **[Amazon S3](outputs/s3/)** | outputs | Store logs in Amazon S3 | GA | [1.6.1](https://github.com/fluent/fluent-plugin-s3/releases/tag/v1.6.1) |
| **[Splunk Hec](outputs/splunk_hec/)** | outputs | Fluent Plugin Splunk Hec Release 1.2.3 | GA | [1.2.7]() |
| **[SQS](outputs/sqs/)** | outputs | Output plugin writes fluent-events as queue messages to Amazon SQS | Testing | [v2.1.0](https://github.com/ixixi/fluent-plugin-sqs) |
| **[SumoLogic](outputs/sumologic/)** | outputs | Send your logs to Sumologic | GA | [2.0.0](https://github.com/SumoLogic/fluentd-output-sumologic/releases/tag/2.0.0) |
| **[Syslog](outputs/syslog/)** | outputs | Output plugin writes events to syslog | GA | [0.9.0.rc.8](https://github.com/cloudfoundry/fluent-plugin-syslog_rfc5424) |
</center>

