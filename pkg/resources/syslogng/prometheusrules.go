// Copyright © 2022 Cisco Systems, Inc. and/or its affiliates
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
	"fmt"

	"github.com/cisco-open/operator-tools/pkg/reconciler"
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"

	prometheus_operator "github.com/kube-logging/logging-operator/pkg/resources/prometheus-operator"
)

func (r *Reconciler) prometheusRules() (runtime.Object, reconciler.DesiredState, error) {
	obj := &v1.PrometheusRule{
		ObjectMeta: r.SyslogNGObjectMeta(ServiceName+"-metrics", ComponentSyslogNG),
	}
	state := reconciler.StateAbsent

	if r.syslogNGSpec.Metrics != nil && r.syslogNGSpec.Metrics.PrometheusRules {
		nsJobLabel := fmt.Sprintf(`job="%s", namespace="%s"`, obj.Name, obj.Namespace)
		state = reconciler.StatePresent
		const ruleGroupName = "syslog-ng"
		obj.Spec.Groups = []v1.RuleGroup{{
			Name: ruleGroupName,
			Rules: []v1.Rule{
				{
					Alert: "SyslogNGNodeDown",
					Expr:  intstr.FromString(fmt.Sprintf("up{%s} == 0", nsJobLabel)),
					For:   prometheus_operator.Duration("10m"),
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "syslog-ng",
						"severity":  "critical",
					},
					Annotations: map[string]string{
						"summary":     `Syslog-NG cannot be scraped`,
						"description": `Prometheus could not scrape "{{ $labels.job }}" for more than 30 minutes.`,
					},
				},
				{
					Alert: "SyslogNGQueueLength",
					Expr:  intstr.FromString(fmt.Sprintf("rate(syslog_ng_status_buffer_queue_length{%s}[5m]) > 0.3", nsJobLabel)),
					For:   prometheus_operator.Duration("1m"),
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "syslog-ng",
						"severity":  "warning",
					},
					Annotations: map[string]string{
						"summary":     `syslog-ng node are failing`,
						"description": `In the last 5 minutes, syslog-ng queues increased 30%. Current value is "{{ $value }}".`,
					},
				},
				{
					Alert: "SyslogNGQueueLength",
					Expr:  intstr.FromString(fmt.Sprintf("rate(syslog_ng_status_buffer_queue_length{%s}[5m]) > 0.5", nsJobLabel)),
					For:   prometheus_operator.Duration("1m"),
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "syslog-ng",
						"severity":  "critical",
					},
					Annotations: map[string]string{
						"summary":     `Syslog-NG nodes buffer queue length are critical`,
						"description": `In the last 5 minutes, Syslog-NG queues increased 50%. Current value is "{{ $value }}".`,
					},
				},
				{
					Alert: "SyslogNGRecordsCountsHigh",
					Expr:  intstr.FromString(fmt.Sprintf("sum(rate(syslog_ng_output_status_emit_records{%[1]s}[5m])) by (job,pod,namespace) > (3 * sum(rate(syslog_ng_output_status_emit_records{%[1]s}[15m])) by (job,pod,namespace))", nsJobLabel)),
					For:   prometheus_operator.Duration("1m"),
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "syslog-ng",
						"severity":  "critical",
					},
					Annotations: map[string]string{
						"summary":     `syslog-ng records count are critical`,
						"description": `In the last 5m, records counts increased 3 times, comparing to the latest 15 min.`,
					},
				},
				{
					Alert: "SyslogNGRetry",
					Expr:  intstr.FromString(fmt.Sprintf("increase(syslog_ng_status_retry_count{%s}[10m]) > 0", nsJobLabel)),
					For:   prometheus_operator.Duration("20m"),
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "syslog-ng",
						"severity":  "warning",
					},
					Annotations: map[string]string{
						"summary":     `Syslog-NG retry count has been "{{ $value }}" for the last 10 minutes.`,
						"description": `Syslog-NG retry count has been "{{ $value }}" for the last 10 minutes.`,
					},
				},
				{
					Alert: "SyslogNGOutputError",
					Expr:  intstr.FromString(fmt.Sprintf("increase(syslog_ng_output_status_num_errors{%s}[10m]) > 0", nsJobLabel)),
					For:   prometheus_operator.Duration("1s"),
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "syslog-ng",
						"severity":  "warning",
					},
					Annotations: map[string]string{
						"summary":     `There have been syslog-ng output error(s) for the last 10 minutes.`,
						"description": `Syslog-NG output error count is "{{ $value }}" for the last 10 minutes.`,
					},
				},
				{
					Alert: "SyslogNGPredictedBufferGrowth",
					Expr:  intstr.FromString(fmt.Sprintf("predict_linear(syslog_ng_output_status_buffer_total_bytes{%[1]s}[10m], 600) > syslog_ng_output_status_buffer_total_bytes{%[1]s}", nsJobLabel)),
					For:   prometheus_operator.Duration("10m"),
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "syslog-ng",
						"severity":  "warning",
					},
					Annotations: map[string]string{
						"summary":     `Syslog-NG buffer size prediction warning.`,
						"description": `Syslog-NG buffer trending watcher.`,
					},
				},
			},
		},
		}
	}
	return obj, state, nil
}
