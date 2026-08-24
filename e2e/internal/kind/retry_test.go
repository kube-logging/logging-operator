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
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetryDelete(t *testing.T) {
	wedged := errors.New("could not kill: tried to kill container, but did not receive an exit event")

	// fails returns a delete that fails n times before succeeding, numbering its
	// errors so the caller can tell which attempt produced the one it kept.
	fails := func(n int, calls *int) func() error {
		return func() error {
			*calls++
			if *calls <= n {
				return fmt.Errorf("attempt %d: %w", *calls, wedged)
			}
			return nil
		}
	}

	t.Run("a delete that works is not repeated", func(t *testing.T) {
		var calls int
		k := &Kind{DeleteAttempts: 3}

		require.NoError(t, k.retryDelete("clean", fails(0, &calls)))
		assert.Equal(t, 1, calls)
	})

	t.Run("retries until the delete lands", func(t *testing.T) {
		var calls int
		k := &Kind{DeleteAttempts: 3}

		require.NoError(t, k.retryDelete("wedged", fails(2, &calls)))
		assert.Equal(t, 3, calls)
	})

	t.Run("gives up with the last error, not the first", func(t *testing.T) {
		var calls int
		k := &Kind{DeleteAttempts: 2}

		err := k.retryDelete("stuck", fails(99, &calls))

		require.ErrorIs(t, err, wedged)
		assert.Contains(t, err.Error(), "attempt 2")
		assert.Equal(t, 2, calls)
	})

	// Three attempts at the full command timeout would outlast the package
	// deadline the caller still has to report inside.
	t.Run("a timeout is not retried", func(t *testing.T) {
		var calls int
		k := &Kind{DeleteAttempts: 3}

		err := k.retryDelete("hung", func() error {
			calls++
			return fmt.Errorf("kind delete cluster %w after 16m", ErrTimeout)
		})

		require.ErrorIs(t, err, ErrTimeout)
		assert.Equal(t, 1, calls)
	})

	t.Run("a zero Kind still deletes once", func(t *testing.T) {
		var calls int
		k := &Kind{}

		require.Error(t, k.retryDelete("once", fails(99, &calls)))
		assert.Equal(t, 1, calls)
	})
}
