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

package fluentd

import (
	"strings"
	"testing"

	"github.com/cisco-open/operator-tools/pkg/reconciler"
	"github.com/cisco-open/operator-tools/pkg/typeoverride"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

// The API server rejects an update whose ipFamilies[0] disagrees with the already-allocated
// clusterIPs[0], and operator-tools cannot recognize that error as immutable, so the Logging
// stops converging until someone deletes the Service by hand.
func TestServiceKeepsExistingPrimaryIPFamily(t *testing.T) {
	tests := []struct {
		name     string
		current  *corev1.Service
		expected []corev1.IPFamily
	}{
		{
			name:     "no existing service keeps the desired IPv6-primary order",
			current:  nil,
			expected: []corev1.IPFamily{corev1.IPv6Protocol, corev1.IPv4Protocol},
		},
		{
			name: "IPv4-primary dual-stack service allocated by an earlier release is left alone",
			current: existingService(
				[]corev1.IPFamily{corev1.IPv4Protocol, corev1.IPv6Protocol},
				"10.96.35.208", "fd00:10:96::989d",
			),
			expected: []corev1.IPFamily{corev1.IPv4Protocol, corev1.IPv6Protocol},
		},
		{
			name: "single-stack IPv4 service cannot gain an IPv6 primary",
			current: existingService(
				[]corev1.IPFamily{corev1.IPv4Protocol}, "10.96.35.208",
			),
			expected: []corev1.IPFamily{corev1.IPv4Protocol},
		},
		{
			name: "already IPv6-primary service stays as desired",
			current: existingService(
				[]corev1.IPFamily{corev1.IPv6Protocol, corev1.IPv4Protocol},
				"fd00:10:96::989d", "10.96.35.208",
			),
			expected: []corev1.IPFamily{corev1.IPv6Protocol, corev1.IPv4Protocol},
		},
		{
			name:     "service with no allocated clusterIPs takes the desired order",
			current:  existingService(nil),
			expected: []corev1.IPFamily{corev1.IPv6Protocol, corev1.IPv4Protocol},
		},
		{
			name: "headless service has nothing allocated to preserve",
			current: existingService(
				[]corev1.IPFamily{corev1.IPv4Protocol}, corev1.ClusterIPNone,
			),
			expected: []corev1.IPFamily{corev1.IPv6Protocol, corev1.IPv4Protocol},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := ipv6Reconciler()

			object, state, err := r.service()
			require.NoError(t, err)

			if test.current != nil {
				hook, ok := state.(reconciler.DesiredStateHook)
				require.True(t, ok, "service() must return a before-update hook")
				require.NoError(t, hook(test.current))
			}

			service, ok := object.(*corev1.Service)
			require.True(t, ok)
			require.Equal(t, test.expected, service.Spec.IPFamilies)
			requireIPFamiliesMatchClusterIPs(t, service)
		})
	}
}

// The policy decides whether the API server may widen the pinned families or must truncate them,
// so pinning it in the wrong direction either strands an existing Service on IPv4 or releases an
// already-allocated address.
func TestServiceIPFamilyPolicyOnlyResistsNarrowing(t *testing.T) {
	singleStack := corev1.IPFamilyPolicySingleStack
	preferDualStack := corev1.IPFamilyPolicyPreferDualStack

	tests := []struct {
		name        string
		enabledIPv6 bool
		overrides   *corev1.ServiceSpec
		current     *corev1.Service
		expected    corev1.IPFamilyPolicy
	}{
		{
			name:        "enabling IPv6 on a single-stack service still widens it",
			enabledIPv6: true,
			current:     serviceWithPolicy(singleStack, []corev1.IPFamily{corev1.IPv4Protocol}, "10.96.35.208"),
			expected:    preferDualStack,
		},
		{
			name:      "narrowing a dual-stack service keeps its policy so no address is released",
			overrides: &corev1.ServiceSpec{IPFamilyPolicy: &singleStack, IPFamilies: []corev1.IPFamily{corev1.IPv4Protocol}},
			current: serviceWithPolicy(preferDualStack,
				[]corev1.IPFamily{corev1.IPv6Protocol, corev1.IPv4Protocol},
				"fd00:10:96::989d", "10.96.35.208"),
			expected: preferDualStack,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := ipv6Reconciler()
			r.fluentdSpec.EnabledIPv6 = test.enabledIPv6
			if test.overrides != nil {
				r.fluentdSpec.ServiceOverrides = &typeoverride.Service{Spec: *test.overrides}
			}

			object, state, err := r.service()
			require.NoError(t, err)

			hook, ok := state.(reconciler.DesiredStateHook)
			require.True(t, ok)
			require.NoError(t, hook(test.current))

			service, ok := object.(*corev1.Service)
			require.True(t, ok)
			require.NotNil(t, service.Spec.IPFamilyPolicy)
			require.Equal(t, test.expected, *service.Spec.IPFamilyPolicy)
		})
	}
}

func serviceWithPolicy(policy corev1.IPFamilyPolicy, families []corev1.IPFamily, clusterIPs ...string) *corev1.Service {
	service := existingService(families, clusterIPs...)
	service.Spec.IPFamilyPolicy = &policy
	return service
}

func requireIPFamiliesMatchClusterIPs(t *testing.T, service *corev1.Service) {
	t.Helper()

	for i, clusterIP := range service.Spec.ClusterIPs {
		if clusterIP == corev1.ClusterIPNone {
			continue
		}
		require.Less(t, i, len(service.Spec.IPFamilies), "clusterIPs is longer than ipFamilies")
		family := corev1.IPv4Protocol
		if strings.Contains(clusterIP, ":") {
			family = corev1.IPv6Protocol
		}
		require.Equal(t, family, service.Spec.IPFamilies[i], "clusterIPs[%d] %q contradicts ipFamilies[%d]", i, clusterIP, i)
	}
}

func ipv6Reconciler() *Reconciler {
	logging := &v1beta1.Logging{
		ObjectMeta: metav1.ObjectMeta{Name: "test"},
		Spec:       v1beta1.LoggingSpec{ControlNamespace: "default"},
	}
	return &Reconciler{
		Logging:        logging,
		fluentdSpec:    &v1beta1.FluentdSpec{EnabledIPv6: true},
		clusterHasIPv6: true,
	}
}

func existingService(families []corev1.IPFamily, clusterIPs ...string) *corev1.Service {
	spec := corev1.ServiceSpec{IPFamilies: families, ClusterIPs: clusterIPs}
	if len(clusterIPs) > 0 {
		spec.ClusterIP = clusterIPs[0]
	}
	return &corev1.Service{Spec: spec}
}
