# -*- mode: dockerfile -*-
#
# A multi-stage Dockerfile that builds a Linux target then creates a small
# final image for deployment.

#
# STAGE 1
#
# Uses a Go image to build a release binary.
#

FROM golang:1.10.2-alpine AS builder
WORKDIR /go/src/github.com/piquette/finance-mock/
ADD ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o finance-mock .

#
# STAGE 2
#
# Use a tiny base image (alpine) and copy in the release target. This produces
# a very small output image for deployment.
#

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /go/src/github.com/piquette/finance-mock/finance-mock .
ENTRYPOINT /finance-mock
EXPOSE 12111
