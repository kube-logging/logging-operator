// Copyright © 2024 Kube logging authors
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

package telemetry_controller

import (
	"fmt"

	"github.com/cisco-open/operator-tools/pkg/utils"
	"github.com/kube-logging/logging-operator/pkg/resources/fluentd"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	telemetryv1alpha1 "github.com/kube-logging/telemetry-controller/api/telemetry/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	k8sMetadataName  = "kubernetes.io/metadata.name"
	tenantKind       = "Tenant"
	subscriptionKind = "Subscription"
	outputKind       = "Output"
)

func CreateTenant(logging *v1beta1.Logging) *telemetryv1alpha1.Tenant {
	tenantBase := &telemetryv1alpha1.Tenant{
		TypeMeta: metav1.TypeMeta{
			APIVersion: telemetryv1alpha1.GroupVersion.String(),
			Kind:       tenantKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   logging.Name,
			Labels: logging.Spec.RouteConfig.TenantLabels,
		},
		Spec: telemetryv1alpha1.TenantSpec{
			SubscriptionNamespaceSelectors: []metav1.LabelSelector{
				{
					MatchLabels: map[string]string{
						k8sMetadataName: logging.Spec.ControlNamespace,
					},
				},
			},
			PersistenceConfig: telemetryv1alpha1.PersistenceConfig{
				EnableFileStorage: true,
			},
		},
	}

	// If no watchNamespaces are specified, the tenant should watch all namespaces
	LogSourceNamespaceSelectors := convertToLabelSelectors(logging.Spec.WatchNamespaces, logging.Spec.WatchNamespaceSelector)
	if LogSourceNamespaceSelectors != nil {
		tenantBase.Spec.LogSourceNamespaceSelectors = LogSourceNamespaceSelectors
	} else {
		tenantBase.Spec.SelectFromAllNamespaces = true
	}

	return tenantBase
}

func CreateSubscription(logging *v1beta1.Logging) *telemetryv1alpha1.Subscription {
	return &telemetryv1alpha1.Subscription{
		TypeMeta: metav1.TypeMeta{
			APIVersion: telemetryv1alpha1.GroupVersion.String(),
			Kind:       subscriptionKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      logging.Name,
			Namespace: logging.Spec.ControlNamespace,
		},
		Spec: telemetryv1alpha1.SubscriptionSpec{
			Condition: "true",
			Outputs: []telemetryv1alpha1.NamespacedName{
				{
					Namespace: logging.Spec.ControlNamespace,
					Name:      logging.Name,
				},
			},
		},
	}
}

func CreateOutput(logging *v1beta1.Logging) *telemetryv1alpha1.Output {
	return &telemetryv1alpha1.Output{
		TypeMeta: metav1.TypeMeta{
			APIVersion: telemetryv1alpha1.GroupVersion.String(),
			Kind:       outputKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      logging.Name,
			Namespace: logging.Spec.ControlNamespace,
		},
		Spec: telemetryv1alpha1.OutputSpec{
			Fluentforward: &telemetryv1alpha1.Fluentforward{
				TCPClientSettings: telemetryv1alpha1.TCPClientSettings{
					Endpoint: &telemetryv1alpha1.Endpoint{
						TCPAddr:               aggregatorEndpoint(logging),
						ValidateTCPResolution: false,
					},
					TLSSetting: &telemetryv1alpha1.TLSClientSetting{
						Insecure: true,
					},
				},
				Tag: utils.StringPointer("otelcol"),
				Kubernetes: &telemetryv1alpha1.KubernetesMetadata{
					Key:              "kubernetes",
					IncludePodLabels: true,
				},
			},
		},
	}
}

func aggregatorEndpoint(l *v1beta1.Logging) *string {
	endpoint := fmt.Sprintf("%s.%s.svc%s:%d", l.QualifiedName(fluentd.ServiceName), l.Spec.ControlNamespace, l.ClusterDomainAsSuffix(), fluentd.ServicePort)
	return &endpoint
}

func convertToLabelSelectors(watchNamespaces []string, watchNamespaceSelector *metav1.LabelSelector) []metav1.LabelSelector {
	if len(watchNamespaces) == 0 && watchNamespaceSelector == nil {
		return nil
	}

	var labelSelectors []metav1.LabelSelector
	for _, ns := range watchNamespaces {
		labelSelectors = append(labelSelectors, metav1.LabelSelector{
			MatchExpressions: []metav1.LabelSelectorRequirement{
				{
					Key:      k8sMetadataName,
					Operator: metav1.LabelSelectorOpIn,
					Values:   []string{ns},
				},
			},
		})
	}

	if watchNamespaceSelector != nil {
		labelSelectors = append(labelSelectors, *watchNamespaceSelector)
	}

	return labelSelectors
}
