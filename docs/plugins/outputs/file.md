# File output plugin for Fluentd
## Overview
This plugin has been designed to output logs or metrics to File.
More info at https://docs.fluentd.org/output/file

 #### Example output configurations
 ```
 spec:
  file:
    path: /tmp/logs/${tag}/%Y/%m/%d.%H.%M
    buffer:
      timekey: 1m
      timekey_wait: 10s
      timekey_use_utc: true
 ```

## Configuration
### FileOutputConfig
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| path | string | Yes | - | The Path of the file. The actual path is path + time + ".log" by default.<br> |
| append | bool | No | - | The flushed chunk is appended to existence file or not. The default is not appended.<br> |
| buffer | *Buffer | No | - | [Buffer](./buffer.md)<br> |
