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

func (r *Reconciler) bufferVolumePrometheusRules() (runtime.Object, reconciler.DesiredState, error) {
	obj := &v1.PrometheusRule{
		ObjectMeta: r.SyslogNGObjectMeta(ServiceName+"-buffer-metrics", ComponentSyslogNG),
	}
	state := reconciler.StateAbsent

	if r.syslogNGSpec.BufferVolumeMetrics != nil && r.syslogNGSpec.BufferVolumeMetrics.PrometheusRules {
		nsJobLabel := fmt.Sprintf(`job="%s", namespace="%s"`, obj.Name, obj.Namespace)
		state = reconciler.StatePresent
		const ruleGroupName = "syslog-ng-buffervolume"
		obj.Spec.Groups = []v1.RuleGroup{{
			Name: ruleGroupName,
			Rules: []v1.Rule{
				{
					Alert: "SyslogNGBufferSize",
					Expr:  intstr.FromString(fmt.Sprintf(`node_filesystem_avail_bytes{mountpoint="/buffers", %[1]s} / node_filesystem_size_bytes{mountpoint="/buffers", %[1]s} * 100 < 10`, nsJobLabel)),
					For:   prometheus_operator.Duration("10m"),
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "syslog-ng",
						"severity":  "warning",
					},
					Annotations: map[string]string{
						"summary":     `Syslog-NG buffer free capacity less than 10%.`,
						"description": `Syslog-NG buffer size capacity is {{ $value }}%.`,
					},
				},
				{
					Alert: "SyslogNGBufferSize",
					Expr:  intstr.FromString(fmt.Sprintf(`node_filesystem_avail_bytes{mountpoint="/buffers", %[1]s} / node_filesystem_size_bytes{mountpoint="/buffers", %[1]s} * 100 < 5`, nsJobLabel)),
					For:   prometheus_operator.Duration("10m"),
					Labels: map[string]string{
						"rulegroup": ruleGroupName,
						"service":   "syslog-ng",
						"severity":  "critical",
					},
					Annotations: map[string]string{
						"summary":     `Syslog-NG buffer free capacity less than 5%.`,
						"description": `Syslog-NG buffer size capacity is {{ $value }}%.`,
					},
				},
			},
		},
		}
	}
	return obj, state, nil
}
