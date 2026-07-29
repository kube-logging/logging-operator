// Copyright © 2025 Kube logging authors
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
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
)

// The two loggingRef values the multi-tenant suites route between.
const (
	InfraRef  = "infra"
	TenantRef = "tenant"
)

// TenantReceiverURL is the cross-namespace form the tenant output posts to,
// where the receiver lives in the infra namespace.
func TenantReceiverURL(release, nsInfra, tag string) string {
	return fmt.Sprintf("http://%s-test-receiver.%s:8080/%s", release, nsInfra, tag)
}

// LoggingInfra builds the infra-side objects: a ClusterOutput and ClusterFlow
// carrying the infra loggingRef, a FluentbitAgent, and the aggregator Logging.
// Returned in creation order; the caller creates them.
func LoggingInfra(nsInfra, release, tag string, buffer *output.Buffer, producerLabels map[string]string) []client.Object {
	out := &v1beta1.ClusterOutput{
		ObjectMeta: metav1.ObjectMeta{Name: "http", Namespace: nsInfra},
		Spec: v1beta1.ClusterOutputSpec{
			OutputSpec: v1beta1.OutputSpec{
				LoggingRef: InfraRef,
				HTTPOutput: &output.HTTPOutputConfig{
					Endpoint:    ReceiverURL(release, tag),
					ContentType: "application/json",
					Buffer:      buffer,
				},
			},
		},
	}

	flow := &v1beta1.ClusterFlow{
		ObjectMeta: metav1.ObjectMeta{Name: "flow", Namespace: nsInfra},
		Spec: v1beta1.ClusterFlowSpec{
			LoggingRef: InfraRef,
			Match: []v1beta1.ClusterMatch{
				{ClusterSelect: &v1beta1.ClusterSelect{Labels: producerLabels}},
			},
			GlobalOutputRefs: []string{out.Name},
		},
	}

	agent := &v1beta1.FluentbitAgent{
		ObjectMeta: metav1.ObjectMeta{Name: InfraRef},
		Spec: v1beta1.FluentbitSpec{
			LoggingRef:        InfraRef,
			ConfigHotReload:   &v1beta1.HotReload{Image: image(ConfigReloaderRepo)},
			BufferVolumeImage: image(NodeExporterRepo),
		},
	}

	logging := &v1beta1.Logging{
		ObjectMeta: metav1.ObjectMeta{Name: InfraRef, Labels: map[string]string{"tenant": InfraRef}},
		Spec: v1beta1.LoggingSpec{
			LoggingRef:       InfraRef,
			ControlNamespace: nsInfra,
			FluentdSpec:      tenancyFluentdSpec(),
		},
	}

	return []client.Object{out, flow, agent, logging}
}

// LoggingTenant builds the tenant-side objects: a namespaced Output and Flow
// carrying the tenant loggingRef, and the aggregator Logging watching the
// tenant namespace. Returned in creation order.
func LoggingTenant(nsTenant, nsInfra, release, tag string, buffer *output.Buffer, producerLabels map[string]string) []client.Object {
	out := &v1beta1.Output{
		ObjectMeta: metav1.ObjectMeta{Name: "http", Namespace: nsTenant},
		Spec: v1beta1.OutputSpec{
			LoggingRef: TenantRef,
			HTTPOutput: &output.HTTPOutputConfig{
				Endpoint:    TenantReceiverURL(release, nsInfra, tag),
				ContentType: "application/json",
				Buffer:      buffer,
			},
		},
	}

	flow := &v1beta1.Flow{
		ObjectMeta: metav1.ObjectMeta{Name: "flow", Namespace: nsTenant},
		Spec: v1beta1.FlowSpec{
			LoggingRef: TenantRef,
			Match: []v1beta1.Match{
				{Select: &v1beta1.Select{Labels: producerLabels}},
			},
			LocalOutputRefs: []string{out.Name},
		},
	}

	logging := &v1beta1.Logging{
		ObjectMeta: metav1.ObjectMeta{Name: TenantRef, Labels: map[string]string{"tenant": TenantRef}},
		Spec: v1beta1.LoggingSpec{
			LoggingRef:       TenantRef,
			ControlNamespace: nsTenant,
			WatchNamespaces:  []string{TenantRef},
			FluentdSpec:      tenancyFluentdSpec(),
		},
	}

	return []client.Object{out, flow, logging}
}

// LoggingRoute builds the route connecting the infra collectors to every
// aggregator carrying a tenant label.
func LoggingRoute() *v1beta1.LoggingRoute {
	return &v1beta1.LoggingRoute{
		ObjectMeta: metav1.ObjectMeta{Name: "tenants"},
		Spec: v1beta1.LoggingRouteSpec{
			Source: InfraRef,
			Targets: metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{Key: "tenant", Operator: metav1.LabelSelectorOpExists},
				},
			},
		},
	}
}

// tenancyFluentdSpec is deliberately not FluentdSpec(): both tenancy
// aggregators disable the PVC and request 50m/50M, where the twelve
// single-tenant call sites use the larger limits and requests pair.
func tenancyFluentdSpec() *v1beta1.FluentdSpec {
	return &v1beta1.FluentdSpec{
		Image:               image(FluentdImageRepo),
		ConfigReloaderImage: image(ConfigReloaderRepo),
		BufferVolumeImage:   image(NodeExporterRepo),
		DisablePvc:          true,
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceCPU:    resource.MustParse("50m"),
				corev1.ResourceMemory: resource.MustParse("50M"),
			},
		},
	}
}
