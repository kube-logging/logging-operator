package plugins

// S3Output CRD name
const S3Output = "s3"

// S3DefaultValues for Amazaon S3 output plugin
var S3DefaultValues = map[string]string{
	"bufferTimeKey":  "3600",
	"bufferTimeWait": "10m",
	"bufferPath":     "/buffers/s3",
	"format":         "json",
}

// S3Template for Amazaon S3 output plugin
const S3Template = `
<match {{ .pattern }}.** >
  @type s3

  aws_key_id {{ .aws_key_id }}
  aws_sec_key {{ .aws_sec_key }}
  s3_bucket {{ .s3_bucket }}
  s3_region {{ .s3_region }}

  path logs/${tag}/%Y/%m/%d/
  s3_object_key_format %{path}%{time_slice}_%{index}.%{file_extension}

  # if you want to use ${tag} or %Y/%m/%d/ like syntax in path / s3_object_key_format,
  # need to specify tag for ${tag} and time for %Y/%m/%d in <buffer> argument.
  <buffer tag,time>
    @type file
    path /buffers/s3
    timekey 3600 # 1 hour partition
    timekey_wait 10m
    timekey_use_utc true # use utc
  </buffer>
  <format>
    @type json
  </format>
</match>
`
