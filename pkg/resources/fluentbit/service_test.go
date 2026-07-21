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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

func TestMetricsServicesIPFamilies(t *testing.T) {
	metricsEnabled := true
	tests := []struct {
		name        string
		enabledIPv6 bool
	}{
		{name: "single-stack"},
		{name: "dual-stack", enabledIPv6: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logging := &v1beta1.Logging{ObjectMeta: metav1.ObjectMeta{Name: "test"}}
			r := &Reconciler{
				Logging: logging,
				fluentbitSpec: &v1beta1.FluentbitSpec{
					EnabledIPv6: test.enabledIPv6,
					Metrics: &v1beta1.Metrics{
						Enabled: &metricsEnabled,
						Port:    2020,
					},
					BufferVolumeMetrics: &v1beta1.Metrics{
						Enabled: &metricsEnabled,
						Port:    9200,
					},
				},
				nameProvider: NewLegacyFluentbitNameProvider(logging),
			}

			metricsService, _, err := r.serviceMetrics()
			require.NoError(t, err)
			bufferMetricsService, _, err := r.serviceBufferMetrics()
			require.NoError(t, err)

			assertIPFamilies(t, metricsService, test.enabledIPv6)
			assertIPFamilies(t, bufferMetricsService, test.enabledIPv6)
		})
	}
}

func assertIPFamilies(t *testing.T, object runtime.Object, enabledIPv6 bool) {
	t.Helper()

	service, ok := object.(*corev1.Service)
	require.True(t, ok)

	if !enabledIPv6 {
		require.Nil(t, service.Spec.IPFamilyPolicy)
		require.Nil(t, service.Spec.IPFamilies)
		return
	}

	require.NotNil(t, service.Spec.IPFamilyPolicy)
	require.Equal(t, corev1.IPFamilyPolicyPreferDualStack, *service.Spec.IPFamilyPolicy)
	require.Equal(t, []corev1.IPFamily{corev1.IPv4Protocol, corev1.IPv6Protocol}, service.Spec.IPFamilies)
}
