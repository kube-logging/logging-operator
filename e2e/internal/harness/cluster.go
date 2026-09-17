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
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"emperror.dev/errors"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/kube-logging/logging-operator/e2e/internal/kind"
)

// controller-runtime prints a stack trace on first use when no logger is set,
// and its cache logs are not the suite's.
func init() {
	ctrllog.SetLogger(logr.Discard())
}

type kindCluster struct {
	cluster.Cluster
	name       string
	kubeconfig string
}

func newCluster(name string, scheme *runtime.Scheme) (*kindCluster, error) {
	kubeconfig, err := createCluster(name)
	if err != nil {
		return nil, errors.WrapIfWithDetails(err, "creating kind cluster", "clusterName", name)
	}
	restCfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, errors.WrapIfWithDetails(err, "reading kubeconfig", "path", kubeconfig)
	}
	c, err := cluster.New(restCfg, func(o *cluster.Options) { o.Scheme = scheme })
	if err != nil {
		return nil, errors.WrapIfWithDetails(err, "creating cluster with rest config", "cfg", restCfg)
	}
	return &kindCluster{Cluster: c, name: name, kubeconfig: kubeconfig}, nil
}

// A delete that fails after the assertions have run is the runner's state, and
// failing here would not remove the leftover the next run recreates anyway.
func deleteClusterOrLog(t *testing.T, name string) {
	t.Helper()
	if err := deleteCluster(name); err != nil {
		t.Logf("deleting cluster %q: %v", name, err)
	}
}

func deleteCluster(name string) error {
	kubeconfig, err := clusterKubeconfigPath(name)
	if err != nil {
		return err
	}
	if err := kindCLI.DeleteCluster(kind.DeleteClusterOptions{
		Name:       name,
		Kubeconfig: kubeconfig,
	}); err != nil {
		return errors.WrapIfWithDetails(err, "deleting kind cluster", "clusterName", name)
	}
	return removeClusterKubeconfig(name)
}

func kubectl(kubeconfig string, args ...string) *exec.Cmd {
	cmd := exec.Command("kubectl", args...)
	cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeconfig)
	cmd.Stderr = os.Stderr
	return cmd
}

func (c *kindCluster) loadImages(images ...string) error {
	return kindCLI.LoadDockerImage(images, kind.LoadDockerImageOptions{Name: c.name})
}

func (c *kindCluster) printLogs(namespaces []string, path string, limit int) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	cmd := exec.Command("stern", "-n", strings.Join(namespaces, ","), ".*", "--no-follow", "--tail", strconv.Itoa(limit), "--kubeconfig", c.kubeconfig)
	cmd.Stdout = f
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *kindCluster) collectCoverage(namespace, operator string) error {
	deployment := "deployment/" + operator
	if out, err := kubectl(c.kubeconfig, "-n", namespace, "exec", deployment, "--", "kill", "-USR1", "1").Output(); err != nil {
		return errors.WrapIfWithDetails(err, "Error in sending signal to logging-operator", out)
	}
	tarball, err := kubectl(c.kubeconfig, "-n", namespace, "exec", deployment, "--", "tar", "-cf", "-", "-C", "/", "covdatafiles").Output()
	if err != nil {
		return errors.WrapIfWithDetails(err, "Error in reading test coverage files", tarball)
	}

	extract := exec.Command("tar", "-xf", "-", "-C", os.Getenv("E2E_TEST_COV_DIR"))
	extract.Stdin = bytes.NewReader(tarball)
	if out, err := extract.CombinedOutput(); err != nil {
		return errors.WrapIfWithDetails(err, "Error in extracting test coverage files", out)
	}
	return nil
}
