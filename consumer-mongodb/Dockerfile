FROM golang:alpine

WORKDIR /go/src/github.com/williamlsh/consumer-mongodb

COPY . .

ENV GO111MODULE=on

RUN go install

FROM alpine:latest

COPY --from=0 /go/bin/consumer-mongodb .

ENTRYPOINT [ "./consumer-mongodb" ]