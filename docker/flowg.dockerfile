# syntax=docker/dockerfile:1.7-labs

ARG UPX_VERSION="4.2.4"
ARG UPX_ARCH="amd64"
ARG UPX_OS="linux"

##############################
## SOURCES FILES
##############################

## Golang sources
FROM scratch AS sources-go

ADD VERSION.txt /src/VERSION.txt
ADD api /src/api
ADD cmd /src/cmd

ADD \
  --exclude=internal/utils/langs/vrl/rust-crate \
  internal/ /src/internal/

ADD \
  --exclude=web/app \
  web /src/web

ADD go.mod go.sum /src/

## JS sources
FROM scratch AS sources-js

ADD VERSION.txt /src/VERSION.txt
ADD web/app /src/web/app

## VRL sources
FROM scratch AS sources-rust-vrl

ADD internal/utils/langs/vrl/rust-crate /src/internal/utils/langs/vrl/rust-crate

##############################
## BUILD RUST DEPENDENCIES
##############################

## VRL
FROM rust:1.93-alpine3.23 AS builder-rust-vrl

RUN apk add --no-cache musl-dev

RUN mkdir -p /workspace/internal/utils/langs/vrl/rust-crate
WORKDIR /workspace/internal/utils/langs/vrl/rust-crate

COPY --from=sources-rust-vrl /src/internal/utils/langs/vrl/rust-crate/Cargo.toml .
COPY --from=sources-rust-vrl /src/internal/utils/langs/vrl/rust-crate/Cargo.lock .
RUN mkdir src \
    && echo "// dummy file" > src/lib.rs \
    && cargo build

COPY --from=sources-rust-vrl /src /workspace
RUN cargo build --release
RUN cargo test

##############################
## BUILD JS DEPENDENCIES
##############################

FROM node:25-alpine3.23 AS builder-js

RUN mkdir -p /workspace/web/app
WORKDIR /workspace/web/app

COPY --from=sources-js /src/web/app/package.json /workspace/web/app
RUN npm i

COPY --from=sources-js /src /workspace
RUN npm run lint
RUN NODE_ENV="production" npm run build

##############################
## BUILD GO CODE
##############################

FROM golang:1.25-alpine3.23 AS builder-go
ARG UPX_VERSION
ARG UPX_ARCH
ARG UPX_OS

RUN apk add --no-cache gcc musl-dev curl

ADD https://github.com/upx/upx/releases/download/v${UPX_VERSION}/upx-${UPX_VERSION}-${UPX_ARCH}_${UPX_OS}.tar.xz /tmp/upx.tar.xz
RUN set -ex && \
    tar -xvf /tmp/upx.tar.xz && \
    mv upx-${UPX_VERSION}-${UPX_ARCH}_${UPX_OS}/upx /usr/local/bin/ && \
    rm -rf /tmp/upx.tar.xz upx-${UPX_VERSION}-${UPX_ARCH}_${UPX_OS}

RUN mkdir -p /workspace
WORKDIR /workspace

COPY --from=sources-go /src/go.mod .
COPY --from=sources-go /src/go.sum .
RUN go mod download

COPY --from=sources-go /src /workspace
COPY --from=builder-rust-vrl /workspace/internal/utils/langs/vrl/rust-crate/target/release/libflowg_vrl.a /workspace/internal/utils/langs/vrl/rust-crate/target/release/libflowg_vrl.a
COPY --from=builder-js /workspace/web/app/dist /workspace/web/public

RUN go generate ./...
RUN go test -timeout 500ms -v ./...
RUN go build -ldflags="-s -w" -o bin/ ./cmd/...
RUN upx bin/*

##############################
## FINAL ARTIFACT
##############################

FROM alpine:3.23 AS runner

RUN apk add --no-cache libgcc su-exec

COPY --from=builder-go /workspace/bin/flowg-server /usr/local/bin/flowg-server
COPY --from=builder-go /workspace/bin/flowg-health /usr/local/bin/flowg-health

ADD docker/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod 0700 /docker-entrypoint.sh

RUN addgroup -S flowg && adduser -S -G flowg -h /app flowg
WORKDIR /app

ENV FLOWG_SECRET_KEY=""
ENV FLOWG_VERBOSE=false

ENV FLOWG_HTTP_BIND_ADDRESS=":5080"
ENV FLOWG_HTTP_TLS_ENABLED=false
ENV FLOWG_HTTP_TLS_CERT=""
ENV FLOWG_HTTP_TLS_KEY=""

ENV FLOWG_MGMT_BIND_ADDRESS=":9113"
ENV FLOWG_MGMT_TLS_ENABLED=false
ENV FLOWG_MGMT_TLS_CERT=""
ENV FLOWG_MGMT_TLS_KEY=""

ENV FLOWG_SYSLOG_PROTOCOL="udp"
ENV FLOWG_SYSLOG_BIND_ADDRESS=":5514"
ENV FLOWG_SYSLOG_ALLOW_ORIGINS=""
ENV FLOWG_SYSLOG_TLS_ENABLED=false
ENV FLOWG_SYSLOG_TLS_CERT=""
ENV FLOWG_SYSLOG_TLS_KEY=""
ENV FLOWG_SYSLOG_TLS_AUTH=false

ENV FLOWG_AUTH_DIR="/data/auth"
ENV FLOWG_CONFIG_DIR="/data/config"
ENV FLOWG_LOG_DIR="/data/logs"
ENV FLOWG_CLUSTER_STATE_DIR="/data/state"

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD []

HEALTHCHECK --interval=5s --timeout=1s --retries=3 CMD ["/usr/local/bin/flowg-health", "--pid", "1"]
