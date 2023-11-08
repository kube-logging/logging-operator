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

func (r *Reconciler) bufferVolumePrometheusRules() (runtime.Object, reconciler.DesiredState, error) {
	obj := &v1.PrometheusRule{
		ObjectMeta: r.FluentbitObjectMeta(fluentbitServiceName + "-buffer-metrics"),
	}
	state := reconciler.StateAbsent

	if r.fluentbitSpec.BufferVolumeMetrics != nil && r.fluentbitSpec.BufferVolumeMetrics.PrometheusRules {
		nsJobLabel := fmt.Sprintf(`job="%s", namespace="%s"`, obj.Name, obj.Namespace)
		state = reconciler.StatePresent
		const ruleGroupName = "fluentbit-buffervolume"
		obj.Spec.Groups = []v1.RuleGroup{{
			Name: ruleGroupName,
			Rules: []v1.Rule{
				{
					Alert: "FluentbitBufferSize",
					Expr:  intstr.FromString(fmt.Sprintf(`node_filesystem_avail_bytes{mountpoint="/buffers", %[1]s} / node_filesystem_size_bytes{mountpoint="/buffers", %[1]s} * 100 < 10`, nsJobLabel)),
					For:   prometheus_operator.Duration("10m"),
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "fluentbit",
						"severity":  "warning",
					},
					Annotations: map[string]string{
						"summary":     `Fluentbit buffer free capacity less than 10%.`,
						"description": `Fluentbit buffer size capacity is {{ $value }}%.`,
					},
				},
				{
					Alert: "FluentbitBufferSize",
					Expr:  intstr.FromString(fmt.Sprintf(`node_filesystem_avail_bytes{mountpoint="/buffers", %[1]s} / node_filesystem_size_bytes{mountpoint="/buffers", %[1]s} * 100 < 5`, nsJobLabel)),
					For:   prometheus_operator.Duration("10m"),
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "fluentbit",
						"severity":  "critical",
					},
					Annotations: map[string]string{
						"summary":     `Fluentbit buffer free capacity less than 5%.`,
						"description": `Fluentbit buffer size capacity is {{ $value }}%.`,
					},
				},
			},
		},
		}
	}
	return obj, state, nil
}
