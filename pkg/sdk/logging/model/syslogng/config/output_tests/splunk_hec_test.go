// Copyright © 2019 Banzai Cloud
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

package test

import (
	"testing"

	"github.com/cisco-open/operator-tools/pkg/secret"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/config"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/output"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSplunkHECOutput(t *testing.T) {
	config.CheckConfigForOutput(t,
		v1beta1.SyslogNGOutput{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "test-splunk-hec-out",
			},
			Spec: v1beta1.SyslogNGOutputSpec{
				SplunkHEC: &output.SplunkHECOutput{
					HTTPOutput: output.HTTPOutput{
						URL: "https://localhost:8088",
					},
					Token: secret.Secret{Value: "secret_value"},
				},
			},
		},
		`
destination "output_default_test-splunk-hec-out" {
	splunk_hec_event(url("https://localhost:8088") persist_name("output_default_test-splunk-hec-out") token("secret_value"));
};
`,
	)
}

func TestSplunkHECOutputAdvancedUsage(t *testing.T) {
	config.CheckConfigForOutput(t,
		v1beta1.SyslogNGOutput{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "test-splunk-hec-out",
			},
			Spec: v1beta1.SyslogNGOutputSpec{
				SplunkHEC: &output.SplunkHECOutput{
					HTTPOutput: output.HTTPOutput{
						URL: "https://localhost:8088",
						Batch: output.Batch{BatchLines: 5000,
							BatchBytes:   4096 * 1024,
							BatchTimeout: 300,
						},
						Workers: 8,
						Timeout: 10,
					},
					Token:        secret.Secret{Value: "secret_value"},
					Host:         `${HOST}`,
					ContentType:  "application/json",
					ExtraHeaders: []string{"a:b", "c:d"},
					Time:         `${S_UNIXTIME}.${S_MSEC}`,
				},
			},
		},
		`
destination "output_default_test-splunk-hec-out" {
	splunk_hec_event(url("https://localhost:8088") batch-lines(5000) batch-bytes(4194304) batch-timeout(300) workers(8) persist_name("output_default_test-splunk-hec-out") timeout(10) token("secret_value") host("${HOST}") time("${S_UNIXTIME}.${S_MSEC}") extra_headers("a:b" "c:d") content_type("application/json"));
};
`,
	)
}
