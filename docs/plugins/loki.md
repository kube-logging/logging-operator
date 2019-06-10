# Plugin loki
## Variables
| Variable name | Default | Applied function |
|---|---|---|
| pattern | - |  |
| url | - | loki url |
| username | - |  |
| password | - |  |
| extra_labels | - |  |
| flush_interval | 10s |  |
| flush_at_shutdown | true |  |
| buffer_chunk_limit | 1m |  |
## Plugin template
```
<match {{ .pattern }}.**>
  @type loki
  url {{ .url }}
  username {{ .username }}
  password {{ .password }}
  extra_labels {{ .extraLabels }}
  flush_interval {{ .flushInterval }}
  flush_at_shutdown true
  buffer_chunk_limit  {{ .bufferChunkLimit }}
</match>
```