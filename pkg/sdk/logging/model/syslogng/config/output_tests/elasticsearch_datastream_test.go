// Copyright © 2024 Kube logging authors
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

func TestElasticsearchDatastreamOutput(t *testing.T) {
	config.CheckConfigForOutput(t,
		v1beta1.SyslogNGOutput{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "test-elasticsearch-datastream-out",
			},
			Spec: v1beta1.SyslogNGOutputSpec{
				ElasticsearchDatastream: &output.ElasticsearchDatastreamOutput{
					HTTPOutput: output.HTTPOutput{
						URL:  "https://elastic-endpoint:9200/my-data-stream/_bulk",
						User: "elastic",
						Password: secret.Secret{
							Value: "elastic-password",
						},
					},
				},
			},
		},
		`
destination "output_default_test-elasticsearch-datastream-out" {
	elasticsearch-datastream(url("https://elastic-endpoint:9200/my-data-stream/_bulk") user("elastic") password("elastic-password") persist_name("output_default_test-elasticsearch-datastream-out"));
};
`,
	)
}
