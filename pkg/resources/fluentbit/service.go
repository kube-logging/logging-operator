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
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/kube-logging/logging-operator/pkg/resources/kubetool"
	"github.com/kube-logging/logging-operator/pkg/resources/model"
)

func (r *Reconciler) serviceMetrics() (runtime.Object, reconciler.DesiredState, error) {
	objectMetadata := r.FluentbitObjectMeta(fluentbitServiceName + "-metrics")

	if r.fluentbitSpec.Metrics != nil && r.fluentbitSpec.Metrics.IsEnabled() {
		ports := []corev1.ServicePort{
			{
				Protocol:   corev1.ProtocolTCP,
				Name:       "http-metrics",
				Port:       r.fluentbitSpec.Metrics.Port,
				TargetPort: intstr.IntOrString{IntVal: r.fluentbitSpec.Metrics.Port},
			},
		}
		// Add config-reloader metrics port if hotreload is configured
		if r.fluentbitSpec.ConfigHotReload != nil {
			ports = append(ports, corev1.ServicePort{
				Protocol:   corev1.ProtocolTCP,
				Name:       model.ConfigReloaderMetricsPortName,
				Port:       model.ConfigReloaderMetricsPort,
				TargetPort: intstr.IntOrString{IntVal: model.ConfigReloaderMetricsPort},
			})
		}
		return &corev1.Service{
			ObjectMeta: objectMetadata,
			Spec: corev1.ServiceSpec{
				Ports:     ports,
				Selector:  r.getFluentBitLabels(),
				Type:      corev1.ServiceTypeClusterIP,
				ClusterIP: corev1.ClusterIPNone,
			},
		}, reconciler.StatePresent, nil
	}
	return &corev1.Service{
		ObjectMeta: objectMetadata,
		Spec:       corev1.ServiceSpec{},
	}, reconciler.StateAbsent, nil
}

func (r *Reconciler) monitorServiceMetrics() (runtime.Object, reconciler.DesiredState, error) {
	objectMetadata := r.FluentbitObjectMeta(fluentbitServiceName + "-metrics")

	if r.fluentbitSpec.Metrics != nil && r.fluentbitSpec.Metrics.IsEnabled() && r.fluentbitSpec.Metrics.ServiceMonitor {
		if r.fluentbitSpec.Metrics.ServiceMonitorConfig.Scheme == "" {
			r.fluentbitSpec.Metrics.ServiceMonitorConfig.Scheme = kubetool.To(monitoringv1.SchemeHTTP).String()
		}

		if r.fluentbitSpec.Metrics.ServiceMonitorConfig.AdditionalLabels != nil {
			for k, v := range r.fluentbitSpec.Metrics.ServiceMonitorConfig.AdditionalLabels {
				objectMetadata.Labels[k] = v
			}
		}

		var SampleLimit uint64 = 0
		endpoints := []monitoringv1.Endpoint{
			{
				Port:                 "http-metrics",
				Path:                 r.fluentbitSpec.Metrics.Path,
				Interval:             monitoringv1.Duration(r.fluentbitSpec.Metrics.Interval),
				ScrapeTimeout:        monitoringv1.Duration(r.fluentbitSpec.Metrics.Timeout),
				HonorLabels:          r.fluentbitSpec.Metrics.ServiceMonitorConfig.HonorLabels,
				RelabelConfigs:       r.fluentbitSpec.Metrics.ServiceMonitorConfig.Relabelings,
				MetricRelabelConfigs: r.fluentbitSpec.Metrics.ServiceMonitorConfig.MetricsRelabelings,
				Scheme:               kubetool.To(monitoringv1.Scheme(r.fluentbitSpec.Metrics.ServiceMonitorConfig.Scheme)),
				HTTPConfigWithProxyAndTLSFiles: monitoringv1.HTTPConfigWithProxyAndTLSFiles{
					HTTPConfigWithTLSFiles: monitoringv1.HTTPConfigWithTLSFiles{
						TLSConfig: r.fluentbitSpec.Metrics.ServiceMonitorConfig.TLSConfig,
					},
				},
			},
		}
		// Add config-reloader metrics endpoint if hotreload is configured
		if r.fluentbitSpec.ConfigHotReload != nil {
			endpoints = append(endpoints, monitoringv1.Endpoint{
				Port:                 model.ConfigReloaderMetricsPortName,
				Path:                 "/metrics",
				Interval:             monitoringv1.Duration(r.fluentbitSpec.Metrics.Interval),
				ScrapeTimeout:        monitoringv1.Duration(r.fluentbitSpec.Metrics.Timeout),
				HonorLabels:          r.fluentbitSpec.Metrics.ServiceMonitorConfig.HonorLabels,
				RelabelConfigs:       r.fluentbitSpec.Metrics.ServiceMonitorConfig.Relabelings,
				MetricRelabelConfigs: r.fluentbitSpec.Metrics.ServiceMonitorConfig.MetricsRelabelings,
				Scheme:               kubetool.To(monitoringv1.Scheme(r.fluentbitSpec.Metrics.ServiceMonitorConfig.Scheme)),
				HTTPConfigWithProxyAndTLSFiles: monitoringv1.HTTPConfigWithProxyAndTLSFiles{
					HTTPConfigWithTLSFiles: monitoringv1.HTTPConfigWithTLSFiles{
						TLSConfig: r.fluentbitSpec.Metrics.ServiceMonitorConfig.TLSConfig,
					},
				},
			})
		}
		return &monitoringv1.ServiceMonitor{
			ObjectMeta: objectMetadata,
			Spec: monitoringv1.ServiceMonitorSpec{
				JobLabel:        "",
				TargetLabels:    nil,
				PodTargetLabels: nil,
				Endpoints:       endpoints,
				Selector: metav1.LabelSelector{
					MatchLabels: util.MergeLabels(r.fluentbitSpec.Labels, r.getFluentBitLabels(), generateLoggingRefLabels(r.Logging.GetName())),
				},
				NamespaceSelector: monitoringv1.NamespaceSelector{MatchNames: []string{r.Logging.Spec.FluentbitAgentNamespace}},
				SampleLimit:       &SampleLimit,
			},
		}, reconciler.StatePresent, nil
	}
	return &monitoringv1.ServiceMonitor{
		ObjectMeta: objectMetadata,
		Spec:       monitoringv1.ServiceMonitorSpec{},
	}, reconciler.StateAbsent, nil
}

func (r *Reconciler) serviceBufferMetrics() (runtime.Object, reconciler.DesiredState, error) {
	objectMetadata := r.FluentbitObjectMeta(fluentbitServiceName + "-buffer-metrics")

	if r.fluentbitSpec.BufferVolumeMetrics != nil && r.fluentbitSpec.BufferVolumeMetrics.IsEnabled() {
		port := int32(defaultBufferVolumeMetricsPort)
		if r.fluentbitSpec.BufferVolumeMetrics.Port != 0 {
			port = r.fluentbitSpec.BufferVolumeMetrics.Port
		}

		return &corev1.Service{
			ObjectMeta: objectMetadata,
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
				ClusterIP: corev1.ClusterIPNone,
			},
		}, reconciler.StatePresent, nil
	}
	return &corev1.Service{
		ObjectMeta: objectMetadata,
		Spec:       corev1.ServiceSpec{},
	}, reconciler.StateAbsent, nil
}

func (r *Reconciler) monitorBufferServiceMetrics() (runtime.Object, reconciler.DesiredState, error) {
	objectMetadata := r.FluentbitObjectMeta(fluentbitServiceName + "-buffer-metrics")

	if r.fluentbitSpec.BufferVolumeMetrics != nil && r.fluentbitSpec.BufferVolumeMetrics.IsEnabled() && r.fluentbitSpec.BufferVolumeMetrics.ServiceMonitor {
		if r.fluentbitSpec.BufferVolumeMetrics.ServiceMonitorConfig.Scheme == "" {
			r.fluentbitSpec.BufferVolumeMetrics.ServiceMonitorConfig.Scheme = kubetool.To(monitoringv1.SchemeHTTP).String()
		}

		objectMetadata.Labels = util.MergeLabels(objectMetadata.Labels, r.fluentbitSpec.BufferVolumeMetrics.ServiceMonitorConfig.AdditionalLabels)

		var SampleLimit uint64 = 0
		return &monitoringv1.ServiceMonitor{
			ObjectMeta: objectMetadata,
			Spec: monitoringv1.ServiceMonitorSpec{
				JobLabel:        "",
				TargetLabels:    nil,
				PodTargetLabels: nil,
				Endpoints: []monitoringv1.Endpoint{{
					Port:                 "buffer-metrics",
					Path:                 r.fluentbitSpec.BufferVolumeMetrics.Path,
					Interval:             monitoringv1.Duration(r.fluentbitSpec.BufferVolumeMetrics.Interval),
					ScrapeTimeout:        monitoringv1.Duration(r.fluentbitSpec.BufferVolumeMetrics.Timeout),
					HonorLabels:          r.fluentbitSpec.BufferVolumeMetrics.ServiceMonitorConfig.HonorLabels,
					RelabelConfigs:       r.fluentbitSpec.BufferVolumeMetrics.ServiceMonitorConfig.Relabelings,
					MetricRelabelConfigs: r.fluentbitSpec.BufferVolumeMetrics.ServiceMonitorConfig.MetricsRelabelings,
					Scheme:               kubetool.To(monitoringv1.Scheme(r.fluentbitSpec.BufferVolumeMetrics.ServiceMonitorConfig.Scheme)),
					HTTPConfigWithProxyAndTLSFiles: monitoringv1.HTTPConfigWithProxyAndTLSFiles{
						HTTPConfigWithTLSFiles: monitoringv1.HTTPConfigWithTLSFiles{
							TLSConfig: r.fluentbitSpec.BufferVolumeMetrics.ServiceMonitorConfig.TLSConfig,
						},
					},
				}},
				Selector:          metav1.LabelSelector{MatchLabels: r.getFluentBitLabels()},
				NamespaceSelector: monitoringv1.NamespaceSelector{MatchNames: []string{r.Logging.Spec.FluentbitAgentNamespace}},
				SampleLimit:       &SampleLimit,
			},
		}, reconciler.StatePresent, nil
	}
	return &monitoringv1.ServiceMonitor{
		ObjectMeta: objectMetadata,
		Spec:       monitoringv1.ServiceMonitorSpec{},
	}, reconciler.StateAbsent, nil
}
