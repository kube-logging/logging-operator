# End-to-end tests

One directory per suite. Each is its own Go package and test binary, and each
provisions its own KIND cluster, installs the operator into it and tears it down
afterwards. Suites do not share a cluster, so they can run concurrently and a
failure in one leaves the others alone.

`common/` and `internal/` are helpers rather than suites, and the Makefile
filters them out of the selection.

## Running them

```bash
make test-e2e                                # everything: builds the images first, ~20 min
make test-e2e-nodeps E2E_TEST=volumedrain    # one suite, reusing the images already built
```

`E2E_TEST` is a **directory path, not a `-run` regex**, so naming a suite still
runs every test in that package.

The loop while working on a suite is `make test-e2e` once, then
`make test-e2e-nodeps` repeatedly — the first builds the six images `image.All()`
names, the second
reuses them. Rebuild only when you change operator code, not test code.

### What has to be installed

`docker`, and `kubectl` and `helm` on `PATH`. `make test-e2e` fetches `kind` and
`stern` into `bin/` for you; `make test-e2e-nodeps` assumes they are already
there.

On Linux, KIND needs more inotify instances than the default 128:

```bash
sudo sysctl -w fs.inotify.max_user_instances=512 fs.inotify.max_user_watches=524288
```

Without this, clusters fail to come up in ways that look like unrelated timeouts.
It resets on reboot.

### Knobs

| variable | default | what it does |
| --- | --- | --- |
| `E2E_TEST` | all | suite directory to run |
| `E2E_TEST_TIMEOUT` | `20m` | per suite binary |
| `E2E_CLUSTERS` | `4` | suite binaries at once (`go test -p`) |
| `E2E_SUITE_PARALLEL` | `2` | clusters one binary builds at once |
| `KIND_COMMAND_TIMEOUT` | derived | per `kind` invocation |

The two parallelism knobs multiply. Raising them starves the aggregators on a
small runner: a suite whose fluentd never finishes its config check usually
means the box is oversubscribed rather than that the test is wrong. A slowdown
in an unrelated suite run at the same time is the cheapest way to confirm that.

### Artifacts

Each suite writes `build/_test/cluster-<TestName>.log` — a `stern` dump of every
watched namespace — and coverage into `build/_test_coverage`. Both survive the
cluster being deleted, and the dump is the first place to look when a wait times
out.

## Writing a suite

`internal/harness` owns the lifecycle. A suite says what it wants and asserts on
it; it does not create clusters, install the operator, or register cleanup.

```go
func TestSomething(t *testing.T) {
	env := harness.New(t).
		WithCluster("something").        // must be unique across all suites
		WithRelease("e2e").
		WithControlNamespace("test").
		Start()                          // marks the test parallel: call it first

	env.Create(&v1beta1.Logging{ /* ... */ })
	env.StartLogProducer("test", map[string]string{"app": "producer"})

	env.WaitFor(
		wait.Operator("e2e"),
		wait.Fluentbit("test"),
		wait.FluentdAggregator("test"),
	)

	env.Receiver.MustReceive("some-tag")
}
```

`Start()` returns an `Env` carrying `T`, `Ctx`, `Client`, `Cluster`, `Release`,
`ControlNamespace` and `Receiver`. Teardown is registered for you and runs in
order: artifacts (the log dump and coverage), the temporary kubeconfig, stopping
the cluster, deleting it. Each step is isolated, so one failing does not strand
the cluster.

### Builder options

| | |
| --- | --- |
| `WithCluster(name)` | KIND cluster name. Required, and see the rule below. |
| `WithRelease(name)` | operator release name; also names the test receiver. |
| `WithControlNamespace(ns)` | where the operator goes. |
| `WithNamespaces(ns...)` | extra namespaces to create and to dump at teardown. |
| `WithOperatorArgs(args...)` | flags for the operator deployment. |
| `WithScheme(add...)` | extra scheme builders, for CRDs outside the default set. |

`WithScheme` has exactly one caller — `logging_metrics_monitoring`, for the
prometheus-operator types whose CRDs only it installs. Everything else is
covered by the defaults.

### Cluster names

Every suite's cluster name must be **distinct**, and KIND only accepts
`[a-z0-9.-]`. `t.Name()` cannot be used directly: test names are mixed-case, and
subtests carry a `/`. Name it after what the suite does, not after the test
function.

### Waiting

`internal/wait` holds the conditions. A `Condition` is a name and a predicate
over a read-only client:

```go
type Condition struct {
	Name string
	Met  func(context.Context, client.Reader) (bool, error)
}
```

`env.WaitFor(conditions...)` polls all of them against the shared budget.
`env.WaitWithin(budget, interval, conditions...)` takes the suite's own, for an
assertion where *how long* is part of the claim — settling within thirty seconds
and settling within five minutes are different statements.

Existing conditions cover pods, deployments, jobs, PVCs, the operator, the
producer, fluentbit, both aggregators, attached and excess configs, and the
`Logging` status. Add one when a suite needs a state none of them names, and
keep it read-only: no `testing.T`, no logging, no mutation. That is what lets
them compose in a single `WaitFor`.

**When a Condition is not enough.** A `Condition` cannot log, so it cannot
report *why* it is still waiting. A suite that needs a per-poll diagnostic
drives the condition from its own `require.Eventually` instead —
`elasticsearch-multiversion` does this to print container restart counts while
it waits. That is the only place in the tree that needs it; if a second suite
wants the same thing, widen the contract rather than copying the pattern.

### The test receiver

The chart installs a receiver that suites send logs to. `env.Receiver` owns it:

```go
env.Receiver.URL("tag")            // address to point an Output at
env.Receiver.MustReceive("tag")    // wait until the tag shows up
env.Receiver.MustNotReceive("tag") // point-in-time check that it did not
env.Receiver.Scale(0)              // take it away, to make an aggregator buffer
```

Prefer these to shelling out to `kubectl` — `#2325` tracks the calls that are
left, and each one is a case where nothing better exists yet.

### Images

`internal/image` names every image the suites load, in one place:

```go
image.Fluentd().Spec()          // v1beta1.ImageSpec
image.ConfigReloader().Basic()  // *v1beta1.BasicImageSpec
image.NodeExporter().Ref()      // "repo:tag"
```

Never write an image reference into a suite. The list `image.All()` returns is
what actually gets loaded into the cluster, so a hardcoded one will not be there.

### What stays in the suite

The `Logging`, `Flow` and `Output` literals. They are each suite's
specification, and defaulted builders for them were tried and removed: no field
turned out to be universal enough to default safely, and the ones that came
close were the ones some suite existed to test.

The rule that replaced them: **extract a value when it is identical across all
callers and no test asserts it.** Image names pass on both counts. `Workers: 2`
does not, because a suite exists to test it.

Repetition *within* one suite is different — three near-identical Deployments
that differ in a name and a version are a table, not a specification.
