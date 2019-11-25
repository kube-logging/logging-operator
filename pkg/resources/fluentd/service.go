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
	"github.com/banzaicloud/logging-operator/pkg/k8sutil"
	v1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *Reconciler) service() (runtime.Object, k8sutil.DesiredState, error) {
	desired := &corev1.Service{
		ObjectMeta: r.FluentdObjectMeta(ServiceName),
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       24240,
					TargetPort: intstr.IntOrString{IntVal: 24240},
				},
			},
			Selector: r.getFluentdLabels(),
			Type:     corev1.ServiceTypeClusterIP,
		},
	}

	beforeUpdateHook := k8sutil.DesiredStateHook(func(current runtime.Object) error {
		if s, ok := current.(*corev1.Service); ok {
			desired.Spec.ClusterIP = s.Spec.ClusterIP
		} else {
			return errors.Errorf("failed to cast service object %+v", current)
		}
		return nil
	})

	return desired, beforeUpdateHook, nil
}

func (r *Reconciler) serviceMetrics() (runtime.Object, k8sutil.DesiredState, error) {
	if r.Logging.Spec.FluentdSpec.Metrics != nil {
		return &corev1.Service{
			ObjectMeta: r.FluentdObjectMeta(ServiceName + "-metrics"),
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Protocol:   corev1.ProtocolTCP,
						Name:       "metrics",
						Port:       r.Logging.Spec.FluentdSpec.Metrics.Port,
						TargetPort: intstr.IntOrString{IntVal: r.Logging.Spec.FluentdSpec.Metrics.Port},
					},
				},
				Selector:  r.getFluentdLabels(),
				Type:      corev1.ServiceTypeClusterIP,
				ClusterIP: "None",
			},
		}, k8sutil.StatePresent, nil
	}
	return &corev1.Service{
		ObjectMeta: r.FluentdObjectMeta(ServiceName + "-monitor"),
		Spec:       corev1.ServiceSpec{}}, k8sutil.StateAbsent, nil
}

func (r *Reconciler) monitorServiceMetrics() (runtime.Object, k8sutil.DesiredState, error) {
	if r.Logging.Spec.FluentdSpec.Metrics != nil && r.Logging.Spec.FluentdSpec.Metrics.ServiceMonitor {
		return &v1.ServiceMonitor{
			ObjectMeta: r.FluentdObjectMeta(ServiceName + "-metrics"),
			Spec: v1.ServiceMonitorSpec{
				JobLabel:        "",
				TargetLabels:    nil,
				PodTargetLabels: nil,
				Endpoints: []v1.Endpoint{{
					Port:          "metrics",
					Path:          r.Logging.Spec.FluentdSpec.Metrics.Path,
					Interval:      r.Logging.Spec.FluentdSpec.Metrics.Interval,
					ScrapeTimeout: r.Logging.Spec.FluentdSpec.Metrics.Timeout,
				}},
				Selector:          v12.LabelSelector{MatchLabels: r.getFluentdLabels()},
				NamespaceSelector: v1.NamespaceSelector{MatchNames: []string{r.Logging.Spec.ControlNamespace}},
				SampleLimit:       0,
			},
		}, k8sutil.StatePresent, nil
	}
	return &v1.ServiceMonitor{
		ObjectMeta: r.FluentdObjectMeta(ServiceName + "-metrics"),
		Spec:       v1.ServiceMonitorSpec{},
	}, k8sutil.StateAbsent, nil
}
