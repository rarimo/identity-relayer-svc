FROM golang:1.19

WORKDIR /go/src/relayer-svc

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/relayer-svc relayer-svc


###

FROM alpine:3.9

COPY --from=0 /usr/local/bin/relayer-svc /usr/local/bin/relayer-svc
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["relayer-svc"]
