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

package config

import (
	"strings"
	"testing"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/config/render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRenderClusterFlow(t *testing.T) {
	testCases := map[string]struct {
		clusterFlow v1beta1.SyslogNGClusterFlow
		expected    string
	}{
		"nil match": {
			clusterFlow: v1beta1.SyslogNGClusterFlow{
				ObjectMeta: v1.ObjectMeta{
					Name:      "test_clusterflow",
					Namespace: "test_ns",
				},
				Spec: v1beta1.SyslogNGClusterFlowSpec{
					Match: nil,
				},
			},
			expected: Untab(`log {
source("test_input");
};
`),
		},
		"empty match": {
			clusterFlow: v1beta1.SyslogNGClusterFlow{
				ObjectMeta: v1.ObjectMeta{
					Name:      "test_clusterflow",
					Namespace: "test_ns",
				},
				Spec: v1beta1.SyslogNGClusterFlowSpec{
					Match: &v1beta1.SyslogNGMatch{},
				},
			},
			expected: Untab(`log {
source("test_input");
};
`),
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			out := strings.Builder{}
			require.NoError(t, renderClusterFlow("test_input", testCase.clusterFlow, nil)(render.RenderContext{
				Out: &out,
			}))
			assert.Equal(t, testCase.expected, out.String())
		})
	}
}
