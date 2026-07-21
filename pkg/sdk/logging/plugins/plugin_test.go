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

func TestRawFilterIsDisabled(t *testing.T) {
	plugins.EnableRawFilter = false
	filter := v1beta1.Filter{
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
	_, err := plugins.CreateFilter(filter, "test", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "raw filter is disabled")

	filter = v1beta1.Filter{
		StdOut: &modelfilter.StdOutFilterConfig{
			OutputType: "json",
		},
	}
	_, err = plugins.CreateFilter(filter, "test", nil)
	require.NoError(t, err)
}

func TestRawFilterIsEnabled(t *testing.T) {
	plugins.EnableRawFilter = true
	filter := v1beta1.Filter{
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
	_, err := plugins.CreateFilter(filter, "test", nil)
	require.NoError(t, err)
}
