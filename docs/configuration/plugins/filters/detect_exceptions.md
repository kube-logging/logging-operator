---
title: Exception Detector
weight: 200
generated_file: true
---

# Exception Detector
## Overview
 This filter plugin consumes a log stream of JSON objects which contain single-line log messages. If a consecutive sequence of log messages form an exception stack trace, they forwarded as a single, combined JSON object. Otherwise, the input log data is forwarded as is.
 More info at https://github.com/GoogleCloudPlatform/fluent-plugin-detect-exceptions

 > Note: As Tag management is not supported yet, this Plugin is **mutually exclusive** with [Tag normaliser](../tagnormaliser)

 ## Example output configurations
 ```yaml
 filters:
   - detectExceptions:
     languages: java, python
     multiline_flush_interval: 0.1

 ```

## Configuration
## DetectExceptions

### message (string, optional) {#detectexceptions-message}

The field which contains the raw message text in the input JSON data.  

Default:  ""

### remove_tag_prefix (string, optional) {#detectexceptions-remove_tag_prefix}

The prefix to be removed from the input tag when outputting a record.  

Default:  kubernetes

### multiline_flush_interval (string, optional) {#detectexceptions-multiline_flush_interval}

The interval of flushing the buffer for multiline format.  

Default:  nil

### languages ([]string, optional) {#detectexceptions-languages}

Programming languages for which to detect exceptions.  

Default:  []

### max_lines (int, optional) {#detectexceptions-max_lines}

Maximum number of lines to flush (0 means no limit)  

Default:  1000

### max_bytes (int, optional) {#detectexceptions-max_bytes}

Maximum number of bytes to flush (0 means no limit)  

Default:  0

### stream (string, optional) {#detectexceptions-stream}

Separate log streams by this field in the input JSON data.  

Default:  ""

### force_line_breaks (bool, optional) {#detectexceptions-force_line_breaks}

Force line breaks between each lines when comibining exception stacks.  

Default:  false

### match_tag (string, optional) {#detectexceptions-match_tag}

Tag used in match directive.  

Default:  kubernetes.**


 ## Example `Exception Detector` filter configurations
 ```yaml
 apiVersion: logging.banzaicloud.io/v1beta1
 kind: Flow
 metadata:

	name: demo-flow

 spec:

	filters:
	  - detectExceptions:
	      multiline_flush_interval: 0.1
	      languages:
	        - java
	        - python
	selectors: {}
	localOutputRefs:
	  - demo-output

 ```

 #### Fluentd Config Result
 ```yaml
 <match kubernetes.**>

	@type detect_exceptions
	@id test_detect_exceptions
	languages ["java","python"]
	multiline_flush_interval 0.1
	remove_tag_prefix kubernetes

 </match>
 ```

---
