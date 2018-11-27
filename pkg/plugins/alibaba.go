package plugins

// AlibabaOutput CRD name
const AlibabaOutput = "alibaba"

// AlibabaDefaultValues for Alibaba OSS output plugin
var AlibabaDefaultValues = map[string]string{
	"buffer_chunk_limit": "256m",
	"buffer_path":        "/buffers/ali",
	"time_slice_format":  "%Y%m%d",
	"time_slice_wait":    "10m",
}

// AlibabaTemplate for Alibaba OSS output plugin
const AlibabaTemplate = `
<match {{ .pattern }}.**>
  @type oss
  oss_key_id {{ .aliKeyId }}
  oss_key_secret {{ .aliKeySecret }}
  oss_bucket {{ .bucket }}
  oss_endpoint {{ .aliBucketEndpoint }}
  oss_object_key_format "%{time_slice}/%{host}-%{uuid}.%{file_ext}"

  buffer_path /buffers/ali
  buffer_chunk_limit 1m
  time_slice_format %Y%m%d
  time_slice_wait 10m
</match>`
