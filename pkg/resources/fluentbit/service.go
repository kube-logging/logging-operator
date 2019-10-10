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
	"github.com/banzaicloud/logging-operator/pkg/k8sutil"
	"github.com/banzaicloud/logging-operator/pkg/resources/templates"
	"github.com/banzaicloud/logging-operator/pkg/util"
	v1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *Reconciler) serviceMetrics() (runtime.Object, k8sutil.DesiredState) {
	if r.Logging.Spec.FluentbitSpec.Metrics != nil {
		return &corev1.Service{
			ObjectMeta: templates.FluentbitObjectMeta(
				r.Logging.QualifiedName(fluentbitServiceName+"-monitor"), util.MergeLabels(r.Logging.Labels, r.getFluentBitLabels()), r.Logging),
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{
					{
						Protocol:   corev1.ProtocolTCP,
						Name:       "monitor",
						Port:       r.Logging.Spec.FluentbitSpec.Metrics.Port,
						TargetPort: intstr.IntOrString{IntVal: r.Logging.Spec.FluentbitSpec.Metrics.Port},
					},
				},
				Selector:  r.getFluentBitLabels(),
				Type:      corev1.ServiceTypeClusterIP,
				ClusterIP: "None",
			},
		}, k8sutil.StatePresent
	}
	return &corev1.Service{
		ObjectMeta: templates.FluentbitObjectMeta(
			r.Logging.QualifiedName(fluentbitServiceName+"-monitor"), util.MergeLabels(r.Logging.Labels, r.getFluentBitLabels()), r.Logging),
		Spec: corev1.ServiceSpec{}}, k8sutil.StateAbsent
}

func (r *Reconciler) monitorServiceMetrics() (runtime.Object, k8sutil.DesiredState) {
	if r.Logging.Spec.FluentbitSpec.Metrics != nil {
		return &v1.ServiceMonitor{
			ObjectMeta: templates.FluentbitObjectMeta(r.Logging.QualifiedName(fluentbitServiceName+"-metrics"), util.MergeLabels(r.Logging.Labels, r.getFluentBitLabels()), r.Logging),
			Spec: v1.ServiceMonitorSpec{
				JobLabel:        "",
				TargetLabels:    nil,
				PodTargetLabels: nil,
				Endpoints: []v1.Endpoint{{
					Port: "metrics",
					Path: r.Logging.Spec.FluentbitSpec.Metrics.Path,
				}},
				Selector:          v12.LabelSelector{MatchLabels: util.MergeLabels(r.Logging.Labels, r.getFluentBitLabels(), generataLoggingRefLabels(r.Logging.ObjectMeta.GetName()))},
				NamespaceSelector: v1.NamespaceSelector{MatchNames: []string{r.Logging.Spec.ControlNamespace}},
				SampleLimit:       0,
			},
		}, k8sutil.StatePresent
	}
	return &v1.ServiceMonitor{ObjectMeta: templates.FluentbitObjectMeta(r.Logging.QualifiedName(fluentbitServiceName+"-metrics"), util.MergeLabels(r.Logging.Labels, r.getFluentBitLabels()), r.Logging), Spec: v1.ServiceMonitorSpec{}}, k8sutil.StateAbsent

}
