# Fluentd concat filter
## Overview
 Fluentd Filter plugin to concatenate multiline log separated in multiple events.
 More information at https://github.com/fluent-plugins-nursery/fluent-plugin-concat

 #### Example record configurations
 ```
 spec:
  filters:
    - concat:
        key: message
        partial_key: "partial_message"
        separator: ""
 ```

## Configuration
### Concat
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| key | string | Yes | - | The key for part of multiline log.<br> |
| separator | string | No |  "\n" | The separator of lines. <br> |
| n_lines | int | No | - | The number of lines. This is exclusive with multiline_start_regex.<br> |
| multiline_start_regexp | string | No | - | The regexp to match beginning of multiline. This is exclusive with n_lines.<br> |
| multiline_end_regexp | string | No | - | The regexp to match ending of multiline. This is exclusive with n_lines.<br> |
| continuous_line_regexp | string | No | - | The regexp to match continuous lines. This is exclusive with n_lines.<br> |
| stream_identity_key | string | No | - | The key to determine which stream an event belongs to.<br> |
| flush_interval | int | No | - | The number of seconds after which the last received event log will be flushed. If specified 0, wait for next line forever.<br> |
| timeout_label | string | No | - | The label name to handle events caused by timeout.<br> |
| use_first_timestamp | bool | No |  False | Use timestamp of first record when buffer is flushed. <br> |
| partial_key | string | No | - | The field name that is the reference to concatenate records<br> |
| partial_value | string | No | - | The value stored in the field specified by partial_key that represent partial log<br> |
| keep_partial_key | bool | No | False | If true, keep partial_key in concatenated records <br> |
| use_partial_metadata | string | No | - | Use partial metadata to concatenate multiple records<br> |
| keep_partial_metadata | string | No | - | If true, keep partial metadata<br> |
