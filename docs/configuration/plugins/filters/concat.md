---
title: Concat
weight: 200
generated_file: true
---

# [Concat Filter](https://github.com/fluent-plugins-nursery/fluent-plugin-concat)
## Overview
 Fluentd Filter plugin to concatenate multiline log separated in multiple events.

## Configuration
## Concat

### key (string, optional) {#concat-key}

Specify field name in the record to parse. If you leave empty the Container Runtime default will be used. 

Default: -

### separator (*string, optional) {#concat-separator}

The separator of lines. (default: "\n") 

Default: \"\\n\"

### n_lines (int, optional) {#concat-n_lines}

The number of lines. This is exclusive with multiline_start_regex. 

Default: -

### multiline_start_regexp (string, optional) {#concat-multiline_start_regexp}

The regexp to match beginning of multiline. This is exclusive with n_lines. 

Default: -

### multiline_end_regexp (string, optional) {#concat-multiline_end_regexp}

The regexp to match ending of multiline. This is exclusive with n_lines. 

Default: -

### continuous_line_regexp (string, optional) {#concat-continuous_line_regexp}

The regexp to match continuous lines. This is exclusive with n_lines. 

Default: -

### stream_identity_key (string, optional) {#concat-stream_identity_key}

The key to determine which stream an event belongs to. 

Default: -

### flush_interval (int, optional) {#concat-flush_interval}

The number of seconds after which the last received event log will be flushed. If specified 0, wait for next line forever. 

Default: -

### timeout_label (string, optional) {#concat-timeout_label}

The label name to handle events caused by timeout. 

Default: -

### use_first_timestamp (bool, optional) {#concat-use_first_timestamp}

Use timestamp of first record when buffer is flushed.  

Default:  False

### partial_key (string, optional) {#concat-partial_key}

The field name that is the reference to concatenate records 

Default: -

### partial_value (string, optional) {#concat-partial_value}

The value stored in the field specified by partial_key that represent partial log 

Default: -

### keep_partial_key (bool, optional) {#concat-keep_partial_key}

If true, keep partial_key in concatenated records  

Default: False

### use_partial_metadata (string, optional) {#concat-use_partial_metadata}

Use partial metadata to concatenate multiple records 

Default: -

### keep_partial_metadata (string, optional) {#concat-keep_partial_metadata}

If true, keep partial metadata 

Default: -

### partial_metadata_format (string, optional) {#concat-partial_metadata_format}

Input format of the partial metadata (fluentd or journald docker log driver)( docker-fluentd, docker-journald, docker-journald-lowercase) 

Default: -

### use_partial_cri_logtag (bool, optional) {#concat-use_partial_cri_logtag}

Use cri log tag to concatenate multiple records 

Default: -

### partial_cri_logtag_key (string, optional) {#concat-partial_cri_logtag_key}

The key name that is referred to concatenate records on cri log 

Default: -

### partial_cri_stream_key (string, optional) {#concat-partial_cri_stream_key}

The key name that is referred to detect stream name on cri log 

Default: -


 ## Example `Concat` filter configurations
 ```yaml
 apiVersion: logging.banzaicloud.io/v1beta1
 kind: Flow
 metadata:

	name: demo-flow

 spec:

	filters:
	  - concat:
	      partial_key: "partial_message"
	      separator: ""
	      n_lines: 10
	selectors: {}
	localOutputRefs:
	  - demo-output

 ```

 #### Fluentd Config Result
 ```yaml
 <filter **>

	@type concat
	@id test_concat
	key message
	n_lines 10
	partial_key partial_message

 </filter>
 ```

---
