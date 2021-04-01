# this makefile was generated by
include Makefile.app

OS = $(shell uname | tr A-Z a-z)

# Image URL to use all building/pushing image targets
IMG ?= controller:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true,preserveUnknownFields=false,maxDescLen=0"

KUBEBUILDER_VERSION = 2.3.1
LICENSEI_VERSION = 0.3.1
VERSION := $(shell git describe --abbrev=0 --tags)
DOCKER_IMAGE = banzaicloud/logging-operator
DOCKER_TAG ?= ${VERSION}
GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./client/*")
GOFILES_NOPLUGINS =  $(shell find . -type f -name '*.go' -not -path "./pkg/sdk/model/filter/*"  -not -path "./pkg/sdk/model/output/*"  -not -path "./pkg/sdk/model/input/*")
PKGS=$(shell go list ./... | grep -v /vendor)

CONTROLLER_GEN_VERSION = v0.2.4
CONTROLLER_GEN = $(PWD)/bin/controller-gen

GOLANGCI_VERSION = 1.33.0
GOBIN_VERSION = 0.0.13

export KUBEBUILDER_ASSETS := $(PWD)/bin
export PATH := $(PWD)/bin:$(PATH)

all: manager

# Generate docs
.PHONY: docs
docs:
	go run cmd/docs.go

bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint
bin/golangci-lint-${GOLANGCI_VERSION}:
	@mkdir -p bin
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | BINARY=golangci-lint bash -s -- v${GOLANGCI_VERSION}
	@mv bin/golangci-lint $@

.PHONY: lint
lint: bin/golangci-lint ## Run linter
	bin/golangci-lint run 
	cd pkg/sdk && ../../bin/golangci-lint run  -c ../../.golangci.yml

bin/licensei: bin/licensei-${LICENSEI_VERSION}
	@ln -sf licensei-${LICENSEI_VERSION} bin/licensei
bin/licensei-${LICENSEI_VERSION}:
	@mkdir -p bin
	curl -sfL https://raw.githubusercontent.com/goph/licensei/master/install.sh | bash -s v${LICENSEI_VERSION}
	@mv bin/licensei $@

.PHONY: license-check
license-check: bin/licensei ## Run license check
	bin/licensei check
	./scripts/check-header.sh

.PHONY: license-cache
license-cache: bin/licensei ## Generate license cache
	bin/licensei cache

.PHONY: lint-fix
lint-fix: bin/golangci-lint ## Run linter
	bin/golangci-lint run --fix
	cd pkg/sdk && golangci-lint run -c ../../.golangci.yml --fix

bin/gobin: bin/gobin-${GOBIN_VERSION}
	@ln -sf gobin-${GOBIN_VERSION} bin/gobin
bin/gobin-${GOBIN_VERSION}:
	@mkdir -p bin
	curl -L https://github.com/myitcv/gobin/releases/download/v${GOBIN_VERSION}/${OS}-amd64 > ./bin/gobin-${GOBIN_VERSION} && chmod +x ./bin/gobin-${GOBIN_VERSION}

.PHONY: bin/kubebuilder_${KUBEBUILDER_VERSION}
bin/kubebuilder_${KUBEBUILDER_VERSION}:
	@ if ! test -L bin/kubebuilder_${KUBEBUILDER_VERSION}; then \
		mkdir -p bin; \
		curl -L https://github.com/kubernetes-sigs/kubebuilder/releases/download/v${KUBEBUILDER_VERSION}/kubebuilder_${KUBEBUILDER_VERSION}_${OS}_amd64.tar.gz | tar xvz -C bin; \
		ln -sf kubebuilder_${KUBEBUILDER_VERSION}_${OS}_amd64/bin bin/kubebuilder_${KUBEBUILDER_VERSION}; \
	fi

bin/kubebuilder: bin/kubebuilder_${KUBEBUILDER_VERSION}
	@ln -sf kubebuilder_${KUBEBUILDER_VERSION}/kubebuilder bin/kubebuilder
	@ln -sf kubebuilder_${KUBEBUILDER_VERSION}/kube-apiserver bin/kube-apiserver
	@ln -sf kubebuilder_${KUBEBUILDER_VERSION}/etcd bin/etcd
	@ln -sf kubebuilder_${KUBEBUILDER_VERSION}/kubectl bin/kubectl

# Run tests
test: generate fmt vet manifests bin/kubebuilder
	@which kubebuilder
	@which etcd
	kubebuilder version
	cd pkg/sdk && go test ./...
	go test ./controllers/... ./pkg/... -coverprofile cover.out -v

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./main.go --verbose --pprof

# remote debug
debug: manager
	dlv --listen=:40000 --log --headless=true --api-version=2 exec bin/manager -- $(ARGS)

# Install CRDs into a cluster
install: manifests
	kubectl apply -f config/crd/bases

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	kubectl apply -f config/crd/bases
	kustomize build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: bin/controller-gen
	cd pkg/sdk && $(CONTROLLER_GEN) $(CRD_OPTIONS) webhook paths="./..." output:crd:artifacts:config=../../config/crd/bases output:webhook:artifacts:config=../../config/webhook
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role paths="./controllers/..." output:rbac:artifacts:config=./config/rbac
	cp config/crd/bases/* charts/logging-operator/crds/
	echo "{{- if .Values.rbac.enabled }}" > ./charts/logging-operator/templates/clusterrole.yaml && cat config/rbac/role.yaml |sed -e 's@manager-role@{{ template "logging-operator.fullname" . }}@' | cat >> ./charts/logging-operator/templates/clusterrole.yaml && echo "{{- end }}" >> ./charts/logging-operator/templates/clusterrole.yaml

# Run go fmt against code
fmt:
	go fmt ./...
	cd pkg/sdk && go fmt ./...

# Run go vet against code
vet:
	go vet ./...
	cd pkg/sdk && go vet ./...

# Generate code
generate: bin/controller-gen
	cd pkg/sdk && $(CONTROLLER_GEN) object:headerFile=./../../hack/boilerplate.go.txt paths=./api/...
	cd pkg/sdk && $(CONTROLLER_GEN) object:headerFile=./../../hack/boilerplate.go.txt paths=./model/...
	cd pkg/sdk && $(CONTROLLER_GEN) object:headerFile=./../../hack/boilerplate.go.txt paths=./resourcebuilder/...
	cd pkg/sdk && go generate ./static

# Build the docker image
docker-build:
	docker build . -t ${IMG}
	@echo "updating kustomize image patch file for manager resource"
	sed -i'' -e 's@image: .*@image: '"${IMG}"'@' ./config/default/manager_image_patch.yaml

# Push the docker image
docker-push:
	docker push ${IMG}

.PHONY: bin/controller-gen
bin/controller-gen:
	@ if ! test -x bin/controller-gen; then \
		set -ex ;\
		CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
		cd $$CONTROLLER_GEN_TMP_DIR ;\
		go mod init tmp ;\
		GOBIN=$(PWD)/bin go get sigs.k8s.io/controller-tools/cmd/controller-gen@${CONTROLLER_GEN_VERSION} ;\
		rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	fi

check-diff: check
	go mod tidy
	$(MAKE) generate manifests
	git diff --exit-code ':(exclude)./ADOPTERS.md' ':(exclude)./docs/*'
