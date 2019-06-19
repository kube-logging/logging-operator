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

// LokiOutput CRD name
const LokiOutput = "loki"

// LokiDefaultValues for Loki output plugin
var LokiDefaultValues = map[string]string{
	"url":             "",
	"username":        "",
	"password":        "",
	"extraLabels":     "",
	"flushInterval":   "10s",
	"chunkLimitSize":  "1m",
	"flushAtShutdown": "true",
}

// LokiTemplate for Loki output plugin
const LokiTemplate = `
<match {{ .pattern }}.**>
  @type kubernetes_loki
  url {{ .url }}
  username {{ .username }}
  password {{ .password }}
  extra_labels {{ .extraLabels }}
  <buffer>
    flush_interval {{ .flushInterval }}
    chunk_limit_size {{ .chunkLimitSize }}
    flush_at_shutdown {{ .flushAtShutdown }}
  </buffer>
</match>`
