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

package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClusterKubeconfigPath(t *testing.T) {
	alpha, err := ClusterKubeconfigPath("alpha")
	require.NoError(t, err)
	beta, err := ClusterKubeconfigPath("beta")
	require.NoError(t, err)

	require.NotEqual(t, alpha, beta)
	require.Equal(t, filepath.Dir(alpha), filepath.Dir(beta))

	again, err := ClusterKubeconfigPath("alpha")
	require.NoError(t, err)
	require.Equal(t, alpha, again, "create and delete have to name the same file")

	dir := filepath.Dir(alpha)
	require.NotEqual(t, os.TempDir(), dir, "a path in the shared temp directory is guessable")

	info, err := os.Stat(dir)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0o700), info.Mode().Perm())
}

func TestRemoveClusterKubeconfig(t *testing.T) {
	path := stubKubeconfig(t, "gamma")
	lock := path + ".lock"
	require.NoError(t, os.WriteFile(lock, nil, 0o600))

	require.NoError(t, RemoveClusterKubeconfig("gamma"))

	require.NoFileExists(t, path)
	require.NoFileExists(t, lock)
}

func TestRemoveClusterKubeconfigToleratesMissingFiles(t *testing.T) {
	require.NoError(t, RemoveClusterKubeconfig("never-created"))
}

// stubKubeconfig recreates the directory a previous case may have removed.
func stubKubeconfig(t *testing.T, name string) string {
	t.Helper()

	path, err := ClusterKubeconfigPath(name)
	require.NoError(t, err)
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o700))
	require.NoError(t, os.WriteFile(path, nil, 0o600))
	return path
}
