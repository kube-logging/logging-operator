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

// SplunkHecOutput CRD name
const SplunkHecOutput = "splunk_hec"

// SplunkHecDefaultValues for Splunk Event Collector output plugin
var SplunkHecDefaultValues = map[string]string{
	"protocol":   "http",
	"host":       "localhost",
	"port":       "9997",
	"sourceType": "_json",
	"index":      "default_index",
}

// SplunkHecTemplate for Splunk Event Collector output plugin
const SplunkHecTemplate = `
 <match {{ .pattern }}.** >
  @type splunk_hec
  protocol {{ .protocol }}
  hec_host {{ .host }}
  hec_port {{ .port }}
  hec_token {{ .token }}
  index {{ .index }}
 </match>`
