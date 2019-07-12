# Plugin forward
## Variables
| Variable name | Default | Applied function |
|---|---|---|
| pattern | - |  |
| clientHostname | fluentd.client |  |
| tlsSharedKey |  |  |
| name | target |  |
| host | - |  |
| port | - |  |
| bufferPath | /buffers/forward |  |
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
<match {{ .pattern }}.** >
  @type forward

  {{ if not (eq .tlsSharedKey "") -}}
  transport tls
  tls_version TLSv1_2
  tls_cert_path                /fluentd/tls/caCert
  tls_client_cert_path         /fluentd/tls/clientCert
  tls_client_private_key_path  /fluentd/tls/clientKey
  <security>
    self_hostname           {{ .clientHostname }}
    shared_key              {{ .tlsSharedKey }}
  </security>
  {{ end -}}

  <server>
    name {{ .name }}
    host {{ .host }}
    port {{ .port }}
  </server>

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