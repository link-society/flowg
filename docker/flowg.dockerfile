# syntax=docker/dockerfile:1.7-labs

##############################
## SOURCES FILES
##############################

## Golang sources
FROM scratch AS sources-go

ADD api /src/api
ADD cmd /src/cmd
ADD internal/app /src/internal/app
ADD internal/data /src/internal/data
ADD internal/integrations /src/internal/integrations
ADD internal/webutils /src/internal/webutils
ADD --exclude=internal/ffi/filterdsl/rust-crate internal/ffi/filterdsl /src/internal/ffi/filterdsl
ADD --exclude=internal/ffi/vrl/rust-crate internal/ffi/vrl /src/internal/ffi/vrl
ADD --exclude=web/components web /src/web
ADD go.mod go.sum /src/

## FlowEditor sources
FROM scratch AS sources-js-floweditor

ADD web/components/floweditor /src/web/components/floweditor

## CodeEditor sources
FROM scratch AS sources-js-code-editor

ADD web/components/code-editor /src/web/components/code-editor

## FilterDSL sources
FROM scratch AS sources-rust-filterdsl

ADD internal/ffi/filterdsl/rust-crate /src/internal/ffi/filterdsl/rust-crate

## VRL sources
FROM scratch AS sources-rust-vrl

ADD internal/ffi/vrl/rust-crate /src/internal/ffi/vrl/rust-crate

##############################
## BUILD RUST DEPENDENCIES
##############################

## FilterDSL
FROM rust:1.81-alpine3.20 AS builder-rust-filterdsl
RUN apk add --no-cache musl-dev

COPY --from=sources-rust-filterdsl /src /workspace
WORKDIR /workspace/internal/ffi/filterdsl/rust-crate

RUN cargo build --release
RUN cargo test

## VRL
FROM rust:1.81-alpine3.20 AS builder-rust-vrl
RUN apk add --no-cache musl-dev

COPY --from=sources-rust-vrl /src /workspace
WORKDIR /workspace/internal/ffi/vrl/rust-crate

RUN cargo build --release
RUN cargo test

##############################
## BUILD JS DEPENDENCIES
##############################

## FlowEditor
FROM node:22-alpine3.20 AS builder-js-floweditor
COPY --from=sources-js-floweditor /src /workspace
WORKDIR /workspace/web/components/floweditor

RUN npm i
RUN npm run build

## CodeEditor
FROM node:22-alpine3.20 AS builder-js-code-editor
COPY --from=sources-js-code-editor /src /workspace
WORKDIR /workspace/web/components/code-editor

RUN npm i
RUN npm run build

##############################
## BUILD GO CODE
##############################

FROM golang:1.23-alpine3.20 AS builder-go

RUN apk add --no-cache gcc musl-dev
RUN go install github.com/a-h/templ/cmd/templ@v0.2.778

COPY --from=sources-go /src /workspace
COPY --from=builder-rust-filterdsl /workspace/internal/ffi/filterdsl/rust-crate/target/release/libflowg_filterdsl.a /workspace/internal/ffi/filterdsl/rust-crate/target/release/libflowg_filterdsl.a
COPY --from=builder-rust-vrl /workspace/internal/ffi/vrl/rust-crate/target/release/libflowg_vrl.a /workspace/internal/ffi/vrl/rust-crate/target/release/libflowg_vrl.a
COPY --from=builder-js-floweditor /workspace/web/static /workspace/web/static
COPY --from=builder-js-code-editor /workspace/web/static /workspace/web/static
WORKDIR /workspace

RUN go generate ./...
RUN go build -o bin/ ./...
RUN go test -v ./...

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

ENV FLOWG_HTTP_BIND_ADDRESS=":5080"
ENV FLOWG_SYSLOG_BIND_ADDRESS=":5514"

ENV FLOWG_AUTH_DIR="/data/auth"
ENV FLOWG_CONFIG_DIR="/data/config"
ENV FLOWG_LOG_DIR="/data/logs"

ENTRYPOINT ["/docker-entrypoint.sh"]
