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

package setup

import (
	"fmt"
	"os"
	"testing"

	"helm.sh/helm/v4/pkg/action"
	"helm.sh/helm/v4/pkg/chart/loader"
	"helm.sh/helm/v4/pkg/cli"
	corev1 "k8s.io/api/core/v1"

	"github.com/kube-logging/logging-operator/e2e/common"
	"github.com/kube-logging/logging-operator/e2e/internal/image"
)

func LoggingOperator(t *testing.T, c common.Cluster, opts ...LoggingOperatorOption) {
	opt := &LoggingOperatorOptions{
		Namespace:    "default",
		NameOverride: "logging-operator",
	}

	for _, o := range opts {
		o.ApplyToLoggingOperatorOptions(opt)
	}

	restClientGetter, err := newRESTClientGetter(c.KubeConfigFilePath(), opt.Namespace)
	if err != nil {
		t.Fatalf("helm rest client getter: %s", err)
	}
	actionConfig := new(action.Configuration)

	if err := actionConfig.Init(restClientGetter, opt.Namespace, "memory", func(format string, v ...any) {
		t.Logf(format, v...)
	}); err != nil {
		t.Fatalf("helm action config init: %s", err)
	}

	installer := action.NewInstall(actionConfig)

	installer.Namespace = opt.Namespace
	installer.CreateNamespace = true
	installer.ReleaseName = "logging-operator"

	projectDir := os.Getenv("PROJECT_DIR")
	if projectDir == "" {
		projectDir = "../.."
	}

	cp, err := installer.LocateChart(fmt.Sprintf("%s/charts/logging-operator", projectDir), cli.New())
	if err != nil {
		t.Fatalf("helm locate chart: %s", err)
	}
	chartReq, err := loader.Load(cp)
	if err != nil {
		t.Fatalf("helm load chart: %s", err)
	}

	loggingOperatorImage := image.Operator()
	images := make([]string, 0, len(image.All()))
	for _, img := range image.All() {
		t.Logf("%s: loading %s", img.Env, img.Ref())
		images = append(images, img.Ref())
	}

	// One invocation, so docker save writes shared layers once instead of once
	// per image.
	if err := c.LoadImages(images...); err != nil {
		t.Fatalf("kind load images: %s", err)
	}

	_, err = installer.Run(chartReq, map[string]any{
		"nameOverride": opt.NameOverride,
		"image": map[string]any{
			"repository": loggingOperatorImage.Repository,
			"tag":        loggingOperatorImage.Tag,
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
		"extraArgs": opt.Args,
	})
	if err != nil {
		t.Fatalf("helm chart install: %s", err)
	}
}

type LoggingOperatorOption interface {
	ApplyToLoggingOperatorOptions(options *LoggingOperatorOptions)
}

type LoggingOperatorOptionFunc func(*LoggingOperatorOptions)

func (fn LoggingOperatorOptionFunc) ApplyToLoggingOperatorOptions(options *LoggingOperatorOptions) {
	fn(options)
}

type LoggingOperatorOptions struct {
	Namespace    string
	NameOverride string
	Args         []string
}
