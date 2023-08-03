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
 kubectl create secret generic sumo-output --from-literal "endpoint=$URL"
 ```

 # Example ClusterOutput

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

export URL='https://endpoint1.collection.eu.sumologic.com/receiver/v1/http/.......'

## Configuration
## Output Config

### data_type (string, optional) {#output config-data_type}

The type of data that will be sent to Sumo Logic, either logs or metrics  

Default:  logs

### endpoint (*secret.Secret, required) {#output config-endpoint}

SumoLogic HTTP Collector URL 

Default: -

### verify_ssl (bool, optional) {#output config-verify_ssl}

Verify ssl certificate.  

Default:  true

### metric_data_format (string, optional) {#output config-metric_data_format}

The format of metrics you will be sending, either graphite or carbon2 or prometheus  

Default:  graphite

### log_format (string, optional) {#output config-log_format}

Format to post logs into Sumo.  

Default:  json

### log_key (string, optional) {#output config-log_key}

Used to specify the key when merging json or sending logs in text format  

Default:  message

### source_category (string, optional) {#output config-source_category}

Set _sourceCategory metadata field within SumoLogic  

Default:  nil

### source_name (string, required) {#output config-source_name}

Set _sourceName metadata field within SumoLogic - overrides source_name_key (default is nil) 

Default: -

### source_name_key (string, optional) {#output config-source_name_key}

Set as source::path_key's value so that the source_name can be extracted from Fluentd's buffer  

Default:  source_name

### source_host (string, optional) {#output config-source_host}

Set _sourceHost metadata field within SumoLogic  

Default:  nil

### open_timeout (int, optional) {#output config-open_timeout}

Set timeout seconds to wait until connection is opened.  

Default:  60

### add_timestamp (bool, optional) {#output config-add_timestamp}

Add timestamp (or timestamp_key) field to logs before sending to sumologic  

Default:  true

### timestamp_key (string, optional) {#output config-timestamp_key}

Field name when add_timestamp is on  

Default:  timestamp

### proxy_uri (string, optional) {#output config-proxy_uri}

Add the uri of the proxy environment if present. 

Default: -

### disable_cookies (bool, optional) {#output config-disable_cookies}

Option to disable cookies on the HTTP Client.  

Default:  false

### delimiter (string, optional) {#output config-delimiter}

Delimiter  

Default:  .

### custom_fields ([]string, optional) {#output config-custom_fields}

Comma-separated key=value list of fields to apply to every log. [more information](https://help.sumologic.com/Manage/Fields#http-source-fields) 

Default: -

### sumo_client (string, optional) {#output config-sumo_client}

Name of sumo client which is send as X-Sumo-Client header  

Default:  fluentd-output

### compress (*bool, optional) {#output config-compress}

Compress payload  

Default:  false

### compress_encoding (string, optional) {#output config-compress_encoding}

Encoding method of compression (either gzip or deflate)  

Default:  gzip

### custom_dimensions (string, optional) {#output config-custom_dimensions}

Dimensions string (eg "cluster=payment, service=credit_card") which is going to be added to every metric record. 

Default: -

### buffer (*Buffer, optional) {#output config-buffer}

[Buffer](../buffer/) 

Default: -

### slow_flush_log_threshold (string, optional) {#output config-slow_flush_log_threshold}

The threshold for chunk flush performance check. Parameter type is float, not time, default: 20.0 (seconds) If chunk flush takes longer time than this threshold, fluentd logs warning message and increases metric fluentd_output_status_slow_flush_count. 

Default: -


