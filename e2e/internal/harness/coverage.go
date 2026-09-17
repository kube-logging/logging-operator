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
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/types"
	corev1 "k8s.io/api/core/v1"
)

const coverageDir = "/covdatafiles"

// collectCoverage has the operator flush its counters, then pulls the coverage
// directory out as a tar stream into E2E_TEST_COV_DIR.
func (c *kindCluster) collectCoverage(ctx context.Context, namespace, release string) error {
	dest := os.Getenv("E2E_TEST_COV_DIR")
	if dest == "" {
		return errors.New("E2E_TEST_COV_DIR is not set")
	}
	pods, err := c.pods(ctx, namespace, map[string]string{types.NameLabel: release})
	if err != nil {
		return errors.WrapIf(err, "listing the operator pods")
	}
	i := slices.IndexFunc(pods, func(p corev1.Pod) bool { return p.Status.Phase == corev1.PodRunning })
	if i < 0 {
		return errors.New("no running operator pod")
	}
	pod := pods[i].Name
	if _, err := c.exec(ctx, namespace, pod, "", "kill", "-USR1", "1"); err != nil {
		return errors.WrapIf(err, "signaling the operator")
	}
	tarball, err := c.exec(ctx, namespace, pod, "", "tar", "-cf", "-", "-C", filepath.Dir(coverageDir), filepath.Base(coverageDir))
	if err != nil {
		return errors.WrapIf(err, "reading the coverage files")
	}
	return untar(bytes.NewReader(tarball), dest)
}

// untar writes a tar stream under dir, refusing members that would land
// outside it.
func untar(r io.Reader, dir string) error {
	root := filepath.Clean(dir) + string(os.PathSeparator)
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		target := filepath.Join(dir, hdr.Name)
		if !strings.HasPrefix(target, root) {
			return fmt.Errorf("tar member %q escapes %s", hdr.Name, dir)
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := writeFile(target, tr, os.FileMode(hdr.Mode).Perm()); err != nil {
				return err
			}
		}
	}
}

func writeFile(path string, r io.Reader, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, r); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}
