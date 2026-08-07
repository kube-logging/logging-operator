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
	"testing"

	"github.com/cisco-open/operator-tools/pkg/reconciler"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
