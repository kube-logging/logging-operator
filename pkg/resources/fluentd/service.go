/*
 * Copyright © 2019 Banzai Cloud
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package fluentd

import (
	"github.com/banzaicloud/logging-operator/pkg/resources/templates"
	"github.com/banzaicloud/logging-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *Reconciler) service() runtime.Object {
	return &corev1.Service{
		ObjectMeta: templates.FluentdObjectMeta(serviceName, util.MergeLabels(r.Fluentd.Labels, labelSelector), r.Fluentd),
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       24240,
					TargetPort: intstr.IntOrString{IntVal: 24240},
				},
			},
			Selector: labelSelector,
			Type:     "ClusterIP",
		},
	}
}
