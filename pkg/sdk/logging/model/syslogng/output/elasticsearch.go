// Copyright © 2023 Kube logging authors
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

import "fmt"

// +name:"Elasticsearch"
// +weight:"200"
type _hugoElasticsearch interface{} //nolint:deadcode,unused

// +docName:"Sending messages over Elasticsearch"
// More info at https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/destination-forward-send-and-store-log-messages/elasticsearch-http-sending-messages-to-elasticsearch-http-bulk-api/
type _docSElasticsearch interface{} //nolint:deadcode,unused

// +name:"Elasticsearch"
// +url:"https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/destination-forward-send-and-store-log-messages/elasticsearch-http-sending-messages-to-elasticsearch-http-bulk-api/"
// +description:"Sending messages over Elasticsearch"
// +status:"Testing"
type _metaElasticsearch interface{} //nolint:deadcode,unused

// +kubebuilder:object:generate=true
type ElasticsearchOutput struct {
	HTTPOutput `json:",inline"`
	// Name of the data stream, index, or index alias to perform the action on.
	Index string `json:"index,omitempty"`
	// The document type associated with the operation. Elasticsearch indices now support a single document type: _doc
	Type *string `json:"type,omitempty"`
	// The document ID. If no ID is specified, a document ID is automatically generated.
	CustomID       string `json:"custom_id,omitempty"`
	LogstashPrefix string `json:"logstash_prefix,omitempty" syslog-ng:"ignore"`
}

func (o *ElasticsearchOutput) BeforeRender() {
	if o.LogstashPrefix != "" {
		o.Index = fmt.Sprintf("%s-${YEAR}.${MONTH}.${DAY}", o.LogstashPrefix)
	}
}
