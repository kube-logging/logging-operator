FROM alpine:3.6

RUN adduser -D logging-operator
USER logging-operator

ADD tmp/_output/bin/logging-operator /usr/local/bin/logging-operator
