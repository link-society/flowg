package rafthttp

import (
	"context"

	"crypto/tls"
	"net"
)

type Dial func(ctx context.Context, addr string) (net.Conn, error)

func NewDialTCP() Dial {
	return func(ctx context.Context, addr string) (net.Conn, error) {
		dialer := &net.Dialer{}
		return dialer.DialContext(ctx, "tcp", addr)
	}
}

func NewDialTLS(tlsConfig *tls.Config) Dial {
	return func(ctx context.Context, addr string) (net.Conn, error) {
		dialer := &tls.Dialer{Config: tlsConfig}
		return dialer.DialContext(ctx, "tcp", addr)
	}
}
