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

package harness

import (
	"path/filepath"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/kube-logging/logging-operator/e2e/internal/image"
)

func installOperator(t *testing.T, c *kindCluster, cfg config) {
	images := make([]string, 0, len(image.All()))
	for _, img := range image.All() {
		t.Logf("%s: loading %s", img.Env, img.Ref())
		images = append(images, img.Ref())
	}
	// One invocation, so docker save writes shared layers once instead of once
	// per image.
	loading := time.Now()
	if err := c.loadImages(images...); err != nil {
		t.Fatalf("kind load images: %s", err)
	}
	t.Logf("images loaded in %s", since(loading))

	installing := time.Now()
	operator := image.Operator()
	requireNoError(t, helmInstall(c.kubeconfig, Chart{
		Release:   "logging-operator",
		Namespace: cfg.controlNamespace,
		Name:      filepath.Join(projectDir(), "charts/logging-operator"),
		Values: map[string]any{
			"nameOverride": cfg.release,
			"image": map[string]any{
				"repository": operator.Repository,
				"tag":        operator.Tag,
				"pullPolicy": corev1.PullNever,
			},
			"testReceiver": map[string]any{
				"enabled": true,
			},
			"volumes": []map[string]any{
				{
					"name":     "coverage-data",
					"emptyDir": map[string]string{},
				},
			},
			"volumeMounts": []map[string]any{
				{
					"mountPath": "/covdatafiles",
					"name":      "coverage-data",
				},
			},
			"env": []map[string]any{
				{
					"name":  "GOCOVERDIR",
					"value": "/covdatafiles",
				},
			},
			"extraArgs": cfg.operatorArgs,
		},
	}))
	t.Logf("operator installed in %s", since(installing))
}
