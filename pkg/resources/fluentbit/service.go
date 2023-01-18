// Copyright © 2019 Banzai Cloud
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
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	util "github.com/banzaicloud/operator-tools/pkg/utils"
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *Reconciler) serviceMetrics() (runtime.Object, reconciler.DesiredState, error) {
	if r.Logging.Spec.FluentbitSpec.Metrics != nil {
		return &corev1.Service{
			ObjectMeta: r.FluentbitObjectMeta(fluentbitServiceName + "-monitor"),
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Protocol:   corev1.ProtocolTCP,
						Name:       "http-metrics",
						Port:       r.Logging.Spec.FluentbitSpec.Metrics.Port,
						TargetPort: intstr.IntOrString{IntVal: r.Logging.Spec.FluentbitSpec.Metrics.Port},
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

func (r *Reconciler) monitorServiceMetrics() (runtime.Object, reconciler.DesiredState, error) {
	if r.Logging.Spec.FluentbitSpec.Metrics != nil && r.Logging.Spec.FluentbitSpec.Metrics.ServiceMonitor {
		objectMetadata := r.FluentbitObjectMeta(fluentbitServiceName + "-metrics")
		if r.Logging.Spec.FluentbitSpec.Metrics.ServiceMonitorConfig.AdditionalLabels != nil {
			for k, v := range r.Logging.Spec.FluentbitSpec.Metrics.ServiceMonitorConfig.AdditionalLabels {
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
					Path:                 r.Logging.Spec.FluentbitSpec.Metrics.Path,
					HonorLabels:          r.Logging.Spec.FluentbitSpec.Metrics.ServiceMonitorConfig.HonorLabels,
					RelabelConfigs:       r.Logging.Spec.FluentbitSpec.Metrics.ServiceMonitorConfig.Relabelings,
					MetricRelabelConfigs: r.Logging.Spec.FluentbitSpec.Metrics.ServiceMonitorConfig.MetricsRelabelings,
					Scheme:               r.Logging.Spec.FluentbitSpec.Metrics.ServiceMonitorConfig.Scheme,
					TLSConfig:            r.Logging.Spec.FluentbitSpec.Metrics.ServiceMonitorConfig.TLSConfig,
				}},
				Selector: v12.LabelSelector{
					MatchLabels: util.MergeLabels(r.Logging.Spec.FluentbitSpec.Labels, r.getFluentBitLabels(), generateLoggingRefLabels(r.Logging.ObjectMeta.GetName())),
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

func (r *Reconciler) serviceBufferMetrics() (runtime.Object, reconciler.DesiredState, error) {
	if r.Logging.Spec.FluentbitSpec.BufferVolumeMetrics != nil {
		port := int32(defaultBufferVolumeMetricsPort)
		if r.Logging.Spec.FluentbitSpec.BufferVolumeMetrics.Port != 0 {
			port = r.Logging.Spec.FluentbitSpec.BufferVolumeMetrics.Port
		}

		return &corev1.Service{
			ObjectMeta: r.FluentbitObjectMeta(fluentbitServiceName + "-buffer-metrics"),

			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Protocol:   corev1.ProtocolTCP,
						Name:       "buffer-metrics",
						Port:       port,
						TargetPort: intstr.IntOrString{IntVal: port},
					},
				},
				Selector:  r.getFluentBitLabels(),
				Type:      corev1.ServiceTypeClusterIP,
				ClusterIP: "None",
			},
		}, reconciler.StatePresent, nil
	}
	return &corev1.Service{
		ObjectMeta: r.FluentbitObjectMeta(fluentbitServiceName + "-buffer-monitor"),
		Spec:       corev1.ServiceSpec{}}, reconciler.StateAbsent, nil
}

func (r *Reconciler) monitorBufferServiceMetrics() (runtime.Object, reconciler.DesiredState, error) {
	if r.Logging.Spec.FluentbitSpec.BufferVolumeMetrics != nil && r.Logging.Spec.FluentbitSpec.BufferVolumeMetrics.ServiceMonitor {
		objectMetadata := r.FluentbitObjectMeta(fluentbitServiceName + "-buffer-metrics")

		objectMetadata.Labels = util.MergeLabels(objectMetadata.Labels, r.Logging.Spec.FluentbitSpec.BufferVolumeMetrics.ServiceMonitorConfig.AdditionalLabels)
		return &v1.ServiceMonitor{
			ObjectMeta: objectMetadata,
			Spec: v1.ServiceMonitorSpec{
				JobLabel:        "",
				TargetLabels:    nil,
				PodTargetLabels: nil,
				Endpoints: []v1.Endpoint{{
					Port:                 "buffer-metrics",
					Path:                 r.Logging.Spec.FluentbitSpec.BufferVolumeMetrics.Path,
					Interval:             v1.Duration(r.Logging.Spec.FluentbitSpec.BufferVolumeMetrics.Interval),
					ScrapeTimeout:        v1.Duration(r.Logging.Spec.FluentbitSpec.BufferVolumeMetrics.Timeout),
					HonorLabels:          r.Logging.Spec.FluentbitSpec.BufferVolumeMetrics.ServiceMonitorConfig.HonorLabels,
					RelabelConfigs:       r.Logging.Spec.FluentbitSpec.BufferVolumeMetrics.ServiceMonitorConfig.Relabelings,
					MetricRelabelConfigs: r.Logging.Spec.FluentbitSpec.BufferVolumeMetrics.ServiceMonitorConfig.MetricsRelabelings,
				}},
				Selector:          v12.LabelSelector{MatchLabels: r.getFluentBitLabels()},
				NamespaceSelector: v1.NamespaceSelector{MatchNames: []string{r.Logging.Spec.ControlNamespace}},
				SampleLimit:       0,
			},
		}, reconciler.StatePresent, nil
	}
	return &v1.ServiceMonitor{
		ObjectMeta: r.FluentbitObjectMeta(fluentbitServiceName + "-buffer-metrics"),
		Spec:       v1.ServiceMonitorSpec{},
	}, reconciler.StateAbsent, nil
}
