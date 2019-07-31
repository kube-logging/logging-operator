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

// SplunkTcpOutput CRD name
const SplunkTcpOutput = "splunk_tcp"

// SplunkTcpDefaultValues for Splunk TCP output plugin
var SplunkTcpDefaultValues = map[string]string{
	"bufferTimeKey":   "3600",
	"bufferTimeWait":  "10m",
	"bufferPath":      "/buffers/splunk",
	"format":          "json",
	"eventKey":        "log",
	"timekey_use_utc": "true",
	"host":            "localhost",
	"port":            "9997",
}

// SplunkTcpTemplate for Splunk TCP output plugin
const SplunkTcpTemplate = `
 <match {{ .pattern }}.** >
   @type splunk_tcp
   host {{ .host }}
   port {{ .port }}
   format {{ .format }}
   event_key {{ .eventKey }}
   <buffer tag,time>
	 @type file
	 path {{ .bufferPath }}
	 timekey {{ .bufferTimeKey }}
	 timekey_wait {{ .bufferTimeWait }}
	 timekey_use_utc {{ .timekey_use_utc }}
   </buffer>
 </match>`
