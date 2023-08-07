---
title: LogZ
weight: 200
generated_file: true
---

# LogZ output plugin for Fluentd
## Overview
 More info at https://github.com/tarokkk/fluent-plugin-logzio

 ## Example output configurations
 ```yaml
 spec:

	logz:
	  endpoint:
	    url: https://listener.logz.io
	    port: 8071
	    token:
	      valueFrom:
	       secretKeyRef:
	   	  name: logz-token
	         key: token
	  output_include_tags: true
	  output_include_time: true
	  buffer:
	    type: file
	    flush_mode: interval
	    flush_thread_count: 4
	    flush_interval: 5s
	    chunk_limit_size: 16m
	    queue_limit_length: 4096

 ```

## Configuration
## Logzio

LogZ Send your logs to LogZ.io

### endpoint (*Endpoint, required) {#logzio-endpoint}

Define LogZ endpoint URL 

Default: -

### output_include_time (bool, optional) {#logzio-output_include_time}

Should the appender add a timestamp to your logs on their process time (recommended). 

Default: -

### output_include_tags (bool, optional) {#logzio-output_include_tags}

Should the appender add the fluentd tag to the document, called "fluentd_tag" 

Default: -

### http_idle_timeout (int, optional) {#logzio-http_idle_timeout}

Timeout in seconds that the http persistent connection will stay open without traffic. 

Default: -

### retry_count (int, optional) {#logzio-retry_count}

How many times to resend failed bulks. 

Default: -

### retry_sleep (int, optional) {#logzio-retry_sleep}

How long to sleep initially between retries, exponential step-off. 

Default: -

### bulk_limit (int, optional) {#logzio-bulk_limit}

Limit to the size of the Logz.io upload bulk. Defaults to 1000000 bytes leaving about 24kB for overhead. 

Default: -

### bulk_limit_warning_limit (int, optional) {#logzio-bulk_limit_warning_limit}

Limit to the size of the Logz.io warning message when a record exceeds bulk_limit to prevent a recursion when Fluent warnings are sent to the Logz.io output. 

Default: -

### gzip (bool, optional) {#logzio-gzip}

Should the plugin ship the logs in gzip compression. Default is false. 

Default: -

### buffer (*Buffer, optional) {#logzio-buffer}

[Buffer](../buffer/) 

Default: -

### slow_flush_log_threshold (string, optional) {#logzio-slow_flush_log_threshold}

The threshold for chunk flush performance check. Parameter type is float, not time, default: 20.0 (seconds) If chunk flush takes longer time than this threshold, fluentd logs warning message and increases metric fluentd_output_status_slow_flush_count. 

Default: -


## Endpoint

Endpoint defines connection details for LogZ.io.

### url (string, optional) {#endpoint-url}

LogZ URL. 

Default: https://listener.logz.io

### port (int, optional) {#endpoint-port}

Port over which to connect to LogZ URL. 

Default: 8071

### token (*secret.Secret, optional) {#endpoint-token}

LogZ API Token. [Secret](../secret/) 

Default: -


