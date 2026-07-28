// Copyright © 2021 Cisco Systems, Inc. and/or its affiliates
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

// fakeKind stands in for the kind binary. It appends every invocation to
// $FAKE_KIND_LOG and then behaves as the remaining variables ask, which lets
// these tests exercise the stall path without building a real cluster.
//
// Killing this script leaves its `sleep` running, so the sleep drops the
// inherited output: otherwise the orphan would hold the test binary's own
// pipes open and `go test` would sit on them long after the run finished.
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

// installFakeKind points KindPath at the stub for the duration of the test and
// returns the path of its invocation log.
func installFakeKind(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	binary := filepath.Join(dir, "kind")
	require.NoError(t, os.WriteFile(binary, []byte(fakeKind), 0o700))

	logPath := filepath.Join(dir, "invocations.log")
	t.Setenv("FAKE_KIND_LOG", logPath)

	originalPath, originalImage := KindPath, KindImage
	// KindImage would otherwise leak into the recorded argument lists.
	KindPath, KindImage = binary, ""
	t.Cleanup(func() { KindPath, KindImage = originalPath, originalImage })

	return logPath
}

// withTimeouts shortens the deadlines so a stall surfaces in test time. The
// wait delay comes down too: the stub leaves a `sleep` holding the output
// pipe, so the production 10s would be paid on every stalling test.
func withTimeouts(t *testing.T, command, cleanup time.Duration) {
	t.Helper()

	originalCommand, originalCleanup, originalDelay := CommandTimeout, CleanupTimeout, waitDelay
	CommandTimeout, CleanupTimeout, waitDelay = command, cleanup, 250*time.Millisecond
	t.Cleanup(func() {
		CommandTimeout, CleanupTimeout, waitDelay = originalCommand, originalCleanup, originalDelay
	})
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

func TestCreateClusterTimesOut(t *testing.T) {
	installFakeKind(t)
	withTimeouts(t, 200*time.Millisecond, time.Minute)
	t.Setenv("FAKE_KIND_CREATE_SLEEP", "60")

	started := time.Now()
	err := CreateCluster(CreateClusterOptions{Name: "stuck"})

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrTimeout)
	// The point of the change: the caller finds out promptly instead of
	// hanging until `go test` panics on its own -timeout.
	assert.Less(t, time.Since(started), 30*time.Second)
}

func TestCreateClusterTimeoutNamesTheCommand(t *testing.T) {
	installFakeKind(t)
	withTimeouts(t, 200*time.Millisecond, time.Minute)
	t.Setenv("FAKE_KIND_CREATE_SLEEP", "60")

	err := CreateCluster(CreateClusterOptions{Name: "stuck"})

	require.Error(t, err)
	// "signal: killed" on its own says nothing about what was stuck.
	assert.Contains(t, err.Error(), "create cluster")
	assert.Contains(t, err.Error(), "--name stuck")
	assert.Contains(t, err.Error(), commandTimeoutEnv)
}

func TestCreateClusterDeletesThePartialClusterOnTimeout(t *testing.T) {
	logPath := installFakeKind(t)
	withTimeouts(t, 200*time.Millisecond, time.Minute)
	t.Setenv("FAKE_KIND_CREATE_SLEEP", "60")

	require.Error(t, CreateCluster(CreateClusterOptions{Name: "stuck"}))

	// kind tears down a half-built cluster when one of its own actions fails,
	// but not when we kill it, so the delete has to be ours.
	assert.Contains(t, invocations(t, logPath), "delete cluster --name stuck")
}

func TestCreateClusterDoesNotDeleteOnSuccess(t *testing.T) {
	logPath := installFakeKind(t)
	withTimeouts(t, time.Minute, time.Minute)

	require.NoError(t, CreateCluster(CreateClusterOptions{Name: "fine"}))

	assert.Equal(t, []string{"create cluster --name fine"}, invocations(t, logPath))
}

// Negative control: a kind failure that is not a stall must keep reporting
// kind's own stderr, not be relabelled as a timeout.
func TestCreateClusterReportsKindErrorsUnchanged(t *testing.T) {
	installFakeKind(t)
	withTimeouts(t, time.Minute, time.Minute)
	t.Setenv("FAKE_KIND_EXIT", "1")
	t.Setenv("FAKE_KIND_STDERR", "ERROR: node(s) already exist for a cluster with the name")

	err := CreateCluster(CreateClusterOptions{Name: "existing"})

	require.Error(t, err)
	assert.NotErrorIs(t, err, ErrTimeout)
	assert.Contains(t, err.Error(), "node(s) already exist for a cluster with the name")
}

func TestCreateClusterReportsAFailedCleanup(t *testing.T) {
	installFakeKind(t)
	withTimeouts(t, 200*time.Millisecond, 200*time.Millisecond)
	t.Setenv("FAKE_KIND_CREATE_SLEEP", "60")
	t.Setenv("FAKE_KIND_DELETE_SLEEP", "60")

	err := CreateCluster(CreateClusterOptions{Name: "stuck"})

	require.Error(t, err)
	// The original stall is still the headline; the failed cleanup is extra.
	assert.ErrorIs(t, err, ErrTimeout)
	assert.Contains(t, err.Error(), "deleting the partial cluster also failed")
}

func TestLoadDockerImageTimesOut(t *testing.T) {
	installFakeKind(t)
	withTimeouts(t, 200*time.Millisecond, time.Minute)
	t.Setenv("FAKE_KIND_LOAD_SLEEP", "60")

	err := LoadDockerImage([]string{"some/image:latest"}, LoadDockerImageOptions{Name: "cluster"})

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrTimeout)
	assert.Contains(t, err.Error(), "load docker-image")
}

func TestLoadDockerImageSkipsAnEmptyList(t *testing.T) {
	logPath := installFakeKind(t)

	require.NoError(t, LoadDockerImage(nil, LoadDockerImageOptions{Name: "cluster"}))

	assert.Empty(t, invocations(t, logPath))
}

func TestGetKubeconfigReturnsStdout(t *testing.T) {
	installFakeKind(t)
	withTimeouts(t, time.Minute, time.Minute)
	t.Setenv("FAKE_KIND_KUBECONFIG", "apiVersion: v1")

	kubeconfig, err := GetKubeconfig(GetKubeconfigOptions{Name: "cluster"})

	require.NoError(t, err)
	assert.Equal(t, "apiVersion: v1", string(kubeconfig))
}

func TestGetKubeconfigTimesOut(t *testing.T) {
	installFakeKind(t)
	withTimeouts(t, 200*time.Millisecond, time.Minute)
	t.Setenv("FAKE_KIND_GET_SLEEP", "60")

	_, err := GetKubeconfig(GetKubeconfigOptions{Name: "cluster"})

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrTimeout)
	assert.Contains(t, err.Error(), "get kubeconfig")
}

func TestDeleteClusterTimesOut(t *testing.T) {
	installFakeKind(t)
	withTimeouts(t, 200*time.Millisecond, time.Minute)
	t.Setenv("FAKE_KIND_DELETE_SLEEP", "60")

	err := DeleteCluster(DeleteClusterOptions{Name: "cluster"})

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrTimeout)
	assert.Contains(t, err.Error(), "delete cluster")
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
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, resolveCommandTimeout(testCase.raw, fallback))
		})
	}
}
