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
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/kube-logging/logging-operator/e2e/internal/kind"
)

// clusterStopTimeout bounds the wait for cluster.Start to return once the
// context is canceled, so a runnable that will not stop is named here instead
// of running the whole binary out of its -timeout.
const clusterStopTimeout = time.Minute

type Cluster interface {
	cluster.Cluster
	LoadImages(images ...string) error
	Cleanup() error
	PrintLogs(config PrintLogConfig) error
	KubeConfigFilePath() string
	CollectTestCoverageFiles(string, string) error
}

type PrintLogConfig struct {
	Namespaces []string
	FilePath   string
	Limit      int
}

func WithCluster(name string, t *testing.T, fn func(*testing.T, Cluster), beforeCleanup func(*testing.T, Cluster) error, opts ...cluster.Option) {
	zapLogger := zap.New(func(o *zap.Options) {
		o.Development = true
		encoder := zap.ConsoleEncoder()
		encoder(o)
	})

	ctrl.SetLogger(zapLogger)

	cluster, err := GetTestCluster(name, opts...)
	if err != nil {
		// The cluster is created before the client can fail, and the deferred
		// teardown below is not registered yet.
		assert.NoError(t, DeleteTestCluster(name))
	}
	RequireNoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	startErr := make(chan error, 1)
	go func() {
		startErr <- cluster.Start(ctx)
	}()

	defer func() {
		assert.NoError(t, beforeCleanup(t, cluster))
		assert.NoError(t, cluster.Cleanup())
		cancel()

		// Start's error is checked here rather than in the goroutine: FailNow is
		// undefined off the test goroutine, so a cluster that never started used
		// to be dropped and resurface as an unrelated Eventually timeout.
		select {
		case err := <-startErr:
			assert.NoError(t, err, "starting the cluster")
		case <-time.After(clusterStopTimeout):
			assert.Fail(t, "cluster.Start did not return after cancellation")
		}

		// A cluster that came down cleanly must not fail the test on kind's own
		// kubeconfig bookkeeping, so this matches the assertions above.
		assert.NoError(t, DeleteTestCluster(name))
	}()

	fn(t, cluster)
}

func GetTestCluster(clusterName string, opts ...cluster.Option) (Cluster, error) {
	var err error
	var kubeconfig []byte
	var kubeconfigFile *os.File
	var clientCfg clientcmd.ClientConfig
	var restCfg *rest.Config
	var c cluster.Cluster

	kubeconfig, err = KindClusterKubeconfig(clusterName)
	if err != nil {
		return nil, errors.WrapIfWithDetails(err, "getting kubeconfig of kind cluster", "clusterName", clusterName)
	}
	kubeconfigFile, err = os.CreateTemp("", "kind-kind-kubeconfig")
	if err != nil {
		return nil, errors.WrapIfWithDetails(err, "unable to create temp file for kubeconfig", "clusterName", clusterName)
	}
	err = os.WriteFile(kubeconfigFile.Name(), kubeconfig, os.FileMode(0o600))
	if err != nil {
		return nil, errors.WrapIfWithDetails(err, "failed to write kubeconfig", "clusterName", clusterName, "path", kubeconfigFile.Name())
	}
	clientCfg, err = clientcmd.NewClientConfigFromBytes(kubeconfig)
	if err != nil {
		return nil, errors.WrapIfWithDetails(err, "creating client config from kubeconfig bytes", "kubeconfig", kubeconfig)
	}
	restCfg, err = clientCfg.ClientConfig()
	if err != nil {
		return nil, errors.WrapIfWithDetails(err, "creating rest config from client config", "cfg", clientCfg)
	}
	c, err = cluster.New(restCfg, opts...)
	return &kindCluster{
		Cluster:            c,
		kubeconfigFilePath: kubeconfigFile.Name(),
		clusterName:        clusterName,
	}, errors.WrapIfWithDetails(err, "creating cluster with rest config", "cfg", restCfg)
}

func DeleteTestCluster(clusterName string) error {
	kubeconfig, err := ClusterKubeconfigPath(clusterName)
	if err != nil {
		return err
	}
	if err := kindCLI.DeleteCluster(kind.DeleteClusterOptions{
		Name:       clusterName,
		Kubeconfig: kubeconfig,
	}); err != nil {
		return errors.WrapIfWithDetails(err, "deleting kind cluster", "clusterName", clusterName)
	}
	return RemoveClusterKubeconfig(clusterName)
}

func CmdEnv(cmd *exec.Cmd, c Cluster) *exec.Cmd {
	cmd.Env = append(os.Environ(), fmt.Sprintf("KUBECONFIG=%s", c.KubeConfigFilePath()))
	cmd.Stderr = os.Stderr
	return cmd
}

// Collects test coverage data files from logging-operator
func (c kindCluster) CollectTestCoverageFiles(ns string, loggingOperatorName string) error {
	cmd := CmdEnv(exec.Command("kubectl", "-n", ns,
		"exec", fmt.Sprintf("deployment/%s", loggingOperatorName), "--",
		"kill", "-USR1", "1"), c)
	cmdOut, err := cmd.Output()
	if err != nil {
		return errors.WrapIfWithDetails(err, "Error in sending signal to logging-operator", cmdOut)
	}
	testCovDir := os.Getenv("E2E_TEST_COV_DIR")
	archive := CmdEnv(exec.Command("kubectl", "-n", ns,
		"exec", fmt.Sprintf("deployment/%s", loggingOperatorName), "--",
		"tar", "-cf", "-", "/covdatafiles"), c)
	tarball, err := archive.Output()
	if err != nil {
		return errors.WrapIfWithDetails(err, "Error in reading test coverage files", tarball)
	}

	extract := exec.Command("tar", "-xf", "-", "-C", testCovDir)
	extract.Stdin = bytes.NewReader(tarball)
	if cmdOut, err := extract.CombinedOutput(); err != nil {
		return errors.WrapIfWithDetails(err, "Error in extracting test coverage files", cmdOut)
	}
	return nil
}

type kindCluster struct {
	cluster.Cluster
	kubeconfigFilePath string
	clusterName        string
}

func (c kindCluster) PrintLogs(config PrintLogConfig) error {
	cmd := exec.Command("stern", "-n", strings.Join(config.Namespaces, ","), ".*", "--no-follow", "--tail", cast.ToString(config.Limit), "--kubeconfig", c.kubeconfigFilePath)
	f, err := os.Create(config.FilePath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	cmd.Stdout = f
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (c kindCluster) Cleanup() error {
	return os.Remove(c.kubeconfigFilePath)
}

func (c kindCluster) LoadImages(images ...string) error {
	return kindCLI.LoadDockerImage(images, kind.LoadDockerImageOptions{
		Name: c.clusterName,
	})
}

func (c kindCluster) KubeConfigFilePath() string {
	return c.kubeconfigFilePath
}
