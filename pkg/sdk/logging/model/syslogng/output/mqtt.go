// Copyright © 2022 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package output

type MQTT struct {
	// Address of the destination host
	Address string `json:"address,omitempty"`
	// Event tag [more information](https://documentation.solarwinds.com/en/success_center/loggly/content/admin/tags.htm)
	Topic         string `json:"topic,omitempty"`
	FallbackTopic string `json:"fallback-topic,omitempty"`
	// Specifies a template defining the logformat to be used in the destination. [more information](https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/56#kanchor905) (default: 0)
	Template string `json:"template,omitempty"`
	QOS      int    `json:"qos,omitempty"`
}
