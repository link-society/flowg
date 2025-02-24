package rafthttp

import (
	"fmt"

	"context"
	"time"

	"bufio"
	"io"

	"net"
	"net/http"
	"net/url"

	"github.com/vladopajic/go-actor/actor"
	
    "github.com/hashicorp/raft"
)

type Layer struct {
	mbox      actor.MailboxReceiver[net.Conn]
	url       *url.URL
	localAddr net.Addr
	dial      Dial
}

func NewLayer(
	mbox actor.MailboxReceiver[net.Conn],
	path string,
	localAddr net.Addr,
	dial Dial,
) *Layer {
	return &Layer{
		mbox:      mbox,
		url:       &url.URL{Path: path},
		localAddr: localAddr,
		dial:      dial,
	}
}

func (l *Layer) Accept() (net.Conn, error) {
	conn, ok := <-l.mbox.ReceiveC()
	if !ok {
		return nil, io.EOF
	}

	return conn, nil
}

func (l *Layer) Close() error {
	return nil
}

func (l *Layer) Addr() net.Addr {
	return l.localAddr
}

func (l *Layer) Dial(addr raft.ServerAddress, timeout time.Duration) (net.Conn, error) {
	req := &http.Request{
		Method:     "GET",
		URL:        l.url,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       l.Addr().String(),
	}
	req.Header.Set("Upgrade", "raft")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	conn, err := l.dial(ctx, string(addr))
	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}

	if err := req.Write(conn); err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}

	resp, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP response: %w", err)
	}

	if resp.StatusCode != http.StatusSwitchingProtocols {
        body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("expected status code 101, got: %d - %v", resp.StatusCode, string(body))
	}

	if resp.Header.Get("Upgrade") != "raft" {
		return nil, fmt.Errorf("expected `Upgrade: raft` header")
	}

	return conn, nil
}
