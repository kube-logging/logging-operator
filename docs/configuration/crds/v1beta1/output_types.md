---
title: OutputSpec
weight: 200
generated_file: true
---

## OutputSpec

OutputSpec defines the desired state of Output

### awsElasticsearch (*output.AwsElasticsearchOutputConfig, optional) {#outputspec-awselasticsearch}

Default: -

### azurestorage (*output.AzureStorage, optional) {#outputspec-azurestorage}

Default: -

### cloudwatch (*output.CloudWatchOutput, optional) {#outputspec-cloudwatch}

Default: -

### datadog (*output.DatadogOutput, optional) {#outputspec-datadog}

Default: -

### elasticsearch (*output.ElasticsearchOutput, optional) {#outputspec-elasticsearch}

Default: -

### file (*output.FileOutputConfig, optional) {#outputspec-file}

Default: -

### forward (*output.ForwardOutput, optional) {#outputspec-forward}

Default: -

### gcs (*output.GCSOutput, optional) {#outputspec-gcs}

Default: -

### gelf (*output.GELFOutputConfig, optional) {#outputspec-gelf}

Default: -

### http (*output.HTTPOutputConfig, optional) {#outputspec-http}

Default: -

### kafka (*output.KafkaOutputConfig, optional) {#outputspec-kafka}

Default: -

### kinesisFirehose (*output.KinesisFirehoseOutputConfig, optional) {#outputspec-kinesisfirehose}

Default: -

### kinesisStream (*output.KinesisStreamOutputConfig, optional) {#outputspec-kinesisstream}

Default: -

### logdna (*output.LogDNAOutput, optional) {#outputspec-logdna}

Default: -

### logz (*output.LogZOutput, optional) {#outputspec-logz}

Default: -

### loggingRef (string, optional) {#outputspec-loggingref}

Default: -

### loki (*output.LokiOutput, optional) {#outputspec-loki}

Default: -

### mattermost (*output.MattermostOutputConfig, optional) {#outputspec-mattermost}

Default: -

### newrelic (*output.NewRelicOutputConfig, optional) {#outputspec-newrelic}

Default: -

### nullout (*output.NullOutputConfig, optional) {#outputspec-nullout}

Default: -

### oss (*output.OSSOutput, optional) {#outputspec-oss}

Default: -

### opensearch (*output.OpenSearchOutput, optional) {#outputspec-opensearch}

Default: -

### redis (*output.RedisOutputConfig, optional) {#outputspec-redis}

Default: -

### relabel (*output.RelabelOutputConfig, optional) {#outputspec-relabel}

Default: -

### s3 (*output.S3OutputConfig, optional) {#outputspec-s3}

Default: -

### sqs (*output.SQSOutputConfig, optional) {#outputspec-sqs}

Default: -

### splunkHec (*output.SplunkHecOutput, optional) {#outputspec-splunkhec}

Default: -

### sumologic (*output.SumologicOutput, optional) {#outputspec-sumologic}

Default: -

### syslog (*output.SyslogOutputConfig, optional) {#outputspec-syslog}

Default: -


## OutputStatus

OutputStatus defines the observed state of Output

### active (*bool, optional) {#outputstatus-active}

Default: -

### problems ([]string, optional) {#outputstatus-problems}

Default: -

### problemsCount (int, optional) {#outputstatus-problemscount}

Default: -


## Output

Output is the Schema for the outputs API

###  (metav1.TypeMeta, required) {#output-}

Default: -

### metadata (metav1.ObjectMeta, optional) {#output-metadata}

Default: -

### spec (OutputSpec, optional) {#output-spec}

Default: -

### status (OutputStatus, optional) {#output-status}

Default: -


## OutputList

OutputList contains a list of Output

###  (metav1.TypeMeta, required) {#outputlist-}

Default: -

### metadata (metav1.ListMeta, optional) {#outputlist-metadata}

Default: -

### items ([]Output, required) {#outputlist-items}

Default: -


