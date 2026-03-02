# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Kubernetes operator that manages a complete logging pipeline. It deploys and configures:
- **Fluent Bit** (DaemonSet) — collects logs from every node
- **Fluentd or syslog-ng** (StatefulSet) — aggregates and routes logs to outputs

Users define `Flow`/`Output` CRDs (or their cluster-scoped equivalents) to route logs to destinations like S3, Elasticsearch, Loki, Kafka, etc.

## Project Structure

```
logging-operator/
├── .github/                         # GitHub Actions workflows
├── charts/logging-operator/         # Helm chart (synced with manifests via make manifests)
├── controllers/
│   ├── logging/                     # Core reconcilers (Logging, LoggingRoute, TelemetryController, AxoSyslog)
│   └── extensions/                  # EventTailer and HostTailer reconcilers
├── e2e/                             # KIND-based e2e test suites
├── images/
│   ├── config-reloader/             # Sidecar that hot-reloads on ConfigMap/Secret change
│   ├── fluentd/                     # Custom Fluentd Docker image (Ruby gems, fluent.conf)
│   ├── fluentd-drain-watch/         # Shell script image for buffer drain
│   ├── syslog-ng-reloader/          # Syslog-ng config reload image
│   └── node-exporter/               # Node exporter wrapper image
├── pkg/
│   ├── resources/                   # Kubernetes resource generators
│   ├── sdk/                         # Separate Go module (pkg/sdk/go.mod)
│   │   ├── logging/
│   │   │   ├── api/v1beta1/         # CRD types (Logging, Flow, Output, FluentBitAgent, …)
│   │   │   └── model/
│   │   │       ├── filter/          # 30+ Fluentd filter implementations
│   │   │       ├── output/          # 25+ Fluentd output implementations
│   │   │       └── syslogng/        # syslog-ng config DSL and renderers
│   │   └── extensions/api/          # EventTailer / HostTailer types
│   └── webhook/podhandler/          # Pod mutation webhook (tailer injection)
└── main.go                          # Operator entry point (controller manager setup)
```

## Module Structure

The repository contains five Go modules:
- **Root module** (`go.mod`) — operator binary, controllers, resource generators
- **SDK module** (`pkg/sdk/`) — CRD API types and configuration models
- **E2E module** (`e2e/`) — e2e tests
- **config-reloader** (`images/config-reloader/`) — sidecar binary that watches ConfigMaps/Secrets and triggers webhook reloads when they change
- The other `images/` subdirectories (fluentd, syslog-ng-reloader, etc.) contain Dockerfiles and scripts but no Go modules

When running `go` commands targeting SDK code, `cd pkg/sdk` first.

## Build & Development Commands

```bash
# Build operator binary
make manager

# Generate CRD manifests and RBAC from API types
make manifests

# Regenerate Go code (deepcopy, etc.)
make codegen

# Full generation cycle after any API type change (codegen + fmt + manifests + docs)
make generate

# Format all code
make fmt

# Run go vet
make vet

# Run all linters (golangci-lint)
make lint
make lint-fix   # auto-fix issues

# Run all tests
make test

# Run a single test
cd pkg/sdk && go test -run TestName -v ./logging/...
# or for controller tests:
go test -run TestName -v ./controllers/logging/...

# Check test coverage meets threshold
make check-coverage

# Run e2e tests on KIND cluster
make test-e2e

# Run all checks (license + lint + test)
make check
```

## Architecture

### CRD API (`pkg/sdk/logging/api/v1beta1/`)

Core resource types:
- **`Logging`** — cluster-scoped root resource; references the FluentBit agent and Fluentd/SyslogNG aggregator configs
- **`Flow` / `ClusterFlow`** — namespace-scoped / cluster-scoped routing rules with match selectors and filter chains
- **`Output` / `ClusterOutput`** — log destinations (25+ types: S3, GCS, Azure, Elasticsearch, Loki, Kafka, CloudWatch, Splunk, etc.)
- **`FluentBitAgent`** — standalone DaemonSet spec for log collection
- **`FluentdConfig` / `SyslogNGConfig`** — detached aggregator configs (v4.5+, decoupled from `Logging`)
- **`LoggingRoute`** — connects collectors in one logging domain to aggregators in another (multi-tenancy)
- **`SyslogNGFlow`, `SyslogNGClusterFlow`, `SyslogNGOutput`, `SyslogNGClusterOutput`** — syslog-ng equivalents

v1alpha1 is the legacy API; conversion functions exist to v1beta1 (which is the storage version).

### Controllers (`controllers/logging/`)

- **`LoggingReconciler`** — primary reconciler; orchestrates Fluent Bit DaemonSet, Fluentd/SyslogNG StatefulSet, ConfigMaps, Secrets, RBAC, Services, and Prometheus integrations
- **`LoggingRouteController`** — manages `LoggingRoute` cross-domain routing
- **`TelemetryControllerController`** — integrates with `kube-logging/telemetry-controller`
- **`AxoSyslogController`** — AxoSyslog-specific handling
- Extensions controllers in `controllers/extensions/` handle `EventTailer` and `HostTailer`

### Resource Generation (`pkg/resources/`)

- `model/` — `LoggingResourceRepository` (gathers all Flows/Outputs for a Logging instance), core reconciliation orchestrator
- `fluentd/` — generates Fluentd config and Kubernetes resources
- `fluentbit/` — generates Fluent Bit DaemonSet and config
- `syslogng/` — generates syslog-ng config and resources
- `telemetry-controller/` — telemetry integration resources

### Configuration Models (`pkg/sdk/logging/model/`)

- `filter/` — 30+ filter implementations (parser, grep, record_transformer, dedot, throttle, etc.)
- `output/` — 25+ output plugins
- `syslogng/` — syslog-ng-specific config generation

### Key Patterns

**Repository Pattern**: `model.NewLoggingResourceRepository()` abstracts querying all related resources for a `Logging` instance (Flows, Outputs, FluentdConfig, SyslogNGConfig, LoggingRoutes).

**Reconciler Chain**: Uses `cisco-open/operator-tools` for chaining multiple reconcilers, immutable field detection, and StatefulSet recreation on config changes.

**Finalizers**: Used on `FluentdConfig` and `SyslogNGConfig` to gracefully remove configurations from aggregators before deletion.

**Multi-tenancy**: `WatchNamespaces` / `WatchNamespaceSelector` on the `Logging` resource scopes which namespaces' Flows and Outputs are processed. `LoggingRoute` connects separate logging domains.

**`loggingRef` label**: When multiple `Logging` resources exist in a cluster, child resources (`Flow`, `Output`, `FluentBitAgent`, etc.) are associated with the right instance via a matching `loggingRef` field. Resources without a `loggingRef` belong to the default (empty-string) `Logging` instance.

**Immutable fields**: `controlNamespace`, `FluentbitAgentNamespace`, and `AllowClusterResourcesFromAllNamespaces` are CRD-enforced immutable on `Logging` resources (via `XValidation`). Changing them requires deleting and recreating the resource.

## Testing

- Unit/integration tests use `envtest` (embedded Kubernetes API server + etcd); no cluster needed
- E2E tests in `e2e/` use KIND and cover scenarios: fluentd-aggregator, fluentbit-multitenant, syslog-ng-aggregator
- Coverage profile config in `.testcoverage.yml`; tool: `go-test-coverage`
- Test framework: Ginkgo + Gomega for BDD-style tests; testify for unit tests

## Key Configuration

Operator flags (set in Deployment args):
- `--watch-namespace` — restrict to specific namespaces
- `--watch-logging-name` — restrict to a specific Logging resource
- `--enable-leader-election` — required for HA
- `--verbose` — debug logging

Environment variables:
- `ENABLE_WEBHOOKS=true` — enable pod mutation webhook
- `WEBHOOK_PORT` — webhook server port
- `GOCOVERDIR` — coverage output for e2e coverage builds
