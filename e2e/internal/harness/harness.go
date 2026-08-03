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
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/cisco-open/operator-tools/pkg/types"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"

	"github.com/kube-logging/logging-operator/e2e/common"
	"github.com/kube-logging/logging-operator/e2e/common/setup"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

// clusterStopTimeout bounds the wait for cluster.Start to return after the
// context is canceled, so a runnable that will not stop is named rather than
// running the binary out of its -timeout.
const clusterStopTimeout = time.Minute

// clusterLogLimit is the tail length used for the archived cluster dump.
const clusterLogLimit = 100 * 1000

// defaultSchemeBuilders covers every type the suites register today except the
// prometheus-operator ones, which only logging_metrics_monitoring needs.
var defaultSchemeBuilders = []func(*runtime.Scheme) error{
	v1beta1.AddToScheme,
	apiextensionsv1.AddToScheme,
	appsv1.AddToScheme,
	batchv1.AddToScheme,
	corev1.AddToScheme,
	rbacv1.AddToScheme,
}

type Config struct {
	// Cluster is the kind cluster name. Named rather than derived from
	// t.Name(): kind accepts only [a-z0-9.-] and none of the test names match.
	Cluster string

	// Release is the helm release, and the operator's nameOverride with it.
	Release string

	// ControlNamespace is where the operator and the test receiver run.
	ControlNamespace string

	// Namespaces are created before the test body. They are dumped at teardown
	// alongside ControlNamespace and default.
	Namespaces []string

	// OperatorArgs are appended to the operator's container arguments.
	OperatorArgs []string

	// ExtraSchemeBuilders are registered on top of the shared set, for the one
	// suite that needs the prometheus-operator types.
	ExtraSchemeBuilders []func(*runtime.Scheme) error
}

type Env struct {
	T   *testing.T
	Ctx context.Context

	Client  client.Client
	Cluster common.Cluster

	Release          string
	ControlNamespace string
}

// Start brings up the cluster, installs the operator and registers teardown.
// It marks the test parallel, so it has to be the first call in the test.
func Start(t *testing.T, cfg Config) *Env {
	common.Initialize(t)

	scheme, err := buildScheme(cfg.ExtraSchemeBuilders)
	common.RequireNoError(t, err)

	c, err := common.GetTestCluster(cfg.Cluster, func(o *cluster.Options) { o.Scheme = scheme })
	if err != nil {
		// The kind cluster is up before the client can fail, and teardown is
		// not registered yet.
		assert.NoError(t, common.DeleteTestCluster(cfg.Cluster))
	}
	common.RequireNoError(t, err)

	// Not t.Context(): that is canceled before the first Cleanup runs, which
	// would stop the cache before the log dump reads through it.
	ctx, cancel := context.WithCancel(context.Background())
	startErr := make(chan error, 1)
	go func() { startErr <- c.Start(ctx) }()

	env := &Env{
		T:                t,
		Ctx:              ctx,
		Client:           c.GetClient(),
		Cluster:          c,
		Release:          cfg.Release,
		ControlNamespace: cfg.ControlNamespace,
	}

	registerTeardown(t.Cleanup,
		func(step string, v any) { assert.Fail(t, fmt.Sprintf("teardown step %q panicked: %v", step, v)) },
		teardownSteps(
			func() { env.collectArtifacts(cfg) },
			func() { assert.NoError(t, c.Cleanup()) },
			func() { stopCluster(t, cancel, startErr) },
			func() { assert.NoError(t, common.DeleteTestCluster(cfg.Cluster)) },
		))

	setup.LoggingOperator(t, c, setup.LoggingOperatorOptionFunc(func(o *setup.LoggingOperatorOptions) {
		o.Namespace = cfg.ControlNamespace
		o.NameOverride = cfg.Release
		o.Args = cfg.OperatorArgs
	}))

	for _, ns := range cfg.Namespaces {
		env.Create(&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}})
	}

	return env
}

func (e *Env) Create(objects ...client.Object) {
	e.T.Helper()
	for _, object := range objects {
		common.RequireNoError(e.T, e.Client.Create(e.Ctx, object))
	}
}

// StartLogProducer runs on the test goroutine: the create errors reach it
// through RequireNoError, which needs FailNow to be defined.
func (e *Env) StartLogProducer(namespace string, labels map[string]string) {
	e.T.Helper()
	setup.LogProducer(e.T, e.Client, setup.LogProducerOptionFunc(func(o *setup.LogProducerOptions) {
		o.Namespace = namespace
		o.Labels = labels
	}))
}

// ReceiverLogs returns the tail of the chart's test receiver, which is where a
// suite looks to see that its logs arrived.
func (e *Env) ReceiverLogs(tail int) (string, error) {
	out, err := common.CmdEnv(exec.Command("kubectl",
		"logs",
		"-n", e.ControlNamespace,
		"--tail", fmt.Sprint(tail),
		"-l", fmt.Sprintf("%s=%s-test-receiver", types.NameLabel, e.Release)), e.Cluster).Output()
	return string(out), err
}

func (e *Env) collectArtifacts(cfg Config) {
	path, err := artifactPath(e.T.Name())
	if err != nil {
		e.T.Logf("Skipping cluster logs: %s", err)
	} else {
		e.T.Logf("Printing cluster logs to %s", path)
		assert.NoError(e.T, e.Cluster.PrintLogs(common.PrintLogConfig{
			Namespaces: dumpNamespaces(cfg),
			FilePath:   path,
			Limit:      clusterLogLimit,
		}))
	}

	operator := "logging-operator-" + e.Release
	e.T.Logf("Collecting coverage files from logging-operator: %s/%s", e.ControlNamespace, operator)
	if err := e.Cluster.CollectTestCoverageFiles(e.ControlNamespace, operator); err != nil {
		// Logged, never fatal: the run's coverage is not the suite's verdict.
		e.T.Logf("Failed collecting coverage files: %s", err)
	}
}

type teardownStep struct {
	name string
	run  func()
}

// teardownSteps lists them in the order they have to run. This order is what
// kept clusters from leaking before it moved here, so it is one list rather
// than four Cleanup calls in reverse.
func teardownSteps(artifacts, kubeconfig, stop, deleteCluster func()) []teardownStep {
	return []teardownStep{
		{"artifacts", artifacts},
		{"kubeconfig", kubeconfig},
		{"stop", stop},
		{"delete", deleteCluster},
	}
}

// registerTeardown hands the steps over back to front, because Cleanup is LIFO,
// and isolates each one: Go abandons the remaining cleanups at the first panic,
// which would strand the cluster.
func registerTeardown(cleanup func(func()), onPanic func(step string, v any), steps []teardownStep) {
	for _, step := range slices.Backward(steps) {
		cleanup(func() {
			defer func() {
				if v := recover(); v != nil {
					onPanic(step.name, v)
				}
			}()
			step.run()
		})
	}
}

func stopCluster(t *testing.T, cancel context.CancelFunc, startErr <-chan error) {
	cancel()
	// Checked here, not in the goroutine: FailNow is undefined off the test one.
	select {
	case err := <-startErr:
		assert.NoError(t, err, "starting the cluster")
	case <-time.After(clusterStopTimeout):
		assert.Fail(t, "cluster.Start did not return after cancellation")
	}
}

func buildScheme(extra []func(*runtime.Scheme) error) (*runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	for _, add := range slices.Concat(defaultSchemeBuilders, extra) {
		if err := add(scheme); err != nil {
			return nil, err
		}
	}
	return scheme, nil
}

// dumpNamespaces is ControlNamespace, the configured ones and default, without
// repeats, so a suite that names one of them does not dump it twice.
func dumpNamespaces(cfg Config) []string {
	out := make([]string, 0, len(cfg.Namespaces)+2)
	seen := map[string]bool{}
	// Concat, not append: appending to cfg.Namespaces could write into the
	// caller's backing array.
	for _, ns := range slices.Concat([]string{cfg.ControlNamespace}, cfg.Namespaces, []string{"default"}) {
		if ns != "" && !seen[ns] {
			seen[ns] = true
			out = append(out, ns)
		}
	}
	return out
}

// artifactPath is build/_test under PROJECT_DIR, the directory every suite
// used to build for itself in an init().
func artifactPath(name string) (string, error) {
	root, ok := os.LookupEnv("PROJECT_DIR")
	if !ok {
		root = "../.."
	}
	dir := filepath.Join(root, "build/_test")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(dir, fmt.Sprintf("cluster-%s.log", name)), nil
}
