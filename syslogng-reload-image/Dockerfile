FROM ghcr.io/banzaicloud/custom-runner AS custom-runner

FROM alpine:latest
RUN apk add socat
WORKDIR /
COPY --from=custom-runner /runner .
ENTRYPOINT ["/runner"]
