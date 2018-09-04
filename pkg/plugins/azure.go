package plugins

// AzureOutput CRD name
const AzureOutput = "azure"

// AzureDefaultValues for Google Cloud Storage output plugin
var AzureDefaultValues = map[string]string{
	"bufferTimeKey":  "3600",
	"bufferTimeWait": "10m",
	"bufferPath":     "/buffers/azure",
	"format":         "json",
}

// AzureTemplate for Google Cloud Storage output plugin
const AzureTemplate = `
<match {{ .pattern }}.**>
  @type azurestorage

  azure_storage_account    {{ .storageAccountName }}
  azure_storage_access_key {{ .storageAccountKey }}
  azure_container          {{ .bucket }}
  azure_storage_type       blob
  store_as                 gzip
  auto_create_container    true
  azure_object_key_format %{path}%{time_slice}_%{index}.%{file_extension}
  path logs/${tag}/%Y/%m/%d/
  time_slice_format        %Y%m%d-%H
  # if you want to use ${tag} or %Y/%m/%d/ like syntax in path / object_key_format,
  # need to specify tag for ${tag} and time for %Y/%m/%d in <buffer> argument.
  <buffer tag,time>
    @type file
    path /buffers/azure
    timekey 1h # 1 hour partition
    timekey_wait 10m
    timekey_use_utc true # use utc
  </buffer>

  <format>
    @type json
  </format>
</match>`
