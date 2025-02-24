package rafthttp_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"time"

	"crypto/tls"
	"net"

	"link-society.com/flowg/internal/cluster/rafthttp"
)

func TestDialTCP(t *testing.T) {
	testDial(t, rafthttp.NewDialTCP(), (*httptest.Server).Start)
}

func TestDialTLS(t *testing.T) {
	testDial(
		t,
		rafthttp.NewDialTLS(&tls.Config{InsecureSkipVerify: true}),
		(*httptest.Server).StartTLS,
	)
}

func TestDialTCP_Timeout(t *testing.T) {
	testDialTimeout(t, rafthttp.NewDialTCP())
}

func TestDialTLS_Timeout(t *testing.T) {
	testDialTimeout(t, rafthttp.NewDialTLS(&tls.Config{InsecureSkipVerify: true}))
}

func testDial(
	t *testing.T,
	dial rafthttp.Dial,
	startServer func(*httptest.Server),
) {
	server := httptest.NewUnstartedServer(nil)
	startServer(server)
	defer server.Close()

	addr := server.Listener.Addr().String()

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()
	conn, err := dial(ctx, addr)
	if err != nil {
		t.Errorf("dial returned error: %v", err)
	}

	if conn == nil {
		t.Errorf("dial returned nil connection")
	}
}

func testDialTimeout(t *testing.T, dial rafthttp.Dial) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
	defer cancel()
	_, err := dial(ctx, "localhost:0")

	if err == nil {
		t.Errorf("dial did not return error")
	}

	if neterr, ok := err.(net.Error); ok {
		if !neterr.Timeout() {
			t.Errorf("dial returned non-timeout error: %v", neterr)
		}
	} else {
		t.Errorf("dial returned non-network error: %v", err)
	}
}
