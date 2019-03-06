# Plugin alibaba
## Variables
| Variable name | Default | Applied function |
|---|---|---|
| pattern | - |  |
| aliKeyId | - |  |
| aliKeySecret | - |  |
| bucket | - |  |
| aliBucketEndpoint | - |  |
| oss_object_key_format | %{time_slice}/%{host}-%{uuid}.%{file_ext} |  |
| buffer_path | /buffers/ali |  |
| buffer_chunk_limit | 1m |  |
| time_slice_format | %Y%m%d |  |
| time_slice_wait | 10m |  |
## Plugin template
```
<match {{ .pattern }}.**>
  @type oss
  oss_key_id {{ .aliKeyId }}
  oss_key_secret {{ .aliKeySecret }}
  oss_bucket {{ .bucket }}
  oss_endpoint {{ .aliBucketEndpoint }}
  oss_object_key_format  {{ .oss_object_key_format }}

  buffer_path {{ .buffer_path }}
  buffer_chunk_limit {{ .buffer_chunk_limit }}
  time_slice_format {{ .time_slice_format }}
  time_slice_wait {{ .time_slice_wait }}
</match>
```