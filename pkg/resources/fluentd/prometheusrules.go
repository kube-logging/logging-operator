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

	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (r *Reconciler) prometheusRules() (runtime.Object, reconciler.DesiredState, error) {
	obj := &v1.PrometheusRule{
		ObjectMeta: r.FluentdObjectMeta(ServiceName+"-metrics", ComponentFluentd),
	}
	state := reconciler.StateAbsent

	if r.Logging.Spec.FluentdSpec.Metrics != nil && r.Logging.Spec.FluentdSpec.Metrics.PrometheusRules {
		nsJobLabel := fmt.Sprintf(`job="%s", namespace="%s"`, obj.Name, obj.Namespace)
		state = reconciler.StatePresent
		const ruleGroupName = "fluentd"
		obj.Spec.Groups = []v1.RuleGroup{{
			Name: ruleGroupName,
			Rules: []v1.Rule{
				{
					Alert: "FluentdNodeDown",
					Expr:  intstr.FromString(fmt.Sprintf("up{%s} == 0", nsJobLabel)),
					For:   "10m",
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "fluentd",
						"severity":  "critical",
					},
					Annotations: map[string]string{
						"summary":     `fluentd cannot be scraped`,
						"description": `Prometheus could not scrape {{ "{{ $labels.job }}" }} for more than 30 minutes`,
					},
				},
				{
					Alert: "FluentdQueueLength",
					Expr:  intstr.FromString(fmt.Sprintf("rate(fluentd_status_buffer_queue_length{%s}[5m]) > 0.3", nsJobLabel)),
					For:   "1m",
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "fluentd",
						"severity":  "warning",
					},
					Annotations: map[string]string{
						"summary":     `fluentd node are failing`,
						"description": `In the last 5 minutes, fluentd queues increased 30%. Current value is {{ "{{ $value }}" }}`,
					},
				},
				{
					Alert: "FluentdQueueLength",
					Expr:  intstr.FromString(fmt.Sprintf("rate(fluentd_status_buffer_queue_length{%s}[5m]) > 0.5", nsJobLabel)),
					For:   "1m",
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "fluentd",
						"severity":  "critical",
					},
					Annotations: map[string]string{
						"summary":     `fluentd nodes buffer queue length are critical`,
						"description": `In the last 5 minutes, fluentd queues increased 50%. Current value is {{ "{{ $value }}" }}`,
					},
				},
				{
					Alert: "FluentdRecordsCountsHigh",
					Expr:  intstr.FromString(fmt.Sprintf("sum(rate(fluentd_output_status_emit_records{%[1]s}[5m])) by (job,pod,namespace) > (3 * sum(rate(fluentd_output_status_emit_records{%[1]s}[15m])) by (job,pod,namespace))", nsJobLabel)),
					For:   "1m",
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
					For:   "20m",
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "fluentd",
						"severity":  "warning",
					},
					Annotations: map[string]string{
						"summary":     `Fluentd retry count has been  {{ "{{ $value }}" }} for the last 10 minutes`,
						"description": `Fluentd retry count has been  {{ "{{ $value }}" }} for the last 10 minutes`,
					},
				},
				{
					Alert: "FluentdOutputError",
					Expr:  intstr.FromString(fmt.Sprintf("increase(fluentd_output_status_num_errors{%s}[10m]) > 0", nsJobLabel)),
					For:   "1s",
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "fluentd",
						"severity":  "warning",
					},
					Annotations: map[string]string{
						"summary":     `There have been Fluentd output error(s) for the last 10 minutes`,
						"description": `Fluentd output error count is {{ "{{ $value }}" }} for the last 10 minutes`,
					},
				},
				{
					Alert: "FluentdBufferSize",
					Expr:  intstr.FromString(fmt.Sprintf("node_filesystem_avail_bytes{mountpoint=\"/buffers\", %s} / node_filesystem_size_bytes{mountpoint=\"/buffers\", %s} * 100 < 10", nsJobLabel, nsJobLabel)),
					For:   "10m",
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "fluentd",
						"severity":  "warning",
					},
					Annotations: map[string]string{
						"summary":     `Fluentd buffer free capacity less than 10%.`,
						"description": `Fluentd buffer size capacity is {{ "{{ $value }}" }}% `,
					},
				},
				{
					Alert: "FluentdBufferSize",
					Expr:  intstr.FromString(fmt.Sprintf("node_filesystem_avail_bytes{mountpoint=\"/buffers\", %s} / node_filesystem_size_bytes{mountpoint=\"/buffers\" ,%s} * 100 < 5", nsJobLabel, nsJobLabel)),
					For:   "10m",
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "fluentd",
						"severity":  "critical",
					},
					Annotations: map[string]string{
						"summary":     `Fluentd buffer free capacity less than 5%`,
						"description": `Fluentd buffer size capacity is {{ "{{ $value }}" }}% `,
					},
				},
				{
					Alert: "FluentdPredictedBufferGrowth",
					Expr:  intstr.FromString(fmt.Sprintf("predict_linear(fluentd_output_status_buffer_total_bytes{%s}[10m], 600) > fluentd_output_status_buffer_total_bytes{%s}", nsJobLabel, nsJobLabel)),
					For:   "10m",
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "fluentd",
						"severity":  "warning",
					},
					Annotations: map[string]string{
						"summary":     `Fluentd buffer size prediction warning`,
						"description": `Fluentd buffer trending watcher`,
					},
				},
			},
		},
		}
	}
	return obj, state, nil
}
