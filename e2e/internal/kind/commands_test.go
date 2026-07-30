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
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeKind stands in for the kind binary: it records every invocation and then
// behaves as the FAKE_KIND_* variables ask, so the stall paths can be exercised
// without building a cluster.
//
// The sleep drops its inherited output because killing this script leaves the
// sleep running, and an orphan holding the test binary's pipes would keep
// `go test` waiting long after the run finished.
const fakeKind = `#!/bin/sh
echo "$*" >> "$FAKE_KIND_LOG"
stall() {
	if [ "${1:-0}" = "0" ]; then
		return
	fi
	sleep "$1" >/dev/null 2>&1
}
case "$1" in
create) stall "${FAKE_KIND_CREATE_SLEEP:-0}" ;;
delete) stall "${FAKE_KIND_DELETE_SLEEP:-0}" ;;
load)   stall "${FAKE_KIND_LOAD_SLEEP:-0}" ;;
get)
	stall "${FAKE_KIND_GET_SLEEP:-0}"
	printf '%s' "${FAKE_KIND_KUBECONFIG:-}"
	;;
esac
if [ -n "${FAKE_KIND_STDERR:-}" ]; then
	echo "$FAKE_KIND_STDERR" >&2
fi
exit "${FAKE_KIND_EXIT:-0}"
`

const (
	// shortTimeout is the deadline used by the tests that stall the stub.
	shortTimeout = 200 * time.Millisecond

	// stubStall is how long the stub sleeps when a test wants it to hang. It
	// only has to outlast shortTimeout: the sleep is orphaned when the stub is
	// killed, so a longer one would just linger on the runner.
	stubStall = "5"
)

// newFakeKind returns a Kind pointed at the stub, and the path of its
// invocation log.
func newFakeKind(t *testing.T) (*Kind, string) {
	t.Helper()

	dir := t.TempDir()
	binary := filepath.Join(dir, "kind")
	require.NoError(t, os.WriteFile(binary, []byte(fakeKind), 0o700))

	logPath := filepath.Join(dir, "invocations.log")
	t.Setenv("FAKE_KIND_LOG", logPath)

	return &Kind{
		Path:           binary,
		CommandTimeout: time.Minute,
		CleanupTimeout: time.Minute,
		// The stub leaves a sleep holding the output pipe, so the production
		// delay would be paid by every test that stalls.
		waitDelay: 250 * time.Millisecond,
	}, logPath
}

func invocations(t *testing.T, logPath string) []string {
	t.Helper()

	contents, err := os.ReadFile(logPath)
	if os.IsNotExist(err) {
		return nil
	}
	require.NoError(t, err)

	return strings.Split(strings.TrimSpace(string(contents)), "\n")
}

func TestCommandsTimeOut(t *testing.T) {
	testCases := map[string]struct {
		sleepEnv string
		invoke   func(*Kind) error
		names    string
	}{
		"create cluster": {
			sleepEnv: "FAKE_KIND_CREATE_SLEEP",
			invoke:   func(k *Kind) error { return k.CreateCluster(CreateClusterOptions{Name: "stuck"}) },
			names:    "create cluster --name stuck",
		},
		"delete cluster": {
			sleepEnv: "FAKE_KIND_DELETE_SLEEP",
			invoke:   func(k *Kind) error { return k.DeleteCluster(DeleteClusterOptions{Name: "stuck"}) },
			names:    "delete cluster --name stuck",
		},
		"load docker-image": {
			sleepEnv: "FAKE_KIND_LOAD_SLEEP",
			invoke: func(k *Kind) error {
				return k.LoadDockerImage([]string{"some/image:latest"}, LoadDockerImageOptions{Name: "stuck"})
			},
			names: "load docker-image --name stuck some/image:latest",
		},
		"get kubeconfig": {
			sleepEnv: "FAKE_KIND_GET_SLEEP",
			invoke: func(k *Kind) error {
				_, err := k.GetKubeconfig(GetKubeconfigOptions{Name: "stuck"})
				return err
			},
			names: "get kubeconfig --name stuck",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			k, _ := newFakeKind(t)
			k.CommandTimeout = shortTimeout
			t.Setenv(testCase.sleepEnv, stubStall)

			started := time.Now()
			err := testCase.invoke(k)

			require.Error(t, err)
			assert.ErrorIs(t, err, ErrTimeout)
			// "signal: killed" alone would not say what was stuck.
			assert.Contains(t, err.Error(), testCase.names)
			assert.Contains(t, err.Error(), commandTimeoutEnv)
			// The point of the change: the caller finds out promptly instead of
			// hanging until go test panics on its own -timeout.
			assert.Less(t, time.Since(started), 30*time.Second)
		})
	}
}

func TestCreateClusterErrorReporting(t *testing.T) {
	const alreadyExists = "ERROR: failed to create cluster: node(s) already exist for a cluster with the name"

	testCases := map[string]struct {
		env         map[string]string
		missingKind bool
		timeouts    time.Duration
		isTimeout   bool
		contains    string
	}{
		// Negative control: a failure that is not a stall keeps reporting kind's
		// own stderr rather than being relabelled as a timeout.
		"kind's own diagnosis is passed through": {
			env:      map[string]string{"FAKE_KIND_EXIT": "1", "FAKE_KIND_STDERR": alreadyExists},
			contains: "failed to create cluster: node(s) already exist for a cluster with the name",
		},
		"the exec error stands in when kind writes nothing": {
			missingKind: true,
			contains:    "kind-is-not-installed-here",
		},
		"a failed cleanup is reported alongside the timeout": {
			env:       map[string]string{"FAKE_KIND_CREATE_SLEEP": stubStall, "FAKE_KIND_DELETE_SLEEP": stubStall},
			timeouts:  shortTimeout,
			isTimeout: true,
			contains:  "deleting the partial cluster also failed",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			k, _ := newFakeKind(t)
			if testCase.timeouts > 0 {
				k.CommandTimeout, k.CleanupTimeout = testCase.timeouts, testCase.timeouts
			}
			if testCase.missingKind {
				k.Path = filepath.Join(t.TempDir(), "kind-is-not-installed-here")
			}
			for key, value := range testCase.env {
				t.Setenv(key, value)
			}

			err := k.CreateCluster(CreateClusterOptions{Name: "cluster"})

			require.Error(t, err)
			assert.NotEmpty(t, err.Error(), "the error must say something")
			assert.Contains(t, err.Error(), testCase.contains)
			if testCase.isTimeout {
				assert.ErrorIs(t, err, ErrTimeout)
			} else {
				assert.NotErrorIs(t, err, ErrTimeout)
			}
		})
	}
}

func TestInvocations(t *testing.T) {
	testCases := map[string]struct {
		sleepEnv string
		timeout  time.Duration
		invoke   func(*Kind) error
		wantErr  bool
		want     []string
	}{
		"a successful create is the only call": {
			invoke:  func(k *Kind) error { return k.CreateCluster(CreateClusterOptions{Name: "fine"}) },
			timeout: time.Minute,
			want:    []string{"create cluster --name fine"},
		},
		// kind tears down a half-built cluster when one of its own actions
		// fails, but not when we kill it, so the delete has to be ours.
		"a stalled create deletes the partial cluster": {
			sleepEnv: "FAKE_KIND_CREATE_SLEEP",
			timeout:  shortTimeout,
			invoke:   func(k *Kind) error { return k.CreateCluster(CreateClusterOptions{Name: "stuck"}) },
			wantErr:  true,
			want:     []string{"create cluster --name stuck", "delete cluster --name stuck"},
		},
		"an empty image list runs nothing": {
			invoke:  func(k *Kind) error { return k.LoadDockerImage(nil, LoadDockerImageOptions{Name: "c"}) },
			timeout: time.Minute,
			want:    nil,
		},
		// One call for every image, so docker save writes shared layers once.
		"several images load in a single call": {
			invoke: func(k *Kind) error {
				return k.LoadDockerImage([]string{"a:local", "b:local", "c:local"}, LoadDockerImageOptions{Name: "c"})
			},
			timeout: time.Minute,
			want:    []string{"load docker-image --name c a:local b:local c:local"},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			k, logPath := newFakeKind(t)
			k.CommandTimeout = testCase.timeout
			if testCase.sleepEnv != "" {
				t.Setenv(testCase.sleepEnv, stubStall)
			}

			err := testCase.invoke(k)

			if testCase.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, testCase.want, invocations(t, logPath))
		})
	}
}

func TestCreateClusterAppliesTheDefaultImage(t *testing.T) {
	testCases := map[string]struct {
		image        string
		optionsImage string
		want         string
	}{
		"the default image is applied": {
			image: "kindest/node:v1.30.0",
			want:  "create cluster --image kindest/node:v1.30.0 --name c",
		},
		"an explicit image wins": {
			image:        "kindest/node:v1.30.0",
			optionsImage: "kindest/node:v1.31.0",
			want:         "create cluster --image kindest/node:v1.31.0 --name c",
		},
		"no image is passed when none is set": {
			want: "create cluster --name c",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			k, logPath := newFakeKind(t)
			k.Image = testCase.image

			require.NoError(t, k.CreateCluster(CreateClusterOptions{Name: "c", Image: testCase.optionsImage}))

			assert.Equal(t, []string{testCase.want}, invocations(t, logPath))
		})
	}
}

func TestGetKubeconfigReturnsStdout(t *testing.T) {
	k, _ := newFakeKind(t)
	t.Setenv("FAKE_KIND_KUBECONFIG", "apiVersion: v1")

	kubeconfig, err := k.GetKubeconfig(GetKubeconfigOptions{Name: "cluster"})

	require.NoError(t, err)
	assert.Equal(t, "apiVersion: v1", string(kubeconfig))
}

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
