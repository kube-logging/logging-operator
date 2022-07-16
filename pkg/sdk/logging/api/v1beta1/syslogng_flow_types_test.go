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

package v1beta1

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/filter"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/render/syslogng"
)

func TestSyslogNGFlow_RenderAsSyslogNGConfig(t *testing.T) {
	ctx := syslogng.Context{
		Indent: "    ",
	}
	tests := map[string]struct {
		flow    SyslogNGFlow
		ctx     syslogng.Context
		wantOut string
		wantErr bool
	}{
		"empty flow": {
			flow: SyslogNGFlow{},
			ctx:  ctx,
			wantOut: UseSpaces(`log flow__ {
	source(the_input);
	filter { (not match("" value(".kubernetes.namespace_name") type("string"))) }; # filter messages from other namespaces
};
`),
		},
		"simple flow": {
			flow: SyslogNGFlow{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "apache",
				},
				Spec: SyslogNGFlowSpec{
					Match: (*SyslogNGMatch)(&filter.MatchConfig{
						Regexp: &filter.RegexpMatchExpr{
							Pattern: "apache",
							Value:   ".kubernetes.labels.app.kubernetes.io/name",
							Type:    "string",
						},
					}),
					Filters: []SyslogNGFilter{
						{
							Rewrite: &filter.RewriteConfig{
								Set: &filter.SetConfig{
									FieldName: ".kubernetes.labels.app",
									Value:     "apache",
								},
							},
						},
					},
					LocalOutputRefs: []string{
						"es",
						"os",
					},
					GlobalOutputRefs: []string{
						"s3-audit",
					},
				},
			},
			ctx: ctx.WithControlNamespace("logging"),
			wantOut: UseSpaces(`log flow_default_apache {
	source(the_input);
	filter { (not match("default" value(".kubernetes.namespace_name") type("string"))) }; # filter messages from other namespaces
	filter { match("apache" value(".kubernetes.labels.app.kubernetes.io/name") type("string")) }; # flow match
	rewrite { set("apache" value(".kubernetes.labels.app")) };
	destination(output_default_es);
	destination(output_default_os);
	destination(clusteroutput_logging_s3-audit);
};
`),
		},
	}
	for name, testCase := range tests {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			b := strings.Builder{}
			testCase.ctx.Out = &b
			err := testCase.flow.RenderAsSyslogNGConfig(testCase.ctx)
			if (err != nil) != testCase.wantErr {
				t.Errorf("SyslogNGFlow.RenderAsSyslogNGConfig() error = %v, wantErr %v", err, testCase.wantErr)
			}
			require.Equal(t, testCase.wantOut, b.String())
		})
	}
}

func withControlNamespace(ctx syslogng.Context, ns string) syslogng.Context {
	ctx.ControlNamespace = ns
	return ctx
}
