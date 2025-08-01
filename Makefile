# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

####
##  Dependency versions
####

# renovate: datasource=github-releases depName=kubernetes-sigs/controller-tools versioning=semver
CONTROLLER_GEN_VERSION := 0.18.0

# renovate: datasource=github-releases depName=golangci/golangci-lint versioning=semver
GOLANGCI_LINT_VERSION := 2.3.0

# renovate: datasource=go depName=github.com/vladopajic/go-test-coverage/v2 versioning=semver
GO_TEST_COVERAGE_VERSION := 2.16.0

# renovate: datasource=github-releases depName=norwoodj/helm-docs versioning=semver
HELM_DOCS_VERSION = 1.14.2

# renovate: datasource=github-releases depName=kubernetes-sigs/kind versioning=semver
KIND_VERSION ?= 0.29.0

# renovate: datasource=github-releases depName=kubernetes-sigs/kubebuilder versioning=semver
KUBEBUILDER_VERSION = 4.7.0

# renovate: datasource=go depName=github.com/goph/licensei versioning=semver
LICENSEI_VERSION = 0.9.0

# renovate: datasource=go depName=github.com/stern/stern versioning=semver
STERN_VERSION := 1.32.0

ENVTEST_K8S_VERSION := 1.33.0

BIN := ${PWD}/bin

export PATH := $(BIN):$(PATH)

OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)

DOCKER ?= docker
GOVERSION := $(shell go env GOVERSION)

# Image name to use for building/pushing image targets
FLUENTD_IMG ?= fluentd-full:local
CONFIG_RELOADER_IMG ?= config-reloader:local
SYSLOG_NG_RELOADER_IMG ?= syslog-ng-reloader:local
FLUENTD_DRAIN_WATCH_IMG ?= fluentd-drain-watch:local
NODE_EXPORTER_IMG ?= node-exporter:local
OPERATOR_IMG ?= controller:local
OPERATOR_IMG_DEBUG ?= controller:debug

# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= crd:maxDescLen=0

VERSION := $(shell git describe --abbrev=0 --tags)

E2E_TEST_TIMEOUT ?= 20m
TEST_COV_DIR := $(shell mkdir -p build/_test_coverage && realpath build/_test_coverage)

CONTROLLER_GEN := ${BIN}/controller-gen

ENVTEST_BIN_DIR := ${BIN}/envtest
ENVTEST_BINARY_ASSETS := ${ENVTEST_BIN_DIR}/bin

GOLANGCI_LINT := ${BIN}/golangci-lint
LINTER_FLAGS := --timeout 10m

GO_TEST_COVERAGE := ${BIN}/go-test-coverage

HELM_DOCS := ${BIN}/helm-docs

KIND := ${BIN}/kind
KIND_IMAGE ?= kindest/node:v1.33.0@sha256:02f73d6ae3f11ad5d543f16736a2cb2a63a300ad60e81dac22099b0b04784a4e
KIND_CLUSTER := kind

KUBEBUILDER := ${BIN}/kubebuilder

LICENSEI := ${BIN}/licensei

SETUP_ENVTEST := ${BIN}/setup-envtest

## =============
## ==  Rules  ==
## =============

.PHONY: all
all: manager

.PHONY: check
check: license-check lint test

.PHONY: generate
generate: codegen fmt manifests docs helm-docs ## Generate code, documentation, etc

.PHONY: check-diff
check-diff: generate
	git diff --exit-code ':(exclude)./ADOPTERS.md'

.PHONY: debug
debug: manager ## Remote debug
	dlv --listen=:40000 --log --headless=true --api-version=2 exec bin/manager -- $(ARGS)

.PHONY: docker-build
docker-build: ## Build the docker image
	${DOCKER} build . -t ${OPERATOR_IMG}
	@echo "updating kustomize image patch file for manager resource"
	sed -i'' -e 's@image: .*@image: '"${OPERATOR_IMG}"'@' ./config/default/manager_image_patch.yaml

.PHONY: docker-build-debug
docker-build-debug: ## Build the debug docker image
	${DOCKER} build --target debug -t ${OPERATOR_IMG_DEBUG} .

.PHONY: docker-build-e2e-test
docker-build-e2e-test: ## Build the coverage docker image
	${DOCKER} build --build-arg GO_BUILD_FLAGS="-cover -covermode=atomic" -t ${OPERATOR_IMG} --target e2e-test .
	sed -i'' -e 's@image: .*@image: '"${OPERATOR_IMG}"'@' ./config/default/manager_image_patch.yaml
	${DOCKER} build -t ${CONFIG_RELOADER_IMG} images/config-reloader
	${DOCKER} build -t ${SYSLOG_NG_RELOADER_IMG} images/syslog-ng-reloader
	${DOCKER} build -t ${FLUENTD_DRAIN_WATCH_IMG} images/fluentd-drain-watch
	${DOCKER} build -t ${NODE_EXPORTER_IMG} images/node-exporter
	${DOCKER} build -t ${FLUENTD_IMG} --target full images/fluentd

.PHONY: docker-build-drain-watch
docker-build-drain-watch: ## Build the drain-watch docker image
	${DOCKER} build drain-watch-image -t ${FLUENTD_DRAIN_WATCH_IMG} images/fluentd-drain-watch

.PHONY: docker-push
docker-push: ## Push the docker image
	${DOCKER} push ${OPERATOR_IMG}

.PHONY: docs
docs: ## Generate docs
	go run cmd/docs.go

.PHONY: fmt
fmt: ## Run go fmt against code
	go fmt ./...
	cd pkg/sdk && go fmt ./...

.PHONY: codegen
codegen: ${CONTROLLER_GEN} tidy ## Generate code
	cd pkg/sdk && $(CONTROLLER_GEN) $(CRD_OPTIONS) output:crd:artifacts:config=../../config/crd/bases object:headerFile=./../../hack/boilerplate.go.txt paths=./logging/api/...
	cd pkg/sdk && $(CONTROLLER_GEN) $(CRD_OPTIONS) output:crd:artifacts:config=../../config/crd/bases object:headerFile=./../../hack/boilerplate.go.txt paths=./logging/model/...
	cd pkg/sdk && $(CONTROLLER_GEN) $(CRD_OPTIONS) output:crd:artifacts:config=../../config/crd/bases object:headerFile=./../../hack/boilerplate.go.txt paths=./extensions/api/...

.PHONY: install
install: manifests ## Install CRDs into the cluster in ~/.kube/config
	kubectl apply -f config/crd/bases --server-side --force-conflicts

.PHONY: license-check
license-check: ${LICENSEI} .licensei.cache ## Run license check
	${LICENSEI} check
	${LICENSEI} header

.PHONY: license-cache
license-cache: ${LICENSEI} ## Generate license cache
	${LICENSEI} cache

.PHONY: lint
lint: ${GOLANGCI_LINT} ## Run linter
	${GOLANGCI_LINT} run ${LINTER_FLAGS}
	cd pkg/sdk && ${GOLANGCI_LINT} run ${LINTER_FLAGS}

.PHONY: lint-fix
lint-fix: ${GOLANGCI_LINT} ## Run linter
	${GOLANGCI_LINT} run --fix
	cd pkg/sdk && ${GOLANGCI_LINT} run --fix

.PHONY: list
list: ## List all make targets
	@${MAKE} -pRrn : -f $(MAKEFILE_LIST) 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | sort

.PHONY: manager
manager: codegen fmt vet ## Build manager binary
	go build -o bin/manager main.go

HELM_MANIFEST_OVERRIDE='s@manager-role@{{ template "logging-operator.fullname" . }}\n  annotations:\n  {{- if .Values.rbac.retainOnDelete }}\n    "helm.sh/resource-policy": keep\n  {{- end }}\n  labels:\n{{ include "logging-operator.labels" . | indent 4 }}@'

.PHONY: manifests
manifests: ${CONTROLLER_GEN} ## Generate manifests e.g. CRD, RBAC etc.
	cd pkg/sdk && $(CONTROLLER_GEN) $(CRD_OPTIONS)  webhook paths="./..." output:crd:artifacts:config=../../config/crd/bases output:webhook:artifacts:config=../../config/webhook
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role paths="./controllers/..." output:rbac:artifacts:config=./config/rbac
	cp config/crd/bases/* charts/logging-operator/crds/
	for f in config/crd/bases/*.yaml; do sed '/controller-gen.kubebuilder.io\/version/ r hack/crds.annotations.snippet.txt' $${f} > charts/logging-operator/charts/logging-operator-crds/templates/$${f##*/}; done
	echo "{{- if .Values.rbac.enabled }}" > ./charts/logging-operator/templates/clusterrole.yaml
	cat config/rbac/role.yaml | sed -e $(HELM_MANIFEST_OVERRIDE) | cat >> ./charts/logging-operator/templates/clusterrole.yaml
	echo "{{- end }}" >> ./charts/logging-operator/templates/clusterrole.yaml

.PHONY: run
run: codegen fmt vet ## Run against the configured Kubernetes cluster in ~/.kube/config
	go run ./main.go --verbose --pprof

.PHONY: test
test: codegen fmt vet manifests ${ENVTEST_BINARY_ASSETS} ${KUBEBUILDER} ## Run tests
	cd pkg/sdk/logging && ENVTEST_BINARY_ASSETS=${ENVTEST_BINARY_ASSETS} go test ./... -coverprofile ${TEST_COV_DIR}/cover_logging.out
	cd pkg/sdk/extensions && go test ./... -coverprofile  ${TEST_COV_DIR}/cover_extensions.out
	cd pkg/sdk/logging/model/syslogng/config && go test ./...  -coverprofile ${TEST_COV_DIR}/cover_syslogng.out
	ENVTEST_BINARY_ASSETS=${ENVTEST_BINARY_ASSETS} go test ./controllers/logging/... ./pkg/...  -coverprofile ${TEST_COV_DIR}/cover_controllers_logging.out
	ENVTEST_BINARY_ASSETS=${ENVTEST_BINARY_ASSETS} go test ./controllers/extensions/... ./pkg/...  -coverprofile ${TEST_COV_DIR}/cover_controllers_extensions.out

.PHONY: generate-test-coverage
generate-test-coverage: test
	rm -f ${TEST_COV_DIR}/coverage_all.out
	echo "mode: set" > ${TEST_COV_DIR}/coverage_all.out
	find -name 'cover_*.out' | xargs cat | grep -v "mode: set" >> ${TEST_COV_DIR}/coverage_all.out

.PHONY: check-coverage
check-coverage: generate-test-coverage
	go-test-coverage --config=./.testcoverage.yml

.PHONY: test-e2e-coverage-report-ci
test-e2e-coverage-report-ci: ${GO_TEST_COVERAGE}
	$(MAKE) test-e2e-coverage-report-no-deps

.PHONY: test-e2e-coverage-report-no-deps
test-e2e-coverage-report-no-deps:
	${GO_TEST_COVERAGE} --profile=${TEST_COV_DIR}/coverage_e2e.out

.PHONY: test-e2e
test-e2e: ${KIND} codegen manifests docker-build-e2e-test stern ## Run E2E tests
	$(MAKE) test-e2e-nodeps E2E_TEST=${E2E_TEST}

.PHONY: test-e2e-ci
test-e2e-ci: ${BIN}
	curl -Lo ./bin/kind https://kind.sigs.k8s.io/dl/v${KIND_VERSION}/kind-linux-amd64
	chmod +x ./bin/kind
	curl -L https://github.com/stern/stern/releases/download/v${STERN_VERSION}/stern_${STERN_VERSION}_linux_amd64.tar.gz | tar xz -C bin stern
	chmod +x ./bin/stern
	$(MAKE) test-e2e-nodeps E2E_TEST=${E2E_TEST}

.PHONY: test-e2e-nodeps
test-e2e-nodeps:
	cd e2e && \
		LOGGING_OPERATOR_IMAGE="${OPERATOR_IMG}" \
		CONFIG_RELOADER_IMAGE="${CONFIG_RELOADER_IMG}" \
		SYSLOG_NG_RELOADER_IMAGE="${SYSLOG_NG_RELOADER_IMG}" \
		FLUENTD_DRAIN_WATCH_IMAGE="${FLUENTD_DRAIN_WATCH_IMG}" \
		NODE_EXPORTER_IMAGE="${NODE_EXPORTER_IMG}" \
		FLUENTD_IMAGE="${FLUENTD_IMG}" \
		KIND_PATH="$(KIND)" \
		KIND_IMAGE="$(KIND_IMAGE)" \
		PROJECT_DIR="$(PWD)" \
		E2E_TEST_COV_DIR=${TEST_COV_DIR} \
		go test -count=1 -v -timeout ${E2E_TEST_TIMEOUT} ./${E2E_TEST}/...
		go tool covdata textfmt -i=${TEST_COV_DIR}/covdatafiles -o ${TEST_COV_DIR}/coverage_e2e.out
	@echo "--- E2E test coverage report"
	go tool covdata percent -i=${TEST_COV_DIR}/covdatafiles

.PHONY: tidy
tidy: ## Tidy Go modules
	find . -iname "go.mod" -not -path "./.devcontainer/*" | xargs -L1 sh -c 'cd $$(dirname $$0); go mod tidy'

.PHONY: vet
vet: ## Run go vet against code
	go vet ./...
	cd pkg/sdk && go vet ./...

.PHONY: kind-cluster
kind-cluster: ${KIND}
	kind create cluster --name $(KIND_CLUSTER) --image $(KIND_IMAGE)

.PHONY: helm-docs
helm-docs: ${HELM_DOCS}
	${HELM_DOCS} -s file -c charts/ -t ../charts-docs/templates/overrides.gotmpl -t README.md.gotmpl

## =========================
## ==  Tool dependencies  ==
## =========================

${CONTROLLER_GEN}: ${CONTROLLER_GEN}_${CONTROLLER_GEN_VERSION}_${GOVERSION} | ${BIN}
	ln -sf $(notdir $<) $@

${CONTROLLER_GEN}_${CONTROLLER_GEN_VERSION}_${GOVERSION}: IMPORT_PATH := sigs.k8s.io/controller-tools/cmd/controller-gen
${CONTROLLER_GEN}_${CONTROLLER_GEN_VERSION}_${GOVERSION}: VERSION := v${CONTROLLER_GEN_VERSION}
${CONTROLLER_GEN}_${CONTROLLER_GEN_VERSION}_${GOVERSION}: | ${BIN}
	${go_install_binary}

${ENVTEST_BINARY_ASSETS}: ${ENVTEST_BINARY_ASSETS}_${ENVTEST_K8S_VERSION}
	ln -sf $(notdir $<) $@

${ENVTEST_BINARY_ASSETS}_${ENVTEST_K8S_VERSION}: | ${SETUP_ENVTEST} ${ENVTEST_BIN_DIR}
	ln -sf $$(${SETUP_ENVTEST} --bin-dir ${ENVTEST_BIN_DIR} use ${ENVTEST_K8S_VERSION} -p path) $@

${GOLANGCI_LINT}: ${GOLANGCI_LINT}_${GOLANGCI_LINT_VERSION}_${GOVERSION} | ${BIN}
	ln -sf $(notdir $<) $@

${GOLANGCI_LINT}_${GOLANGCI_LINT_VERSION}_${GOVERSION}: IMPORT_PATH := github.com/golangci/golangci-lint/v2/cmd/golangci-lint
${GOLANGCI_LINT}_${GOLANGCI_LINT_VERSION}_${GOVERSION}: VERSION := v${GOLANGCI_LINT_VERSION}
${GOLANGCI_LINT}_${GOLANGCI_LINT_VERSION}_${GOVERSION}: | ${BIN}
	${go_install_binary}

${GO_TEST_COVERAGE}: ${GO_TEST_COVERAGE}_${GO_TEST_COVERAGE_VERSION}_${GOVERSION} | ${BIN}

${GO_TEST_COVERAGE}_${GO_TEST_COVERAGE_VERSION}_${GOVERSION}: IMPORT_PATH := github.com/vladopajic/go-test-coverage/v2
${GO_TEST_COVERAGE}_${GO_TEST_COVERAGE_VERSION}_${GOVERSION}: VERSION := v${GO_TEST_COVERAGE_VERSION}
${GO_TEST_COVERAGE}_${GO_TEST_COVERAGE_VERSION}_${GOVERSION}: | ${BIN}
	GOBIN=${BIN} go install ${IMPORT_PATH}@${VERSION}

${KIND}: ${KIND}_${KIND_VERSION}_${GOVERSION} | ${BIN}
	ln -sf $(notdir $<) $@

${KIND}_${KIND_VERSION}_${GOVERSION}: IMPORT_PATH := sigs.k8s.io/kind
${KIND}_${KIND_VERSION}_${GOVERSION}: VERSION := v${KIND_VERSION}
${KIND}_${KIND_VERSION}_${GOVERSION}: | ${BIN}
	${go_install_binary}

${KUBEBUILDER}: ${KUBEBUILDER}_$(KUBEBUILDER_VERSION) | ${BIN}
	ln -sf $(notdir $<) $@

${KUBEBUILDER}_$(KUBEBUILDER_VERSION): | ${BIN}
	curl -sL https://github.com/kubernetes-sigs/kubebuilder/releases/download/v${KUBEBUILDER_VERSION}/kubebuilder_${OS}_${ARCH} -o $@
	chmod +x $@

${LICENSEI}: ${LICENSEI}_${LICENSEI_VERSION}_${GOVERSION} | ${BIN}
	ln -sf $(notdir $<) $@

${LICENSEI}_${LICENSEI_VERSION}_${GOVERSION}: IMPORT_PATH := github.com/goph/licensei/cmd/licensei
${LICENSEI}_${LICENSEI_VERSION}_${GOVERSION}: VERSION := v${LICENSEI_VERSION}
${LICENSEI}_${LICENSEI_VERSION}_${GOVERSION}: | ${BIN}
	${go_install_binary}

.licensei.cache: ${LICENSEI}
ifndef GITHUB_TOKEN
	@>&2 echo "WARNING: building licensei cache without Github token, rate limiting might occur."
	@>&2 echo "(Hint: If too many licenses are missing, try specifying a Github token via the environment variable GITHUB_TOKEN.)"
endif
	${LICENSEI} cache

${SETUP_ENVTEST}: IMPORT_PATH := sigs.k8s.io/controller-runtime/tools/setup-envtest
${SETUP_ENVTEST}: VERSION := latest
${SETUP_ENVTEST}: | ${BIN}
	GOBIN=${BIN} go install ${IMPORT_PATH}@${VERSION}

stern: | ${BIN}
	GOBIN=${BIN} go install github.com/stern/stern@latest

${ENVTEST_BIN_DIR}: | ${BIN}
	mkdir -p $@

${HELM_DOCS}: ${HELM_DOCS}-${HELM_DOCS_VERSION}
	@ln -sf ${HELM_DOCS}-${HELM_DOCS_VERSION} ${HELM_DOCS}
${HELM_DOCS}-${HELM_DOCS_VERSION}:
	@mkdir -p bin
	curl -L https://github.com/norwoodj/helm-docs/releases/download/v${HELM_DOCS_VERSION}/helm-docs_${HELM_DOCS_VERSION}_$(shell uname)_x86_64.tar.gz | tar -zOxf - helm-docs > ${HELM_DOCS}-${HELM_DOCS_VERSION} && chmod +x ${HELM_DOCS}-${HELM_DOCS_VERSION}

${BIN}:
	mkdir -p bin

define go_install_binary
find ${BIN} -name '$(notdir ${IMPORT_PATH})_*' -exec rm {} +
GOBIN=${BIN} go install ${IMPORT_PATH}@${VERSION}
mv ${BIN}/$(notdir ${IMPORT_PATH}) $@
endef

# Self-documenting Makefile
.DEFAULT_GOAL = help
.PHONY: help
help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
