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

package harness

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

// cleanupRecorder stands in for testing.T's Cleanup, so the order the steps end
// up running in can be asserted without a cluster.
type cleanupRecorder struct {
	registered []func()
}

func (c *cleanupRecorder) Cleanup(fn func()) {
	c.registered = append(c.registered, fn)
}

// run does what testing does at the end of a test: LIFO.
func (c *cleanupRecorder) run() {
	for _, fn := range slices.Backward(c.registered) {
		fn()
	}
}

// The order here is what kept clusters from leaking: the log dump needs the
// cache and the kubeconfig, both of which the later steps take away.
func TestTeardownRunsInDeclaredOrder(t *testing.T) {
	var ran []string
	record := func(name string) func() { return func() { ran = append(ran, name) } }

	recorder := &cleanupRecorder{}
	registerTeardown(recorder.Cleanup, failOnPanic(t), teardownSteps(
		record("artifacts"), record("kubeconfig"), record("stop"), record("delete"),
	))
	recorder.run()

	assert.Equal(t, []string{"artifacts", "kubeconfig", "stop", "delete"}, ran)
}

// Go abandons the remaining cleanups at the first panic, so without the recover
// in registerTeardown a panicking log dump would leave the cluster running.
func TestTeardownDeletesTheClusterAfterAnEarlierPanic(t *testing.T) {
	var ran, panicked []string
	record := func(name string) func() { return func() { ran = append(ran, name) } }

	recorder := &cleanupRecorder{}
	registerTeardown(recorder.Cleanup,
		func(step string, _ any) { panicked = append(panicked, step) },
		teardownSteps(
			func() { ran = append(ran, "artifacts"); panic("boom") },
			record("kubeconfig"), record("stop"), record("delete"),
		))
	recorder.run()

	assert.Equal(t, []string{"artifacts", "kubeconfig", "stop", "delete"}, ran)
	assert.Equal(t, []string{"artifacts"}, panicked, "the panic is reported, not swallowed")
}

func failOnPanic(t *testing.T) func(string, any) {
	return func(step string, v any) { assert.Fail(t, "unexpected panic", "%s: %v", step, v) }
}

func TestBuildScheme(t *testing.T) {
	t.Run("the default set covers what the suites register", func(t *testing.T) {
		scheme, err := buildScheme(nil)
		require.NoError(t, err)

		for _, obj := range []runtime.Object{
			&v1beta1.Logging{}, &v1beta1.ClusterFlow{}, &corev1.Pod{},
		} {
			assert.True(t, scheme.Recognizes(mustGVK(t, scheme, obj)), "%T", obj)
		}
	})

	t.Run("extras are added on top of the default", func(t *testing.T) {
		scheme, err := buildScheme([]func(*runtime.Scheme) error{monitoringv1.AddToScheme})
		require.NoError(t, err)

		assert.True(t, scheme.Recognizes(mustGVK(t, scheme, &monitoringv1.ServiceMonitor{})))
		assert.True(t, scheme.Recognizes(mustGVK(t, scheme, &v1beta1.Logging{})), "the shared set stays")
	})

	t.Run("the default set is not mutated by extras", func(t *testing.T) {
		before := len(defaultSchemeBuilders)
		_, err := buildScheme([]func(*runtime.Scheme) error{monitoringv1.AddToScheme})
		require.NoError(t, err)

		assert.Len(t, defaultSchemeBuilders, before)
	})

	t.Run("a failing builder is reported", func(t *testing.T) {
		want := errors.New("boom")
		_, err := buildScheme([]func(*runtime.Scheme) error{
			func(*runtime.Scheme) error { return want },
		})
		assert.ErrorIs(t, err, want)
	})
}

func TestDumpNamespaces(t *testing.T) {
	testCases := map[string]struct {
		config config
		want   []string
	}{
		"control namespace first, default last": {
			config: config{controlNamespace: "infra", namespaces: []string{"tenant"}},
			want:   []string{"infra", "tenant", "default"},
		},
		// A suite that runs in default, or names it explicitly, would otherwise
		// have stern dump it twice.
		"repeats are dropped": {
			config: config{controlNamespace: "default", namespaces: []string{"tenant", "tenant"}},
			want:   []string{"default", "tenant"},
		},
		"an unset control namespace is skipped": {
			config: config{namespaces: []string{"tenant"}},
			want:   []string{"tenant", "default"},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, testCase.want, dumpNamespaces(testCase.config))
		})
	}
}

// The cap is what the suites spend by hand today; the deadline is what stops
// one wait taking the whole package budget with it.
func TestWaitBudgetStaysUnderTheCapAndTheDeadline(t *testing.T) {
	env := &Env{T: t}

	budget := env.waitBudget()

	assert.Positive(t, budget)
	assert.LessOrEqual(t, budget, waitBudget)
	if deadline, ok := t.Deadline(); ok {
		assert.Less(t, budget, time.Until(deadline), "a wait must not outlive the binary")
	}
}

func TestArtifactPath(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PROJECT_DIR", root)

	path, err := artifactPath("TestSomething")
	require.NoError(t, err)

	assert.Equal(t, filepath.Join(root, "build/_test", "cluster-TestSomething.log"), path)
	// The directory each suite used to build in its own init().
	info, err := os.Stat(filepath.Dir(path))
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func mustGVK(t *testing.T, scheme *runtime.Scheme, obj runtime.Object) schema.GroupVersionKind {
	t.Helper()
	gvks, _, err := scheme.ObjectKinds(obj)
	require.NoError(t, err)
	require.NotEmpty(t, gvks)
	return gvks[0]
}
