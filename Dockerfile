FROM --platform=$BUILDPLATFORM golang:1.26.0-alpine3.22@sha256:169d3991a4f795124a88b33c73549955a3d856e26e8504b5530c30bd245f9f1b AS builder

RUN apk add --update --no-cache ca-certificates make git curl

ARG TARGETOS
ARG TARGETARCH
ARG TARGETPLATFORM
ARG GO_BUILD_FLAGS

WORKDIR /usr/local/src/logging-operator

ARG GOPROXY

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
COPY pkg/sdk/go.mod pkg/sdk/go.mod
COPY pkg/sdk/go.sum pkg/sdk/go.sum
COPY pkg/sdk/logging/model/syslogng/config/go.mod pkg/sdk/logging/model/syslogng/config/go.mod
COPY pkg/sdk/logging/model/syslogng/config/go.sum pkg/sdk/logging/model/syslogng/config/go.sum

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY controllers/ controllers/
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build $GO_BUILD_FLAGS -o /usr/local/bin/manager

FROM builder AS debug

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go install github.com/go-delve/delve/cmd/dlv@latest

CMD ["/go/bin/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/usr/local/bin/manager"]


FROM gcr.io/distroless/static:debug AS e2e-test

COPY --from=builder /usr/local/bin/manager /manager

ENTRYPOINT ["/manager"]


FROM gcr.io/distroless/static:latest@sha256:d90359c7a3ad67b3c11ca44fd5f3f5208cbef546f2e692b0dc3410a869de46bf

COPY --from=builder /usr/local/bin/manager /manager

ENTRYPOINT ["/manager"]
