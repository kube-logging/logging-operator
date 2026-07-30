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

package fixture

import (
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

func TestLoggingInfraShape(t *testing.T) {
	buf := Buffer("time")
	labels := map[string]string{"my-unique-label": "log-producer"}
	objs := LoggingInfra("infra-ns", "e2e", "test.tag", buf, labels)

	require.Len(t, objs, 4, "ClusterOutput, ClusterFlow, FluentbitAgent, Logging — in creation order")

	out := objs[0].(*v1beta1.ClusterOutput)
	require.Equal(t, "http", out.Name)
	require.Equal(t, "infra-ns", out.Namespace)
	require.Equal(t, InfraRef, out.Spec.LoggingRef)
	require.Equal(t, "http://e2e-test-receiver:8080/test.tag", out.Spec.HTTPOutput.Endpoint)
	require.Same(t, buf, out.Spec.HTTPOutput.Buffer)

	flow := objs[1].(*v1beta1.ClusterFlow)
	require.Equal(t, InfraRef, flow.Spec.LoggingRef)
	require.Equal(t, labels, flow.Spec.Match[0].ClusterSelect.Labels)
	require.Equal(t, []string{"http"}, flow.Spec.GlobalOutputRefs)

	agent := objs[2].(*v1beta1.FluentbitAgent)
	require.Equal(t, InfraRef, agent.Spec.LoggingRef)
	require.Equal(t, "config-reloader", agent.Spec.ConfigHotReload.Image.Repository)

	lg := objs[3].(*v1beta1.Logging)
	require.Equal(t, InfraRef, lg.Spec.LoggingRef)
	require.Equal(t, "infra-ns", lg.Spec.ControlNamespace)
	require.Equal(t, map[string]string{"tenant": InfraRef}, lg.Labels)
	require.Empty(t, lg.Spec.WatchNamespaces, "only the tenant side watches a namespace")
}

func TestLoggingTenantShape(t *testing.T) {
	buf := Buffer("time")
	labels := map[string]string{"my-unique-label": "log-producer"}
	objs := LoggingTenant("tenant-ns", "infra-ns", "e2e", "test.tag", buf, labels)

	require.Len(t, objs, 3, "Output, Flow, Logging — no agent on the tenant side")

	out := objs[0].(*v1beta1.Output)
	require.Equal(t, TenantRef, out.Spec.LoggingRef)
	// Cross-namespace: the receiver lives in the infra namespace.
	require.Equal(t, "http://e2e-test-receiver.infra-ns:8080/test.tag", out.Spec.HTTPOutput.Endpoint)

	flow := objs[1].(*v1beta1.Flow)
	require.Equal(t, TenantRef, flow.Spec.LoggingRef)
	require.Equal(t, []string{"http"}, flow.Spec.LocalOutputRefs)

	lg := objs[2].(*v1beta1.Logging)
	require.Equal(t, TenantRef, lg.Spec.LoggingRef)
	require.Equal(t, "tenant-ns", lg.Spec.ControlNamespace)
	// "tenant-ns" differs from TenantRef on purpose: were the two equal, the
	// assertion below would also pass against a hard-coded ref.
	require.NotEqual(t, TenantRef, flow.Namespace)
	require.Equal(t, []string{flow.Namespace}, lg.Spec.WatchNamespaces)
}

func TestTenancyAggregatorDivergesFromTheDefault(t *testing.T) {
	infra := LoggingInfra("i", "e2e", "t", Buffer("time"), nil)[3].(*v1beta1.Logging)
	tenant := LoggingTenant("t", "i", "e2e", "t", Buffer("time"), nil)[2].(*v1beta1.Logging)

	for _, lg := range []*v1beta1.Logging{infra, tenant} {
		fd := lg.Spec.FluentdSpec
		require.True(t, fd.DisablePvc)
		require.Equal(t, resource.MustParse("50m"), fd.Resources.Requests[corev1.ResourceCPU])
		require.Equal(t, resource.MustParse("50M"), fd.Resources.Requests[corev1.ResourceMemory])
		require.Empty(t, fd.Resources.Limits, "the tenancy specs set no limits")
		require.Zero(t, fd.Workers, "tenancy does not set Workers")
		require.Nil(t, fd.Scaling)
	}

	// The single-tenant builder must still carry its own, different values.
	def := FluentdSpec()
	require.False(t, def.DisablePvc)
	require.Equal(t, resource.MustParse("500m"), def.Resources.Limits[corev1.ResourceCPU])
}

func TestLoggingRoute(t *testing.T) {
	r := LoggingRoute()
	require.Equal(t, "tenants", r.Name)
	require.Equal(t, InfraRef, r.Spec.Source)
	require.Equal(t, metav1.LabelSelectorOpExists, r.Spec.Targets.MatchExpressions[0].Operator)
	require.Equal(t, "tenant", r.Spec.Targets.MatchExpressions[0].Key)
}
