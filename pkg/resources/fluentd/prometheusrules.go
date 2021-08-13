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
	if r.Logging.Spec.FluentdSpec.Metrics != nil && r.Logging.Spec.FluentdSpec.Metrics.PrometheusRules {
		objectMetadata := r.FluentdObjectMeta(ServiceName+"-metrics", ComponentFluentd)

		return &v1.PrometheusRule{
			ObjectMeta: objectMetadata,
			Spec: v1.PrometheusRuleSpec{
				Groups: []v1.RuleGroup{{
					Name: "fluentd",
					Rules: []v1.Rule{
						{
							Alert: "FluentdNodeDown",
							Expr: intstr.IntOrString{
								Type:   intstr.String,
								StrVal: fmt.Sprintf("up{job=\"%s\"} == 0", objectMetadata.Name),
							},
							For: "10m",
							Labels: map[string]string{
								"service":  "fluentd",
								"severity": "warning",
							},
							Annotations: map[string]string{
								"summary":     `fluentd cannot be scraped`,
								"description": `Prometheus could not scrape {{ "{{ $labels.job }}" }} for more than 10 minutes`,
							},
						},
						{
							Alert: "FluentdNodeDown",
							Expr: intstr.IntOrString{
								Type:   intstr.String,
								StrVal: fmt.Sprintf("up{job=\"%s\"} == 0", objectMetadata.Name),
							},
							For: "30m",
							Labels: map[string]string{
								"service":  "fluentd",
								"severity": "critical",
							},
							Annotations: map[string]string{
								"summary":     `fluentd cannot be scraped`,
								"description": `Prometheus could not scrape {{ "{{ $labels.job }}" }} for more than 30 minutes`,
							},
						},
						{
							Alert: "FluentdQueueLength",
							Expr: intstr.IntOrString{
								Type:   intstr.String,
								StrVal: `rate(fluentd_status_buffer_queue_length[5m]) > 0.3`,
							},
							For: "1m",
							Labels: map[string]string{
								"service":  "fluentd",
								"severity": "warning",
							},
							Annotations: map[string]string{
								"summary":     `fluentd node are failing`,
								"description": `In the last 5 minutes, fluentd queues increased 30%. Current value is {{ "{{ $value }}" }}`,
							},
						},
						{
							Alert: "FluentdQueueLength",
							Expr: intstr.IntOrString{
								Type:   intstr.String,
								StrVal: `rate(fluentd_status_buffer_queue_length[5m]) > 0.5`,
							},
							For: "1m",
							Labels: map[string]string{
								"service":  "fluentd",
								"severity": "critical",
							},
							Annotations: map[string]string{
								"summary":     `fluentd node are critical`,
								"description": `In the last 5 minutes, fluentd queues increased 50%. Current value is {{ "{{ $value }}" }}`,
							},
						},
						{
							Alert: "FluentdRecordsCountsHigh",
							Expr: intstr.IntOrString{
								Type: intstr.String,
								StrVal: `sum(rate(fluentd_output_status_emit_records{job="{{ tpl .Release.Name . }}"}[5m]))
      BY (instance) >  (3 * sum(rate(fluentd_output_status_emit_records{job="{{ tpl .Release.Name . }}"}[15m]))
      BY (instance))`,
							},
							For: "1m",
							Labels: map[string]string{
								"service":  "fluentd",
								"severity": "critical",
							},
							Annotations: map[string]string{
								"summary":     `fluentd records count are critical`,
								"description": `In the last 5m, records counts increased 3 times, comparing to the latest 15 min.`,
							},
						},
						{
							Alert: "FluentdRetry",
							Expr: intstr.IntOrString{
								Type:   intstr.String,
								StrVal: `increase(fluentd_status_retry_count[10m]) > 0`,
							},
							For: "20m",
							Labels: map[string]string{
								"service":  "fluentd",
								"severity": "warning",
							},
							Annotations: map[string]string{
								"summary":     `Fluentd retry count has been  {{ "{{ $value }}" }} for the last 10 minutes`,
								"description": `Fluentd retry count has been  {{ "{{ $value }}" }} for the last 10 minutes`,
							},
						},
						{
							Alert: "FluentdOutputError",
							Expr: intstr.IntOrString{
								Type:   intstr.String,
								StrVal: `increase(fluentd_output_status_num_errors[10m]) > 0`,
							},
							For: "1s",
							Labels: map[string]string{
								"service":  "fluentd",
								"severity": "warning",
							},
							Annotations: map[string]string{
								"summary":     `There have been Fluentd output error(s) for the last 10 minutes`,
								"description": `Fluentd output error count is {{ "{{ $value }}" }} for the last 10 minutes`,
							},
						},
						{
							Alert: "FluentdBufferSize",
							Expr: intstr.IntOrString{
								Type:   intstr.String,
								StrVal: `node_filesystem_avail_bytes{mountpoint="/buffers"} / node_filesystem_size_bytes{mountpoint="/buffers"} * 100 < 10`,
							},
							For: "10m",
							Labels: map[string]string{
								"service":  "fluentd",
								"severity": "warning",
							},
							Annotations: map[string]string{
								"summary":     `There have been Fluentd output error(s) for the last 10 minutes`,
								"description": `Fluentd buffer size capacity is {{ "{{ $value }}" }}% `,
							},
						},
						{
							Alert: "FluentdBufferSize",
							Expr: intstr.IntOrString{
								Type:   intstr.String,
								StrVal: `node_filesystem_avail_bytes{mountpoint="/buffers"} / node_filesystem_size_bytes{mountpoint="/buffers"} * 100 < 5`,
							},
							For: "10m",
							Labels: map[string]string{
								"service":  "fluentd",
								"severity": "critical",
							},
							Annotations: map[string]string{
								"summary":     `There have been Fluentd output error(s) for the last 10 minutes`,
								"description": `Fluentd buffer size capacity is {{ "{{ $value }}" }}% `,
							},
						},
					},
				},
				},
			},
		}, reconciler.StatePresent, nil
	}
	return &v1.PrometheusRule{
		ObjectMeta: r.FluentdObjectMeta(ServiceName, ComponentFluentd),
		Spec:       v1.PrometheusRuleSpec{},
	}, reconciler.StateAbsent, nil
}
