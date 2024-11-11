# syntax=docker/dockerfile:1.7-labs

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
FROM rust:1.82-alpine3.20 AS builder-rust-filterdsl
RUN apk add --no-cache musl-dev

COPY --from=sources-rust-filterdsl /src /workspace
WORKDIR /workspace/internal/utils/ffi/filterdsl/rust-crate

RUN cargo build --release
RUN cargo test

## VRL
FROM rust:1.82-alpine3.20 AS builder-rust-vrl
RUN apk add --no-cache musl-dev

COPY --from=sources-rust-vrl /src /workspace
WORKDIR /workspace/internal/utils/ffi/vrl/rust-crate

RUN cargo build --release
RUN cargo test

##############################
## BUILD JS DEPENDENCIES
##############################

FROM node:23-alpine3.20 AS builder-js
COPY --from=sources-js /src /workspace
WORKDIR /workspace/web/app

RUN npm i
RUN NODE_ENV="production" npm run build

##############################
## BUILD GO CODE
##############################

FROM golang:1.23-alpine3.20 AS builder-go

RUN apk add --no-cache gcc musl-dev
RUN go install github.com/a-h/templ/cmd/templ@v0.2.778

COPY --from=sources-go /src /workspace
COPY --from=builder-rust-filterdsl /workspace/internal/utils/ffi/filterdsl/rust-crate/target/release/libflowg_filterdsl.a /workspace/internal/utils/ffi/filterdsl/rust-crate/target/release/libflowg_filterdsl.a
COPY --from=builder-rust-vrl /workspace/internal/utils/ffi/vrl/rust-crate/target/release/libflowg_vrl.a /workspace/internal/utils/ffi/vrl/rust-crate/target/release/libflowg_vrl.a
COPY --from=builder-js /workspace/web/app/dist /workspace/web/public
WORKDIR /workspace

RUN go generate ./...
RUN go build -o bin/ ./...
RUN go test -timeout 500ms -v ./...

##############################
## FINAL ARTIFACT
##############################

FROM alpine:3.20 AS runner

RUN apk add --no-cache libgcc su-exec

COPY --from=builder-go /workspace/bin/ /usr/local/bin/

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

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["serve"]
