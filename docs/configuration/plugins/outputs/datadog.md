---
title: Datadog
weight: 200
generated_file: true
---

# Datadog output plugin for Fluentd
## Overview

It mainly contains a proper JSON formatter and a socket handler that streams logs directly to Datadog - so no need to use a log shipper if you don't want to.
For details, see [https://github.com/DataDog/fluent-plugin-datadog](https://github.com/DataDog/fluent-plugin-datadog).

## Example
```yaml
spec:
  datadog:
    api_key:
      value: '<YOUR_API_KEY>' # For referencing a secret, see https://kube-logging.dev/docs/configuration/plugins/outputs/secret/
    dd_source: '<INTEGRATION_NAME>'
    dd_tags: '<KEY1:VALUE1>,<KEY2:VALUE2>'
    dd_sourcecategory: '<YOUR_SOURCE_CATEGORY>'
```


## Configuration
## Output Config

### api_key (*secret.Secret, required) {#output config-api_key}

This parameter is required in order to authenticate your fluent agent.  +docLink:"Secret,../secret/"

Default: nil

### buffer (*Buffer, optional) {#output config-buffer}

[Buffer](../buffer/) 


### compression_level (string, optional) {#output config-compression_level}

Set the log compression level for HTTP (1 to 9, 9 being the best ratio)

Default: "6"

### dd_hostname (string, optional) {#output config-dd_hostname}

Used by Datadog to identify the host submitting the logs.

Default: "hostname -f"

### dd_source (string, optional) {#output config-dd_source}

This tells Datadog what integration it is

Default: nil

### dd_sourcecategory (string, optional) {#output config-dd_sourcecategory}

Multiple value attribute. Can be used to refine the source attribute

Default: nil

### dd_tags (string, optional) {#output config-dd_tags}

Custom tags with the following format "key1:value1, key2:value2"

Default: nil

### host (string, optional) {#output config-host}

Proxy endpoint when logs are not directly forwarded to Datadog

Default: "http-intake.logs.datadoghq.com"

### include_tag_key (bool, optional) {#output config-include_tag_key}

Automatically include the Fluentd tag in the record.

Default: false

### max_backoff (string, optional) {#output config-max_backoff}

The maximum time waited between each retry in seconds

Default: "30"

### max_retries (string, optional) {#output config-max_retries}

The number of retries before the output plugin stops. Set to -1 for unlimited retries

Default: "-1"

### no_ssl_validation (bool, optional) {#output config-no_ssl_validation}

Disable SSL validation (useful for proxy forwarding)

Default: false

### port (string, optional) {#output config-port}

Proxy port when logs are not directly forwarded to Datadog and ssl is not used

Default: "80"

### service (string, optional) {#output config-service}

Used by Datadog to correlate between logs, traces and metrics.

Default: nil

### slow_flush_log_threshold (string, optional) {#output config-slow_flush_log_threshold}

The threshold for chunk flush performance check. Parameter type is float, not time, default: 20.0 (seconds) If chunk flush takes longer time than this threshold, fluentd logs warning message and increases metric fluentd_output_status_slow_flush_count. 


### ssl_port (string, optional) {#output config-ssl_port}

Port used to send logs over a SSL encrypted connection to Datadog. If use_http is disabled, use 10516 for the US region and 443 for the EU region.

Default: "443"

### tag_key (string, optional) {#output config-tag_key}

Where to store the Fluentd tag.

Default: "tag"

### timestamp_key (string, optional) {#output config-timestamp_key}

Name of the attribute which will contain timestamp of the log event. If nil, timestamp attribute is not added.

Default: "@timestamp"

### use_compression (bool, optional) {#output config-use_compression}

Enable log compression for HTTP

Default: true

### use_http (bool, optional) {#output config-use_http}

Enable HTTP forwarding. If you disable it, make sure to change the port to 10514 or ssl_port to 10516

Default: true

### use_json (bool, optional) {#output config-use_json}

Event format, if true, the event is sent in json format. Othwerwise, in plain text.

Default: true

### use_ssl (bool, optional) {#output config-use_ssl}

If true, the agent initializes a secure connection to Datadog. In clear TCP otherwise.

Default: true


