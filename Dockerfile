FROM --platform=$BUILDPLATFORM golang:1.27.1-alpine3.23@sha256:0908ac9b9319e09d7c238aabe914e0395c51d63c4e3d0ae8c554fda9158a5769 AS builder

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


FROM gcr.io/distroless/static:latest@sha256:58133991db06659feaabe0f4e97a35cebf15ef4ea08f8a4c6d2ee5f75e4aa6a0

COPY --from=builder /usr/local/bin/manager /manager

ENTRYPOINT ["/manager"]
