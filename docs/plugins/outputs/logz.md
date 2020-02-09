# LogZ output plugin for Fluentd
## Overview
More info at https://github.com/logzio/fluent-plugin-logzio
>Example Deployment: [Save all logs to LogZ](../../../docs/example-logz.md)

 #### Example output configurations
 ```
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
       flush_thread_count: 4
       flush_interval: 3s
       chunk_limit_size: 16m
       queue_limit_length: 4096
 ```

## Configuration
### Logzio
#### LogZ Send your logs to LogZ.io

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| endpoint | *Endpoint | Yes | - | [Shared Credentials](#Shared-Credentials)<br> |
| output_include_time | bool | No | - | Token (LogZ API Token).<br>[Secret](./secret.md"<br>Token             *secret.Secret `json:"token)`<br> |
| output_include_tags | bool | No | - |  |
| http_idle_timeout | int | No | - |  |
| retry_count | int | No | - |  |
| retry_sleep | int | No | - |  |
| gzip | bool | No | - |  |
| buffer | *Buffer | No | - | [Buffer](./buffer.md)<br> |
### Endpoint
#### Endpoint defines connection details for LogZ.io.

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| url | string | No | "https://listener.logz.io |  |
| port | int | No | 8071 |  |
| token | *secret.Secret | No | - | LogZ API Token.<br>[Secret](./secret.md)<br> |
