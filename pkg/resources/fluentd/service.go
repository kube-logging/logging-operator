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

package fluentd

import (
	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/merge"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	"github.com/kube-logging/logging-operator/pkg/resources/kubetool"
	"github.com/kube-logging/logging-operator/pkg/resources/model"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *Reconciler) service() (runtime.Object, reconciler.DesiredState, error) {
	desired := &corev1.Service{
		ObjectMeta: r.FluentdObjectMeta(ServiceName, ComponentFluentd),
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "tcp-fluentd",
					Protocol:   corev1.ProtocolTCP,
					Port:       ServicePort,
					TargetPort: intstr.IntOrString{IntVal: r.fluentdSpec.Port},
				},
				{
					Name:       "udp-fluentd",
					Protocol:   corev1.ProtocolUDP,
					Port:       ServicePort,
					TargetPort: intstr.IntOrString{IntVal: r.fluentdSpec.Port},
				},
			},
			Selector: r.Logging.GetFluentdLabels(ComponentFluentd, *r.fluentdSpec),
			Type:     corev1.ServiceTypeClusterIP,
		},
	}

	err := merge.Merge(desired, r.fluentdSpec.ServiceOverrides)
	if err != nil {
		return desired, reconciler.StatePresent, errors.WrapIf(err, "unable to merge overrides to base object")
	}

	if r.fluentdSpec.EnabledIPv6 {
		v1beta1.EnableIPv6Options(&desired.Spec)
	}

	beforeUpdateHook := reconciler.DesiredStateHook(func(current runtime.Object) error {
		if s, ok := current.(*corev1.Service); ok {
			desired.Spec.ClusterIP = s.Spec.ClusterIP
			// Preserve ClusterIPs for dual-stack configuration
			if len(s.Spec.ClusterIPs) > 0 {
				desired.Spec.ClusterIPs = s.Spec.ClusterIPs
			}
		} else {
			return errors.Errorf("failed to cast service object %+v", current)
		}
		return nil
	})

	return desired, beforeUpdateHook, nil
}

func (r *Reconciler) serviceMetrics() (runtime.Object, reconciler.DesiredState, error) {
	objectMetadata := r.FluentdObjectMeta(ServiceName+"-metrics", ComponentFluentd)

	if r.fluentdSpec.Metrics != nil && r.fluentdSpec.Metrics.IsEnabled() {
		desired := &corev1.Service{
			ObjectMeta: objectMetadata,
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Protocol:   corev1.ProtocolTCP,
						Name:       "http-metrics",
						Port:       r.fluentdSpec.Metrics.Port,
						TargetPort: intstr.IntOrString{IntVal: r.fluentdSpec.Metrics.Port},
					},
					{
						Protocol:   corev1.ProtocolTCP,
						Name:       model.ConfigReloaderMetricsPortName,
						Port:       model.ConfigReloaderMetricsPort,
						TargetPort: intstr.IntOrString{IntVal: model.ConfigReloaderMetricsPort},
					},
				},
				Selector:  r.Logging.GetFluentdLabels(ComponentFluentd, *r.fluentdSpec),
				Type:      corev1.ServiceTypeClusterIP,
				ClusterIP: corev1.ClusterIPNone,
			},
		}

		if r.fluentdSpec.EnabledIPv6 {
			v1beta1.EnableIPv6Options(&desired.Spec)
		}

		return desired, reconciler.StatePresent, nil
	}
	return &corev1.Service{
		ObjectMeta: objectMetadata,
		Spec:       corev1.ServiceSpec{}}, reconciler.StateAbsent, nil
}

func (r *Reconciler) monitorServiceMetrics() (runtime.Object, reconciler.DesiredState, error) {
	objectMetadata := r.FluentdObjectMeta(ServiceName+"-metrics", ComponentFluentd)

	if r.fluentdSpec.Metrics != nil && r.fluentdSpec.Metrics.IsEnabled() && r.fluentdSpec.Metrics.ServiceMonitor {
		if r.fluentdSpec.Metrics.ServiceMonitorConfig.Scheme == "" {
			r.fluentdSpec.Metrics.ServiceMonitorConfig.Scheme = kubetool.To(v1.SchemeHTTP).String()
		}

		if r.fluentdSpec.Metrics.ServiceMonitorConfig.AdditionalLabels != nil {
			for k, v := range r.fluentdSpec.Metrics.ServiceMonitorConfig.AdditionalLabels {
				objectMetadata.Labels[k] = v
			}
		}

		var SampleLimit uint64 = 0
		return &v1.ServiceMonitor{
			ObjectMeta: objectMetadata,
			Spec: v1.ServiceMonitorSpec{
				JobLabel:        "",
				TargetLabels:    nil,
				PodTargetLabels: nil,
				Endpoints: []v1.Endpoint{
					{
						Port:                 "http-metrics",
						Path:                 r.fluentdSpec.GetFluentdMetricsPath(),
						Interval:             v1.Duration(r.fluentdSpec.Metrics.Interval),
						ScrapeTimeout:        v1.Duration(r.fluentdSpec.Metrics.Timeout),
						HonorLabels:          r.fluentdSpec.Metrics.ServiceMonitorConfig.HonorLabels,
						RelabelConfigs:       r.fluentdSpec.Metrics.ServiceMonitorConfig.Relabelings,
						MetricRelabelConfigs: r.fluentdSpec.Metrics.ServiceMonitorConfig.MetricsRelabelings,
						Scheme:               kubetool.To(v1.Scheme(r.fluentdSpec.Metrics.ServiceMonitorConfig.Scheme)),
						TLSConfig:            r.fluentdSpec.Metrics.ServiceMonitorConfig.TLSConfig,
					},
					{
						Port:                 model.ConfigReloaderMetricsPortName,
						Path:                 "/metrics",
						Interval:             v1.Duration(r.fluentdSpec.Metrics.Interval),
						ScrapeTimeout:        v1.Duration(r.fluentdSpec.Metrics.Timeout),
						HonorLabels:          r.fluentdSpec.Metrics.ServiceMonitorConfig.HonorLabels,
						RelabelConfigs:       r.fluentdSpec.Metrics.ServiceMonitorConfig.Relabelings,
						MetricRelabelConfigs: r.fluentdSpec.Metrics.ServiceMonitorConfig.MetricsRelabelings,
						Scheme:               kubetool.To(v1.Scheme(r.fluentdSpec.Metrics.ServiceMonitorConfig.Scheme)),
						TLSConfig:            r.fluentdSpec.Metrics.ServiceMonitorConfig.TLSConfig,
					},
				},
				Selector:          v12.LabelSelector{MatchLabels: r.Logging.GetFluentdLabels(ComponentFluentd, *r.fluentdSpec)},
				NamespaceSelector: v1.NamespaceSelector{MatchNames: []string{r.Logging.Spec.ControlNamespace}},
				SampleLimit:       &SampleLimit,
			},
		}, reconciler.StatePresent, nil
	}
	return &v1.ServiceMonitor{
		ObjectMeta: objectMetadata,
		Spec:       v1.ServiceMonitorSpec{},
	}, reconciler.StateAbsent, nil
}

func (r *Reconciler) serviceBufferMetrics() (runtime.Object, reconciler.DesiredState, error) {
	objectMetadata := r.FluentdObjectMeta(ServiceName+"-buffer-metrics", ComponentFluentd)

	if r.fluentdSpec.BufferVolumeMetrics != nil && r.fluentdSpec.BufferVolumeMetrics.IsEnabled() {
		port := int32(defaultBufferVolumeMetricsPort)
		if r.fluentdSpec.BufferVolumeMetrics.Port != 0 {
			port = r.fluentdSpec.BufferVolumeMetrics.Port
		}

		desired := &corev1.Service{
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
				Selector:  r.Logging.GetFluentdLabels(ComponentFluentd, *r.fluentdSpec),
				Type:      corev1.ServiceTypeClusterIP,
				ClusterIP: corev1.ClusterIPNone,
			},
		}

		if r.fluentdSpec.EnabledIPv6 {
			v1beta1.EnableIPv6Options(&desired.Spec)
		}

		return desired, reconciler.StatePresent, nil
	}
	return &corev1.Service{
		ObjectMeta: objectMetadata,
		Spec:       corev1.ServiceSpec{}}, reconciler.StateAbsent, nil
}

func (r *Reconciler) monitorBufferServiceMetrics() (runtime.Object, reconciler.DesiredState, error) {
	objectMetadata := r.FluentdObjectMeta(ServiceName+"-buffer-metrics", ComponentFluentd)

	if r.fluentdSpec.BufferVolumeMetrics != nil && r.fluentdSpec.BufferVolumeMetrics.IsEnabled() && r.fluentdSpec.BufferVolumeMetrics.ServiceMonitor {
		if r.fluentdSpec.BufferVolumeMetrics.ServiceMonitorConfig.Scheme == "" {
			r.fluentdSpec.BufferVolumeMetrics.ServiceMonitorConfig.Scheme = kubetool.To(v1.SchemeHTTP).String()
		}

		if r.fluentdSpec.BufferVolumeMetrics.ServiceMonitorConfig.AdditionalLabels != nil {
			for k, v := range r.fluentdSpec.BufferVolumeMetrics.ServiceMonitorConfig.AdditionalLabels {
				objectMetadata.Labels[k] = v
			}
		}

		var SampleLimit uint64 = 0
		return &v1.ServiceMonitor{
			ObjectMeta: objectMetadata,
			Spec: v1.ServiceMonitorSpec{
				JobLabel:        "",
				TargetLabels:    nil,
				PodTargetLabels: nil,
				Endpoints: []v1.Endpoint{{
					Port:                 "buffer-metrics",
					Path:                 r.fluentdSpec.BufferVolumeMetrics.Path,
					Interval:             v1.Duration(r.fluentdSpec.BufferVolumeMetrics.Interval),
					ScrapeTimeout:        v1.Duration(r.fluentdSpec.BufferVolumeMetrics.Timeout),
					HonorLabels:          r.fluentdSpec.BufferVolumeMetrics.ServiceMonitorConfig.HonorLabels,
					RelabelConfigs:       r.fluentdSpec.BufferVolumeMetrics.ServiceMonitorConfig.Relabelings,
					MetricRelabelConfigs: r.fluentdSpec.BufferVolumeMetrics.ServiceMonitorConfig.MetricsRelabelings,
					Scheme:               kubetool.To(v1.Scheme(r.fluentdSpec.BufferVolumeMetrics.ServiceMonitorConfig.Scheme)),
					TLSConfig:            r.fluentdSpec.BufferVolumeMetrics.ServiceMonitorConfig.TLSConfig,
				}},
				Selector:          v12.LabelSelector{MatchLabels: r.Logging.GetFluentdLabels(ComponentFluentd, *r.fluentdSpec)},
				NamespaceSelector: v1.NamespaceSelector{MatchNames: []string{r.Logging.Spec.ControlNamespace}},
				SampleLimit:       &SampleLimit,
			},
		}, reconciler.StatePresent, nil
	}
	return &v1.ServiceMonitor{
		ObjectMeta: objectMetadata,
		Spec:       v1.ServiceMonitorSpec{},
	}, reconciler.StateAbsent, nil
}

func (r *Reconciler) headlessService() (runtime.Object, reconciler.DesiredState, error) {
	desired := &corev1.Service{
		ObjectMeta: r.FluentdObjectMeta(ServiceName+"-headless", ComponentFluentd),
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "tcp-fluentd",
					Protocol: corev1.ProtocolTCP,
					// This port should match the containerport and targetPort will be automatically set to the same
					// https://github.com/kubernetes/kubernetes/issues/20488
					Port: r.fluentdSpec.Port,
				},
				{
					Name:     "udp-fluentd",
					Protocol: corev1.ProtocolUDP,
					// This port should match the containerport and targetPort will be automatically set to the same
					// https://github.com/kubernetes/kubernetes/issues/20488
					Port: r.fluentdSpec.Port,
				},
			},
			Selector:  r.Logging.GetFluentdLabels(ComponentFluentd, *r.fluentdSpec),
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: corev1.ClusterIPNone,
		},
	}

	if r.fluentdSpec.EnabledIPv6 {
		v1beta1.EnableIPv6Options(&desired.Spec)
	}

	return desired, reconciler.StatePresent, nil
}
