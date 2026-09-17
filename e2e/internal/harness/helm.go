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
	"fmt"
	"strings"

	"emperror.dev/errors"
	"helm.sh/helm/v4/pkg/action"
	"helm.sh/helm/v4/pkg/chart/loader"
	"helm.sh/helm/v4/pkg/cli"
)

// Chart is a helm chart to install: Name is a path when Repo is empty, and a
// chart in that repository otherwise.
type Chart struct {
	Release   string
	Namespace string
	Repo      string
	Name      string
	Values    map[string]any
}

func (e *Env) InstallChart(chart Chart) {
	e.T.Helper()
	requireNoError(e.T, helmInstall(e.cluster.kubeconfig, chart))
}

// helmInstall goes through the SDK rather than the helm binary, and keeps
// helm's debug log for the error: it narrates every resource, which only says
// anything when the install fails.
func helmInstall(kubeconfig string, chart Chart) error {
	getter, err := newRESTClientGetter(kubeconfig, chart.Namespace)
	if err != nil {
		return errors.WrapIf(err, "helm rest client getter")
	}
	var helmLog strings.Builder
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(getter, chart.Namespace, "memory", func(format string, v ...any) {
		fmt.Fprintf(&helmLog, format+"\n", v...)
	}); err != nil {
		return errors.WrapIf(err, "helm action config init")
	}

	installer := action.NewInstall(actionConfig)
	installer.Namespace = chart.Namespace
	installer.CreateNamespace = true
	installer.ReleaseName = chart.Release
	installer.RepoURL = chart.Repo

	path, err := installer.LocateChart(chart.Name, cli.New())
	if err != nil {
		return errors.WrapIfWithDetails(err, "helm locate chart", "chart", chart.Name, "repo", chart.Repo)
	}
	loaded, err := loader.Load(path)
	if err != nil {
		return errors.WrapIfWithDetails(err, "helm load chart", "path", path)
	}
	if _, err := installer.Run(loaded, chart.Values); err != nil {
		return errors.WrapIfWithDetails(err, "helm install", "release", chart.Release, "helm", helmLog.String())
	}
	return nil
}
