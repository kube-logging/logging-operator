// Copyright © 2024 Kube logging authors
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

package v1beta1

import (
	"reflect"
	"testing"

	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func intstrRef(val string) *intstr.IntOrString {
	x := intstr.FromString(val)
	return &x
}

func TestMerge(t *testing.T) {
	tests := []struct {
		name     string
		rule     v1.Rule
		override PrometheusRulesOverride
		expected v1.Rule
	}{
		{
			name: "SeverityOverride",
			rule: v1.Rule{
				Alert:  "TestAlert",
				Labels: map[string]string{"severity": "critical"},
			},
			override: PrometheusRulesOverride{
				Alert:  "TestAlert",
				Labels: map[string]string{"severity": "none"},
			},
			expected: v1.Rule{
				Alert:  "TestAlert",
				Labels: map[string]string{"severity": "none"},
			},
		},
		{
			name: "OverrideExpr",
			rule: v1.Rule{
				Alert:  "TestAlert",
				Labels: map[string]string{"severity": "critical"},
				Expr:   intstr.FromString("up > 0"),
			},
			override: PrometheusRulesOverride{
				Alert: "TestAlert",
				Expr:  intstrRef("up > 1"),
			},
			expected: v1.Rule{
				Alert:  "TestAlert",
				Labels: map[string]string{"severity": "critical"},
				Expr:   intstr.FromString("up > 1"),
			},
		},
	}

	for _, tt := range tests {
		ttp := tt
		t.Run(ttp.name, func(t *testing.T) {
			actual := *(ttp.override.Override(&ttp.rule))
			if !reflect.DeepEqual(actual, ttp.expected) {
				t.Fatalf("expected: %v, got: %v", ttp.expected, actual)
			}
		})
	}
}

func TestListMerge(t *testing.T) {
	tests := []struct {
		name          string
		rules         []v1.Rule
		overrides     []PrometheusRulesOverride
		expectedRules []v1.Rule
	}{
		{
			name: "Alert2CriticalToNone",
			rules: []v1.Rule{
				{
					Alert:  "TestAlert",
					Labels: map[string]string{"severity": "critical"},
				},
				{
					Alert:  "TestAlert2",
					Labels: map[string]string{"severity": "critical"},
				},
			},
			overrides: []PrometheusRulesOverride{
				{
					Alert:  "TestAlert2",
					Labels: map[string]string{"severity": "none"},
				},
			},
			expectedRules: []v1.Rule{
				{
					Alert:  "TestAlert",
					Labels: map[string]string{"severity": "critical"},
				},
				{
					Alert:  "TestAlert2",
					Labels: map[string]string{"severity": "none"},
				},
			},
		},
		{
			name: "OverrideAlert2Mismatch",
			rules: []v1.Rule{
				{
					Alert:  "TestAlert",
					Labels: map[string]string{"severity": "critical"},
				},
			},
			overrides: []PrometheusRulesOverride{
				{
					Alert:  "TestAlert2",
					Labels: map[string]string{"severity": "none"},
				},
			},
			expectedRules: []v1.Rule{
				{
					Alert:  "TestAlert",
					Labels: map[string]string{"severity": "critical"},
				},
			},
		},
		{
			name: "MultipleOverridesAppliedCorrectly",
			rules: []v1.Rule{
				{
					Alert: "FluentdRetry",
					Expr:  intstr.FromString("increase(fluentd_status_retry_count[10m]) > 5"),
				},
				{
					Alert: "FluentdOutputError",
					Expr:  intstr.FromString("increase(fluentd_output_status_num_errors[10m]) > 2"),
				},
			},
			overrides: []PrometheusRulesOverride{
				{
					Alert: "FluentdRetry",
					Expr:  intstrRef("increase(fluentd_status_retry_count[10m]) > 10"),
				},
				{
					Alert: "FluentdOutputError",
					Expr:  intstrRef("increase(fluentd_output_status_num_errors[10m]) > 5"),
				},
			},
			expectedRules: []v1.Rule{
				{
					Alert: "FluentdRetry",
					Expr:  intstr.FromString("increase(fluentd_status_retry_count[10m]) > 10"),
				},
				{
					Alert: "FluentdOutputError",
					Expr:  intstr.FromString("increase(fluentd_output_status_num_errors[10m]) > 5"),
				},
			},
		},
	}

	for _, tt := range tests {
		ttp := tt
		t.Run(ttp.name, func(t *testing.T) {
			actual := ttp.rules
			for _, o := range ttp.overrides {
				actual = o.ListOverride(actual)
			}

			if !reflect.DeepEqual(actual, ttp.expectedRules) {
				t.Fatalf("expected: %v, got: %v", ttp.expectedRules, actual)
			}
		})
	}
}
