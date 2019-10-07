FROM golang:alpine

WORKDIR /go/src/github.com/williamlsh/producer-random

COPY . .

ENV GO111MODULE=on

RUN go install

FROM alpine:latest

COPY --from=0 /go/bin/producer-random .

ENTRYPOINT [ "./producer-random" ]