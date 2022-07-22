// Copyright © 2020 Banzai Cloud
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

package config

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func checkError(t *testing.T, expected interface{}, actual error, msgAndArgs ...interface{}) {
	t.Helper()
	switch expected := expected.(type) {
	case nil:
		require.NoError(t, actual, msgAndArgs...)
	case bool:
		if expected {
			require.Error(t, actual, msgAndArgs...)
		} else {
			require.NoError(t, actual, msgAndArgs...)
		}
	case func(error) bool:
		require.True(t, expected(actual), msgAndArgs...)
	default:
		require.Equal(t, expected, actual, msgAndArgs...)
	}
}

var leadingTabs = regexp.MustCompile("(?m:^\t+)")

func untab(s string) string {
	return leadingTabs.ReplaceAllStringFunc(s, func(match string) string {
		return strings.Repeat("    ", len(match))
	})
}
