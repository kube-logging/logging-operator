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

func TestLogglyOutput(t *testing.T) {
	config.CheckConfigForOutput(t,
		v1beta1.SyslogNGOutput{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "test-loggly-out",
			},
			Spec: v1beta1.SyslogNGOutputSpec{
				Loggly: &output.Loggly{
					SyslogOutput: output.SyslogOutput{Host: "localhost"},
					Token: &secret.Secret{
						Value: "asd",
					},
					Tag: "test-tag",
				},
			},
		},
		`
destination "output_default_test-loggly-out" {
	loggly("localhost" tag("test-tag") token("asd") persist_name("output_default_test-loggly-out"));
};
`,
	)
}
