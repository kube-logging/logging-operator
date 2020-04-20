---
title: LogZ
weight: 200
---

# LogZ output plugin for Fluentd
## Overview
More info at https://github.com/logzio/fluent-plugin-logzio

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
| endpoint | *Endpoint | Yes | - | Define LogZ endpoint URL<br> |
| output_include_time | bool | No | - | Should the appender add a timestamp to your logs on their process time (recommended).<br> |
| output_include_tags | bool | No | - | Should the appender add the fluentd tag to the document, called "fluentd_tag"<br> |
| http_idle_timeout | int | No | - | Timeout in seconds that the http persistent connection will stay open without traffic.<br> |
| retry_count | int | No | - | How many times to resend failed bulks.<br> |
| retry_sleep | int | No | - | How long to sleep initially between retries, exponential step-off.<br> |
| gzip | bool | No | - | Should the plugin ship the logs in gzip compression. Default is false.<br> |
| buffer | *Buffer | No | - | [Buffer](../buffer/)<br> |
### Endpoint
#### Endpoint defines connection details for LogZ.io.

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| url | string | No | https://listener.logz.io | LogZ URL.<br> |
| port | int | No | 8071 | Port over which to connect to LogZ URL.<br> |
| token | *secret.Secret | No | - | LogZ API Token.<br>[Secret](../secret/)<br> |
