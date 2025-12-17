FROM ghcr.io/kube-logging/custom-runner:v0.14.0 AS custom-runner

FROM alpine:3.22.2@sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412

RUN apk add socat

COPY --from=custom-runner /runner /

WORKDIR /

ENTRYPOINT ["/runner"]
