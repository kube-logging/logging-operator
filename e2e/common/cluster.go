// Copyright © 2021 Banzai Cloud
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
	"context"
	"testing"

	"emperror.dev/errors"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/cluster"

	"github.com/banzaicloud/logging-operator/e2e/common/kind"
)

const defaultClusterName = "e2e-test"

type Cluster interface {
	cluster.Cluster
	LoadImages(images ...string) error
}

func WithCluster(t *testing.T, fn func(*testing.T, Cluster), opts ...cluster.Option) {
	cluster, err := GetTestCluster(defaultClusterName, opts...)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		require.NoError(t, cluster.Start(ctx))
	}()

	defer func() {
		cancel()
		require.NoError(t, DeleteTestCluster(defaultClusterName))
	}()

	fn(t, cluster)
}

func GetTestCluster(clusterName string, opts ...cluster.Option) (Cluster, error) {
	kubeconfig, err := KindClusterKubeconfig(clusterName)
	if err != nil {
		return nil, errors.WrapIfWithDetails(err, "getting kubeconfig of kind cluster", "clusterName", clusterName)
	}
	clientCfg, err := clientcmd.NewClientConfigFromBytes(kubeconfig)
	if err != nil {
		return nil, errors.WrapIfWithDetails(err, "creating client config from kubeconfig bytes", "kubeconfig", kubeconfig)
	}
	restCfg, err := clientCfg.ClientConfig()
	if err != nil {
		return nil, errors.WrapIfWithDetails(err, "creating rest config from client config", "cfg", clientCfg)
	}
	c, err := cluster.New(restCfg, opts...)
	return &kindCluster{
		Cluster:     c,
		clusterName: clusterName,
	}, errors.WrapIfWithDetails(err, "creating cluster with rest config", "cfg", restCfg)
}

func DeleteTestCluster(clusterName string) error {
	return errors.WrapIfWithDetails(kind.DeleteCluster(kind.DeleteClusterOptions{
		Name: clusterName,
	}), "deleting kind cluster", "clusterName", clusterName)
}

type kindCluster struct {
	cluster.Cluster
	clusterName string
}

func (c kindCluster) LoadImages(images ...string) error {
	return kind.LoadDockerImage(images, kind.LoadDockerImageOptions{
		Name: c.clusterName,
	})
}
