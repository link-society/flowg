package cluster

import (
	"fmt"
	"log/slog"

	"crypto/tls"
	"net"
	"net/url"

	"github.com/hashicorp/go-sockaddr"
)

type Listener struct {
	socket    net.Listener
	tlsConfig *tls.Config
}

func NewListener(
	logger *slog.Logger,
	bindAddress string,
	tlsConfig *tls.Config,
) (*Listener, error) {
	logger.Info("Listen on Management interface")

	socket, err := net.Listen("tcp", bindAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on Management interface: %w", err)
	}

	return &Listener{
		tlsConfig: tlsConfig,
		socket:    socket,
	}, nil
}

func (l *Listener) ResolveLocalEndpoint() (*url.URL, error) {
	host, port, err := net.SplitHostPort(l.socket.Addr().String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse listener address: %w", err)
	}

	if host == "0.0.0.0" || host == "::" {
		ip, err := sockaddr.GetPrivateIP()
		if err != nil {
			return nil, fmt.Errorf("failed to get private IP: %w", err)
		}
		if ip == "" {
			return nil, fmt.Errorf("no private IP found")
		}

		host = ip
	}

	localEndpoint := &url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(host, port),
	}

	if l.tlsConfig != nil {
		localEndpoint.Scheme = "https"
	}

	return localEndpoint, nil
}

func (l *Listener) Socket() net.Listener {
	return l.socket
}
