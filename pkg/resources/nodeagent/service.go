// Copyright © 2021 Banzai Cloud
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

package nodeagent

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	util "github.com/banzaicloud/operator-tools/pkg/utils"
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (n *nodeAgentInstance) serviceMetrics() (runtime.Object, reconciler.DesiredState, error) {
	if n.nodeAgent.FluentbitSpec.Metrics != nil {
		return &corev1.Service{
			ObjectMeta: r.FluentbitObjectMeta(fluentbitServiceName + "-monitor"),
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Protocol:   corev1.ProtocolTCP,
						Name:       "http-metrics",
						Port:       n.nodeAgent.FluentbitSpec.Metrics.Port,
						TargetPort: intstr.IntOrString{IntVal: n.nodeAgent.FluentbitSpec.Metrics.Port},
					},
				},
				Selector:  r.getFluentBitLabels(),
				Type:      corev1.ServiceTypeClusterIP,
				ClusterIP: "None",
			},
		}, reconciler.StatePresent, nil
	}
	return &corev1.Service{
		ObjectMeta: r.FluentbitObjectMeta(fluentbitServiceName + "-monitor"),
		Spec:       corev1.ServiceSpec{}}, reconciler.StateAbsent, nil
}

func (n *nodeAgentInstance) monitorServiceMetrics() (runtime.Object, reconciler.DesiredState, error) {
	if n.nodeAgent.FluentbitSpec.Metrics != nil && n.nodeAgent.FluentbitSpec.Metrics.ServiceMonitor {
		objectMetadata := r.FluentbitObjectMeta(fluentbitServiceName + "-metrics")
		if n.nodeAgent.FluentbitSpec.Metrics.ServiceMonitorConfig.AdditionalLabels != nil {
			for k, v := range n.nodeAgent.FluentbitSpec.Metrics.ServiceMonitorConfig.AdditionalLabels {
				objectMetadata.Labels[k] = v
			}
		}
		return &v1.ServiceMonitor{
			ObjectMeta: objectMetadata,
			Spec: v1.ServiceMonitorSpec{
				JobLabel:        "",
				TargetLabels:    nil,
				PodTargetLabels: nil,
				Endpoints: []v1.Endpoint{{
					Port:                 "http-metrics",
					Path:                 n.nodeAgent.FluentbitSpec.Metrics.Path,
					HonorLabels:          n.nodeAgent.FluentbitSpec.Metrics.ServiceMonitorConfig.HonorLabels,
					RelabelConfigs:       n.nodeAgent.FluentbitSpec.Metrics.ServiceMonitorConfig.Relabelings,
					MetricRelabelConfigs: n.nodeAgent.FluentbitSpec.Metrics.ServiceMonitorConfig.MetricsRelabelings,
				}},
				Selector: v12.LabelSelector{
					MatchLabels: util.MergeLabels(n.nodeAgent.FluentbitSpec.Labels, r.getFluentBitLabels(), generateLoggingRefLabels(r.Logging.ObjectMeta.GetName())),
				},
				NamespaceSelector: v1.NamespaceSelector{MatchNames: []string{r.Logging.Spec.ControlNamespace}},
				SampleLimit:       0,
			},
		}, reconciler.StatePresent, nil
	}
	return &v1.ServiceMonitor{
		ObjectMeta: r.FluentbitObjectMeta(fluentbitServiceName + "-metrics"),
		Spec:       v1.ServiceMonitorSpec{},
	}, reconciler.StateAbsent, nil
}
