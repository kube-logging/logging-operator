---
title: OutputSpec
weight: 200
generated_file: true
---

## OutputSpec

OutputSpec defines the desired state of Output

### awsElasticsearch (*output.AwsElasticsearchOutputConfig, optional) {#outputspec-awselasticsearch}


### azurestorage (*output.AzureStorage, optional) {#outputspec-azurestorage}


### cloudwatch (*output.CloudWatchOutput, optional) {#outputspec-cloudwatch}


### datadog (*output.DatadogOutput, optional) {#outputspec-datadog}


### elasticsearch (*output.ElasticsearchOutput, optional) {#outputspec-elasticsearch}


### file (*output.FileOutputConfig, optional) {#outputspec-file}


### forward (*output.ForwardOutput, optional) {#outputspec-forward}


### gcs (*output.GCSOutput, optional) {#outputspec-gcs}


### gelf (*output.GELFOutputConfig, optional) {#outputspec-gelf}


### http (*output.HTTPOutputConfig, optional) {#outputspec-http}


### kafka (*output.KafkaOutputConfig, optional) {#outputspec-kafka}


### kinesisFirehose (*output.KinesisFirehoseOutputConfig, optional) {#outputspec-kinesisfirehose}


### kinesisStream (*output.KinesisStreamOutputConfig, optional) {#outputspec-kinesisstream}


### logdna (*output.LogDNAOutput, optional) {#outputspec-logdna}


### logz (*output.LogZOutput, optional) {#outputspec-logz}


### loggingRef (string, optional) {#outputspec-loggingref}


### loki (*output.LokiOutput, optional) {#outputspec-loki}


### mattermost (*output.MattermostOutputConfig, optional) {#outputspec-mattermost}


### newrelic (*output.NewRelicOutputConfig, optional) {#outputspec-newrelic}


### nullout (*output.NullOutputConfig, optional) {#outputspec-nullout}


### oss (*output.OSSOutput, optional) {#outputspec-oss}


### opensearch (*output.OpenSearchOutput, optional) {#outputspec-opensearch}


### redis (*output.RedisOutputConfig, optional) {#outputspec-redis}


### relabel (*output.RelabelOutputConfig, optional) {#outputspec-relabel}


### s3 (*output.S3OutputConfig, optional) {#outputspec-s3}


### sqs (*output.SQSOutputConfig, optional) {#outputspec-sqs}


### splunkHec (*output.SplunkHecOutput, optional) {#outputspec-splunkhec}


### sumologic (*output.SumologicOutput, optional) {#outputspec-sumologic}


### syslog (*output.SyslogOutputConfig, optional) {#outputspec-syslog}


### vmwareLogInsight (*output.VMwareLogInsightOutput, optional) {#outputspec-vmwareloginsight}


### vmwareLogIntelligence (*output.VMwareLogIntelligenceOutputConfig, optional) {#outputspec-vmwarelogintelligence}



## OutputStatus

OutputStatus defines the observed state of Output

### active (*bool, optional) {#outputstatus-active}


### problems ([]string, optional) {#outputstatus-problems}


### problemsCount (int, optional) {#outputstatus-problemscount}



## Output

Output is the Schema for the outputs API

###  (metav1.TypeMeta, required) {#output-}


### metadata (metav1.ObjectMeta, optional) {#output-metadata}


### spec (OutputSpec, optional) {#output-spec}


### status (OutputStatus, optional) {#output-status}



## OutputList

OutputList contains a list of Output

###  (metav1.TypeMeta, required) {#outputlist-}


### metadata (metav1.ListMeta, optional) {#outputlist-metadata}


### items ([]Output, required) {#outputlist-items}



