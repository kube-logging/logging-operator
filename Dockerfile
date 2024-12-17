FROM --platform=$BUILDPLATFORM golang:1.23.4-alpine3.20@sha256:9a31ef0803e6afdf564edc8ba4b4e17caed22a0b1ecd2c55e3c8fdd8d8f68f98 AS builder

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


FROM gcr.io/distroless/static:latest@sha256:5c7e2b465ac6a2a4e5f4f7f722ce43b147dabe87cb21ac6c4007ae5178a1fa58

COPY --from=builder /usr/local/bin/manager /manager

ENTRYPOINT ["/manager"]
