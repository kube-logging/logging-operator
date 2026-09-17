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
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// dumpLogs writes the tail of every container in the namespaces to path, each
// line prefixed with namespace/pod container. A container whose log cannot be
// read gets the error in its place rather than aborting the dump.
func (c *kindCluster) dumpLogs(ctx context.Context, namespaces []string, path string, tail int64) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	for _, ns := range namespaces {
		pods, err := c.clientset.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			_, _ = fmt.Fprintf(w, "%s: listing pods: %v\n", ns, err)
			continue
		}
		for _, pod := range pods.Items {
			for _, container := range slices.Concat(pod.Spec.InitContainers, pod.Spec.Containers) {
				prefix := fmt.Sprintf("%s/%s %s", ns, pod.Name, container.Name)
				logs, err := c.podLogs(ctx, ns, pod.Name, container.Name, tail)
				if err != nil {
					_, _ = fmt.Fprintf(w, "%s: %v\n", prefix, err)
					continue
				}
				writePrefixed(w, prefix, logs)
			}
		}
	}
	if err := w.Flush(); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

func writePrefixed(w io.Writer, prefix, logs string) {
	for line := range strings.Lines(logs) {
		_, _ = fmt.Fprintf(w, "%s %s\n", prefix, strings.TrimSuffix(line, "\n"))
	}
}
