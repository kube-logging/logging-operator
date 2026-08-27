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
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

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
	// Names a runnable that will not stop, rather than letting it run the
	// binary out of its -timeout.
	clusterStopTimeout = time.Minute

	clusterLogLimit = 100 * 1000

	// waitBudget and waitInterval are what the suites spend by hand today.
	waitBudget   = 5 * time.Minute
	waitInterval = 3 * time.Second

	// Leaves a wait that runs to the end time to tear down.
	teardownMargin = 90 * time.Second
)

// Everything the suites register except the prometheus-operator types, which
// only logging_metrics_monitoring needs.
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

// WithCluster is named rather than derived from t.Name(): kind accepts only
// [a-z0-9.-] and none of the test names match.
func (b *Builder) WithCluster(name string) *Builder {
	b.cfg.cluster = name
	return b
}

func (b *Builder) WithRelease(name string) *Builder {
	b.cfg.release = name
	return b
}

func (b *Builder) WithControlNamespace(namespace string) *Builder {
	b.cfg.controlNamespace = namespace
	return b
}

// WithNamespaces are also dumped at teardown.
func (b *Builder) WithNamespaces(namespaces ...string) *Builder {
	b.cfg.namespaces = append(b.cfg.namespaces, namespaces...)
	return b
}

func (b *Builder) WithOperatorArgs(args ...string) *Builder {
	b.cfg.operatorArgs = append(b.cfg.operatorArgs, args...)
	return b
}

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
	Receiver         Receiver

	dumpNamespaces []string
}

// Start marks the test parallel, so it has to be the first call in the test.
func (b *Builder) Start() *Env {
	t := b.t
	common.Initialize(t)

	scheme, err := buildScheme(b.cfg.schemeBuilders)
	common.RequireNoError(t, err)

	c, err := common.GetTestCluster(b.cfg.cluster, func(o *cluster.Options) { o.Scheme = scheme })
	if err != nil {
		// The kind cluster is up before the client can fail, and teardown is
		// not registered yet.
		common.DeleteTestClusterOrLog(t, b.cfg.cluster)
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
	env.Receiver = Receiver{env: env}

	teardown{
		{"artifacts", env.collectArtifacts},
		{"kubeconfig", func() { assert.NoError(t, c.Cleanup()) }},
		{"stop", func() { stopCluster(t, cancel, startErr) }},
		{"delete", func() { common.DeleteTestClusterOrLog(t, b.cfg.cluster) }},
	}.register(t)

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

// StartLogProducer runs on the test goroutine, where FailNow is defined.
func (e *Env) StartLogProducer(namespace string, labels map[string]string) {
	e.T.Helper()
	setup.LogProducer(e.T, e.Client, setup.LogProducerOptionFunc(func(o *setup.LogProducerOptions) {
		o.Namespace = namespace
		o.Labels = labels
	}))
}

func (e *Env) WaitFor(conditions ...wait.Condition) {
	e.T.Helper()
	e.WaitWithin(e.waitBudget(), waitInterval, conditions...)
}

// WaitWithin is WaitFor with the suite's own budget, for an assertion where how
// long it takes is part of what is being tested: settling within thirty seconds
// and settling within the shared five minutes are different claims.
//
// What neither takes is a diagnostic on each failed poll, since a Condition is
// read-only by construction and cannot log. A suite that needs one drives the
// Condition from its own require.Eventually instead, as
// elasticsearch-multiversion does to report container restarts while it waits.
// That is the one escape hatch the thirteen suites have needed; widen the
// Condition contract only if a second suite wants the same thing.
func (e *Env) WaitWithin(budget, interval time.Duration, conditions ...wait.Condition) {
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
	}, budget, interval, "still waiting for %s", &outstanding)
}

// waitBudget holds a wait under what is left of the binary's deadline, so one
// cannot take the package's whole budget with it.
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

type step struct {
	name string
	run  func()
}

type teardown []step

type cleanupT interface {
	Cleanup(func())
	Errorf(format string, args ...any)
}

// register goes back to front because Cleanup is LIFO, and isolates each step:
// Go abandons the rest at the first panic, which would strand the cluster.
func (td teardown) register(t cleanupT) {
	for _, s := range slices.Backward(td) {
		t.Cleanup(func() {
			defer func() {
				if v := recover(); v != nil {
					assert.Fail(t, fmt.Sprintf("teardown step %q panicked: %v", s.name, v))
				}
			}()
			s.run()
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
		// Logged, never fatal: coverage is not the suite's verdict.
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

func dumpNamespaces(cfg config) []string {
	out := make([]string, 0, len(cfg.namespaces)+2)
	seen := map[string]bool{}
	// Concat, not append: appending to cfg.namespaces could write into the
	// caller's backing array.
	for _, ns := range slices.Concat([]string{cfg.controlNamespace}, cfg.namespaces, []string{"default"}) {
		if ns != "" && !seen[ns] {
			seen[ns] = true
			out = append(out, ns)
		}
	}
	return out
}

func artifactPath(name string) (string, error) {
	root, ok := os.LookupEnv("PROJECT_DIR")
	if !ok {
		root = "../.."
	}
	dir := filepath.Join(root, "build/_test")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	// A subtest's name carries a slash, which Join would read as a directory
	// that nothing creates, so the dump would fail to open.
	return filepath.Join(dir, fmt.Sprintf("cluster-%s.log", strings.ReplaceAll(name, "/", "_"))), nil
}
