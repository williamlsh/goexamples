FROM ghcr.io/williamlsh/srt-docker:latest AS libsrt

FROM golang:alpine

COPY --from=libsrt /usr/local /usr/local

RUN apk update && apk add --no-cache \
    build-base

WORKDIR /app

COPY . .

RUN go mod tidy \
    && go mod verify

ENTRYPOINT [ "go" ]

CMD [ "run", "main.go" ]
