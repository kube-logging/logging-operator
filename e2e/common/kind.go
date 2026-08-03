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
	"strings"

	"emperror.dev/errors"

	"github.com/kube-logging/logging-operator/e2e/internal/kind"
)

// KindClusterCreationTimeout is passed to `kind create cluster --wait`, which
// bounds only the last of kind's actions, waiting for control plane readiness.
// The whole invocation is bounded by Kind.CommandTimeout instead.
const KindClusterCreationTimeout = "3m"

var kindCLI = kind.New()

func KindClusterKubeconfig(name string) ([]byte, error) {
	create := kind.CreateClusterOptions{
		Name: name,
		Wait: KindClusterCreationTimeout,
	}

	err := kindCLI.CreateCluster(create)
	if err != nil && isClusterAlreadyExistsError(err) {
		// Adopting a leftover would hand the suite an unknown operator and data.
		fmt.Printf("kind cluster %q already exists, recreating it\n", name)
		if err := kindCLI.DeleteCluster(kind.DeleteClusterOptions{Name: name}); err != nil {
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
