---
title: LogDNA
weight: 200
generated_file: true
---

# [LogDNA Output](https://github.com/logdna/fluent-plugin-logdna)
## Overview
 This plugin has been designed to output logs to LogDNA.

## Configuration
## LogDNA

Send your logs to LogDNA

### api_key (string, required) {#logdna-api_key}

LogDNA Api key 

Default: -

### hostname (string, required) {#logdna-hostname}

Hostname 

Default: -

### app (string, optional) {#logdna-app}

Application name 

Default: -

### tags (string, optional) {#logdna-tags}

Comma-Separated List of Tags, Optional 

Default: -

### request_timeout (string, optional) {#logdna-request_timeout}

HTTPS POST Request Timeout, Optional. Supports s and ms Suffices  

Default:  30 s

### ingester_domain (string, optional) {#logdna-ingester_domain}

Custom Ingester URL, Optional  

Default:  https://logs.logdna.com

### ingester_endpoint (string, optional) {#logdna-ingester_endpoint}

Custom Ingester Endpoint, Optional  

Default:  /logs/ingest

### buffer (*Buffer, optional) {#logdna-buffer}

[Buffer](../buffer/) 

Default: -

### slow_flush_log_threshold (string, optional) {#logdna-slow_flush_log_threshold}

The threshold for chunk flush performance check. Parameter type is float, not time, default: 20.0 (seconds) If chunk flush takes longer time than this threshold, fluentd logs warning message and increases metric fluentd_output_status_slow_flush_count. 

Default: -


 ## Example `LogDNA` filter configurations
 ```yaml
 apiVersion: logging.banzaicloud.io/v1beta1
 kind: Output
 metadata:

	name: logdna-output-sample

 spec:

	logdna:
	  api_key: xxxxxxxxxxxxxxxxxxxxxxxxxxx
	  hostname: logging-operator
	  app: my-app
	  tags: web,dev
	  ingester_domain https://logs.logdna.com
	  ingester_endpoint /logs/ingest

 ```

 #### Fluentd Config Result
 ```
 <match **>

	@type logdna
	@id test_logdna
	api_key xxxxxxxxxxxxxxxxxxxxxxxxxxy
	app my-app
	hostname logging-operator

 </match>
 ```

---
