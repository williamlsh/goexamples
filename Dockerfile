FROM golang:1.15-alpine

RUN apk update && apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    && update-ca-certificates

WORKDIR /src

COPY go.mod .

RUN go env -w GOPROXY="https://goproxy.io,direct"; \
    go mod download; \
    go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o hotrod .

FROM scratch

COPY --from=0 /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /src/hotrod /

EXPOSE 8080 8081 8082 8083

ENTRYPOINT [ "/hotrod" ]

CMD ["all"]