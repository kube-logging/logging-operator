
GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./client/*")
PKGS=$(shell go list ./... | grep -v /vendor)

build:
	@go build $(PKGS)

vet:
	@go vet -composites=false ./...

check-fmt:
	PKGS="${GOFILES_NOVENDOR}" GOFMT="gofmt" ./scripts/fmt-check.sh

fmt:
	@gofmt -w ${GOFILES_NOVENDOR}

lint: install-golint
	golint -min_confidence 0.9 -set_exit_status $(PKGS)

install-golint:
	GOLINT_CMD=$(shell command -v golint 2> /dev/null)
ifndef GOLINT_CMD
	go get golang.org/x/lint/golint
endif

check-misspell: install-misspell
	PKGS="${GOFILES_NOVENDOR}" MISSPELL="misspell" ./scripts/misspell-check.sh

misspell: install-misspell
	misspell -w ${GOFILES_NOVENDOR}

install-misspell:
	MISSPELL_CMD=$(shell command -v misspell 2> /dev/null)
ifndef MISSPELL_CMD
	go get -u github.com/client9/misspell/cmd/misspell
endif

ineffassign: install-ineffassign
	ineffassign ${GOFILES_NOVENDOR}

install-ineffassign:
	INEFFASSIGN_CMD=$(shell command -v ineffassign 2> /dev/null)
ifndef INEFFASSIGN_CMD
	go get -u github.com/gordonklaus/ineffassign
endif

