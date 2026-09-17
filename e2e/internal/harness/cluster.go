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
	"context"
	"testing"

	"emperror.dev/errors"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
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
	restConfig *rest.Config
	clientset  kubernetes.Interface
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
	clientset, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return nil, errors.WrapIf(err, "creating clientset")
	}
	c, err := cluster.New(restCfg, func(o *cluster.Options) { o.Scheme = scheme })
	if err != nil {
		return nil, errors.WrapIfWithDetails(err, "creating cluster with rest config", "cfg", restCfg)
	}
	return &kindCluster{Cluster: c, name: name, kubeconfig: kubeconfig, restConfig: restCfg, clientset: clientset}, nil
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

func (c *kindCluster) loadImages(images ...string) error {
	return kindCLI.LoadDockerImage(images, kind.LoadDockerImageOptions{Name: c.name})
}

func (c *kindCluster) pods(ctx context.Context, namespace string, selector map[string]string) ([]corev1.Pod, error) {
	list, err := c.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(selector).String(),
	})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

// podLogs returns the last tail lines of one container, or all of them when
// tail is negative.
func (c *kindCluster) podLogs(ctx context.Context, namespace, pod, container string, tail int64) (string, error) {
	opts := &corev1.PodLogOptions{Container: container}
	if tail >= 0 {
		opts.TailLines = &tail
	}
	out, err := c.clientset.CoreV1().Pods(namespace).GetLogs(pod, opts).DoRaw(ctx)
	return string(out), err
}

// exec runs command in a container and returns its stdout; an empty container
// means the pod's only one.
func (c *kindCluster) exec(ctx context.Context, namespace, pod, container string, command ...string) ([]byte, error) {
	req := c.clientset.CoreV1().RESTClient().Post().
		Resource("pods").Namespace(namespace).Name(pod).SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: container,
			Command:   command,
			Stdout:    true,
			Stderr:    true,
		}, clientgoscheme.ParameterCodec)
	executor, err := remotecommand.NewSPDYExecutor(c.restConfig, "POST", req.URL())
	if err != nil {
		return nil, errors.WrapIf(err, "creating the exec request")
	}
	var stdout, stderr bytes.Buffer
	if err := executor.StreamWithContext(ctx, remotecommand.StreamOptions{Stdout: &stdout, Stderr: &stderr}); err != nil {
		return stdout.Bytes(), errors.WrapIfWithDetails(err, "exec", "pod", namespace+"/"+pod, "command", command, "stderr", stderr.String())
	}
	return stdout.Bytes(), nil
}
