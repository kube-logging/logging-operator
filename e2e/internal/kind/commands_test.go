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

package kind

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandTimeoutIsResolvedOnce(t *testing.T) {
	t.Run("an explicit value is kept", func(t *testing.T) {
		k := &Kind{CommandTimeout: 42 * time.Second}

		assert.Equal(t, 42*time.Second, k.commandTimeout())
		assert.Equal(t, 42*time.Second, k.commandTimeout(), "resolving twice must not change it")
	})

	t.Run("zero is derived from the enclosing deadline", func(t *testing.T) {
		k := &Kind{}
		enclosing := enclosingTestTimeout()
		require.Positive(t, enclosing)

		resolved := k.commandTimeout()

		assert.Positive(t, resolved)
		assert.Less(t, resolved, enclosing)
		assert.Equal(t, resolved, k.commandTimeout(), "resolving twice must not change it")
	})
}

func TestDeriveCommandTimeout(t *testing.T) {
	testCases := map[string]struct {
		enclosing time.Duration
		expected  time.Duration
	}{
		"the suite's own 20m budget": {enclosing: 20 * time.Minute, expected: 16 * time.Minute},
		"a short budget":             {enclosing: 30 * time.Second, expected: 24 * time.Second},
		// -timeout 0 means the binary has no deadline, so there is nothing to
		// derive from and a fixed backstop is all that is left.
		"no deadline falls back": {enclosing: 0, expected: fallbackTimeout},
		"a negative one too":     {enclosing: -1, expected: fallbackTimeout},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, deriveCommandTimeout(testCase.enclosing))
		})
	}
}

// The invariant that matters. A cap at or above the deadline it sits inside
// would never fire in time to be useful: the binary would panic on its own
// deadline first, with nothing named and no cluster cleaned up.
func TestDerivedTimeoutStaysBelowTheEnclosingDeadline(t *testing.T) {
	for _, enclosing := range []time.Duration{
		30 * time.Second, 2 * time.Minute, 10 * time.Minute,
		20 * time.Minute, time.Hour,
	} {
		derived := deriveCommandTimeout(enclosing)
		assert.Less(t, derived, enclosing, "a %s budget derived %s", enclosing, derived)
		assert.Positive(t, derived)
	}
}

// Guards the plumbing: the value really does come from this binary's -timeout.
func TestEnclosingTestTimeoutIsReadable(t *testing.T) {
	enclosing := enclosingTestTimeout()
	t.Logf("this binary was started with -timeout %s", enclosing)
	assert.Positive(t, enclosing, "go test always sets a deadline unless -timeout 0")
}

func TestResolveCommandTimeout(t *testing.T) {
	fallback := 10 * time.Minute

	testCases := map[string]struct {
		raw      string
		expected time.Duration
	}{
		"unset keeps the default":         {raw: "", expected: fallback},
		"a duration is honored":           {raw: "90s", expected: 90 * time.Second},
		"a bare number is not a duration": {raw: "600", expected: fallback},
		"nonsense keeps the default":      {raw: "soon", expected: fallback},
		// A zero deadline has already expired, so honoring these would fail
		// every invocation instantly instead of disabling the timeout.
		"zero keeps the default":         {raw: "0", expected: fallback},
		"zero seconds keeps the default": {raw: "0s", expected: fallback},
		"negative keeps the default":     {raw: "-1s", expected: fallback},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, resolveCommandTimeout(testCase.raw, fallback))
		})
	}
}
