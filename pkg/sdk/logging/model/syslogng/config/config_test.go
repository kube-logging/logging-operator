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

package config

import (
	"strings"
	"testing"

	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/filter"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/output"
	"github.com/banzaicloud/operator-tools/pkg/secret"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRenderConfigInto(t *testing.T) {
	testCases := map[string]struct {
		input   Input
		wantOut string
		wantErr any
	}{
		"empty input": {
			input: Input{
				SecretLoaderFactory: secretLoaderFactory{},
			},
			wantErr: true,
		},
		"no syslog-ng spec": {
			input: Input{
				Logging: v1beta1.Logging{
					Spec: v1beta1.LoggingSpec{
						SyslogNGSpec: nil,
					},
				},
				SecretLoaderFactory: secretLoaderFactory{},
			},
			wantErr: true,
		},
		"single flow with single output": {
			input: Input{
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
							Name:      "test-syslog-out",
						},
						Spec: v1beta1.SyslogNGOutputSpec{
							Syslog: &output.SyslogNGSyslogOutput{
								Host:      "test.local",
								Transport: "tcp",
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
									Value:   ".kubernetes.labels.app",
								},
							},
							Filters: []v1beta1.SyslogNGFilter{
								{
									Rewrite: &filter.RewriteConfig{
										Set: &filter.SetConfig{
											FieldName: "cluster",
											Value:     "test-cluster",
										},
									},
								},
							},
							LocalOutputRefs: []string{"test-syslog-out"},
						},
					},
				},
				SecretLoaderFactory: secretLoaderFactory{},
			},
			wantOut: untab(`@version: 3.37

source main_input {
	network(transport(tcp) port(2000) flags(no-parse));
};

destination output_default_test-syslog-out {
	syslog("test.local" transport(tcp));
};

log {
	source(main_input)
	parser {
		json-parser();
	};
	filter {
		match("default" value(.kubernetes.namespace_name) type(string));
	};
	filter {
		match("nginx" value(.kubernetes.labels.app));
	};
	rewrite {
		set("test-cluster" value(cluster));
	};
	destination(output_default_test-syslog-out)
};
`),
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			var buf strings.Builder
			err := RenderConfigInto(testCase.input, &buf)
			checkError(t, testCase.wantErr, err)
			require.Equal(t, testCase.wantOut, buf.String())
		})
	}
}

type secretLoaderFactory map[string]secret.SecretLoader

func (f secretLoaderFactory) OutputSecretLoaderForNamespace(ns string) secret.SecretLoader {
	return f[ns]
}

type secretLoader struct {
}

func (l *secretLoader) Load(s *secret.Secret) (string, error) {
	return "", nil // TODO
}

var _ secret.SecretLoader = (*secretLoader)(nil)
