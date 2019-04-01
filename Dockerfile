FROM golang:1.11-alpine as golang

RUN apk add --update --no-cache ca-certificates curl git make
RUN go get -u github.com/golang/dep/cmd/dep

ADD Gopkg.toml /go/src/github.com/banzaicloud/logging-operator/Gopkg.toml
ADD Gopkg.lock /go/src/github.com/banzaicloud/logging-operator/Gopkg.lock

WORKDIR /go/src/github.com/banzaicloud/logging-operator
RUN dep ensure -v -vendor-only
ADD . /go/src/github.com/banzaicloud/logging-operator
RUN go install ./cmd/manager


FROM alpine:3.8

RUN apk add --no-cache ca-certificates

COPY --from=golang /go/bin/manager /usr/local/bin/logging-operator

RUN adduser -D logging-operator
USER logging-operator

ENTRYPOINT ["/usr/local/bin/logging-operator"]