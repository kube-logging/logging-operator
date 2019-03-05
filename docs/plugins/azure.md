# Plugin azure
## Variables
| Variable name | Default | Applied function |
|---|---|---|
| pattern | - |  |
| storageAccountName | - |  |
| storageAccountKey | - |  |
| bucket | - |  |
| azure_object_key_format | %{path}%{time_slice}_%{index}.%{file_extension} |  |
| path | logs/${tag}/%Y/%m/%d/ |  |
| time_slice_format | %Y%m%d-%H |  |
| bufferPath | /buffers/azure |  |
| timekey | 1h |  |
| timekey_wait | 10m |  |
| timekey_use_utc | true |  |
| format | json |  |
## Plugin template
```
<match {{ .pattern }}.**>
  @type azurestorage

  azure_storage_account    {{ .storageAccountName }}
  azure_storage_access_key {{ .storageAccountKey }}
  azure_container          {{ .bucket }}
  azure_storage_type       blob
  store_as                 gzip
  auto_create_container    true
  azure_object_key_format {{ .azure_object_key_format }}
  path {{ .path }}
  time_slice_format {{ .time_slice_format }}  
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