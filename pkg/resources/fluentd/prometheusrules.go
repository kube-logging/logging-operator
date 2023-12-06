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
	"fmt"

	"github.com/cisco-open/operator-tools/pkg/reconciler"
	prometheus_operator "github.com/kube-logging/logging-operator/pkg/resources/prometheus-operator"
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *Reconciler) prometheusRules() (runtime.Object, reconciler.DesiredState, error) {
	obj := &v1.PrometheusRule{
		ObjectMeta: r.FluentdObjectMeta(ServiceName+"-metrics", ComponentFluentd),
	}
	state := reconciler.StateAbsent

	if r.fluentdSpec.Metrics != nil && r.fluentdSpec.Metrics.PrometheusRules {
		nsJobLabel := fmt.Sprintf(`job="%s", namespace="%s"`, obj.Name, obj.Namespace)
		state = reconciler.StatePresent
		const ruleGroupName = "fluentd"
		obj.Spec.Groups = []v1.RuleGroup{{
			Name: ruleGroupName,
			Rules: []v1.Rule{
				{
					Alert: "FluentdNodeDown",
					Expr:  intstr.FromString(fmt.Sprintf("up{%s} == 0", nsJobLabel)),
					For:   prometheus_operator.Duration("10m"),
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "fluentd",
						"severity":  "critical",
					},
					Annotations: map[string]string{
						"summary":     `fluentd cannot be scraped`,
						"description": `Prometheus could not scrape "{{ $labels.job }}" for more than 30 minutes.`,
					},
				},
				{
					Alert: "FluentdQueueLength",
					Expr:  intstr.FromString(fmt.Sprintf("rate(fluentd_status_buffer_queue_length{%s}[5m]) > 0.3", nsJobLabel)),
					For:   prometheus_operator.Duration("1m"),
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "fluentd",
						"severity":  "warning",
					},
					Annotations: map[string]string{
						"summary":     `fluentd node are failing`,
						"description": `In the last 5 minutes, fluentd queues increased 30%. Current value is "{{ $value }}".`,
					},
				},
				{
					Alert: "FluentdQueueLength",
					Expr:  intstr.FromString(fmt.Sprintf("rate(fluentd_status_buffer_queue_length{%s}[5m]) > 0.5", nsJobLabel)),
					For:   prometheus_operator.Duration("1m"),
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "fluentd",
						"severity":  "critical",
					},
					Annotations: map[string]string{
						"summary":     `fluentd nodes buffer queue length are critical`,
						"description": `In the last 5 minutes, fluentd queues increased 50%. Current value is "{{ $value }}".`,
					},
				},
				{
					Alert: "FluentdRecordsCountsHigh",
					Expr:  intstr.FromString(fmt.Sprintf("sum(rate(fluentd_output_status_emit_records{%[1]s}[5m])) by (job,pod,namespace) > (3 * sum(rate(fluentd_output_status_emit_records{%[1]s}[15m])) by (job,pod,namespace))", nsJobLabel)),
					For:   prometheus_operator.Duration("1m"),
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "fluentd",
						"severity":  "critical",
					},
					Annotations: map[string]string{
						"summary":     `fluentd records count are critical`,
						"description": `In the last 5m, records counts increased 3 times, comparing to the latest 15 min.`,
					},
				},
				{
					Alert: "FluentdRetry",
					Expr:  intstr.FromString(fmt.Sprintf("increase(fluentd_status_retry_count{%s}[10m]) > 0", nsJobLabel)),
					For:   prometheus_operator.Duration("20m"),
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "fluentd",
						"severity":  "warning",
					},
					Annotations: map[string]string{
						"summary":     `Fluentd retry count has been "{{ $value }}" for the last 10 minutes.`,
						"description": `Fluentd retry count has been "{{ $value }}" for the last 10 minutes.`,
					},
				},
				{
					Alert: "FluentdOutputError",
					Expr:  intstr.FromString(fmt.Sprintf("increase(fluentd_output_status_num_errors{%s}[10m]) > 0", nsJobLabel)),
					For:   prometheus_operator.Duration("1s"),
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "fluentd",
						"severity":  "warning",
					},
					Annotations: map[string]string{
						"summary":     `There have been Fluentd output error(s) for the last 10 minutes.`,
						"description": `Fluentd output error count is "{{ $value }}" for the last 10 minutes.`,
					},
				},
				{
					Alert: "FluentdPredictedBufferGrowth",
					Expr:  intstr.FromString(fmt.Sprintf("sum(predict_linear(fluentd_output_status_buffer_total_bytes{%[1]s}[10m], 600)) > sum(fluentd_output_status_buffer_total_bytes{%[1]s}) * 1.5", nsJobLabel)),
					For:   prometheus_operator.Duration("10m"),
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "fluentd",
						"severity":  "warning",
					},
					Annotations: map[string]string{
						"summary":     `Fluentd buffer size is predicted to increase more than 50% in the next 10 minutes.`,
						"description": `Fluentd buffer trending watcher.`,
					},
				},
			},
		},
		}
	}
	return obj, state, nil
}
