FROM golang:alpine AS builder

WORKDIR /src


COPY go.mod .

RUN go mod download; \
    go mod verify

COPY . .

RUN go build -o demo .

FROM alpine AS bin

COPY --from=builder /src/demo /demo

ENTRYPOINT [ "/demo" ]
