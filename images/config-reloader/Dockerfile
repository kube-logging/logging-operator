FROM golang:1.26.2-alpine3.22@sha256:c259ff7ffa06f1fd161a6abfa026573cf00f64cfd959c6d2a9d43e3ff63e8729 AS builder

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
COPY pkg pkg

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download
RUN go mod tidy

# Copy the go source
COPY /cmd/configreloader/main.go main.go

# Build
RUN CGO_ENABLED=0 go build -a -o config-reloader main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:latest

WORKDIR /

COPY --from=builder /workspace/config-reloader .

ENTRYPOINT ["/config-reloader"]
