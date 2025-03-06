package cluster

import (
	"io"
	"testing"

	"log/slog"

	"context"
	"time"

	"net"
	"net/http"
	"net/url"

	"github.com/jarcoal/httpmock"

	"github.com/vladopajic/go-actor/actor"

	"github.com/hashicorp/memberlist"
)

func TestFinalAdvertiseAddr_WithValidIPv4(t *testing.T) {
	transport := &httpTransport{
		delegate: &delegate{
			logger: nil,
			localEndpoint: &url.URL{
				Scheme: "http",
				Host:   "1.2.3.4:5678",
			},
			endpoints: make(map[string]*url.URL),
		},
		cookie: "test",

		connM:   nil,
		packetM: nil,
	}

	ip, port, err := transport.FinalAdvertiseAddr("", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ip.String() != "1.2.3.4" {
		t.Fatalf("unexpected IPv4: %s", ip.String())
	}

	if port != 5678 {
		t.Fatalf("unexpected port: %d", port)
	}
}

func TestFinalAdvertiseAddr_WithValidIPv6(t *testing.T) {
	transport := &httpTransport{
		delegate: &delegate{
			logger: nil,
			localEndpoint: &url.URL{
				Scheme: "http",
				Host:   "[::1]:5678",
			},
			endpoints: make(map[string]*url.URL),
		},
		cookie: "test",

		connM:   nil,
		packetM: nil,
	}

	ip, port, err := transport.FinalAdvertiseAddr("", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ip.String() != "::1" {
		t.Fatalf("unexpected IPv6: %s", ip.String())
	}

	if port != 5678 {
		t.Fatalf("unexpected port: %d", port)
	}
}

func TestFinalAdvertiseAddr_WithInvalidHost(t *testing.T) {
	transport := &httpTransport{
		delegate: &delegate{
			logger: nil,
			localEndpoint: &url.URL{
				Scheme: "http",
				Host:   "invalid",
			},
			endpoints: make(map[string]*url.URL),
		},
		cookie: "test",

		connM:   nil,
		packetM: nil,
	}

	_, _, err := transport.FinalAdvertiseAddr("", 0)
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}
}

func TestWriteToAddress_EmptyNodeName(t *testing.T) {
	transport := &httpTransport{
		delegate: &delegate{
			logger: nil,
			localEndpoint: &url.URL{
				Scheme: "http",
				Host:   "invalid",
			},
			endpoints: make(map[string]*url.URL),
		},
		cookie: "test",

		connM:   nil,
		packetM: nil,
	}

	_, err := transport.WriteToAddress(nil, memberlist.Address{})
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}
}

func TestWriteToAddress_EndpointNotFound(t *testing.T) {
	transport := &httpTransport{
		delegate: &delegate{
			logger: nil,
			localEndpoint: &url.URL{
				Scheme: "http",
				Host:   "invalid",
			},
			endpoints: make(map[string]*url.URL),
		},
		cookie: "test",

		connM:   nil,
		packetM: nil,
	}

	_, err := transport.WriteToAddress(nil, memberlist.Address{Name: "node"})
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}
}

func TestWriteToAddress_Valid(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST", "http://1.2.3.4:5678/cluster/gossip",
		httpmock.NewStringResponder(http.StatusAccepted, ""),
	)

	transport := &httpTransport{
		delegate: &delegate{
			logger: nil,
			localEndpoint: &url.URL{
				Scheme: "http",
				Host:   "invalid",
			},
			endpoints: map[string]*url.URL{
				"node": {
					Scheme: "http",
					Host:   "1.2.3.4:5678",
				},
			},
		},
	}

	_, err := transport.WriteToAddress(nil, memberlist.Address{Name: "node"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGossipStream(t *testing.T) {
	delegate := &delegate{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		localEndpoint: &url.URL{
			Scheme: "http",
			Host:   "127.0.0.1:9113",
		},
		endpoints: map[string]*url.URL{
			"node": {
				Scheme: "http",
				Host:   "127.0.0.1:9113",
			},
		},
	}

	connM := actor.NewMailbox[net.Conn]()
	connM.Start()
	defer connM.Stop()

	packetM := actor.NewMailbox[*memberlist.Packet]()
	packetM.Start()
	defer packetM.Stop()

	transport := &httpTransport{
		delegate: delegate,
		cookie:   "",
		connM:    connM,
		packetM:  packetM,
	}

	readyC := make(chan error, 1)

	server := &http.Server{
		Addr:    "127.0.0.1:9113",
		Handler: transport,
	}
	go func() {
		l, err := net.Listen("tcp", server.Addr)
		if err != nil {
			readyC <- err
			close(readyC)
			return
		}

		readyC <- nil
		close(readyC)
		server.Serve(l)
	}()

	err := <-readyC
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	sender, err := transport.DialAddressTimeout(memberlist.Address{Name: "node"}, time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer sender.Close()

	receiver := <-connM.ReceiveC()

	_, err = sender.Write([]byte("test"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	buf := make([]byte, 4)
	_, err = receiver.Read(buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(buf) != "test" {
		t.Fatalf("unexpected message: %s", string(buf))
	}
}

func TestGossipPacket(t *testing.T) {
	delegate := &delegate{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		localEndpoint: &url.URL{
			Scheme: "http",
			Host:   "127.0.0.1:9113",
		},
		endpoints: map[string]*url.URL{
			"node": {
				Scheme: "http",
				Host:   "127.0.0.1:9113",
			},
		},
	}

	connM := actor.NewMailbox[net.Conn]()
	connM.Start()
	defer connM.Stop()

	packetM := actor.NewMailbox[*memberlist.Packet]()
	packetM.Start()
	defer packetM.Stop()

	transport := &httpTransport{
		delegate: delegate,
		cookie:   "",
		connM:    connM,
		packetM:  packetM,
	}

	readyC := make(chan error, 1)

	server := &http.Server{
		Addr:    "127.0.0.1:9113",
		Handler: transport,
	}
	go func() {
		l, err := net.Listen("tcp", server.Addr)
		if err != nil {
			readyC <- err
			close(readyC)
			return
		}

		readyC <- nil
		close(readyC)

		server.Serve(l)
	}()

	err := <-readyC
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	_, err = transport.WriteToAddress([]byte("test"), memberlist.Address{Name: "node"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	packet := <-packetM.ReceiveC()
	if string(packet.Buf) != "test" {
		t.Fatalf("unexpected message: %s", string(packet.Buf))
	}
}
