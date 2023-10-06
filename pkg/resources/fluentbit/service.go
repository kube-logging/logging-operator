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
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	util "github.com/cisco-open/operator-tools/pkg/utils"
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *Reconciler) serviceMetrics() (runtime.Object, reconciler.DesiredState, error) {
	if r.fluentbitSpec.Metrics != nil {
		return &corev1.Service{
			ObjectMeta: r.FluentbitObjectMeta(fluentbitServiceName + "-monitor"),
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Protocol:   corev1.ProtocolTCP,
						Name:       "http-metrics",
						Port:       r.fluentbitSpec.Metrics.Port,
						TargetPort: intstr.IntOrString{IntVal: r.fluentbitSpec.Metrics.Port},
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
	var SampleLimit uint64 = 0
	if r.fluentbitSpec.Metrics != nil && r.fluentbitSpec.Metrics.ServiceMonitor {
		objectMetadata := r.FluentbitObjectMeta(fluentbitServiceName + "-metrics")
		if r.fluentbitSpec.Metrics.ServiceMonitorConfig.AdditionalLabels != nil {
			for k, v := range r.fluentbitSpec.Metrics.ServiceMonitorConfig.AdditionalLabels {
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
					Path:                 r.fluentbitSpec.Metrics.Path,
					HonorLabels:          r.fluentbitSpec.Metrics.ServiceMonitorConfig.HonorLabels,
					RelabelConfigs:       r.fluentbitSpec.Metrics.ServiceMonitorConfig.Relabelings,
					MetricRelabelConfigs: r.fluentbitSpec.Metrics.ServiceMonitorConfig.MetricsRelabelings,
					Scheme:               r.fluentbitSpec.Metrics.ServiceMonitorConfig.Scheme,
					TLSConfig:            r.fluentbitSpec.Metrics.ServiceMonitorConfig.TLSConfig,
				}},
				Selector: v12.LabelSelector{
					MatchLabels: util.MergeLabels(r.fluentbitSpec.Labels, r.getFluentBitLabels(), generateLoggingRefLabels(r.Logging.GetName())),
				},
				NamespaceSelector: v1.NamespaceSelector{MatchNames: []string{r.Logging.Spec.ControlNamespace}},
				SampleLimit:       &SampleLimit,
			},
		}, reconciler.StatePresent, nil
	}
	return &v1.ServiceMonitor{
		ObjectMeta: r.FluentbitObjectMeta(fluentbitServiceName + "-metrics"),
		Spec:       v1.ServiceMonitorSpec{},
	}, reconciler.StateAbsent, nil
}

func (r *Reconciler) serviceBufferMetrics() (runtime.Object, reconciler.DesiredState, error) {
	if r.fluentbitSpec.BufferVolumeMetrics != nil {
		port := int32(defaultBufferVolumeMetricsPort)
		if r.fluentbitSpec.BufferVolumeMetrics.Port != 0 {
			port = r.fluentbitSpec.BufferVolumeMetrics.Port
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
	var SampleLimit uint64 = 0
	if r.fluentbitSpec.BufferVolumeMetrics != nil && r.fluentbitSpec.BufferVolumeMetrics.ServiceMonitor {
		objectMetadata := r.FluentbitObjectMeta(fluentbitServiceName + "-buffer-metrics")

		objectMetadata.Labels = util.MergeLabels(objectMetadata.Labels, r.fluentbitSpec.BufferVolumeMetrics.ServiceMonitorConfig.AdditionalLabels)
		return &v1.ServiceMonitor{
			ObjectMeta: objectMetadata,
			Spec: v1.ServiceMonitorSpec{
				JobLabel:        "",
				TargetLabels:    nil,
				PodTargetLabels: nil,
				Endpoints: []v1.Endpoint{{
					Port:                 "buffer-metrics",
					Path:                 r.fluentbitSpec.BufferVolumeMetrics.Path,
					Interval:             v1.Duration(r.fluentbitSpec.BufferVolumeMetrics.Interval),
					ScrapeTimeout:        v1.Duration(r.fluentbitSpec.BufferVolumeMetrics.Timeout),
					HonorLabels:          r.fluentbitSpec.BufferVolumeMetrics.ServiceMonitorConfig.HonorLabels,
					RelabelConfigs:       r.fluentbitSpec.BufferVolumeMetrics.ServiceMonitorConfig.Relabelings,
					MetricRelabelConfigs: r.fluentbitSpec.BufferVolumeMetrics.ServiceMonitorConfig.MetricsRelabelings,
					Scheme:               r.fluentbitSpec.Metrics.ServiceMonitorConfig.Scheme,
					TLSConfig:            r.fluentbitSpec.Metrics.ServiceMonitorConfig.TLSConfig,
				}},
				Selector:          v12.LabelSelector{MatchLabels: r.getFluentBitLabels()},
				NamespaceSelector: v1.NamespaceSelector{MatchNames: []string{r.Logging.Spec.ControlNamespace}},
				SampleLimit:       &SampleLimit,
			},
		}, reconciler.StatePresent, nil
	}
	return &v1.ServiceMonitor{
		ObjectMeta: r.FluentbitObjectMeta(fluentbitServiceName + "-buffer-metrics"),
		Spec:       v1.ServiceMonitorSpec{},
	}, reconciler.StateAbsent, nil
}
