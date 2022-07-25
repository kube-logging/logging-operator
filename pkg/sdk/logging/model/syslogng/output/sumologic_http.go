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

import "github.com/banzaicloud/operator-tools/pkg/secret"

// +name:"Sumo Logic HTTP"
// +weight:"200"
type _hugoSumologicHTTP interface{} //nolint:deadcode,unused

// +docName:"Storing messages in Sumo Logic over http"
//More info at https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/55
type _docSumologicHTTP interface{} //nolint:deadcode,unused

// +name:"Sumo Logic HTTP"
// +url:"https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/55"
// +description:"Storing messages in Sumo Logic over http"
// +status:"Testing"
type _metaSumologicHTTP interface{} //nolint:deadcode,unused

// +kubebuilder:object:generate=true
type SumologicHTTPOutput struct {
	// The name of a directory that contains a set of trusted CA certificates in PEM format.
	// +docLink:"Secret,../secret/" (default: none)
	CaDir *secret.Secret `json:"ca_dir,omitempty"`
	// The name of a file that contains a set of trusted CA certificates in PEM format. The syslog-ng OSE application uses the CA certificates in this file to validate the certificate of the peer.
	// +docLink:"Secret,../secret/" (default: empty)
	CaFile *secret.Secret `json:"ca_file,omitempty"`
	// The Cloud Syslog Cloud Token that you received from the Sumo Logic service while configuring your cloud syslog source. (default: empty)
	Collector string `json:"collector,omitempty"`
	// This option specifies your Sumo Logic deployment.https://help.sumologic.com/APIs/General-API-Information/Sumo-Logic-Endpoints-by-Deployment-and-Firewall-Security  (default: empty)
	Deployment string `json:"deployment,omitempty"`
	// Custom HTTP headers to include in the request, for example, headers("HEADER1: header1", "HEADER2: header2").  (default: empty)
	Headers string `json:"headers,omitempty"`
	// The time to wait in seconds before a dead connection is reestablished. (default: 60)
	TimeReopen int `json:"time_reopen,omitempty"`
	// This option sets various options related to TLS encryption, for example, key/certificate files and trusted CA locations. TLS can be used only with tcp-based transport protocols. For details, see https://www.syslog-ng.com/technical-documents/doc/syslog-ng-open-source-edition/3.37/administration-guide/73#TOPIC-1829193
	Tls string `json:"tls,omitempty"`
	// This option enables putting outgoing messages into the disk buffer of the destination to avoid message loss in case of a system failure on the destination side.  (default: false)
	DiskBuffer *DiskBuffer `json:"disk_buffer,omitempty"`
}
