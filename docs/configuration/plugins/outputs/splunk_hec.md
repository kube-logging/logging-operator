---
title: Splunk
weight: 200
generated_file: true
---

# Splunk via Hec output plugin for Fluentd
## Overview

For details, see [https://github.com/splunk/fluent-plugin-splunk-hec](https://github.com/splunk/fluent-plugin-splunk-hec).


## Example output configurations

```yaml
spec:
  splunkHec:
    hec_host: splunk.default.svc.cluster.local
    hec_port: 8088
    protocol: http
```


## Configuration
## SplunkHecOutput

SplunkHecOutput sends your logs to Splunk via Hec

### buffer (*Buffer, optional) {#splunkhecoutput-buffer}

[Buffer](../buffer/) 


### ca_file (*secret.Secret, optional) {#splunkhecoutput-ca_file}

The path to a file containing a PEM-format CA certificate. [Secret](../secret/) 


### ca_path (*secret.Secret, optional) {#splunkhecoutput-ca_path}

The path to a directory containing CA certificates in PEM format. [Secret](../secret/) 


### client_cert (*secret.Secret, optional) {#splunkhecoutput-client_cert}

The path to a file containing a PEM-format CA certificate for this client. [Secret](../secret/) 


### client_key (*secret.Secret, optional) {#splunkhecoutput-client_key}

The private key for this client.' [Secret](../secret/) 


### coerce_to_utf8 (*bool, optional) {#splunkhecoutput-coerce_to_utf8}

Indicates whether to allow non-UTF-8 characters in user logs. If set to true, any non-UTF-8 character is replaced by the string specified in non_utf8_replacement_string. If set to false, the Ingest API errors out any non-UTF-8 characters. .

Default: true

### data_type (string, optional) {#splunkhecoutput-data_type}

The type of data that will be sent to Sumo Logic, either event or metric

Default: event

### fields (Fields, optional) {#splunkhecoutput-fields}

In this case, parameters inside `<fields>` are used as indexed fields and removed from the original input events 


### format (*Format, optional) {#splunkhecoutput-format}

[Format](../format/) 


### hec_host (string, required) {#splunkhecoutput-hec_host}

You can specify SplunkHec host by this parameter. 


### hec_port (int, optional) {#splunkhecoutput-hec_port}

The port number for the Hec token or the Hec load balancer.

Default: 8088

### hec_token (*secret.Secret, required) {#splunkhecoutput-hec_token}

Identifier for the Hec token. [Secret](../secret/) 


### host (string, optional) {#splunkhecoutput-host}

The host location for events. Cannot set both host and host_key parameters at the same time. (Default:hostname) 


### host_key (string, optional) {#splunkhecoutput-host_key}

Key for the host location. Cannot set both host and host_key parameters at the same time. 


### idle_timeout (int, optional) {#splunkhecoutput-idle_timeout}

If a connection has not been used for this number of seconds it will automatically be reset upon the next use to avoid attempting to send to a closed connection. nil means no timeout. 


### index (string, optional) {#splunkhecoutput-index}

Identifier for the Splunk index to be used for indexing events. If this parameter is not set, the indexer is chosen by HEC. Cannot set both index and index_key parameters at the same time. 


### index_key (string, optional) {#splunkhecoutput-index_key}

The field name that contains the Splunk index name. Cannot set both index and index_key parameters at the same time. 


### insecure_ssl (*bool, optional) {#splunkhecoutput-insecure_ssl}

Indicates if insecure SSL connection is allowed

Default: false

### keep_keys (bool, optional) {#splunkhecoutput-keep_keys}

By default, all the fields used by the *_key parameters are removed from the original input events. To change this behavior, set this parameter to true. This parameter is set to false by default. When set to true, all fields defined in `index_key`, `host_key`, `source_key`, `sourcetype_key`, `metric_name_key`, and `metric_value_key` are saved in the original event. 


### metric_name_key (string, optional) {#splunkhecoutput-metric_name_key}

Field name that contains the metric name. This parameter only works in conjunction with the metrics_from_event parameter. When this prameter is set, the `metrics_from_event` parameter is automatically set to false.

Default: true

### metric_value_key (string, optional) {#splunkhecoutput-metric_value_key}

Field name that contains the metric value, this parameter is required when `metric_name_key` is configured. 


### metrics_from_event (*bool, optional) {#splunkhecoutput-metrics_from_event}

When data_type is set to "metric", the ingest API will treat every key-value pair in the input event as a metric name-value pair. Set metrics_from_event to false to disable this behavior and use `metric_name_key` and `metric_value_key` to define metrics. (Default:true) 


### non_utf8_replacement_string (string, optional) {#splunkhecoutput-non_utf8_replacement_string}

If coerce_to_utf8 is set to true, any non-UTF-8 character is replaced by the string you specify in this parameter. .

Default: ' '

### open_timeout (int, optional) {#splunkhecoutput-open_timeout}

The amount of time to wait for a connection to be opened. 


### protocol (string, optional) {#splunkhecoutput-protocol}

This is the protocol to use for calling the Hec API. Available values are: http, https.

Default: https

### read_timeout (int, optional) {#splunkhecoutput-read_timeout}

The amount of time allowed between reading two chunks from the socket. 


### ssl_ciphers (string, optional) {#splunkhecoutput-ssl_ciphers}

List of SSL ciphers allowed. 


### slow_flush_log_threshold (string, optional) {#splunkhecoutput-slow_flush_log_threshold}

The threshold for chunk flush performance check. Parameter type is float, not time, default: 20.0 (seconds) If chunk flush takes longer time than this threshold, fluentd logs warning message and increases metric fluentd_output_status_slow_flush_count. 


### source (string, optional) {#splunkhecoutput-source}

The source field for events. If this parameter is not set, the source will be decided by HEC. Cannot set both source and source_key parameters at the same time. 


### source_key (string, optional) {#splunkhecoutput-source_key}

Field name to contain source. Cannot set both source and source_key parameters at the same time. 


### sourcetype (string, optional) {#splunkhecoutput-sourcetype}

The sourcetype field for events. When not set, the sourcetype is decided by HEC. Cannot set both source and `source_key` parameters at the same time. 


### sourcetype_key (string, optional) {#splunkhecoutput-sourcetype_key}

Field name that contains the sourcetype. Cannot set both source and source_key parameters at the same time. 



