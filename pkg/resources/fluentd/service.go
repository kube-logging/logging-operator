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
	"context"

	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *Reconciler) service() (runtime.Object, reconciler.DesiredState, error) {
	ctx := context.TODO()
	fluentdSpec := r.GetFluentdSpec(ctx)
	desired := &corev1.Service{
		ObjectMeta: r.FluentdObjectMeta(ServiceName, ComponentFluentd),
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "tcp-fluentd",
					Protocol:   corev1.ProtocolTCP,
					Port:       ServicePort,
					TargetPort: intstr.IntOrString{IntVal: fluentdSpec.Port},
				},
				{
					Name:       "udp-fluentd",
					Protocol:   corev1.ProtocolUDP,
					Port:       ServicePort,
					TargetPort: intstr.IntOrString{IntVal: fluentdSpec.Port},
				},
			},
			Selector: r.Logging.GetFluentdLabels(ComponentFluentd, *fluentdSpec),
			Type:     corev1.ServiceTypeClusterIP,
		},
	}

	beforeUpdateHook := reconciler.DesiredStateHook(func(current runtime.Object) error {
		if s, ok := current.(*corev1.Service); ok {
			desired.Spec.ClusterIP = s.Spec.ClusterIP
		} else {
			return errors.Errorf("failed to cast service object %+v", current)
		}
		return nil
	})

	return desired, beforeUpdateHook, nil
}

func (r *Reconciler) serviceMetrics() (runtime.Object, reconciler.DesiredState, error) {
	ctx := context.TODO()
	fluentdSpec := r.GetFluentdSpec(ctx)
	if fluentdSpec.Metrics != nil {
		return &corev1.Service{
			ObjectMeta: r.FluentdObjectMeta(ServiceName+"-metrics", ComponentFluentd),
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Protocol:   corev1.ProtocolTCP,
						Name:       "http-metrics",
						Port:       fluentdSpec.Metrics.Port,
						TargetPort: intstr.IntOrString{IntVal: fluentdSpec.Metrics.Port},
					},
				},
				Selector:  r.Logging.GetFluentdLabels(ComponentFluentd, *fluentdSpec),
				Type:      corev1.ServiceTypeClusterIP,
				ClusterIP: "None",
			},
		}, reconciler.StatePresent, nil
	}
	return &corev1.Service{
		ObjectMeta: r.FluentdObjectMeta(ServiceName+"-monitor", ComponentFluentd),
		Spec:       corev1.ServiceSpec{}}, reconciler.StateAbsent, nil
}

func (r *Reconciler) monitorServiceMetrics() (runtime.Object, reconciler.DesiredState, error) {
	var SampleLimit uint64 = 0
	ctx := context.TODO()
	fluentdSpec := r.GetFluentdSpec(ctx)
	if fluentdSpec.Metrics != nil && fluentdSpec.Metrics.ServiceMonitor {
		objectMetadata := r.FluentdObjectMeta(ServiceName+"-metrics", ComponentFluentd)
		if fluentdSpec.Metrics.ServiceMonitorConfig.AdditionalLabels != nil {
			for k, v := range fluentdSpec.Metrics.ServiceMonitorConfig.AdditionalLabels {
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
					Path:                 fluentdSpec.GetFluentdMetricsPath(),
					Interval:             v1.Duration(fluentdSpec.Metrics.Interval),
					ScrapeTimeout:        v1.Duration(fluentdSpec.Metrics.Timeout),
					HonorLabels:          fluentdSpec.Metrics.ServiceMonitorConfig.HonorLabels,
					RelabelConfigs:       fluentdSpec.Metrics.ServiceMonitorConfig.Relabelings,
					MetricRelabelConfigs: fluentdSpec.Metrics.ServiceMonitorConfig.MetricsRelabelings,
					Scheme:               fluentdSpec.Metrics.ServiceMonitorConfig.Scheme,
					TLSConfig:            fluentdSpec.Metrics.ServiceMonitorConfig.TLSConfig,
				}},
				Selector:          v12.LabelSelector{MatchLabels: r.Logging.GetFluentdLabels(ComponentFluentd, *fluentdSpec)},
				NamespaceSelector: v1.NamespaceSelector{MatchNames: []string{r.Logging.Spec.ControlNamespace}},
				SampleLimit:       &SampleLimit,
			},
		}, reconciler.StatePresent, nil
	}
	return &v1.ServiceMonitor{
		ObjectMeta: r.FluentdObjectMeta(ServiceName+"-metrics", ComponentFluentd),
		Spec:       v1.ServiceMonitorSpec{},
	}, reconciler.StateAbsent, nil
}

func (r *Reconciler) serviceBufferMetrics() (runtime.Object, reconciler.DesiredState, error) {
	ctx := context.TODO()
	fluentdSpec := r.GetFluentdSpec(ctx)
	if fluentdSpec.BufferVolumeMetrics != nil {
		port := int32(defaultBufferVolumeMetricsPort)
		if fluentdSpec.BufferVolumeMetrics != nil && fluentdSpec.BufferVolumeMetrics.Port != 0 {
			port = fluentdSpec.BufferVolumeMetrics.Port
		}

		return &corev1.Service{
			ObjectMeta: r.FluentdObjectMeta(ServiceName+"-buffer-metrics", ComponentFluentd),
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Protocol:   corev1.ProtocolTCP,
						Name:       "buffer-metrics",
						Port:       port,
						TargetPort: intstr.IntOrString{IntVal: port},
					},
				},
				Selector:  r.Logging.GetFluentdLabels(ComponentFluentd, *fluentdSpec),
				Type:      corev1.ServiceTypeClusterIP,
				ClusterIP: "None",
			},
		}, reconciler.StatePresent, nil
	}
	return &corev1.Service{
		ObjectMeta: r.FluentdObjectMeta(ServiceName+"-buffer-monitor", ComponentFluentd),
		Spec:       corev1.ServiceSpec{}}, reconciler.StateAbsent, nil
}

func (r *Reconciler) monitorBufferServiceMetrics() (runtime.Object, reconciler.DesiredState, error) {
	var SampleLimit uint64 = 0
	ctx := context.TODO()
	fluentdSpec := r.GetFluentdSpec(ctx)
	if fluentdSpec.BufferVolumeMetrics != nil && fluentdSpec.BufferVolumeMetrics.ServiceMonitor {
		objectMetadata := r.FluentdObjectMeta(ServiceName+"-buffer-metrics", ComponentFluentd)
		if fluentdSpec.BufferVolumeMetrics.ServiceMonitorConfig.AdditionalLabels != nil {
			for k, v := range fluentdSpec.BufferVolumeMetrics.ServiceMonitorConfig.AdditionalLabels {
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
					Port:                 "buffer-metrics",
					Path:                 fluentdSpec.BufferVolumeMetrics.Path,
					Interval:             v1.Duration(fluentdSpec.BufferVolumeMetrics.Interval),
					ScrapeTimeout:        v1.Duration(fluentdSpec.BufferVolumeMetrics.Timeout),
					HonorLabels:          fluentdSpec.BufferVolumeMetrics.ServiceMonitorConfig.HonorLabels,
					RelabelConfigs:       fluentdSpec.BufferVolumeMetrics.ServiceMonitorConfig.Relabelings,
					MetricRelabelConfigs: fluentdSpec.BufferVolumeMetrics.ServiceMonitorConfig.MetricsRelabelings,
					Scheme:               fluentdSpec.BufferVolumeMetrics.ServiceMonitorConfig.Scheme,
					TLSConfig:            fluentdSpec.BufferVolumeMetrics.ServiceMonitorConfig.TLSConfig,
				}},
				Selector:          v12.LabelSelector{MatchLabels: r.Logging.GetFluentdLabels(ComponentFluentd, *fluentdSpec)},
				NamespaceSelector: v1.NamespaceSelector{MatchNames: []string{r.Logging.Spec.ControlNamespace}},
				SampleLimit:       &SampleLimit,
			},
		}, reconciler.StatePresent, nil
	}
	return &v1.ServiceMonitor{
		ObjectMeta: r.FluentdObjectMeta(ServiceName+"-buffer-metrics", ComponentFluentd),
		Spec:       v1.ServiceMonitorSpec{},
	}, reconciler.StateAbsent, nil
}

func (r *Reconciler) headlessService() (runtime.Object, reconciler.DesiredState, error) {
	ctx := context.TODO()
	fluentdSpec := r.GetFluentdSpec(ctx)
	desired := &corev1.Service{
		ObjectMeta: r.FluentdObjectMeta(ServiceName+"-headless", ComponentFluentd),
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "tcp-fluentd",
					Protocol: corev1.ProtocolTCP,
					// This port should match the containerport and targetPort will be automatically set to the same
					// https://github.com/kubernetes/kubernetes/issues/20488
					Port: fluentdSpec.Port,
				},
				{
					Name:     "udp-fluentd",
					Protocol: corev1.ProtocolUDP,
					// This port should match the containerport and targetPort will be automatically set to the same
					// https://github.com/kubernetes/kubernetes/issues/20488
					Port: fluentdSpec.Port,
				},
			},
			Selector:  r.Logging.GetFluentdLabels(ComponentFluentd, *fluentdSpec),
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: corev1.ClusterIPNone,
		},
	}
	return desired, reconciler.StatePresent, nil
}
