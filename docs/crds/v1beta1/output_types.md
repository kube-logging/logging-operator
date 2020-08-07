---
title: OutputSpec
weight: 200
generated_file: true
---

### OutputSpec
#### OutputSpec defines the desired state of Output

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| loggingRef | string | No | - |  |
| s3 | *output.S3OutputConfig | No | - |  |
| azurestorage | *output.AzureStorage | No | - |  |
| gcs | *output.GCSOutput | No | - |  |
| oss | *output.OSSOutput | No | - |  |
| elasticsearch | *output.ElasticsearchOutput | No | - |  |
| logz | *output.LogZOutput | No | - |  |
| loki | *output.LokiOutput | No | - |  |
| sumologic | *output.SumologicOutput | No | - |  |
| datadog | *output.DatadogOutput | No | - |  |
| forward | *output.ForwardOutput | No | - |  |
| file | *output.FileOutputConfig | No | - |  |
| nullout | *output.NullOutputConfig | No | - |  |
| kafka | *output.KafkaOutputConfig | No | - |  |
| cloudwatch | *output.CloudWatchOutput | No | - |  |
| kinesisStream | *output.KinesisStreamOutputConfig | No | - |  |
| logdna | *output.LogDNAOutput | No | - |  |
| newrelic | *output.NewRelicOutputConfig | No | - |  |
| splunkHec | *output.SplunkHecOutput | No | - |  |
| http | *output.HTTPOutputConfig | No | - |  |
| awsElasticsearch | *output.AwsElasticsearchOutputConfig | No | - |  |
| redis | *output.RedisOutputConfig | No | - |  |
### OutputStatus
#### OutputStatus defines the observed state of Output

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
### Output
#### Output is the Schema for the outputs API

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ObjectMeta | No | - |  |
| spec | OutputSpec | No | - |  |
| status | OutputStatus | No | - |  |
### OutputList
#### OutputList contains a list of Output

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
|  | metav1.TypeMeta | Yes | - |  |
| metadata | metav1.ListMeta | No | - |  |
| items | []Output | Yes | - |  |
