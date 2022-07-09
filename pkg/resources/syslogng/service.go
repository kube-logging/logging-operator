// Copyright © 2022 Banzai Cloud
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

package syslogng

import (
	"emperror.dev/errors"
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *Reconciler) service() (runtime.Object, reconciler.DesiredState, error) {
	desired := &corev1.Service{
		ObjectMeta: r.SyslogNGObjectMeta(ServiceName, ComponentSyslogNG),
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "tcp-syslog-ng",
					Protocol:   corev1.ProtocolTCP,
					Port:       601,
					TargetPort: intstr.IntOrString{IntVal: 601},
				},
				{
					Name:       "udp-syslog-ng",
					Protocol:   corev1.ProtocolUDP,
					Port:       514,
					TargetPort: intstr.IntOrString{IntVal: 514},
				},
			},
			Selector: r.Logging.GetSyslogNGLabels(ComponentSyslogNG),
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
	if r.Logging.Spec.SyslogNGSpec.Metrics != nil {
		return &corev1.Service{
			ObjectMeta: r.SyslogNGObjectMeta(ServiceName+"-metrics", ComponentSyslogNG),
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Protocol:   corev1.ProtocolTCP,
						Name:       "http-metrics",
						Port:       r.Logging.Spec.SyslogNGSpec.Metrics.Port,
						TargetPort: intstr.IntOrString{IntVal: r.Logging.Spec.SyslogNGSpec.Metrics.Port},
					},
				},
				Selector:  r.Logging.GetSyslogNGLabels(ComponentSyslogNG),
				Type:      corev1.ServiceTypeClusterIP,
				ClusterIP: "None",
			},
		}, reconciler.StatePresent, nil
	}
	return &corev1.Service{
		ObjectMeta: r.SyslogNGObjectMeta(ServiceName+"-monitor", ComponentSyslogNG),
		Spec:       corev1.ServiceSpec{}}, reconciler.StateAbsent, nil
}

func (r *Reconciler) monitorServiceMetrics() (runtime.Object, reconciler.DesiredState, error) {
	if r.Logging.Spec.SyslogNGSpec.Metrics != nil && r.Logging.Spec.SyslogNGSpec.Metrics.ServiceMonitor {
		objectMetadata := r.SyslogNGObjectMeta(ServiceName+"-metrics", ComponentSyslogNG)
		if r.Logging.Spec.SyslogNGSpec.Metrics.ServiceMonitorConfig.AdditionalLabels != nil {
			for k, v := range r.Logging.Spec.SyslogNGSpec.Metrics.ServiceMonitorConfig.AdditionalLabels {
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
					Path:                 r.Logging.Spec.SyslogNGSpec.Metrics.Path,
					Interval:             r.Logging.Spec.SyslogNGSpec.Metrics.Interval,
					ScrapeTimeout:        r.Logging.Spec.SyslogNGSpec.Metrics.Timeout,
					HonorLabels:          r.Logging.Spec.SyslogNGSpec.Metrics.ServiceMonitorConfig.HonorLabels,
					RelabelConfigs:       r.Logging.Spec.SyslogNGSpec.Metrics.ServiceMonitorConfig.Relabelings,
					MetricRelabelConfigs: r.Logging.Spec.SyslogNGSpec.Metrics.ServiceMonitorConfig.MetricsRelabelings,
					Scheme:               r.Logging.Spec.SyslogNGSpec.Metrics.ServiceMonitorConfig.Scheme,
					TLSConfig:            r.Logging.Spec.SyslogNGSpec.Metrics.ServiceMonitorConfig.TLSConfig,
				}},
				Selector:          v12.LabelSelector{MatchLabels: r.Logging.GetSyslogNGLabels(ComponentSyslogNG)},
				NamespaceSelector: v1.NamespaceSelector{MatchNames: []string{r.Logging.Spec.ControlNamespace}},
				SampleLimit:       0,
			},
		}, reconciler.StatePresent, nil
	}
	return &v1.ServiceMonitor{
		ObjectMeta: r.SyslogNGObjectMeta(ServiceName+"-metrics", ComponentSyslogNG),
		Spec:       v1.ServiceMonitorSpec{},
	}, reconciler.StateAbsent, nil
}

func (r *Reconciler) serviceBufferMetrics() (runtime.Object, reconciler.DesiredState, error) {
	if r.Logging.Spec.SyslogNGSpec.BufferVolumeMetrics != nil {
		port := int32(defaultBufferVolumeMetricsPort)
		if r.Logging.Spec.SyslogNGSpec.BufferVolumeMetrics != nil && r.Logging.Spec.SyslogNGSpec.BufferVolumeMetrics.Port != 0 {
			port = r.Logging.Spec.SyslogNGSpec.BufferVolumeMetrics.Port
		}

		return &corev1.Service{
			ObjectMeta: r.SyslogNGObjectMeta(ServiceName+"-buffer-metrics", ComponentSyslogNG),
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Protocol:   corev1.ProtocolTCP,
						Name:       "buffer-metrics",
						Port:       port,
						TargetPort: intstr.IntOrString{IntVal: port},
					},
				},
				Selector:  r.Logging.GetSyslogNGLabels(ComponentSyslogNG),
				Type:      corev1.ServiceTypeClusterIP,
				ClusterIP: "None",
			},
		}, reconciler.StatePresent, nil
	}
	return &corev1.Service{
		ObjectMeta: r.SyslogNGObjectMeta(ServiceName+"-buffer-monitor", ComponentSyslogNG),
		Spec:       corev1.ServiceSpec{}}, reconciler.StateAbsent, nil
}

func (r *Reconciler) monitorBufferServiceMetrics() (runtime.Object, reconciler.DesiredState, error) {
	if r.Logging.Spec.SyslogNGSpec.BufferVolumeMetrics != nil && r.Logging.Spec.SyslogNGSpec.BufferVolumeMetrics.ServiceMonitor {
		objectMetadata := r.SyslogNGObjectMeta(ServiceName+"-buffer-metrics", ComponentSyslogNG)
		if r.Logging.Spec.SyslogNGSpec.BufferVolumeMetrics.ServiceMonitorConfig.AdditionalLabels != nil {
			for k, v := range r.Logging.Spec.SyslogNGSpec.BufferVolumeMetrics.ServiceMonitorConfig.AdditionalLabels {
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
					Path:                 r.Logging.Spec.SyslogNGSpec.BufferVolumeMetrics.Path,
					Interval:             r.Logging.Spec.SyslogNGSpec.BufferVolumeMetrics.Interval,
					ScrapeTimeout:        r.Logging.Spec.SyslogNGSpec.BufferVolumeMetrics.Timeout,
					HonorLabels:          r.Logging.Spec.SyslogNGSpec.BufferVolumeMetrics.ServiceMonitorConfig.HonorLabels,
					RelabelConfigs:       r.Logging.Spec.SyslogNGSpec.BufferVolumeMetrics.ServiceMonitorConfig.Relabelings,
					MetricRelabelConfigs: r.Logging.Spec.SyslogNGSpec.BufferVolumeMetrics.ServiceMonitorConfig.MetricsRelabelings,
				}},
				Selector:          v12.LabelSelector{MatchLabels: r.Logging.GetSyslogNGLabels(ComponentSyslogNG)},
				NamespaceSelector: v1.NamespaceSelector{MatchNames: []string{r.Logging.Spec.ControlNamespace}},
				SampleLimit:       0,
			},
		}, reconciler.StatePresent, nil
	}
	return &v1.ServiceMonitor{
		ObjectMeta: r.SyslogNGObjectMeta(ServiceName+"-buffer-metrics", ComponentSyslogNG),
		Spec:       v1.ServiceMonitorSpec{},
	}, reconciler.StateAbsent, nil
}

func (r *Reconciler) headlessService() (runtime.Object, reconciler.DesiredState, error) {
	desired := &corev1.Service{
		ObjectMeta: r.SyslogNGObjectMeta(ServiceName+"-headless", ComponentSyslogNG),
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "tcp-syslog-ng",
					Protocol:   corev1.ProtocolTCP,
					Port:       601,
					TargetPort: intstr.IntOrString{IntVal: 601},
				},
				{
					Name:       "udp-syslog-ng",
					Protocol:   corev1.ProtocolUDP,
					Port:       514,
					TargetPort: intstr.IntOrString{IntVal: 514},
				},
			},
			Selector:  r.Logging.GetSyslogNGLabels(ComponentSyslogNG),
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: corev1.ClusterIPNone,
		},
	}
	return desired, reconciler.StatePresent, nil
}
