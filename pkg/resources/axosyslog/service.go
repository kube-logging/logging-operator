// Copyright © 2025 Kube logging authors
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

package axosyslog

import (
	"fmt"

	"github.com/cisco-open/operator-tools/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

func Service(object any) (runtime.Object, reconciler.DesiredState, error) {
	axoSyslog, ok := object.(*v1beta1.AxoSyslog)
	if !ok {
		return nil, reconciler.StateAbsent, fmt.Errorf("expected *v1beta1.AxoSyslog, got %T", axoSyslog)
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      commonAxoSyslogObjectValue,
			Namespace: axoSyslog.Namespace,
			Labels: map[string]string{
				LabelAppName:      commonAxoSyslogObjectValue,
				LabelAppComponent: commonAxoSyslogObjectValue,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "otlp-grpc",
					Protocol:   corev1.ProtocolTCP,
					Port:       4317,
					TargetPort: intstr.IntOrString{IntVal: 4317},
				},
				{
					Name:       "otlp-http",
					Protocol:   corev1.ProtocolTCP,
					Port:       4318,
					TargetPort: intstr.IntOrString{IntVal: 4318},
				},
			},
			Selector: map[string]string{
				LabelAppName:      commonAxoSyslogObjectValue,
				LabelAppComponent: commonAxoSyslogObjectValue,
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	beforeUpdateHook := reconciler.DesiredStateHook(func(current runtime.Object) error {
		if s, ok := current.(*corev1.Service); ok {
			service.Spec.ClusterIP = s.Spec.ClusterIP
		} else {
			return fmt.Errorf("failed to cast service object %+v", current)
		}
		return nil
	})

	return service, beforeUpdateHook, nil
}

func HeadlessService(object any) (runtime.Object, reconciler.DesiredState, error) {
	axoSyslog, ok := object.(*v1beta1.AxoSyslog)
	if !ok {
		return nil, reconciler.StateAbsent, fmt.Errorf("expected *v1beta1.AxoSyslog, got %T", axoSyslog)
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      commonAxoSyslogObjectValue + "-headless",
			Namespace: axoSyslog.Namespace,
			Labels: map[string]string{
				LabelAppName:      commonAxoSyslogObjectValue,
				LabelAppComponent: commonAxoSyslogObjectValue,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "otlp-grpc",
					Protocol:   corev1.ProtocolTCP,
					Port:       4317,
					TargetPort: intstr.IntOrString{IntVal: 4317},
				},
				{
					Name:       "otlp-http",
					Protocol:   corev1.ProtocolTCP,
					Port:       4318,
					TargetPort: intstr.IntOrString{IntVal: 4318},
				},
			},
			Selector: map[string]string{
				LabelAppName:      commonAxoSyslogObjectValue,
				LabelAppComponent: commonAxoSyslogObjectValue,
			},
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: corev1.ClusterIPNone,
		},
	}, reconciler.StatePresent, nil
}

func ServiceMetrics(object any) (runtime.Object, reconciler.DesiredState, error) {
	// TODO: Implement metrics service
	return &corev1.Service{}, reconciler.StateAbsent, nil
}

func ServiceBufferMetrics(object any) (runtime.Object, reconciler.DesiredState, error) {
	// TODO: Implement buffer metrics service
	return &corev1.Service{}, reconciler.StateAbsent, nil
}
