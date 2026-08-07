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
	"testing"

	"github.com/cisco-open/operator-tools/pkg/reconciler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

// newCheckPodReconciler builds a Reconciler that is just complete enough for
// newCheckPod, which only needs the logging resource and the fluentd spec.
func newCheckPodReconciler(t *testing.T, fluentdSpec *v1beta1.FluentdSpec) *Reconciler {
	t.Helper()

	logging := &v1beta1.Logging{}
	logging.Name = "test"
	logging.Spec.ControlNamespace = "logging"
	logging.Spec.FluentdSpec = fluentdSpec
	require.NoError(t, logging.SetDefaults())

	config := "<system>\n</system>\n"

	return New(nil, log.Log, logging, logging.Spec.FluentdSpec, nil, &config, nil, reconciler.ReconcilerOpts{})
}

func TestNewCheckPodDNSSettings(t *testing.T) {
	dnsConfig := &corev1.PodDNSConfig{
		Nameservers: []string{"10.0.0.53"},
		Searches:    []string{"logging.svc.cluster.local"},
	}

	tests := []struct {
		name              string
		dnsPolicy         corev1.DNSPolicy
		dnsConfig         *corev1.PodDNSConfig
		expectedDNSPolicy corev1.DNSPolicy
	}{
		{
			// Left empty the API server applies ClusterFirst, which is what the
			// configcheck pod did unconditionally before dnsPolicy was propagated.
			name:              "Unset",
			expectedDNSPolicy: "",
		},
		{
			// The interesting case: with None the pod resolves *only* through
			// dnsConfig, so dropping the policy while keeping dnsConfig gives the
			// configcheck pod a different resolver than the aggregator it checks for.
			name:              "NoneWithDNSConfig",
			dnsPolicy:         corev1.DNSNone,
			dnsConfig:         dnsConfig,
			expectedDNSPolicy: corev1.DNSNone,
		},
		{
			name:              "ClusterFirstWithHostNet",
			dnsPolicy:         corev1.DNSClusterFirstWithHostNet,
			expectedDNSPolicy: corev1.DNSClusterFirstWithHostNet,
		},
		{
			name:              "Default",
			dnsPolicy:         corev1.DNSDefault,
			expectedDNSPolicy: corev1.DNSDefault,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := &v1beta1.FluentdSpec{
				DNSPolicy: tt.dnsPolicy,
				DNSConfig: tt.dnsConfig,
			}
			r := newCheckPodReconciler(t, spec)

			pod := r.newCheckPod("deadbeef", *r.fluentdSpec)

			assert.Equal(t, tt.expectedDNSPolicy, pod.Spec.DNSPolicy)
			assert.Equal(t, tt.dnsConfig, pod.Spec.DNSConfig)
		})
	}
}

// TestNewCheckPodDNSSettingsMatchStatefulSet pins the configcheck pod and the
// aggregator to the same DNS settings, since a configcheck that resolves names
// differently from the aggregator can pass or fail for the wrong reason.
func TestNewCheckPodDNSSettingsMatchStatefulSet(t *testing.T) {
	ndots := "2"
	spec := &v1beta1.FluentdSpec{
		DNSPolicy: corev1.DNSNone,
		DNSConfig: &corev1.PodDNSConfig{
			Nameservers: []string{"10.0.0.53"},
			Options:     []corev1.PodDNSConfigOption{{Name: "ndots", Value: &ndots}},
		},
	}
	r := newCheckPodReconciler(t, spec)

	checkPod := r.newCheckPod("deadbeef", *r.fluentdSpec)

	obj, _, err := r.statefulset()
	require.NoError(t, err)
	sts, ok := obj.(*appsv1.StatefulSet)
	require.True(t, ok)

	assert.Equal(t, sts.Spec.Template.Spec.DNSPolicy, checkPod.Spec.DNSPolicy)
	assert.Equal(t, sts.Spec.Template.Spec.DNSConfig, checkPod.Spec.DNSConfig)
}

// TestNewCheckPodSidecarContainers pins the configcheck pod to carry the same
// sidecarContainers as the aggregator StatefulSet, since a sidecar that only
// mutates shared config (e.g. refreshing a GeoIP database via extraVolumes)
// needs to run before the check as well, or the check validates against stale
// input.
func TestNewCheckPodSidecarContainers(t *testing.T) {
	sidecar := corev1.Container{
		Name:  "fluentd-sidecar",
		Image: "busybox:1.37",
	}
	spec := &v1beta1.FluentdSpec{
		SidecarContainers: []corev1.Container{sidecar},
	}
	r := newCheckPodReconciler(t, spec)

	checkPod := r.newCheckPod("deadbeef", *r.fluentdSpec)

	assert.Contains(t, checkPod.Spec.Containers, sidecar)
}
