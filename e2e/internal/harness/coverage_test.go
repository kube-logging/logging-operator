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
	"archive/tar"
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func tarball(t *testing.T, members map[string]string) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	for name, body := range members {
		hdr := &tar.Header{Name: name, Mode: 0o644, Size: int64(len(body)), Typeflag: tar.TypeReg}
		if body == "" {
			hdr.Typeflag, hdr.Mode = tar.TypeDir, 0o755
		}
		require.NoError(t, tw.WriteHeader(hdr))
		_, err := tw.Write([]byte(body))
		require.NoError(t, err)
	}
	require.NoError(t, tw.Close())
	return &buf
}

// The stream the operator produces names its members covdatafiles/..., and the
// directory member may or may not precede the files.
func TestUntarWritesMembersUnderDir(t *testing.T) {
	dir := t.TempDir()

	require.NoError(t, untar(tarball(t, map[string]string{
		"covdatafiles/covmeta.abc": "meta",
		"covdatafiles/":            "",
		"covdatafiles/covcounters": "counters",
	}), dir))

	got, err := os.ReadFile(filepath.Join(dir, "covdatafiles", "covcounters"))
	require.NoError(t, err)
	assert.Equal(t, "counters", string(got))
}

func TestUntarRefusesAMemberOutsideDir(t *testing.T) {
	dir := t.TempDir()

	err := untar(tarball(t, map[string]string{"../escaped": "x"}), dir)

	require.ErrorContains(t, err, "escapes")
	assert.NoFileExists(t, filepath.Join(filepath.Dir(dir), "escaped"))
}
