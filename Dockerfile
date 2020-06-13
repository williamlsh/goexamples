FROM --platform=${BUILDPLATFORM} golang:1.14-alpine AS build

WORKDIR /src

ENV CGO_ENABLED=0

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/example .

FROM scratch AS bin

COPY --from=build /out/example /