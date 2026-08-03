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

package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"emperror.dev/errors"

	"github.com/kube-logging/logging-operator/e2e/internal/kind"
)

// KindClusterCreationTimeout is passed to `kind create cluster --wait`, which
// bounds only the last of kind's actions, waiting for control plane readiness.
// The whole invocation is bounded by Kind.CommandTimeout instead.
const KindClusterCreationTimeout = "3m"

var kindCLI = kind.New()

// kubeconfigDir is one 0700 directory per run, so the paths in it are unguessable.
var kubeconfigDir = sync.OnceValues(func() (string, error) {
	return os.MkdirTemp("", "e2e-kubeconfig-*")
})

// clusterKubeconfigPath keeps each cluster's kind bookkeeping in its own file.
// kind locks the kubeconfig it updates and the lock is non-blocking, so
// clusters sharing one fail outright rather than wait.
func clusterKubeconfigPath(name string) (string, error) {
	dir, err := kubeconfigDir()
	if err != nil {
		return "", errors.WrapIf(err, "creating the kubeconfig directory")
	}
	// Reasserted on every lookup, not left to MkdirTemp: kind recreates a missing
	// parent itself, at 0755, so a directory that goes away comes back readable.
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", errors.WrapIfWithDetails(err, "creating the kubeconfig directory", "path", dir)
	}
	return filepath.Join(dir, "kind-"+name+".kubeconfig"), nil
}

// removeClusterKubeconfig drops the file and the lock kind leaves beside it,
// which kind itself does not. The directory stays for the run: removing it with
// one cluster only let kind recreate it at 0755 for the next.
func removeClusterKubeconfig(name string) error {
	path, err := clusterKubeconfigPath(name)
	if err != nil {
		return err
	}
	if err := removeIfExists(path); err != nil {
		return err
	}
	return removeIfExists(path + ".lock")
}

func removeIfExists(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return errors.WrapIfWithDetails(err, "removing kubeconfig", "path", path)
	}
	return nil
}

func KindClusterKubeconfig(name string) ([]byte, error) {
	kubeconfig, err := clusterKubeconfigPath(name)
	if err != nil {
		return nil, err
	}

	create := kind.CreateClusterOptions{
		Name:       name,
		Wait:       KindClusterCreationTimeout,
		Kubeconfig: kubeconfig,
	}

	err = kindCLI.CreateCluster(create)
	if err != nil && isClusterAlreadyExistsError(err) {
		// Adopting a leftover would hand the suite an unknown operator and data.
		fmt.Printf("kind cluster %q already exists, recreating it\n", name)
		if err := kindCLI.DeleteCluster(kind.DeleteClusterOptions{
			Name:       name,
			Kubeconfig: kubeconfig,
		}); err != nil {
			return nil, errors.WrapIfWithDetails(err, "deleting a leftover kind cluster", "clusterName", name)
		}
		err = kindCLI.CreateCluster(create)
	}
	if err != nil {
		return nil, err
	}

	return kindCLI.GetKubeconfig(kind.GetKubeconfigOptions{
		Name: name,
	})
}

func isClusterAlreadyExistsError(err error) bool {
	return strings.Contains(err.Error(), "failed to create cluster: node(s) already exist for a cluster with the name")
}
