// Copyright © 2026 Kube logging authors
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

package fluentbit

import (
	"testing"

	"github.com/stretchr/testify/require"
	rbacv1 "k8s.io/api/rbac/v1"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

func TestClusterRoleKubeletResources(t *testing.T) {
	tests := []struct {
		name                    string
		useKubelet              string
		kubeletFineGrainedAuthz bool
		expected                []string
	}{
		{
			name:     "kubelet not used",
			expected: []string{"pods", "namespaces"},
		},
		{
			name:       "kubelet used",
			useKubelet: "On",
			expected:   []string{"pods", "namespaces", "nodes", "nodes/pods", "nodes/proxy"},
		},
		{
			name:                    "kubelet used with fine-grained kubelet authorization",
			useKubelet:              "On",
			kubeletFineGrainedAuthz: true,
			expected:                []string{"pods", "namespaces", "nodes", "nodes/pods"},
		},
	}

	rbacCreate := true
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logging := &v1beta1.Logging{Name: "test"}
			r := &Reconciler{
				Logging: logging,
				fluentbitSpec: &v1beta1.FluentbitSpec{
					Security:                &v1beta1.Security{RoleBasedAccessControlCreate: &rbacCreate},
					FilterKubernetes:        v1beta1.FilterKubernetes{UseKubelet: test.useKubelet},
					KubeletFineGrainedAuthz: test.kubeletFineGrainedAuthz,
				},
				nameProvider: NewLegacyFluentbitNameProvider(logging),
			}

			object, _, err := r.clusterRole()
			require.NoError(t, err)
			clusterRole, ok := object.(*rbacv1.ClusterRole)
			require.True(t, ok)
			require.Len(t, clusterRole.Rules, 1)
			require.Equal(t, test.expected, clusterRole.Rules[0].Resources)
		})
	}
}
