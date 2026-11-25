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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/input"
)

// fluent-plugin-prometheus makes worker N listen on the configured metrics port + N,
// so every one of those ports has to stay clear of the forward input.
func TestFluentdSpecSetDefaultsWorkerMetricsPorts(t *testing.T) {
	tests := []struct {
		name             string
		spec             FluentdSpec
		inputPort        int32
		wantMetricsPort  int32
		wantErrorMessage string
	}{
		{
			name:            "single worker keeps the default",
			spec:            FluentdSpec{Metrics: &Metrics{}},
			inputPort:       24240,
			wantMetricsPort: 24231,
		},
		{
			name:            "range ending right below the input port keeps the default",
			spec:            FluentdSpec{Workers: 9, Metrics: &Metrics{}},
			inputPort:       24240,
			wantMetricsPort: 24231,
		},
		{
			name:            "range reaching the input port moves above it",
			spec:            FluentdSpec{Workers: 10, Metrics: &Metrics{}},
			inputPort:       24240,
			wantMetricsPort: 24241,
		},
		{
			name:            "worker count reported in issue 1904",
			spec:            FluentdSpec{Workers: 14, Metrics: &Metrics{}},
			inputPort:       24240,
			wantMetricsPort: 24241,
		},
		{
			name: "custom forward input port inside the range",
			spec: FluentdSpec{
				Workers:            5,
				Metrics:            &Metrics{},
				ForwardInputConfig: &input.ForwardInputConfig{Port: "24233"},
			},
			inputPort:       24233,
			wantMetricsPort: 24241,
		},
		{
			name: "custom forward input port above the default input port",
			spec: FluentdSpec{
				Workers:            14,
				Metrics:            &Metrics{},
				ForwardInputConfig: &input.ForwardInputConfig{Port: "24243"},
			},
			inputPort:       24243,
			wantMetricsPort: 24244,
		},
		{
			name:            "explicit non-overlapping port is preserved",
			spec:            FluentdSpec{Workers: 14, Metrics: &Metrics{Port: 24220}},
			inputPort:       24240,
			wantMetricsPort: 24220,
		},
		{
			name:             "explicit overlapping port is rejected",
			spec:             FluentdSpec{Workers: 14, Metrics: &Metrics{Port: 24231}},
			inputPort:        24240,
			wantErrorMessage: "metrics ports 24231-24244 of 14 fluentd worker(s) overlap with port 24240",
		},
		{
			name:             "explicit port equal to the input port is rejected",
			spec:             FluentdSpec{Metrics: &Metrics{Port: 24240}},
			inputPort:        24240,
			wantErrorMessage: "metrics ports 24240-24240 of 1 fluentd worker(s) overlap with port 24240",
		},
		{
			name:            "overlap is ignored while metrics are disabled",
			spec:            FluentdSpec{Workers: 14, Metrics: &Metrics{Port: 24231, Enabled: new(false)}},
			inputPort:       24240,
			wantMetricsPort: 24231,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := tt.spec.DeepCopy()

			err := spec.SetDefaults()
			if tt.wantErrorMessage != "" {
				require.ErrorContains(t, err, tt.wantErrorMessage)
				return
			}
			require.NoError(t, err)

			if spec.Metrics.IsEnabled() {
				var workerPorts []int32
				for worker := range spec.Workers {
					workerPorts = append(workerPorts, spec.Metrics.Port+worker)
				}
				assert.NotContains(t, workerPorts, tt.inputPort, "a worker's metrics port collides with the forward input")
			}
			assert.Equal(t, tt.wantMetricsPort, spec.Metrics.Port)

			require.NoError(t, spec.SetDefaults(), "defaulting twice must be stable")
			assert.Equal(t, tt.wantMetricsPort, spec.Metrics.Port)
		})
	}
}

func TestFluentdSpecSetDefaultsMetricsAnnotationFollowsPort(t *testing.T) {
	spec := &FluentdSpec{Workers: 14, Metrics: &Metrics{PrometheusAnnotations: true}}

	require.NoError(t, spec.SetDefaults())

	assert.Equal(t, fmt.Sprint(spec.Metrics.Port), spec.Annotations["prometheus.io/port"])
	assert.Equal(t, "24241", spec.Annotations["prometheus.io/port"])
}
