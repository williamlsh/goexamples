FROM golang:alpine

WORKDIR /go/src/github.com/williamlsh/producer-api

COPY . .

ENV GO111MODULE=on

RUN go install

FROM alpine:latest

COPY --from=0 /go/bin/producer-api .

EXPOSE 8080

ENTRYPOINT [ "./producer-api" ]