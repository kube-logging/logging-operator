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

// AzureOutput CRD name
const AzureOutput = "azure"

// AzureDefaultValues for Azure ObjectStore output plugin
var AzureDefaultValues = map[string]string{
	"bufferTimeKey":           "3600",
	"bufferTimeWait":          "10m",
	"bufferPath":              "/buffers/azure",
	"format":                  "json",
	"timekey":                 "1h",
	"timekey_wait":            "10m",
	"timekey_use_utc":         "true",
	"time_slice_format":       "%Y%m%d-%H",
	"azure_object_key_format": "%{path}%{time_slice}_%{index}.%{file_extension}",
	"path":                    "logs/${tag}/%Y/%m/%d/",
}

// AzureTemplate for Azure ObjectStore output plugin
const AzureTemplate = `
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
</match>`
