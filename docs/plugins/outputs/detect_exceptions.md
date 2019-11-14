# Exception Detector
## Overview
This output plugin consumes a log stream of JSON objects which contain single-line log messages. If a consecutive sequence of log messages form an exception stack trace, they forwarded as a single, combined JSON object. Otherwise, the input log data is forwarded as is.
More info at https://github.com/GoogleCloudPlatform/fluent-plugin-detect-exceptions

 #### Example output configurations
 ```
 spec:
   exceptionDetector:
     remove_tag_prefix: foo
     languages: java, python
     multiline_flush_interval: 0.1
     buffer:
       timekey: 1m
       timekey_wait: 30s
       timekey_use_utc: true
 ```

## Configuration
### Output Config
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| message | string | No |  "" | The field which contains the raw message text in the input JSON data. <br> |
| remove_tag_prefix | string | No |  "" | The prefix to be removed from the input tag when outputting a record. <br> |
| multiline_flush_interval | string | No |  nil | The interval of flushing the buffer for multiline format. <br> |
| languages | []string | No |  [] | Programming languages for which to detect exceptions. <br> |
| max_lines | int | No |  1000 | Maximum number of lines to flush (0 means no limit) <br> |
| max_bytes | int | No |  0 | Maximum number of bytes to flush (0 means no limit) <br> |
| stream | string | No |  "" | Separate log streams by this field in the input JSON data. <br> |
| buffer | *Buffer | No | - | [Buffer](./buffer.md)<br> |
