FROM golang:alpine

WORKDIR /go/src/github.com/williamlsh/consumer-logger

COPY . .

ENV GO111MODULE=on

RUN go install

FROM alpine:latest

COPY --from=0 /go/bin/consumer-logger .

ENTRYPOINT [ "./consumer-logger" ]