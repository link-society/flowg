FROM golang:1.23-alpine3.20 AS builder

RUN apk add --no-cache npm cargo gcc musl-dev
RUN go install github.com/go-task/task/v3/cmd/task@v3.38.0

ADD . /workspace
WORKDIR /workspace

RUN task build

FROM alpine:3.20 AS runner

RUN apk add --no-cache libgcc

COPY --from=builder /workspace/bin/ /usr/local/bin/

ADD docker/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]
