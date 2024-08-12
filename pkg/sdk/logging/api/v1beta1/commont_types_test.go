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

var overrideTests = []struct {
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
		name: "OverrideAlert2Mismatch",
		rule: v1.Rule{
			Alert:  "TestAlert",
			Labels: map[string]string{"severity": "critical"},
		},
		override: PrometheusRulesOverride{
			Alert:  "TestAlert2",
			Labels: map[string]string{"severity": "none"},
		},
		expected: v1.Rule{
			Alert:  "TestAlert",
			Labels: map[string]string{"severity": "critical"},
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

var overrideListTests = []struct {
	name          string
	rules         []v1.Rule
	override      PrometheusRulesOverride
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
		override: PrometheusRulesOverride{
			Alert:  "TestAlert2",
			Labels: map[string]string{"severity": "none"},
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
}

func TestMerge(t *testing.T) {
	for _, tt := range overrideTests {
		t.Run(tt.name, func(t *testing.T) {
			actual := *(tt.override.Override(&tt.rule))
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Fatalf("%v is not equal to %v", actual, tt.expected)
			}
		})
	}
}

func TestListMerge(t *testing.T) {
	for _, tt := range overrideListTests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.override.ListOverride(tt.rules)
			if !reflect.DeepEqual(actual, tt.expectedRules) {
				t.Fatalf("%v is not equal to %v", actual, tt.expectedRules)
			}
		})
	}
}
