---
title: SumoLogic
weight: 200
generated_file: true
---

# SumoLogic output plugin for Fluentd
## Overview
This plugin has been designed to output logs or metrics to SumoLogic via a HTTP collector endpoint
More info at https://github.com/SumoLogic/fluentd-output-sumologic

 Example secret for HTTP input URL
 ```
export URL='https://endpoint1.collection.eu.sumologic.com/receiver/v1/http/.......'
kubectl create secret generic sumo-output --from-literal "endpoint=$URL"
```

 Example ClusterOutput

```
apiVersion: logging.banzaicloud.io/v1beta1
kind: ClusterOutput
metadata:
  name: sumo-output
spec:
  sumologic:
    buffer:
      flush_interval: 10s
      flush_mode: interval
    compress: true
    endpoint:
      valueFrom:
        secretKeyRef:
          key: endpoint
          name: sumo-output
    source_name: test1
```

## Configuration
### Output Config
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| data_type | string | No |  logs | The type of data that will be sent to Sumo Logic, either logs or metrics <br> |
| endpoint | *secret.Secret | Yes | - | SumoLogic HTTP Collector URL<br> |
| verify_ssl | bool | No |  true | Verify ssl certificate. <br> |
| metric_data_format | string | No |  graphite | The format of metrics you will be sending, either graphite or carbon2 or prometheus <br> |
| log_format | string | No |  json | Format to post logs into Sumo. <br> |
| log_key | string | No |  message | Used to specify the key when merging json or sending logs in text format <br> |
| source_category | string | No |  nil | Set _sourceCategory metadata field within SumoLogic <br> |
| source_name | string | Yes | - | Set _sourceName metadata field within SumoLogic - overrides source_name_key (default is nil)<br> |
| source_name_key | string | No |  source_name | Set as source::path_key's value so that the source_name can be extracted from Fluentd's buffer <br> |
| source_host | string | No |  nil | Set _sourceHost metadata field within SumoLogic <br> |
| open_timeout | int | No |  60 | Set timeout seconds to wait until connection is opened. <br> |
| add_timestamp | bool | No |  true | Add timestamp (or timestamp_key) field to logs before sending to sumologic <br> |
| timestamp_key | string | No |  timestamp | Field name when add_timestamp is on <br> |
| proxy_uri | string | No | - | Add the uri of the proxy environment if present.<br> |
| disable_cookies | bool | No |  false | Option to disable cookies on the HTTP Client. <br> |
| delimiter | string | No |  . | Delimiter <br> |
| custom_fields | []string | No | - | Comma-separated key=value list of fields to apply to every log. [more information](https://help.sumologic.com/Manage/Fields#http-source-fields)<br> |
| sumo_client | string | No |  fluentd-output | Name of sumo client which is send as X-Sumo-Client header <br> |
| compress | *bool | No |  false | Compress payload <br> |
| compress_encoding | string | No |  gzip | Encoding method of compression (either gzip or deflate) <br> |
| custom_dimensions | string | No | - | Dimensions string (eg "cluster=payment, service=credit_card") which is going to be added to every metric record.<br> |
| buffer | *Buffer | No | - | [Buffer](../buffer/)<br> |
