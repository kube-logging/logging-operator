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

// AlibabaOutput CRD name
const AlibabaOutput = "alibaba"

// AlibabaDefaultValues for Alibaba OSS output plugin
var AlibabaDefaultValues = map[string]string{
	"buffer_chunk_limit": "1m",
	"buffer_path":        "/buffers/ali",
	"time_slice_format":  "%Y%m%d",
	"time_slice_wait":    "10m",
	"oss_object_key_format": "%{time_slice}/%{host}-%{uuid}.%{file_ext}",
}

// AlibabaTemplate for Alibaba OSS output plugin
const AlibabaTemplate = `
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
</match>`
