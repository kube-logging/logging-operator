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
	"strings"
	"testing"

	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/syslogng/config"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/syslogng/filter"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/syslogng/output"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestHTTP(t *testing.T) {
	testCases := map[string]struct {
		input   config.Input
		wantOut string
		wantErr any
	}{
		"http output test": {
			input: config.Input{
				SourcePort: 601,
				Logging: v1beta1.Logging{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "config-test",
						Name:      "test",
					},
					Spec: v1beta1.LoggingSpec{
						SyslogNGSpec: &v1beta1.SyslogNGSpec{},
					},
				},
				Outputs: []v1beta1.SyslogNGOutput{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "test-http-out",
						},
						Spec: v1beta1.SyslogNGOutputSpec{
							HTTP: &output.HTTPOutput{
								URL:     "test.local",
								Headers: []string{"a:b", "c:d"},
								Batch: output.Batch{
									BatchLines: 2000,
								},
								Workers: 3,
							},
						},
					},
				},
				Flows: []v1beta1.SyslogNGFlow{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "test-flow",
						},
						Spec: v1beta1.SyslogNGFlowSpec{
							Match: &v1beta1.SyslogNGMatch{
								Regexp: &filter.RegexpMatchExpr{
									Pattern: "nginx",
									Value:   "kubernetes.labels.app",
								},
							},
							Filters: []v1beta1.SyslogNGFilter{
								{
									Rewrite: []filter.RewriteConfig{
										{
											Set: &filter.SetConfig{
												FieldName: "cluster",
												Value:     "test-cluster",
											},
										},
									},
								},
							},
							LocalOutputRefs: []string{"test-http-out"},
						},
					},
				},
				SecretLoaderFactory: &config.TestSecretLoaderFactory{},
			},
			wantOut: config.Untab(`@version: 3.37

@include "scl.conf"

source "main_input" {
    channel {
        source {
            network(flags("no-parse") port(601) transport("tcp"));
        };
        parser {
            json-parser(prefix("json."));
        };
    };
};

destination "output_default_test-http-out" {
	http(url("test.local") headers("a:b" "c:d") batch-lines(2000) workers(3) persist_name("output_default_test-http-out"));
};

filter "flow_default_test-flow_ns_filter" {
	match("default" value("json.kubernetes.namespace_name") type("string"));
};
filter "flow_default_test-flow_match" {
	match("nginx" value("kubernetes.labels.app"));
};
rewrite "flow_default_test-flow_filters_0" {
	set("test-cluster" value("cluster"));
};
log {
	source("main_input");
	filter("flow_default_test-flow_ns_filter");
	filter("flow_default_test-flow_match");
	rewrite("flow_default_test-flow_filters_0");
	destination("output_default_test-http-out");
};
`),
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			var buf strings.Builder
			err := config.RenderConfigInto(testCase.input, &buf)
			config.CheckError(t, testCase.wantErr, err)
			require.Equal(t, testCase.wantOut, buf.String())
		})
	}
}
