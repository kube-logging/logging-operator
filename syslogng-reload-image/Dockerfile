FROM ghcr.io/banzaicloud/custom-runner:0.1.0 AS custom-runner

FROM alpine:latest
RUN apk add socat
WORKDIR /
COPY --from=custom-runner /runner .
ENTRYPOINT ["/runner"]
