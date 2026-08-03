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
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cisco-open/operator-tools/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	"github.com/kube-logging/logging-operator/e2e/internal/wait"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

const (
	// clusterStopTimeout bounds the wait for cluster.Start to return after the
	// context is canceled, so a runnable that will not stop is named rather
	// than running the binary out of its -timeout.
	clusterStopTimeout = time.Minute

	clusterLogLimit = 100 * 1000

	// receiverLogTail is the tail length the suites read from the test receiver.
	receiverLogTail = 30

	// waitBudget and waitInterval are what the suites spend by hand today.
	waitBudget   = 5 * time.Minute
	waitInterval = 3 * time.Second

	// teardownMargin is held back from the budget so a wait that runs to the
	// end still leaves time to dump logs and delete the cluster.
	teardownMargin = 90 * time.Second
)

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

type config struct {
	cluster          string
	release          string
	controlNamespace string
	namespaces       []string
	operatorArgs     []string
	schemeBuilders   []func(*runtime.Scheme) error
}

type Builder struct {
	t   *testing.T
	cfg config
}

func New(t *testing.T) *Builder {
	return &Builder{t: t}
}

// WithCluster names the kind cluster. Named rather than derived from t.Name():
// kind accepts only [a-z0-9.-] and none of the test names match.
func (b *Builder) WithCluster(name string) *Builder {
	b.cfg.cluster = name
	return b
}

// WithRelease names the helm release, and the operator's nameOverride with it.
func (b *Builder) WithRelease(name string) *Builder {
	b.cfg.release = name
	return b
}

// WithControlNamespace sets where the operator and the test receiver run.
func (b *Builder) WithControlNamespace(namespace string) *Builder {
	b.cfg.controlNamespace = namespace
	return b
}

// WithNamespaces are created before the test body and dumped at teardown
// alongside the control namespace and default.
func (b *Builder) WithNamespaces(namespaces ...string) *Builder {
	b.cfg.namespaces = append(b.cfg.namespaces, namespaces...)
	return b
}

func (b *Builder) WithOperatorArgs(args ...string) *Builder {
	b.cfg.operatorArgs = append(b.cfg.operatorArgs, args...)
	return b
}

// WithScheme registers types on top of the shared set.
func (b *Builder) WithScheme(add ...func(*runtime.Scheme) error) *Builder {
	b.cfg.schemeBuilders = append(b.cfg.schemeBuilders, add...)
	return b
}

type Env struct {
	T   *testing.T
	Ctx context.Context

	Client  client.Client
	Cluster common.Cluster

	Release          string
	ControlNamespace string

	dumpNamespaces []string
}

// Start brings up the cluster, installs the operator and registers teardown.
// It marks the test parallel, so it has to be the first call in the test.
func (b *Builder) Start() *Env {
	t := b.t
	common.Initialize(t)

	scheme, err := buildScheme(b.cfg.schemeBuilders)
	common.RequireNoError(t, err)

	c, err := common.GetTestCluster(b.cfg.cluster, func(o *cluster.Options) { o.Scheme = scheme })
	if err != nil {
		// The kind cluster is up before the client can fail, and teardown is
		// not registered yet.
		assert.NoError(t, common.DeleteTestCluster(b.cfg.cluster))
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
		Release:          b.cfg.release,
		ControlNamespace: b.cfg.controlNamespace,
		dumpNamespaces:   dumpNamespaces(b.cfg),
	}

	registerTeardown(t.Cleanup,
		func(step string, v any) { assert.Fail(t, fmt.Sprintf("teardown step %q panicked: %v", step, v)) },
		teardownSteps(
			env.collectArtifacts,
			func() { assert.NoError(t, c.Cleanup()) },
			func() { stopCluster(t, cancel, startErr) },
			func() { assert.NoError(t, common.DeleteTestCluster(b.cfg.cluster)) },
		))

	setup.LoggingOperator(t, c, setup.LoggingOperatorOptionFunc(func(o *setup.LoggingOperatorOptions) {
		o.Namespace = b.cfg.controlNamespace
		o.NameOverride = b.cfg.release
		o.Args = b.cfg.operatorArgs
	}))

	for _, ns := range b.cfg.namespaces {
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

// WaitForRunning blocks until every condition holds, and names the one that
// did not if it runs out of budget.
func (e *Env) WaitForRunning(conditions ...wait.Condition) {
	e.T.Helper()

	var outstanding pending
	require.Eventuallyf(e.T, func() bool {
		for _, condition := range conditions {
			met, err := condition.Met(e.Ctx, e.Client)
			if err != nil {
				e.T.Logf("checking %s: %v", condition.Name, err)
			}
			if !met {
				outstanding.set(condition.Name)
				return false
			}
		}
		return true
	}, e.waitBudget(), waitInterval, "still waiting for %s", &outstanding)
}

// WaitForReceiverLogs blocks until the test receiver has logged every tag, and
// names the one still missing if it runs out of budget. The tail itself is not
// echoed each poll: the archived cluster dump already has the receiver's log.
func (e *Env) WaitForReceiverLogs(tags ...string) {
	e.T.Helper()

	var outstanding pending
	require.Eventuallyf(e.T, func() bool {
		logs, err := e.ReceiverLogs(receiverLogTail)
		if err != nil {
			e.T.Logf("reading the test receiver: %v", err)
			return false
		}
		for _, tag := range tags {
			if !strings.Contains(logs, tag) {
				outstanding.set(tag)
				return false
			}
		}
		return true
	}, e.waitBudget(), waitInterval, "the test receiver never logged %s", &outstanding)
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

// waitBudget is what the suites spend today, held under what is left of the
// binary's deadline so one wait cannot take the package's whole budget with it.
func (e *Env) waitBudget() time.Duration {
	deadline, ok := e.T.Deadline()
	if !ok {
		return waitBudget
	}
	return min(waitBudget, max(time.Until(deadline)-teardownMargin, time.Second))
}

// pending is written by the condition, which testify runs on its own goroutine,
// and read when the assertion message is formatted on the test's.
type pending struct {
	mu   sync.Mutex
	name string
}

func (p *pending) set(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.name = name
}

func (p *pending) String() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.name == "" {
		return "the first condition to be checked"
	}
	return p.name
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

func (e *Env) collectArtifacts() {
	path, err := artifactPath(e.T.Name())
	if err != nil {
		e.T.Logf("Skipping cluster logs: %s", err)
	} else {
		e.T.Logf("Printing cluster logs to %s", path)
		assert.NoError(e.T, e.Cluster.PrintLogs(common.PrintLogConfig{
			Namespaces: e.dumpNamespaces,
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

func buildScheme(extra []func(*runtime.Scheme) error) (*runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	for _, add := range slices.Concat(defaultSchemeBuilders, extra) {
		if err := add(scheme); err != nil {
			return nil, err
		}
	}
	return scheme, nil
}

// dumpNamespaces is the control namespace, the configured ones and default,
// without repeats, so a suite that names one does not dump it twice.
func dumpNamespaces(cfg config) []string {
	out := make([]string, 0, len(cfg.namespaces)+2)
	seen := map[string]bool{}
	for _, ns := range slices.Concat([]string{cfg.controlNamespace}, cfg.namespaces, []string{"default"}) {
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
