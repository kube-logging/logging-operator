// Copyright © 2021 Cisco Systems, Inc. and/or its affiliates
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
	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/merge"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	util "github.com/cisco-open/operator-tools/pkg/utils"
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (n *nodeAgentInstance) serviceMetrics() (runtime.Object, reconciler.DesiredState, error) {
	if n.nodeAgent.FluentbitSpec.Metrics != nil {
		desired := &corev1.Service{
			ObjectMeta: n.NodeAgentObjectMeta(fluentbitServiceName + "-monitor"),
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Protocol:   corev1.ProtocolTCP,
						Name:       "http-metrics",
						Port:       n.nodeAgent.FluentbitSpec.Metrics.Port,
						TargetPort: intstr.IntOrString{IntVal: n.nodeAgent.FluentbitSpec.Metrics.Port},
					},
				},
				Selector:  n.getFluentBitLabels(),
				Type:      corev1.ServiceTypeClusterIP,
				ClusterIP: "None",
			},
		}
		err := merge.Merge(desired, n.nodeAgent.FluentbitSpec.MetricsService)
		if err != nil {
			return desired, reconciler.StatePresent, errors.WrapIf(err, "unable to merge overrides to base object")
		}
		return desired, reconciler.StatePresent, nil
	}
	return &corev1.Service{
		ObjectMeta: n.NodeAgentObjectMeta(fluentbitServiceName + "-monitor"),
		Spec:       corev1.ServiceSpec{}}, reconciler.StateAbsent, nil
}

func (n *nodeAgentInstance) monitorServiceMetrics() (runtime.Object, reconciler.DesiredState, error) {
	var SampleLimit uint64 = 0
	if n.nodeAgent.FluentbitSpec.Metrics != nil && n.nodeAgent.FluentbitSpec.Metrics.ServiceMonitor {
		objectMetadata := n.NodeAgentObjectMeta(fluentbitServiceName + "-metrics")
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
					Interval:             v1.Duration(n.nodeAgent.FluentbitSpec.Metrics.Interval),
					ScrapeTimeout:        v1.Duration(n.nodeAgent.FluentbitSpec.Metrics.Timeout),
					Scheme:               n.nodeAgent.FluentbitSpec.Metrics.ServiceMonitorConfig.Scheme,
					TLSConfig:            n.nodeAgent.FluentbitSpec.Metrics.ServiceMonitorConfig.TLSConfig,
				}},
				Selector: v12.LabelSelector{
					MatchLabels: util.MergeLabels(n.getFluentBitLabels(), generateLoggingRefLabels(n.logging.ObjectMeta.GetName())),
				},
				NamespaceSelector: v1.NamespaceSelector{MatchNames: []string{n.logging.Spec.ControlNamespace}},
				SampleLimit:       &SampleLimit,
			},
		}, reconciler.StatePresent, nil
	}
	return &v1.ServiceMonitor{
		ObjectMeta: n.NodeAgentObjectMeta(fluentbitServiceName + "-metrics"),
		Spec:       v1.ServiceMonitorSpec{},
	}, reconciler.StateAbsent, nil
}
