FROM ubuntu:latest

ENV DEBIAN_FRONTEND=noninteractive

RUN set -eux \
    && apt-get update \
    && apt-get install -y --no-install-recommends \
    git \
    build-essential \
    libssl-dev \
    tclsh \
    pkg-config \
    cmake \
    software-properties-common \
    && rm -rf /var/lib/apt/lists/*

RUN set -eux; \
    add-apt-repository ppa:longsleep/golang-backports -y

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    golang-go

WORKDIR /tmp

RUN git clone --depth=1 --branch master --single-branch https://github.com/Haivision/srt.git \
    && cd srt \
    && ./configure \
    && make \
    && make install

WORKDIR /app

ENV LD_LIBRARY_PATH /lib:/usr/lib:/usr/local/lib64

COPY main.go go.mod go.sum ./

RUN go mod tidy

ENTRYPOINT [ "go" ]

CMD [ "run", "main.go" ]
