FROM --platform=$BUILDPLATFORM golang:1.26.3-alpine3.22@sha256:be93003ee861b3b91b6ebcb22678524947e0cd786c2df3f32af520006b1e54f5 AS builder

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


FROM gcr.io/distroless/static:latest@sha256:3592aa8171c77482f62bbc4164e6a2d141c6122554ace66e5cc910cadb961ff0

COPY --from=builder /usr/local/bin/manager /manager

ENTRYPOINT ["/manager"]
