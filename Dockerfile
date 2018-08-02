FROM golang:1.10-alpine as golang

ADD . /go/src/github.com/banzaicloud/logging-operator
WORKDIR /go/src/github.com/banzaicloud/logging-operator

RUN go install ./cmd/logging-operator


FROM alpine:3.7

RUN apk add --no-cache ca-certificates

COPY --from=golang /go/bin/logging-operator /usr/local/bin/logging-operator

RUN adduser -D logging-operator
USER logging-operator

ENTRYPOINT ["/usr/local/bin/logging-operator"]