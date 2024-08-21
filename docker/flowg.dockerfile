FROM golang:1.23-alpine AS builder

RUN apk add --no-cache npm cargo gcc musl-dev
RUN go install github.com/go-task/task/v3/cmd/task@latest

ADD . /workspace
WORKDIR /workspace

RUN task build

FROM alpine:latest AS runner
COPY --from=builder /workspace/bin/ /usr/local/bin/

ADD docker/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

ENTRYPOINT ["/docker-entrypoint.sh"]
