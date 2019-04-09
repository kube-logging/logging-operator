# Plugin elasticsearch
## Variables
| Variable name | Default | Applied function |
|---|---|---|
| pattern | - |  |
| logLevel | info |  |
| host | - |  |
| port | - |  |
| scheme | scheme |  |
| sslVerify | true |  |
| logstashFormat | true |  |
| logstashPrefix | logstash |  |
| bufferPath | /buffers/elasticsearch |  |
| timekey | 1h |  |
| timekey_wait | 10m |  |
| timekey_use_utc | true |  |
| flush_thread_count | 2 |  |
| flush_interval | 5s |  |
| retry_forever | true |  |
| retry_max_interval | 30 |  |
| chunkLimit | 2M |  |
| queueLimit | 8 |  |
## Plugin template
```
<match {{ .pattern }}.**>
  @type elasticsearch
  @log_level {{ .logLevel }}
  include_tag_key true
  type_name fluentd
  host {{ .host }}
  port {{ .port }}
  scheme  {{ .scheme }}
  {{- if .sslVerify }}
  ssl_verify {{ .sslVerify }}
  {{- end}}
  logstash_format {{ .logstashFormat }}
  logstash_prefix {{ .logstashPrefix }}
  reconnect_on_error true
  {{- if .user }}
  user {{ .user }}
  {{- end}}
  {{- if .password }}
  password {{ .password }}
  {{- end}}
  <buffer tag, time>
    @type file
    path {{ .bufferPath }}
    timekey {{ .timekey }}
    timekey_wait {{ .timekey_wait }}
    timekey_use_utc {{ .timekey_use_utc }}
    flush_mode interval
    retry_type exponential_backoff
    flush_thread_count {{ .flush_thread_count }}
    flush_interval {{ .flush_interval }}
    retry_forever {{ .retry_forever }}
    retry_max_interval {{ .retry_max_interval }}
    chunk_limit_size {{ .chunkLimit }}
    queue_limit_length {{ .queueLimit }}
    overflow_action block
  </buffer>
</match>
```