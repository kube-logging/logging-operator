/*
 * Copyright © 2019 Banzai Cloud
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package plugins

// S3Output CRD name
const S3Output = "s3"

// S3DefaultValues for Amazaon S3 output plugin
var S3DefaultValues = map[string]string{
	"bufferTimeKey":        "3600",
	"bufferTimeWait":       "10m",
	"bufferPath":           "/buffers/s3",
	"format":               "json",
	"timekey_use_utc":      "true",
	"s3_object_key_format": "%{path}%{time_slice}_%{index}.%{file_extension}",
}

// S3Template for Amazaon S3 output plugin
const S3Template = `
<match {{ .pattern }}.** >
  @type s3

  aws_key_id {{ .aws_key_id }}
  aws_sec_key {{ .aws_sec_key }}
  s3_bucket {{ .s3_bucket }}
  s3_region {{ .s3_region }}
  store_as gzip_command

  path logs/${tag}/%Y/%m/%d/
  s3_object_key_format {{ .s3_object_key_format }}

  # if you want to use ${tag} or %Y/%m/%d/ like syntax in path / s3_object_key_format,
  # need to specify tag for ${tag} and time for %Y/%m/%d in <buffer> argument.
  <buffer tag,time>
    @type file
    path {{ .bufferPath }}
    timekey {{ .bufferTimeKey }}
    timekey_wait {{ .bufferTimeWait }}
    timekey_use_utc {{ .timekey_use_utc }}
  </buffer>
  <format>
    @type {{ .format }}
  </format>
</match>`
