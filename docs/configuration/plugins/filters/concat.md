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

### continuous_line_regexp (string, optional) {#concat-continuous_line_regexp}

The regexp to match continuous lines. This is exclusive with n_lines. 


### flush_interval (int, optional) {#concat-flush_interval}

The number of seconds after which the last received event log is flushed. If set to 0, flushing is disabled (wait for next line forever). 


### keep_partial_key (bool, optional) {#concat-keep_partial_key}

If true, keep partial_key in concatenated records

Default: False

### keep_partial_metadata (string, optional) {#concat-keep_partial_metadata}

If true, keep partial metadata 


### key (string, optional) {#concat-key}

Specify field name in the record to parse. If you leave empty the Container Runtime default will be used. 


### multiline_end_regexp (string, optional) {#concat-multiline_end_regexp}

The regexp to match ending of multiline. This is exclusive with n_lines. 


### multiline_start_regexp (string, optional) {#concat-multiline_start_regexp}

The regexp to match beginning of multiline. This is exclusive with n_lines. 


### n_lines (int, optional) {#concat-n_lines}

The number of lines. This is exclusive with multiline_start_regex. 


### partial_cri_logtag_key (string, optional) {#concat-partial_cri_logtag_key}

The key name that is referred to concatenate records on cri log 


### partial_cri_stream_key (string, optional) {#concat-partial_cri_stream_key}

The key name that is referred to detect stream name on cri log 


### partial_key (string, optional) {#concat-partial_key}

The field name that is the reference to concatenate records 


### partial_metadata_format (string, optional) {#concat-partial_metadata_format}

Input format of the partial metadata (fluentd or journald docker log driver)( docker-fluentd, docker-journald, docker-journald-lowercase) 


### partial_value (string, optional) {#concat-partial_value}

The value stored in the field specified by partial_key that represent partial log 


### separator (*string, optional) {#concat-separator}

The separator of lines. (default: "\n") 

Default: \"\\n\"

### stream_identity_key (string, optional) {#concat-stream_identity_key}

The key to determine which stream an event belongs to. 


### timeout_label (string, optional) {#concat-timeout_label}

The label name to handle events caused by timeout. 


### use_first_timestamp (bool, optional) {#concat-use_first_timestamp}

Use timestamp of first record when buffer is flushed.

Default: False

### use_partial_cri_logtag (bool, optional) {#concat-use_partial_cri_logtag}

Use cri log tag to concatenate multiple records 


### use_partial_metadata (string, optional) {#concat-use_partial_metadata}

Use partial metadata to concatenate multiple records 





## Example `Concat` filter configurations

{{< highlight yaml >}}
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
{{</ highlight >}}

Fluentd config result:

{{< highlight xml >}}
<filter **>
  @type concat
  @id test_concat
  key message
  n_lines 10
  partial_key partial_message
</filter>
{{</ highlight >}}


---
