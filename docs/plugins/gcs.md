# Plugin gcs
## Variables
| Variable name | Default | Applied function |
|---|---|---|
| pattern | - |  |
| project | - |  |
| private_key | - | toJson |
| client_email | - |  |
| bucket | - |  |
| object_key_format | %{path}%{time_slice}_%{index}.%{file_extension} |  |
| path | logs/${tag}/%Y/%m/%d/ |  |
| bufferPath | /buffers/gcs |  |
| timekey | 1h |  |
| timekey_wait | 10m |  |
| timekey_use_utc | true |  |
| format | json |  |
## Plugin template
```
<match {{ .pattern }}.**>
  @type gcs

  project {{ .project }}
  credentialsJson { "private_key": {{ toJson .private_key }}, "client_email": "{{ .client_email }}" }
  bucket {{ .bucket }}
  object_key_format {{ .object_key_format }}
  path  {{ .path }}

  # if you want to use ${tag} or %Y/%m/%d/ like syntax in path / object_key_format,
  # need to specify tag for ${tag} and time for %Y/%m/%d in <buffer> argument.
  <buffer tag,time>
    @type file
    path {{ .bufferPath }}
    timekey {{ .timekey }}
    timekey_wait {{ .timekey_wait }}
    timekey_use_utc {{ .timekey_use_utc }}
  </buffer>

  <format>
    @type {{ .format }}
  </format>
</match>
```