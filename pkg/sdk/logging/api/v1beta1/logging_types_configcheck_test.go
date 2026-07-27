// Copyright © 2026 Kube logging authors
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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregatorLevelConfigCheck(t *testing.T) {
	loggingLevel := ConfigCheck{
		Strategy:       ConfigCheckStrategyTimeout,
		TimeoutSeconds: 42,
		Labels:         map[string]string{"from": "logging"},
	}

	tests := []struct {
		name       string
		aggregator *ConfigCheck
		expected   ConfigCheck
	}{
		{
			// Nothing set on the aggregator, so the logging level settings stand.
			name:       "NilKeepsLoggingLevel",
			aggregator: nil,
			expected:   loggingLevel,
		},
		{
			// An empty struct still means "override nothing".
			name:       "EmptyKeepsLoggingLevel",
			aggregator: &ConfigCheck{},
			expected:   loggingLevel,
		},
		{
			// The regression: overriding only the strategy used to drop the
			// timeout and the labels set on the logging resource.
			name:       "StrategyOnlyKeepsTheRest",
			aggregator: &ConfigCheck{Strategy: ConfigCheckStrategyDryRun},
			expected: ConfigCheck{
				Strategy:       ConfigCheckStrategyDryRun,
				TimeoutSeconds: 42,
				Labels:         map[string]string{"from": "logging"},
			},
		},
		{
			name:       "TimeoutOnlyKeepsTheRest",
			aggregator: &ConfigCheck{TimeoutSeconds: 7},
			expected: ConfigCheck{
				Strategy:       ConfigCheckStrategyTimeout,
				TimeoutSeconds: 7,
				Labels:         map[string]string{"from": "logging"},
			},
		},
		{
			name:       "LabelsOnlyKeepsTheRest",
			aggregator: &ConfigCheck{Labels: map[string]string{"from": "aggregator"}},
			expected: ConfigCheck{
				Strategy:       ConfigCheckStrategyTimeout,
				TimeoutSeconds: 42,
				Labels:         map[string]string{"from": "aggregator"},
			},
		},
		{
			name: "AllFieldsOverride",
			aggregator: &ConfigCheck{
				Strategy:       ConfigCheckStrategyDryRun,
				TimeoutSeconds: 7,
				Labels:         map[string]string{"from": "aggregator"},
			},
			expected: ConfigCheck{
				Strategy:       ConfigCheckStrategyDryRun,
				TimeoutSeconds: 7,
				Labels:         map[string]string{"from": "aggregator"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logging{}
			l.Spec.ConfigCheck = *loggingLevel.DeepCopy()

			l.AggregatorLevelConfigCheck(tt.aggregator)

			assert.Equal(t, tt.expected, l.Spec.ConfigCheck)
		})
	}
}

// TestAggregatorLevelConfigCheckDefaultsTimeout keeps the defaulting behavior
// intact: when neither level sets a timeout it still falls back to 10.
func TestAggregatorLevelConfigCheckDefaultsTimeout(t *testing.T) {
	l := &Logging{}
	l.Spec.ConfigCheck = ConfigCheck{Strategy: ConfigCheckStrategyTimeout}

	l.AggregatorLevelConfigCheck(&ConfigCheck{Strategy: ConfigCheckStrategyDryRun})

	assert.Equal(t, ConfigCheckStrategyDryRun, l.Spec.ConfigCheck.Strategy)
	assert.Equal(t, 10, l.Spec.ConfigCheck.TimeoutSeconds)
}
