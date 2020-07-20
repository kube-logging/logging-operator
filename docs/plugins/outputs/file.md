---
title: File
weight: 200
generated_file: true
---

# [File Output](https://docs.fluentd.org/output/file)
## Overview
 This plugin has been designed to output logs or metrics to File.

## Configuration
### FileOutputConfig
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| path | string | Yes | - | The Path of the file. The actual path is path + time + ".log" by default.<br> |
| append | bool | No | - | The flushed chunk is appended to existence file or not. The default is not appended.<br> |
| add_path_suffix | *bool | No | true | Add path suffix(default: true)<br> |
| path_suffix | string | No |  ".log" | The suffix of output result.<br> |
| symlink_path | bool | No |  false | Create symlink to temporary buffered file when buffer_type is file. This is useful for tailing file content to check logs.<br> |
| format | *Format | No | - | [Format](../format/)<br> |
| buffer | *Buffer | No | - | [Buffer](../buffer/)<br> |
 #### Example `File` output configurations
 ```
apiVersion: logging.banzaicloud.io/v1beta1
kind: Output
metadata:
  name: demo-output
spec:
  file:
    path: /tmp/logs/${tag}/%Y/%m/%d.%H.%M
    append: true
    buffer:
      timekey: 1m
      timekey_wait: 10s
      timekey_use_utc: true
 ```

 #### Fluentd Config Result
 ```
  <match **>
	@type file
	@id test_file
	add_path_suffix true
	append true
	path /tmp/logs/${tag}/%Y/%m/%d.%H.%M
	<buffer tag,time>
	  @type file
	  path /buffers/test_file.*.buffer
	  retry_forever true
	  timekey 1m
	  timekey_use_utc true
	  timekey_wait 30s
	</buffer>
  </match>
 ```

---
