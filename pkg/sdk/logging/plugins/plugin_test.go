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

package plugins_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	modelfilter "github.com/kube-logging/logging-operator/pkg/sdk/logging/model/filter"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/plugins"
)

func rawFilter() v1beta1.Filter {
	return v1beta1.Filter{
		Raw: &modelfilter.Raw{
			Config: `
@type my_filter
<my_section>
  foo bar
  tags ["web", "api", "db"]
</my_section>
			`,
		},
	}
}

func TestCreateFilterRawGate(t *testing.T) {
	tests := []struct {
		name    string
		filter  v1beta1.Filter
		options *plugins.CreateFilterOptions
		wantErr string
	}{
		{
			name:    "raw filter rejected without options",
			filter:  rawFilter(),
			options: nil,
			wantErr: "raw filter is disabled",
		},
		{
			name:    "raw filter rejected when disabled",
			filter:  rawFilter(),
			options: &plugins.CreateFilterOptions{RawFilterEnabled: false},
			wantErr: "raw filter is disabled",
		},
		{
			name:    "raw filter allowed when enabled",
			filter:  rawFilter(),
			options: &plugins.CreateFilterOptions{RawFilterEnabled: true},
		},
		{
			name: "non-raw filter unaffected without options",
			filter: v1beta1.Filter{
				StdOut: &modelfilter.StdOutFilterConfig{
					OutputType: "json",
				},
			},
			options: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := plugins.CreateFilterWithOptions(tt.filter, "test", nil, tt.options)
			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestCreateFilterDefaultsToRawDisabled(t *testing.T) {
	_, err := plugins.CreateFilter(rawFilter(), "test", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "raw filter is disabled")
}
