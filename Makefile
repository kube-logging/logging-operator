OS = $(shell uname | tr A-Z a-z)
SHELL := /bin/bash

# Image URL to use all building/pushing image targets
IMG ?= controller:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true,preserveUnknownFields=false,maxDescLen=0"
DRAIN_WATCH_IMAGE_TAG_NAME ?= ghcr.io/banzaicloud/fluentd-drain-watch
DRAIN_WATCH_IMAGE_TAG_VERSION ?= latest

CONTROLLER_GEN_VERSION = v0.5.0
GOLANGCI_VERSION = v1.42.1
KIND_VERSION = v0.11.1
LICENSEI_VERSION = v0.3.1
ENVTEST_CTRL_VERSION = v0.8.3

VERSION := $(shell git describe --abbrev=0 --tags)
DOCKER_IMAGE = banzaicloud/logging-operator
DOCKER_TAG ?= ${VERSION}

CONTROLLER_GEN = $(PWD)/bin/controller-gen
ENVTEST_ASSETS_DIR=${PWD}/testbin
export PATH := $(PWD)/bin:$(PATH)

E2E_TEST_TIMEOUT ?= 20m

.PHONY: all
all: manager

.PHONY: check
check: license-check lint test

.PHONY: check-diff
check-diff: generate fmt manifests
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
	docker build . -t ${IMG}
	@echo "updating kustomize image patch file for manager resource"
	sed -i'' -e 's@image: .*@image: '"${IMG}"'@' ./config/default/manager_image_patch.yaml

.PHONY: docker-build-drain-watch
docker-build-drain-watch: ## Build the drain-watch docker image
	docker build drain-watch-image -t ${DRAIN_WATCH_IMAGE_TAG_NAME}:${DRAIN_WATCH_IMAGE_TAG_VERSION}

.PHONY: docker-push
docker-push: ## Push the docker image
	docker push ${IMG}

.PHONY: docs
docs: ## Generate docs
	go run cmd/docs.go

.PHONY: fmt
fmt: ## Run go fmt against code
	go fmt ./...
	cd pkg/sdk && go fmt ./...

.PHONY: generate
generate: bin/controller-gen tidy ## Generate code
	cd pkg/sdk && $(CONTROLLER_GEN) object:headerFile=./../../hack/boilerplate.go.txt paths=./logging/api/...
	cd pkg/sdk && $(CONTROLLER_GEN) object:headerFile=./../../hack/boilerplate.go.txt paths=./logging/model/...
	cd pkg/sdk && $(CONTROLLER_GEN) object:headerFile=./../../hack/boilerplate.go.txt paths=./extensions/api/...
	cd pkg/sdk && $(CONTROLLER_GEN) object:headerFile=./../../hack/boilerplate.go.txt paths=./resourcebuilder/...
	cd pkg/sdk && go generate ./static

.PHONY: install
install: manifests ## Install CRDs into the cluster in ~/.kube/config
	kubectl create -f config/crd/bases || kubectl replace -f config/crd/bases

.PHONY: license-check
license-check: bin/licensei .licensei.cache ## Run license check
	bin/licensei check
	./scripts/check-header.sh

.PHONY: license-cache
license-cache: bin/licensei ## Generate license cache
	bin/licensei cache

.PHONY: lint
lint: bin/golangci-lint ## Run linter
	bin/golangci-lint run --timeout 2m
	cd pkg/sdk && ../../bin/golangci-lint run  -c ../../.golangci.yml

.PHONY: lint-fix
lint-fix: bin/golangci-lint ## Run linter
	bin/golangci-lint run --fix
	cd pkg/sdk && golangci-lint run -c ../../.golangci.yml --fix

.PHONY: list
list: ## List all make targets
	@${MAKE} -pRrn : -f $(MAKEFILE_LIST) 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | sort

.PHONY: manager
manager: generate fmt vet ## Build manager binary
	go build -o bin/manager main.go

.PHONY: manifests
manifests: bin/controller-gen ## Generate manifests e.g. CRD, RBAC etc.
	cd pkg/sdk && $(CONTROLLER_GEN) $(CRD_OPTIONS) webhook paths="./..." output:crd:artifacts:config=../../config/crd/bases output:webhook:artifacts:config=../../config/webhook
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role paths="./controllers/..." output:rbac:artifacts:config=./config/rbac
	cp config/crd/bases/* charts/logging-operator/crds/
	echo "{{- if .Values.rbac.enabled }}" > ./charts/logging-operator/templates/clusterrole.yaml && cat config/rbac/role.yaml |sed -e 's@manager-role@{{ template "logging-operator.fullname" . }}@' | cat >> ./charts/logging-operator/templates/clusterrole.yaml && echo "{{- end }}" >> ./charts/logging-operator/templates/clusterrole.yaml

.PHONY: run
run: generate fmt vet ## Run against the configured Kubernetes cluster in ~/.kube/config
	go run ./main.go --verbose --pprof

test: generate fmt vet manifests | ${ENVTEST_ASSETS_DIR}/setup-envtest.sh ## Run tests
	source ${ENVTEST_ASSETS_DIR}/setup-envtest.sh; fetch_envtest_tools ${ENVTEST_ASSETS_DIR}; setup_envtest_env ${ENVTEST_ASSETS_DIR}
	cd pkg/sdk/logging && go test ./...
	cd pkg/sdk/extensions && go test ./...
	go test ./controllers/logging/... ./pkg/... -coverprofile cover.out
	go test ./controllers/extensions/... ./pkg/... -coverprofile cover.out

.PHONY: test-e2e
test-e2e: bin/kind docker-build generate fmt vet manifests ## Run E2E tests
	cd e2e && LOGGING_OPERATOR_IMAGE="${IMG}" go test -timeout ${E2E_TEST_TIMEOUT} ./...

.PHONY: tidy
tidy: ## Tidy Go modules
	find . -iname "go.mod" | xargs -L1 sh -c 'cd $$(dirname $$0); go mod tidy -compat=1.17'

.PHONY: vet
vet: ## Run go vet against code
	go vet ./...
	cd pkg/sdk && go vet ./...

.licensei.cache: bin/licensei
ifndef GITHUB_TOKEN
	@>&2 echo "WARNING: building licensei cache without Github token, rate limiting might occur."
	@>&2 echo "(Hint: If too many licenses are missing, try specifying a Github token via the environment variable GITHUB_TOKEN.)"
endif
	bin/licensei cache

bin/controller-gen: | bin/controller-gen_${CONTROLLER_GEN_VERSION} bin
	ln -sf controller-gen_${CONTROLLER_GEN_VERSION} $@

bin/controller-gen_${CONTROLLER_GEN_VERSION}: | bin
	find $(PWD)/bin -name 'controller-gen*' -exec rm {} +
	GOBIN=$(PWD)/bin go install sigs.k8s.io/controller-tools/cmd/controller-gen@${CONTROLLER_GEN_VERSION}
	mv $(PWD)/bin/controller-gen $@

bin/golangci-lint: | bin/golangci-lint_${GOLANGCI_VERSION} bin
	ln -sf golangci-lint_${GOLANGCI_VERSION} bin/golangci-lint

bin/golangci-lint_${GOLANGCI_VERSION}: | bin
	find $(PWD)/bin -name 'golangci-lint*' -exec rm {} +
	GOBIN=$(PWD)/bin go install github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_VERSION}
	mv bin/golangci-lint $@

bin/kind: | bin/kind_${KIND_VERSION} bin
	ln -sf kind_${KIND_VERSION} $@

bin/kind_${KIND_VERSION}: | bin
	find $(PWD)/bin -name 'kind*' -exec rm {} +
	GOBIN=$(PWD)/bin go install sigs.k8s.io/kind@${KIND_VERSION}
	mv bin/kind $@

bin/licensei: | bin/licensei_${LICENSEI_VERSION}
	ln -sf licensei_${LICENSEI_VERSION} $@

bin/licensei_${LICENSEI_VERSION}: | bin
	find $(PWD)/bin -name 'licensei*' -exec rm {} +
	curl -sfL https://raw.githubusercontent.com/goph/licensei/master/install.sh | bash -s ${LICENSEI_VERSION}
	mv bin/licensei $@

bin:
	mkdir -p bin

${ENVTEST_ASSETS_DIR}:
	mkdir -p ${ENVTEST_ASSETS_DIR}

${ENVTEST_ASSETS_DIR}/setup-envtest.sh: | ${ENVTEST_ASSETS_DIR}
	curl -sSLo ${ENVTEST_ASSETS_DIR}/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/$(ENVTEST_CTRL_VERSION)/hack/setup-envtest.sh