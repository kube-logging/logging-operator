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

package syslogng

import (
	"fmt"
	"testing"

	"github.com/cisco-open/operator-tools/pkg/reconciler"
	"github.com/cisco-open/operator-tools/pkg/typeoverride"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kube-logging/logging-operator/pkg/resources/model"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

func TestServiceKeepsExistingPrimaryIPFamily(t *testing.T) {
	r := &Reconciler{
		Logging: &v1beta1.Logging{
			ObjectMeta: metav1.ObjectMeta{Name: "test"},
			Spec:       v1beta1.LoggingSpec{ControlNamespace: "default"},
		},
		syslogNGSpec:   &v1beta1.SyslogNGSpec{EnabledIPv6: true},
		clusterHasIPv6: true,
	}

	object, state, err := r.service()
	require.NoError(t, err)

	hook, ok := state.(reconciler.DesiredStateHook)
	require.True(t, ok, "service() must return a before-update hook")
	require.NoError(t, hook(&corev1.Service{Spec: corev1.ServiceSpec{
		ClusterIP:  "10.96.35.208",
		ClusterIPs: []string{"10.96.35.208", "fd00:10:96::989d"},
		IPFamilies: []corev1.IPFamily{corev1.IPv4Protocol, corev1.IPv6Protocol},
	}}))

	service, ok := object.(*corev1.Service)
	require.True(t, ok)
	require.Equal(t, []corev1.IPFamily{corev1.IPv4Protocol, corev1.IPv6Protocol}, service.Spec.IPFamilies)
}

// The metrics service overrides were declared on the CRD but never read, so anything a user set
// under spec.syslogNG.metricsService was silently dropped.
func TestMetricsServiceOverridesAreApplied(t *testing.T) {
	metricsEnabled := true
	r := &Reconciler{
		Logging: &v1beta1.Logging{
			ObjectMeta: metav1.ObjectMeta{Name: "test"},
			Spec:       v1beta1.LoggingSpec{ControlNamespace: "default"},
		},
		syslogNGSpec: &v1beta1.SyslogNGSpec{
			Metrics:             &v1beta1.Metrics{Enabled: &metricsEnabled, Port: 9577},
			BufferVolumeMetrics: &v1beta1.BufferMetrics{Metrics: v1beta1.Metrics{Enabled: &metricsEnabled, Port: 9578}},
			MetricsServiceOverrides: &typeoverride.Service{
				ObjectMeta: typeoverride.ObjectMeta{Annotations: map[string]string{"metrics": "yes"}},
			},
			BufferVolumeMetricsServiceOverrides: &typeoverride.Service{
				ObjectMeta: typeoverride.ObjectMeta{Annotations: map[string]string{"buffer": "yes"}},
			},
		},
	}

	metrics, _, err := r.serviceMetrics()
	require.NoError(t, err)
	require.Equal(t, "yes", metrics.(*corev1.Service).Annotations["metrics"])

	buffer, _, err := r.serviceBufferMetrics()
	require.NoError(t, err)
	require.Equal(t, "yes", buffer.(*corev1.Service).Annotations["buffer"])
}

// The reloader binary defaults to its own port, which is not the one the ServiceMonitor scrapes.
func TestConfigReloaderListensOnTheScrapedPort(t *testing.T) {
	container := configReloadContainer(&v1beta1.SyslogNGSpec{ConfigReloadImage: &v1beta1.BasicImageSpec{}})
	require.Contains(t, container.Args, "-port")
	require.Contains(t, container.Args, fmt.Sprint(model.ConfigReloaderMetricsPort))
}
