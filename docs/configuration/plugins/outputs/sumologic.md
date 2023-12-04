---
title: SumoLogic
weight: 200
generated_file: true
---

# SumoLogic output plugin for Fluentd
## Overview

This plugin has been designed to output logs or metrics to SumoLogic via a HTTP collector endpoint
For details, see [https://github.com/SumoLogic/fluentd-output-sumologic](https://github.com/SumoLogic/fluentd-output-sumologic).

## Example secret for HTTP input URL:

```
export URL='https://endpoint1.collection.eu.sumologic.com/receiver/v1/http/'
kubectl create secret generic sumo-output --from-literal "endpoint=$URL"
```

## Example ClusterOutput

```yaml
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
## Output Config

### add_timestamp (bool, optional) {#output config-add_timestamp}

Add timestamp (or timestamp_key) field to logs before sending to SumoLogic

Default: true

### buffer (*Buffer, optional) {#output config-buffer}

[Buffer](../buffer/) 


### compress (*bool, optional) {#output config-compress}

Compress payload

Default: false

### compress_encoding (string, optional) {#output config-compress_encoding}

Encoding method of compression (either gzip or deflate)

Default: gzip

### custom_dimensions (string, optional) {#output config-custom_dimensions}

Dimensions string (eg "cluster=payment, service=credit_card") which is going to be added to every metric record. 


### custom_fields ([]string, optional) {#output config-custom_fields}

Comma-separated key=value list of fields to apply to every log. [More information](https://help.sumologic.com/Manage/Fields#http-source-fields) 


### data_type (string, optional) {#output config-data_type}

The type of data that will be sent to Sumo Logic, either logs or metrics

Default: logs

### delimiter (string, optional) {#output config-delimiter}

Delimiter

Default: .

### disable_cookies (bool, optional) {#output config-disable_cookies}

Option to disable cookies on the HTTP Client.

Default: false

### endpoint (*secret.Secret, required) {#output config-endpoint}

SumoLogic HTTP Collector URL 


### log_format (string, optional) {#output config-log_format}

Format to post logs into Sumo.

Default: json

### log_key (string, optional) {#output config-log_key}

Used to specify the key when merging json or sending logs in text format

Default: message

### metric_data_format (string, optional) {#output config-metric_data_format}

The format of metrics you will be sending, either graphite or carbon2 or prometheus

Default: graphite

### open_timeout (int, optional) {#output config-open_timeout}

Set timeout seconds to wait until connection is opened.

Default: 60

### proxy_uri (string, optional) {#output config-proxy_uri}

Add the uri of the proxy environment if present. 


### slow_flush_log_threshold (string, optional) {#output config-slow_flush_log_threshold}

The threshold for chunk flush performance check. Parameter type is float, not time, default: 20.0 (seconds) If chunk flush takes longer time than this threshold, fluentd logs warning message and increases metric fluentd_output_status_slow_flush_count. 


### source_category (string, optional) {#output config-source_category}

Set _sourceCategory metadata field within SumoLogic

Default: nil

### source_host (string, optional) {#output config-source_host}

Set _sourceHost metadata field within SumoLogic

Default: nil

### source_name (string, required) {#output config-source_name}

Set _sourceName metadata field within SumoLogic - overrides source_name_key (default is nil) 


### source_name_key (string, optional) {#output config-source_name_key}

Set as source::path_key's value so that the source_name can be extracted from Fluentd's buffer

Default: source_name

### sumo_client (string, optional) {#output config-sumo_client}

Name of sumo client which is send as X-Sumo-Client header

Default: fluentd-output

### timestamp_key (string, optional) {#output config-timestamp_key}

Field name when add_timestamp is on

Default: timestamp

### verify_ssl (bool, optional) {#output config-verify_ssl}

Verify ssl certificate.

Default: true


