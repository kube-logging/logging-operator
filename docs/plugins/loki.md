# Plugin loki
## Variables
| Variable name | Default | Applied function |
|---|---|---|
| pattern | - |  |
| url |  |  |
| username |  |  |
| password |  |  |
| extraLabels |  |  |
| flushInterval | 10s |  |
| chunkLimitSize | 1m |  |
| flushAtShutdown | true |  |
## Plugin template
```
<match {{ .pattern }}.**>
  @type kubernetes_loki
  url {{ .url }}
  username {{ .username }}
  password {{ .password }}
  extra_labels {{ .extraLabels }}
  <buffer>
    flush_interval {{ .flushInterval }}
    chunk_limit_size {{ .chunkLimitSize }}
    flush_at_shutdown {{ .flushAtShutdown }}
  </buffer>
</match>
```