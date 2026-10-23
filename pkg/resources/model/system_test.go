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

package model

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/filter"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/plugins"
)

func TestCreateSystemDetectExceptionsMultiWorker(t *testing.T) {
	detectExceptions := []v1beta1.Filter{{DetectExceptions: &filter.DetectExceptions{}}}

	tests := []struct {
		name                 string
		workers              int32
		skipInvalidResources bool
		globalFilters        []v1beta1.Filter
		flowFilters          []v1beta1.Filter
		wantErr              bool
		wantFlows            []string
	}{
		{
			name:        "flow with detectExceptions and a single worker",
			workers:     1,
			flowFilters: detectExceptions,
			wantFlows:   []string{"flow:test:detect", "flow:test:plain"},
		},
		{
			name:        "flow with detectExceptions and multiple workers",
			workers:     2,
			flowFilters: detectExceptions,
			wantErr:     true,
		},
		{
			name:                 "flow with detectExceptions and multiple workers is skipped when invalid resources are skipped",
			workers:              2,
			skipInvalidResources: true,
			flowFilters:          detectExceptions,
			wantFlows:            []string{"flow:test:plain"},
		},
		{
			name:          "global detectExceptions filter and multiple workers",
			workers:       2,
			globalFilters: detectExceptions,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logging := loggingFor()
			logging.Spec.FluentdSpec = &v1beta1.FluentdSpec{Workers: tt.workers}
			logging.Spec.SkipInvalidResources = tt.skipInvalidResources
			logging.Spec.GlobalFilters = tt.globalFilters

			outputs := Outputs{{
				Name: "null", Namespace: "test",
				Spec: v1beta1.OutputSpec{NullOutputConfig: output.NewNullOutputConfig()},
			}}
			flows := []v1beta1.Flow{
				{
					Name: "detect", Namespace: "test",
					Spec: v1beta1.FlowSpec{Filters: tt.flowFilters, LocalOutputRefs: []string{"null"}},
				},
				{
					Name: "plain", Namespace: "test",
					Spec: v1beta1.FlowSpec{LocalOutputRefs: []string{"null"}},
				},
			}

			resources := LoggingResources{
				Logging: logging,
				Fluentd: FluentdLoggingResources{Flows: flows, Outputs: outputs},
			}

			secrets := stubSecretLoaderFactory{client: fake.NewClientBuilder().Build()}
			system, err := CreateSystem(resources, secrets, logr.Discard())
			if tt.wantErr {
				require.ErrorIs(t, err, plugins.ErrDetectExceptionsMultiWorker)
				return
			}
			require.NoError(t, err)

			var ids []string
			for _, flow := range system.Flows {
				// Skip the @ERROR flow that is always registered
				if flow.FlowLabel != "@ERROR" {
					ids = append(ids, flow.FlowID)
				}
			}
			require.ElementsMatch(t, tt.wantFlows, ids)
		})
	}
}
