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
  --exclude=internal/utils/ffi/filterdsl/rust-crate \
  --exclude=internal/utils/ffi/vrl/rust-crate \
  internal/ /src/internal/

ADD \
  --exclude=web/app \
  web /src/web

ADD go.mod go.sum /src/

## JS sources
FROM scratch AS sources-js

ADD VERSION.txt /src/VERSION.txt
ADD web/app /src/web/app

## FilterDSL sources
FROM scratch AS sources-rust-filterdsl

ADD internal/utils/ffi/filterdsl/rust-crate /src/internal/utils/ffi/filterdsl/rust-crate

## VRL sources
FROM scratch AS sources-rust-vrl

ADD internal/utils/ffi/vrl/rust-crate /src/internal/utils/ffi/vrl/rust-crate

##############################
## BUILD RUST DEPENDENCIES
##############################

## FilterDSL
FROM rust:1.88-alpine3.22 AS builder-rust-filterdsl
RUN apk add --no-cache musl-dev

COPY --from=sources-rust-filterdsl /src /workspace
WORKDIR /workspace/internal/utils/ffi/filterdsl/rust-crate

RUN cargo build --release
RUN cargo test

## VRL
FROM rust:1.88-alpine3.22 AS builder-rust-vrl
RUN apk add --no-cache musl-dev

COPY --from=sources-rust-vrl /src /workspace
WORKDIR /workspace/internal/utils/ffi/vrl/rust-crate

RUN cargo build --release
RUN cargo test

##############################
## BUILD JS DEPENDENCIES
##############################

FROM node:24-alpine3.22 AS builder-js
COPY --from=sources-js /src /workspace
WORKDIR /workspace/web/app

RUN npm i
RUN npm run lint
RUN NODE_ENV="production" npm run build

##############################
## BUILD GO CODE
##############################

FROM golang:1.24-alpine3.22 AS builder-go
ARG UPX_VERSION
ARG UPX_ARCH
ARG UPX_OS

RUN apk add --no-cache gcc musl-dev curl

ADD https://github.com/upx/upx/releases/download/v${UPX_VERSION}/upx-${UPX_VERSION}-${UPX_ARCH}_${UPX_OS}.tar.xz /tmp/upx.tar.xz
RUN set -ex && \
    tar -xvf /tmp/upx.tar.xz && \
    mv upx-${UPX_VERSION}-${UPX_ARCH}_${UPX_OS}/upx /usr/local/bin/ && \
    rm -rf /tmp/upx.tar.xz upx-${UPX_VERSION}-${UPX_ARCH}_${UPX_OS}

COPY --from=sources-go /src /workspace
COPY --from=builder-rust-filterdsl /workspace/internal/utils/ffi/filterdsl/rust-crate/target/release/libflowg_filterdsl.a /workspace/internal/utils/ffi/filterdsl/rust-crate/target/release/libflowg_filterdsl.a
COPY --from=builder-rust-vrl /workspace/internal/utils/ffi/vrl/rust-crate/target/release/libflowg_vrl.a /workspace/internal/utils/ffi/vrl/rust-crate/target/release/libflowg_vrl.a
COPY --from=builder-js /workspace/web/app/dist /workspace/web/public
WORKDIR /workspace

RUN go generate ./...
RUN go test -timeout 500ms -v ./...
RUN go build -ldflags="-s -w" -o bin/ ./cmd/...
RUN upx bin/*

##############################
## FINAL ARTIFACT
##############################

FROM alpine:3.22 AS runner

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
