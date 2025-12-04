FROM alpine:3.23.0@sha256:51183f2cfa6320055da30872f211093f9ff1d3cf06f39a0bdb212314c5dc7375

RUN apk add curl

COPY drain-watch.sh /bin/drain-watch.sh

ENTRYPOINT ["/bin/drain-watch.sh"]
