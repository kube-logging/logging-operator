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

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/config"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/output"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestElasticsearchOutput(t *testing.T) {
	var emptyString string
	config.CheckConfigForOutput(t,
		v1beta1.SyslogNGOutput{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "test-elasticsearch-out",
			},
			Spec: v1beta1.SyslogNGOutputSpec{
				Elasticsearch: &output.ElasticsearchOutput{
					HTTPOutput: output.HTTPOutput{
						URL:     "test.local",
						Headers: []string{"a:b", "c:d"},
						Batch: output.Batch{
							BatchLines: 2000,
						},
						Workers: 3,
					},
					LogstashPrefix: "xxx",
					Type:           &emptyString,
				},
			},
		},
		`
destination "output_default_test-elasticsearch-out" {
	elasticsearch-http(url("test.local") headers("a:b" "c:d") batch-lines(2000) workers(3) persist_name("output_default_test-elasticsearch-out") index("xxx-${YEAR}.${MONTH}.${DAY}") type(""));
};
`,
	)
}

func TestElasticsearchOutputWithCustomLogstashFormat(t *testing.T) {
	var emptyString string
	config.CheckConfigForOutput(t,
		v1beta1.SyslogNGOutput{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "test-elasticsearch-out",
			},
			Spec: v1beta1.SyslogNGOutputSpec{
				Elasticsearch: &output.ElasticsearchOutput{
					HTTPOutput: output.HTTPOutput{
						URL:     "test.local",
						Headers: []string{"a:b", "c:d"},
						Batch: output.Batch{
							BatchLines: 2000,
						},
						Workers: 3,
					},
					LogstashPrefix:          "my-prefix",
					LogstashPrefixSeparator: "+",
					LogStashSuffix:          "my-suffix",
					Type:                    &emptyString,
				},
			},
		},
		`
destination "output_default_test-elasticsearch-out" {
    elasticsearch-http(url("test.local") headers("a:b" "c:d") batch-lines(2000) workers(3) persist_name("output_default_test-elasticsearch-out") index("my-prefix+my-suffix") type(""));
};
`,
	)
}
