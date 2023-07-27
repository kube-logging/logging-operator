# A Self-Documenting Makefile: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html

BIN := ${PWD}/bin

export PATH := $(BIN):$(PATH)

OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)

DOCKER ?= docker
GOVERSION := $(shell go env GOVERSION)

# Image name to use for building/pushing image targets
IMG ?= controller:local

# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true,preserveUnknownFields=false,maxDescLen=0"

DRAIN_WATCH_IMAGE_TAG_NAME ?= ghcr.io/kube-logging/fluentd-drain-watch
DRAIN_WATCH_IMAGE_TAG_VERSION ?= latest

VERSION := $(shell git describe --abbrev=0 --tags)

E2E_TEST_TIMEOUT ?= 20m

CONTROLLER_GEN := ${BIN}/controller-gen
CONTROLLER_GEN_VERSION := v0.6.0

ENVTEST_BIN_DIR := ${BIN}/envtest
ENVTEST_K8S_VERSION := 1.24.1
ENVTEST_BINARY_ASSETS := ${ENVTEST_BIN_DIR}/bin

GOLANGCI_LINT := ${BIN}/golangci-lint
GOLANGCI_LINT_VERSION := v1.51.2

HELM_DOCS := ${BIN}/helm-docs
HELM_DOCS_VERSION = 1.11.0

KIND := ${BIN}/kind
KIND_VERSION ?= v0.20.0
KIND_IMAGE ?= kindest/node:v1.23.17@sha256:f77f8cf0b30430ca4128cc7cfafece0c274a118cd0cdb251049664ace0dee4ff
KIND_CLUSTER := kind

KUBEBUILDER := ${BIN}/kubebuilder
KUBEBUILDER_VERSION = v3.1.0

LICENSEI := ${BIN}/licensei
LICENSEI_VERSION = v0.8.0

STERN_VERSION := 1.25.0

SETUP_ENVTEST := ${BIN}/setup-envtest

## =============
## ==  Rules  ==
## =============

.PHONY: all
all: manager

.PHONY: check
check: license-check lint test

.PHONY: check-diff
check-diff: generate fmt manifests helm-docs
	git diff --exit-code ':(exclude)./ADOPTERS.md' ':(exclude)./docs/*'

.PHONY: debug
debug: manager ## Remote debug
	dlv --listen=:40000 --log --headless=true --api-version=2 exec bin/manager -- $(ARGS)

.PHONY: deploy
deploy: manifests ## Deploy controller in the configured Kubernetes cluster in ~/.kube/config
	kubectl apply -f config/crd/bases
	kustomize build config/default | kubectl apply -f -

.PHONY: docker-build
docker-build: ## Build the docker image
	docker image inspect ${IMG} && echo "docker image ${IMG} already built" || ${DOCKER} build . -t ${IMG}
	@echo "updating kustomize image patch file for manager resource"
	sed -i'' -e 's@image: .*@image: '"${IMG}"'@' ./config/default/manager_image_patch.yaml

.PHONY: docker-build-drain-watch
docker-build-drain-watch: ## Build the drain-watch docker image
	${DOCKER} build drain-watch-image -t ${DRAIN_WATCH_IMAGE_TAG_NAME}:${DRAIN_WATCH_IMAGE_TAG_VERSION}

.PHONY: docker-push
docker-push: ## Push the docker image
	${DOCKER} push ${IMG}

.PHONY: docs
docs: ## Generate docs
	go run cmd/docs.go

.PHONY: fmt
fmt: ## Run go fmt against code
	go fmt ./...
	cd pkg/sdk && go fmt ./...

.PHONY: generate
generate: ${CONTROLLER_GEN} tidy ## Generate code
	cd pkg/sdk && $(CONTROLLER_GEN) object:headerFile=./../../hack/boilerplate.go.txt paths=./logging/api/...
	cd pkg/sdk && $(CONTROLLER_GEN) object:headerFile=./../../hack/boilerplate.go.txt paths=./logging/model/...
	cd pkg/sdk && $(CONTROLLER_GEN) object:headerFile=./../../hack/boilerplate.go.txt paths=./extensions/api/...
	cd pkg/sdk && $(CONTROLLER_GEN) object:headerFile=./../../hack/boilerplate.go.txt paths=./resourcebuilder/...
	cd pkg/sdk && go generate ./static

.PHONY: install
install: manifests ## Install CRDs into the cluster in ~/.kube/config
	kubectl create -f config/crd/bases || kubectl replace -f config/crd/bases

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
manager: generate fmt vet ## Build manager binary
	go build -o bin/manager main.go

.PHONY: manifests
manifests: ${CONTROLLER_GEN} ## Generate manifests e.g. CRD, RBAC etc.
	cd pkg/sdk && $(CONTROLLER_GEN) $(CRD_OPTIONS) webhook paths="./..." output:crd:artifacts:config=../../config/crd/bases output:webhook:artifacts:config=../../config/webhook
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role paths="./controllers/..." output:rbac:artifacts:config=./config/rbac
	cp config/crd/bases/* charts/logging-operator/crds/
	echo "{{- if .Values.rbac.enabled }}" > ./charts/logging-operator/templates/clusterrole.yaml
	cat config/rbac/role.yaml | sed -e 's@manager-role@{{ template "logging-operator.fullname" . }}@' | sed -e '/creationTimestamp/d' | cat >> ./charts/logging-operator/templates/clusterrole.yaml
	echo "{{- end }}" >> ./charts/logging-operator/templates/clusterrole.yaml

.PHONY: run
run: generate fmt vet ## Run against the configured Kubernetes cluster in ~/.kube/config
	go run ./main.go --verbose --pprof

test: generate fmt vet manifests ${ENVTEST_BINARY_ASSETS} ${KUBEBUILDER} ## Run tests
	cd pkg/sdk/logging && ENVTEST_BINARY_ASSETS=${ENVTEST_BINARY_ASSETS} go test ./...
	cd pkg/sdk/extensions && go test ./...
	cd pkg/sdk/logging/model/syslogng/config && go test ./...
	ENVTEST_BINARY_ASSETS=${ENVTEST_BINARY_ASSETS} go test ./controllers/logging/... ./pkg/... -coverprofile cover.out
	ENVTEST_BINARY_ASSETS=${ENVTEST_BINARY_ASSETS} go test ./controllers/extensions/... ./pkg/... -coverprofile cover.out

.PHONY: test-e2e
test-e2e: ${KIND} generate manifests docker-build stern ## Run E2E tests
	$(MAKE) test-e2e-nodeps

.PHONY: test-e2e-ci
test-e2e-ci: ${BIN}
	curl -Lo ./bin/kind https://kind.sigs.k8s.io/dl/${KIND_VERSION}/kind-linux-amd64
	chmod +x ./bin/kind
	curl -L https://github.com/stern/stern/releases/download/v${STERN_VERSION}/stern_${STERN_VERSION}_linux_amd64.tar.gz | tar xz -C bin stern
	chmod +x ./bin/stern
	cd e2e && \
		LOGGING_OPERATOR_IMAGE="${IMG}" \
		KIND_PATH="$(KIND)" \
		KIND_IMAGE="$(KIND_IMAGE)" \
		PROJECT_DIR="$(PWD)" \
		go test -v -timeout ${E2E_TEST_TIMEOUT} ./...

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
${CONTROLLER_GEN}_${CONTROLLER_GEN_VERSION}_${GOVERSION}: VERSION := ${CONTROLLER_GEN_VERSION}
${CONTROLLER_GEN}_${CONTROLLER_GEN_VERSION}_${GOVERSION}: | ${BIN}
	${go_install_binary}

${ENVTEST_BINARY_ASSETS}: ${ENVTEST_BINARY_ASSETS}_${ENVTEST_K8S_VERSION}
	ln -sf $(notdir $<) $@

${ENVTEST_BINARY_ASSETS}_${ENVTEST_K8S_VERSION}: | ${SETUP_ENVTEST} ${ENVTEST_BIN_DIR}
	ln -sf $$(${SETUP_ENVTEST} --bin-dir ${ENVTEST_BIN_DIR} use ${ENVTEST_K8S_VERSION} -p path) $@

${GOLANGCI_LINT}: ${GOLANGCI_LINT}_${GOLANGCI_LINT_VERSION}_${GOVERSION} | ${BIN}
	ln -sf $(notdir $<) $@

${GOLANGCI_LINT}_${GOLANGCI_LINT_VERSION}_${GOVERSION}: IMPORT_PATH := github.com/golangci/golangci-lint/cmd/golangci-lint
${GOLANGCI_LINT}_${GOLANGCI_LINT_VERSION}_${GOVERSION}: VERSION := ${GOLANGCI_LINT_VERSION}
${GOLANGCI_LINT}_${GOLANGCI_LINT_VERSION}_${GOVERSION}: | ${BIN}
	${go_install_binary}

${KIND}: ${KIND}_${KIND_VERSION}_${GOVERSION} | ${BIN}
	ln -sf $(notdir $<) $@

${KIND}_${KIND_VERSION}_${GOVERSION}: IMPORT_PATH := sigs.k8s.io/kind
${KIND}_${KIND_VERSION}_${GOVERSION}: VERSION := ${KIND_VERSION}
${KIND}_${KIND_VERSION}_${GOVERSION}: | ${BIN}
	${go_install_binary}

${KUBEBUILDER}: ${KUBEBUILDER}_$(KUBEBUILDER_VERSION) | ${BIN}
	ln -sf $(notdir $<) $@

${KUBEBUILDER}_$(KUBEBUILDER_VERSION): | ${BIN}
	curl -sL https://github.com/kubernetes-sigs/kubebuilder/releases/download/${KUBEBUILDER_VERSION}/kubebuilder_${OS}_${ARCH} -o $@
	chmod +x $@

${LICENSEI}: ${LICENSEI}_${LICENSEI_VERSION}_${GOVERSION} | ${BIN}
	ln -sf $(notdir $<) $@

${LICENSEI}_${LICENSEI_VERSION}_${GOVERSION}: IMPORT_PATH := github.com/goph/licensei/cmd/licensei
${LICENSEI}_${LICENSEI_VERSION}_${GOVERSION}: VERSION := ${LICENSEI_VERSION}
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
