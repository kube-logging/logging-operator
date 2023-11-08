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
	"fmt"

	"github.com/cisco-open/operator-tools/pkg/reconciler"
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"

	prometheus_operator "github.com/kube-logging/logging-operator/pkg/resources/prometheus-operator"
)

func (r *Reconciler) prometheusRules() (runtime.Object, reconciler.DesiredState, error) {
	obj := &v1.PrometheusRule{
		ObjectMeta: r.FluentbitObjectMeta(fluentbitServiceName + "-monitor"),
	}
	state := reconciler.StateAbsent

	if r.fluentbitSpec.Metrics != nil && r.fluentbitSpec.Metrics.PrometheusRules {
		nsJobLabel := fmt.Sprintf(`job="%s", namespace="%s"`, obj.Name, obj.Namespace)
		state = reconciler.StatePresent
		obj.Spec.Groups = []v1.RuleGroup{{
			Name: "fluentbit",
			Rules: []v1.Rule{
				{
					Alert: "FluentbitTooManyErrors",
					Expr:  intstr.FromString(fmt.Sprintf("rate(fluentbit_output_retries_failed_total{%s}[10m]) > 0", nsJobLabel)),
					For:   prometheus_operator.Duration("10m"),
					Labels: map[string]string{
						"service":  "fluentbit",
						"severity": "warning",
					},
					Annotations: map[string]string{
						"summary":     `Fluentbit too many errors.`,
						"description": `Fluentbit ({{ $labels.instance }}) is erroring.`,
					},
				},
			},
		},
		}
	}
	return obj, state, nil
}
