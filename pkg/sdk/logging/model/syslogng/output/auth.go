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

// +name:"Authentication config for syslog-ng outputs"
// +weight:"200"
type _hugoAuth interface{} //nolint:deadcode,unused

// +docName:"Authentication config for syslog-ng outputs"
// More info at TODO
type _docAuth interface{} //nolint:deadcode,unused

// +name:"Authentication config for syslog-ng outputs"
// +url:"https://axoflow.com/docs/axosyslog-core/chapter-destinations/syslog-ng-with-loki/#auth"
// +description:"Authentication config for syslog-ng outputs"
// +status:"Testing"
type _metaAuth interface{} //nolint:deadcode,unused

// +kubebuilder:object:generate=true
type Auth struct {
	// Application Layer Transport Security (ALTS) is a simple to use authentication, only available within Google’s infrastructure.
	ALTS *ALTS `json:"alts,omitempty"`
	// Application Default Credentials (ADC).
	ADC *ADC `json:"adc,omitempty"`
	// This is the default method, authentication is disabled (auth(insecure())).
	Insecure *Insecure `json:"insecure,omitempty"`
	// This option sets various options related to TLS encryption, for example, key/certificate files and trusted CA locations. TLS can be used only with tcp-based transport protocols. For details, see [TLS for syslog-ng outputs](../tls/) and the [syslog-ng documentation](https://axoflow.com/docs/axosyslog-core/chapter-encrypted-transport-tls/tlsoptions).
	TLS *TLS `json:"tls,omitempty"`
}

type ADC struct{}

type Insecure struct{}

// +kubebuilder:object:generate=true
type ALTS struct {
	TargetServiceAccounts []string `json:"target-service-accounts,omitempty"`
}
