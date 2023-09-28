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
// +url:"TODO"
// +description:"Authentication config for syslog-ng outputs"
// +status:"Testing"
type _metaAuth interface{} //nolint:deadcode,unused

type Auth struct {
	ALTS     *ALTS     `json:"alts,omitempty"`
	ADC      *ADC      `json:"adc,omitempty"`
	Insecure *Insecure `json:"insecure,omitempty"`
	TLS      *TLS      `json:"tls,omitempty"`
}

type ADC struct{}

type Insecure struct{}

type ALTS struct {
	TargetServiceAccounts []string `json:"target-service-accounts,omitempty"`
}
