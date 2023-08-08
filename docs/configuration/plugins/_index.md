---
title: Supported Plugins
generated_file: true
---

For more information please click on the plugin name
<center>

| Name | Profile | Description | Status |Version |
|:---|---|:---|:---:|---:|
| **[Security](common/security/)** | common |  |  | []() |
| **[Transport](common/transport/)** | common |  |  | []() |
| **[Concat](filters/concat/)** | filters | Fluentd Filter plugin to concatenate multiline log separated in multiple events. | GA | [2.5.0](https://github.com/fluent-plugins-nursery/fluent-plugin-concat) |
| **[Dedot](filters/dedot/)** | filters | Concatenate multiline log separated in multiple events | GA | [1.0.0](https://github.com/lunardial/fluent-plugin-dedot_filter) |
| **[Exception Detector](filters/detect_exceptions/)** | filters | Exception Detector | GA | [0.0.14](https://github.com/GoogleCloudPlatform/fluent-plugin-detect-exceptions) |
| **[ElasticsearchGenId](filters/elasticsearch_genid/)** | filters |  |  | []() |
| **[Enhance K8s Metadata](filters/enhance_k8s/)** | filters | Fluentd output plugin to add extra Kubernetes metadata to the events. | GA | [2.0.0](https://github.com/SumoLogic/sumologic-kubernetes-collection/tree/main/fluent-plugin-enhance-k8s-metadata) |
| **[Geo IP](filters/geoip/)** | filters | Fluentd GeoIP filter | GA | [1.3.2](https://github.com/y-ken/fluent-plugin-geoip) |
| **[Grep](filters/grep/)** | filters | Grep events by the values | GA | [more info](https://docs.fluentd.org/filter/grep) |
| **[Kubernetes Events Timestamp](filters/kube_events_timestamp/)** | filters | Fluentd Filter plugin to select particular timestamp into an additional field | GA | [0.1.4](https://github.com/kube-logging/fluentd-filter-kube-events-timestamp) |
| **[Parser](filters/parser/)** | filters | Parses a string field in event records and mutates its event record with the parsed result. | GA | [more info](https://docs.fluentd.org/filter/parser) |
| **[Prometheus](filters/prometheus/)** | filters | Prometheus Filter Plugin to count Incoming Records | GA | [2.0.2](https://github.com/fluent/fluent-plugin-prometheus#prometheus-outputfilter-plugin) |
| **[Record Modifier](filters/record_modifier/)** | filters | Modify each event record. | GA | [2.1.0](https://github.com/repeatedly/fluent-plugin-record-modifier) |
| **[Record Transformer](filters/record_transformer/)** | filters | Mutates/transforms incoming event streams. | GA | [more info](https://docs.fluentd.org/filter/record_transformer) |
| **[Stdout](filters/stdout/)** | filters | Prints events to stdout | GA | [more info](https://docs.fluentd.org/filter/stdout) |
| **[SumoLogic](filters/sumologic/)** | filters | Sumo Logic collection solution for Kubernetes | GA | [2.3.1](https://github.com/SumoLogic/sumologic-kubernetes-collection) |
| **[Tag Normaliser](filters/tagnormaliser/)** | filters | Re-tag based on log metadata | GA | [0.1.1](https://github.com/kube-logging/fluent-plugin-tag-normaliser) |
| **[Throttle](filters/throttle/)** | filters | A sentry plugin to throttle logs. Logs are grouped by a configurable key. When a group exceeds a configuration rate, logs are dropped for this group. | GA | [0.0.5](https://github.com/rubrikinc/fluent-plugin-throttle) |
| **[Amazon Elasticsearch](outputs/aws_elasticsearch/)** | outputs | Fluent plugin for Amazon Elasticsearch | Testing | [2.4.1](https://github.com/atomita/fluent-plugin-aws-elasticsearch-service) |
| **[Azure Storage](outputs/azurestore/)** | outputs | Store logs in Azure Storage | GA | [0.2.1](https://github.com/microsoft/fluent-plugin-azure-storage-append-blob) |
| **[Buffer](outputs/buffer/)** | outputs | Fluentd event buffer | GA | [mode info](https://docs.fluentd.org/configuration/buffer-section) |
| **[Amazon CloudWatch](outputs/cloudwatch/)** | outputs | Send your logs to AWS CloudWatch | GA | [0.14.2](https://github.com/fluent-plugins-nursery/fluent-plugin-cloudwatch-logs/releases/tag/v0.14.2) |
| **[Datadog](outputs/datadog/)** | outputs | Send your logs to Datadog | Testing | [0.14.1](https://github.com/DataDog/fluent-plugin-datadog/releases/tag/v0.14.1) |
| **[Elasticsearch](outputs/elasticsearch/)** | outputs | Send your logs to Elasticsearch | GA | [5.1.1](https://github.com/uken/fluent-plugin-elasticsearch/releases/tag/v5.1.4) |
| **[File](outputs/file/)** | outputs | Output plugin writes events to files | GA | [more info](https://docs.fluentd.org/output/file) |
| **[Format](outputs/format/)** | outputs | Specify how to format output record. | GA | [more info](https://docs.fluentd.org/configuration/format-section) |
| **[Format rfc5424](outputs/format_rfc5424/)** | outputs | Specify how to format output record. | GA | [more info](https://github.com/cloudfoundry/fluent-plugin-syslog_rfc5424#format-section) |
| **[Forward](outputs/forward/)** | outputs | Forwards events to other fluentd nodes. | GA | [more info](https://docs.fluentd.org/output/forward) |
| **[Google Cloud Storage](outputs/gcs/)** | outputs | Store logs in Google Cloud Storage | GA | [0.4.0](https://github.com/kube-logging/fluent-plugin-gcs) |
| **[Gelf](outputs/gelf/)** | outputs | Output plugin writes events to GELF | Testing | [1.0.8](https://github.com/hotschedules/fluent-plugin-gelf-hs) |
| **[Http](outputs/http/)** | outputs | Sends logs to HTTP/HTTPS endpoints. | GA | [more info](https://docs.fluentd.org/output/http) |
| **[Kafka](outputs/kafka/)** | outputs | Send your logs to Kafka | GA | [0.17.5](https://github.com/fluent/fluent-plugin-kafka/releases/tag/v0.17.5) |
| **[Amazon Kinesis Firehose](outputs/kinesis_firehose/)** | outputs | Fluent plugin for Amazon Kinesis | Testing | [3.4.2](https://github.com/awslabs/aws-fluent-plugin-kinesis/releases/tag/v3.4.2) |
| **[Amazon Kinesis Stream](outputs/kinesis_stream/)** | outputs | Fluent plugin for Amazon Kinesis | GA | [3.4.2](https://github.com/awslabs/aws-fluent-plugin-kinesis/releases/tag/v3.4.2) |
| **[LogDNA](outputs/logdna/)** | outputs | Send your logs to LogDNA | GA | [0.4.0](https://github.com/logdna/fluent-plugin-logdna/releases/tag/v0.4.0) |
| **[LogZ](outputs/logz/)** | outputs | Store logs in LogZ.io | GA | [0.0.21](https://github.com/logzio/fluent-plugin-logzio/releases/tag/v0.0.21) |
| **[Grafana Loki](outputs/loki/)** | outputs | Transfer logs to Loki | GA | [1.2.19](https://github.com/grafana/loki/tree/master/fluentd/fluent-plugin-grafana-loki) |
| **[Mattermost](outputs/mattermost/)** | outputs | Sends logs to Mattermost via webhooks. | GA | [0.2.2](https://github.com/levigo-systems/fluent-plugin-mattermost) |
| **[NewRelic Logs](outputs/newrelic/)** | outputs | Send logs to New Relic Logs | GA | [1.2.1](https://github.com/newrelic/newrelic-fluentd-output) |
| **[OpenSearch](outputs/opensearch/)** | outputs | Send your logs to OpenSearch | GA | [1.0.5](https://github.com/fluent/fluent-plugin-opensearch/releases/tag/v1.0.5) |
| **[Alibaba Cloud Storage](outputs/oss/)** | outputs | Store logs the Alibaba Cloud Object Storage Service | GA | [0.0.2](https://github.com/aliyun/fluent-plugin-oss) |
| **[Redis](outputs/redis/)** | outputs | Sends logs to Redis endpoints. | GA | [0.3.5](https://github.com/fluent-plugins-nursery/fluent-plugin-redis) |
| **[Relabel](outputs/relabel/)** | outputs | Relabel output plugin re-labels events. | GA | [more info](https://docs.fluentd.org/output/relabel) |
| **[Amazon S3](outputs/s3/)** | outputs | Store logs in Amazon S3 | GA | [1.6.1](https://github.com/fluent/fluent-plugin-s3/releases/tag/v1.6.1) |
| **[Splunk Hec](outputs/splunk_hec/)** | outputs | Fluent Plugin Splunk Hec Release | GA | [1.2.9]() |
| **[SQS](outputs/sqs/)** | outputs | Output plugin writes fluent-events as queue messages to Amazon SQS | Testing | [v2.1.0](https://github.com/ixixi/fluent-plugin-sqs) |
| **[SumoLogic](outputs/sumologic/)** | outputs | Send your logs to Sumologic | GA | [1.8.0](https://github.com/SumoLogic/fluentd-output-sumologic/releases/tag/1.8.0) |
| **[Syslog](outputs/syslog/)** | outputs | Output plugin writes events to syslog | GA | [0.9.0.rc.8](https://github.com/cloudfoundry/fluent-plugin-syslog_rfc5424) |
| **[Syslog-NG Match](syslogng-filters/match/)** | syslogng-filters | Selectively keep records | GA | [more info](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/65#TOPIC-1829159) |
| **[Syslog-NG Parser](syslogng-filters/parser/)** | syslogng-filters | Parse data from records | GA | [more info](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.36/administration-guide/90) |
| **[Syslog-NG Rewrite](syslogng-filters/rewrite/)** | syslogng-filters | Rewrite parts of the message | GA | [more info](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/77) |
| **[disk-buffer configuration](syslogng-outputs/disk_buffer/)** | syslogng-outputs | disk-buffer configuration | Testing | [](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/32#kanchor2338) |
| **[File](syslogng-outputs/file/)** | syslogng-outputs | SStoring messages in plain-text files | Testing | [](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.17/administration-guide/32) |
| **[HTTP](syslogng-outputs/http/)** | syslogng-outputs | Sending messages over HTTP | Testing | [](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/40#TOPIC-1829058) |
| **[Loggly](syslogng-outputs/loggly/)** | syslogng-outputs | Send your logs to loggly | Testing | [](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/43#TOPIC-1829072) |
| **[Falcon LogScale](syslogng-outputs/logscale/)** | syslogng-outputs | Storing messages in Falcon's LogScale over http | Testing | [](https://library.humio.com/falcon-logscale/api-ingest.html#api-ingest-structured-data) |
| **[MQTT Destination](syslogng-outputs/mqtt/)** | syslogng-outputs | Sending messages over MQTT Protocol | Testing | [](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/45#TOPIC-1829079) |
| **[Sumo Logic HTTP](syslogng-outputs/sumologic_http/)** | syslogng-outputs | Storing messages in Sumo Logic over http | Testing | [](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/55) |
| **[Sumo Logic Syslog](syslogng-outputs/sumologic_syslog/)** | syslogng-outputs | Storing messages in Sumo Logic over syslog | Testing | [](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/56#TOPIC-1829122) |
| **[Syslog output configuration](syslogng-outputs/syslog/)** | syslogng-outputs | Syslog output configuration | Testing | [](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/32#kanchor2338) |
| **[TLS config for syslog-ng outputs](syslogng-outputs/tls/)** | syslogng-outputs | TLS config for syslog-ng outputs | Testing | [](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/32#kanchor2338) |
</center>

