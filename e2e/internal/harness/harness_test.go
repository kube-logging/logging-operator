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
	"fmt"
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

// Registered in a subtest so the order comes from testing's own LIFO rather
// than a stand-in for it.
func TestTeardownRunsInDeclaredOrder(t *testing.T) {
	var ran []string
	record := func(name string) step { return step{name, func() { ran = append(ran, name) }} }

	t.Run("teardown", func(t *testing.T) {
		teardown{record("artifacts"), record("kubeconfig"), record("stop"), record("delete")}.register(t)
	})

	assert.Equal(t, []string{"artifacts", "kubeconfig", "stop", "delete"}, ran)
}

// A real subtest cannot be used here: register reports the panic, which would
// fail the subtest and this test with it.
type reportingT struct {
	cleanups []func()
	reported []string
}

func (r *reportingT) Cleanup(fn func()) { r.cleanups = append(r.cleanups, fn) }
func (r *reportingT) Errorf(format string, args ...any) {
	r.reported = append(r.reported, fmt.Sprintf(format, args...))
}

func (r *reportingT) runCleanups() {
	for _, fn := range slices.Backward(r.cleanups) {
		fn()
	}
}

// Go abandons the remaining cleanups at the first panic, so without the recover
// a panicking log dump would leave the cluster running.
func TestTeardownDeletesTheClusterAfterAnEarlierPanic(t *testing.T) {
	var ran []string
	record := func(name string) step { return step{name, func() { ran = append(ran, name) }} }

	recorder := &reportingT{}
	teardown{
		{"artifacts", func() { ran = append(ran, "artifacts"); panic("boom") }},
		record("kubeconfig"), record("stop"), record("delete"),
	}.register(recorder)
	recorder.runCleanups()

	assert.Equal(t, []string{"artifacts", "kubeconfig", "stop", "delete"}, ran)
	assert.Len(t, recorder.reported, 1, "the panic is reported, not swallowed")
	assert.Contains(t, recorder.reported[0], `"artifacts" panicked`)
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

func TestWaitBudgetStaysUnderTheCapAndTheDeadline(t *testing.T) {
	budget := (&Env{T: t}).waitBudget()

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
