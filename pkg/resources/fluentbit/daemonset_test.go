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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

func TestDaemonSetSidecarContainers(t *testing.T) {
	sidecar := corev1.Container{
		Name:    "liveness-sidecar",
		Image:   "alpine",
		Command: []string{"sleep", "infinity"},
	}

	tests := []struct {
		name     string
		spec     v1beta1.FluentbitSpec
		expected []string
	}{
		{
			name:     "no sidecars",
			expected: []string{containerName},
		},
		{
			name: "sidecars are appended after the fluent-bit container",
			spec: v1beta1.FluentbitSpec{
				SidecarContainers: []corev1.Container{sidecar},
			},
			expected: []string{containerName, sidecar.Name},
		},
		{
			name: "sidecars are appended after the operator managed containers",
			spec: v1beta1.FluentbitSpec{
				ConfigHotReload:     &v1beta1.HotReload{},
				BufferVolumeMetrics: &v1beta1.Metrics{Enabled: new(true)},
				SidecarContainers:   []corev1.Container{sidecar},
			},
			expected: []string{containerName, "config-reloader", "buffer-metrics-sidecar", sidecar.Name},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			spec := test.spec
			require.NoError(t, v1beta1.FluentBitDefaults(&spec))

			logging := &v1beta1.Logging{Name: "test"}
			r := &Reconciler{
				Logging:       logging,
				fluentbitSpec: &spec,
				nameProvider:  NewLegacyFluentbitNameProvider(logging),
			}

			object, _, err := r.daemonSet()
			require.NoError(t, err)

			daemonSet, ok := object.(*appsv1.DaemonSet)
			require.True(t, ok)

			var names []string
			for _, c := range daemonSet.Spec.Template.Spec.Containers {
				names = append(names, c.Name)
			}
			require.Equal(t, test.expected, names)

			if len(spec.SidecarContainers) != 0 {
				require.Equal(t, sidecar, daemonSet.Spec.Template.Spec.Containers[len(names)-1])
			}
		})
	}
}
