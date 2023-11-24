FROM golang:1.19-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/rarimo/identity-relayer-svc
COPY vendor .
COPY . .

ENV GO111MODULE="on"
ENV CGO_ENABLED=1
ENV GOOS="linux"

RUN GOOS=linux go build  -o /usr/local/bin/relayer-svc /go/src/github.com/rarimo/identity-relayer-svc


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/relayer-svc /usr/local/bin/relayer-svc
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["relayer-svc"]
