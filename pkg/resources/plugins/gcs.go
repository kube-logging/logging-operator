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

// GCSOutput CRD name
const GCSOutput = "gcs"

// GCSDefaultValues for Google Cloud Storage output plugin
var GCSDefaultValues = map[string]string{
	"bufferTimeKey":  "3600",
	"bufferTimeWait": "10m",
	"bufferPath":     "/buffers/gcs",
	"format":         "json",
}

// GCSTemplate for Google Cloud Storage output plugin
const GCSTemplate = `
<match {{ .pattern }}.**>
  @type gcs

  project {{ .project }}
  credentialsJson { "private_key": {{ toJson .private_key }}, "client_email": "{{ .client_email }}" }
  bucket {{ .bucket }}
  object_key_format %{path}%{time_slice}_%{index}.%{file_extension}
  path logs/${tag}/%Y/%m/%d/

  # if you want to use ${tag} or %Y/%m/%d/ like syntax in path / object_key_format,
  # need to specify tag for ${tag} and time for %Y/%m/%d in <buffer> argument.
  <buffer tag,time>
    @type file
    path /buffers/gcs
    timekey 1h # 1 hour partition
    timekey_wait 10m
    timekey_use_utc true # use utc
  </buffer>

  <format>
    @type json
  </format>
</match>`
