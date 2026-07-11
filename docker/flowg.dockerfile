# syntax=docker/dockerfile:1.7-labs

ARG DEBIAN_RELEASE="trixie"

ARG UPX_VERSION="4.2.4"
ARG UPX_OS="linux"

##############################
## SOURCES FILES
##############################

## Golang sources
FROM scratch AS sources-go

ADD third-party /src/third-party

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
FROM rust:1.96-slim-${DEBIAN_RELEASE} AS builder-rust-vrl

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

FROM node:26-${DEBIAN_RELEASE}-slim AS builder-js

RUN mkdir -p /workspace/web/app
WORKDIR /workspace/web/app

COPY --from=sources-js /src/web/app/package.json /workspace/web/app
COPY --from=sources-js /src/web/app/package-lock.json /workspace/web/app
RUN npm i

COPY --from=sources-js /src /workspace
RUN npm run lint
RUN NODE_ENV="production" npm run build

##############################
## BUILD GO CODE
##############################

FROM golang:1.26-${DEBIAN_RELEASE} AS builder-go
ARG TARGETARCH
ARG UPX_VERSION
ARG UPX_OS

RUN set -ex && \
    apt update && \
    apt install -y --no-install-recommends \
        build-essential \
        curl \
        ca-certificates \
      && \
    rm -rf /var/lib/apt/lists/*

ADD https://github.com/upx/upx/releases/download/v${UPX_VERSION}/upx-${UPX_VERSION}-${TARGETARCH}_${UPX_OS}.tar.xz /tmp/upx.tar.xz
RUN set -ex && \
    tar -xvf /tmp/upx.tar.xz && \
    mv upx-${UPX_VERSION}-${TARGETARCH}_${UPX_OS}/upx /usr/local/bin/ && \
    rm -rf /tmp/upx.tar.xz upx-${UPX_VERSION}-${TARGETARCH}_${UPX_OS}

RUN mkdir -p /workspace
WORKDIR /workspace

COPY --from=sources-go /src/go.mod .
COPY --from=sources-go /src/go.sum .
RUN go mod download

COPY --from=sources-go /src /workspace
COPY --from=builder-rust-vrl /workspace/internal/utils/langs/vrl/rust-crate/target/release/libflowg_vrl.a /workspace/internal/utils/langs/vrl/rust-crate/target/release/libflowg_vrl.a
COPY --from=builder-js /workspace/web/app/dist /workspace/web/public

RUN go generate ./...
RUN CGO_ENABLED=1 \
    CGO_CFLAGS="-I/workspace/third-party/foundationdb/7.3.77/include" \
    CGO_LDFLAGS="-L/workspace/third-party/foundationdb/7.3.77/lib/linux/$(go env GOARCH)" \
    go test -timeout 500ms -v ./...
RUN CGO_ENABLED=1 \
    CGO_CFLAGS="-I/workspace/third-party/foundationdb/7.3.77/include" \
    CGO_LDFLAGS="-L/workspace/third-party/foundationdb/7.3.77/lib/linux/$(go env GOARCH)" \
    go build -ldflags="-s -w" -o bin/ ./cmd/...
RUN upx bin/*

##############################
## FINAL ARTIFACT
##############################

FROM debian:${DEBIAN_RELEASE}-slim AS runner
ARG TARGETARCH

COPY --from=builder-go /workspace/bin/flowg-server /usr/local/bin/flowg-server
COPY --from=builder-go /workspace/bin/flowg-health /usr/local/bin/flowg-health

COPY --from=builder-go /workspace/third-party/foundationdb/7.3.77/lib/linux/${TARGETARCH}/libfdb_c.so /usr/local/lib/libfdb_c.so
RUN set -ex && \
    chmod 0755 /usr/local/lib/libfdb_c.so && \
    ldconfig

RUN set -ex && \
    mkdir -p /data && \
    chown 10001:10001 /data
VOLUME /data
USER 10001:10001

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

ENV FLOWG_BADGER_AUTH_DIR="/data/auth"
ENV FLOWG_BADGER_CONFIG_DIR="/data/config"
ENV FLOWG_BADGER_LOG_DIR="/data/logs"

ENTRYPOINT ["/usr/local/bin/flowg-server"]
CMD []

HEALTHCHECK --interval=5s --timeout=1s --retries=3 CMD ["/usr/local/bin/flowg-health", "--pid", "1"]
