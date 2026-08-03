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
	"testing"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

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

	// Concat, not append: a caller's slice must not gain the extras.
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
		config Config
		want   []string
	}{
		"control namespace first, default last": {
			config: Config{ControlNamespace: "infra", Namespaces: []string{"tenant"}},
			want:   []string{"infra", "tenant", "default"},
		},
		// A suite that runs in default, or names it explicitly, would otherwise
		// have stern dump it twice.
		"repeats are dropped": {
			config: Config{ControlNamespace: "default", Namespaces: []string{"tenant", "tenant"}},
			want:   []string{"default", "tenant"},
		},
		"an unset control namespace is skipped": {
			config: Config{Namespaces: []string{"tenant"}},
			want:   []string{"tenant", "default"},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, testCase.want, dumpNamespaces(testCase.config))
		})
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
