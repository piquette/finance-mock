FROM golang:1.9-alpine

ADD . /go/src/github.com/piquette/finance-mock

RUN go install github.com/piquette/finance-mock

EXPOSE 12111

ENTRYPOINT /go/bin/finance-mock