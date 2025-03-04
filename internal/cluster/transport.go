package cluster

import (
	"fmt"
	"log/slog"

	"time"

	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"strconv"

	"crypto/tls"
	"net"
	"net/http"
	"net/url"

	"github.com/vladopajic/go-actor/actor"

	"github.com/hashicorp/memberlist"
)

const COOKIE_HEADER_NAME = "X-FlowG-ClusterKey"

type httpTransport struct {
	delegate *delegate
	cookie   string

	connM   actor.Mailbox[net.Conn]
	packetM actor.Mailbox[*memberlist.Packet]
}

func (t *httpTransport) FinalAdvertiseAddr(string, int) (net.IP, int, error) {
	host, port, err := net.SplitHostPort(t.delegate.localEndpoint.Host)
	if err != nil {
		return nil, 0, err
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return nil, 0, fmt.Errorf("failed to parse IP from %s", host)
	}

	portNumber, err := strconv.Atoi(port)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse port: %w", err)
	}

	return ip, portNumber, nil
}

func (t *httpTransport) PacketCh() <-chan *memberlist.Packet {
	return t.packetM.ReceiveC()
}

func (t *httpTransport) StreamCh() <-chan net.Conn {
	return t.connM.ReceiveC()
}

func (t *httpTransport) DialTimeout(addr string, timeout time.Duration) (net.Conn, error) {
	return t.DialAddressTimeout(memberlist.Address{Addr: addr}, timeout)
}

func (t *httpTransport) WriteTo(b []byte, addr string) (time.Time, error) {
	return t.WriteToAddress(b, memberlist.Address{Addr: addr})
}

func (t *httpTransport) WriteToAddress(b []byte, addr memberlist.Address) (time.Time, error) {
	if addr.Name == "" {
		return time.Time{}, fmt.Errorf("empty node name")
	}

	endpoint, ok := t.delegate.endpoints[addr.Name]
	if !ok {
		return time.Time{}, fmt.Errorf("endpoint not found for %s", addr.Name)
	}

	req := &http.Request{
		Method: http.MethodPost,
		URL: &url.URL{
			Scheme: endpoint.Scheme,
			Host:   endpoint.Host,
			Path:   "/cluster/gossip",
		},
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       endpoint.Host,
		Body:       io.NopCloser(bytes.NewReader(b)),
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Origin", t.delegate.localEndpoint.String())

	if t.cookie != "" {
		req.Header.Set(COOKIE_HEADER_NAME, t.cookie)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return time.Time{}, err
	}

	if resp.StatusCode != http.StatusAccepted {
		return time.Time{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return time.Now(), nil
}

func (t *httpTransport) DialAddressTimeout(addr memberlist.Address, timeout time.Duration) (net.Conn, error) {
	if addr.Name == "" {
		return nil, fmt.Errorf("empty node name")
	}

	endpoint, ok := t.delegate.endpoints[addr.Name]
	if !ok {
		return nil, fmt.Errorf("endpoint not found for %s", addr.Name)
	}

	var dialer func() (net.Conn, error)

	switch endpoint.Scheme {
	case "http":
		dialer = func() (net.Conn, error) {
			dialer := net.Dialer{Timeout: timeout}
			return dialer.Dial("tcp", endpoint.Host)
		}

	case "https":
		dialer = func() (net.Conn, error) {
			dialer := tls.Dialer{NetDialer: &net.Dialer{Timeout: timeout}}
			return dialer.Dial("tcp", endpoint.Host)
		}

	default:
		return nil, fmt.Errorf("unsupported scheme: %s", endpoint.Scheme)
	}

	conn, err := dialer()
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	req := &http.Request{
		Method: http.MethodPost,
		URL: &url.URL{
			Scheme: endpoint.Scheme,
			Host:   endpoint.Host,
			Path:   "/cluster/gossip",
		},
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       endpoint.Host,
	}
	req.Header.Set("Upgrade", "flowg")

	if t.cookie != "" {
		req.Header.Set(COOKIE_HEADER_NAME, t.cookie)
	}

	if err := req.Write(conn); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to write request: %w", err)
	}

	resp, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusSwitchingProtocols {
		conn.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if resp.Header.Get("Upgrade") != "flowg" {
		conn.Close()
		return nil, fmt.Errorf("unexpected upgrade header: %s", resp.Header.Get("Upgrade"))
	}

	return conn, nil
}

func (t *httpTransport) Shutdown() error {
	return nil
}

func (t *httpTransport) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /cluster/nodes", t.handleStatus)
	mux.HandleFunc("POST /cluster/gossip", t.handleGossip)

	mux.ServeHTTP(w, r)
}

func (t *httpTransport) handleStatus(w http.ResponseWriter, r *http.Request) {
	type nodeInfo struct {
		NodeID   string `json:"node-id"`
		Endpoint string `json:"endpoint"`
	}

	var payload struct {
		Nodes []nodeInfo `json:"nodes"`
	}

	for nodeID, endpoint := range t.delegate.endpoints {
		payload.Nodes = append(payload.Nodes, nodeInfo{
			NodeID:   nodeID,
			Endpoint: endpoint.String(),
		})
	}

	body, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (t *httpTransport) handleGossip(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Upgrade") == "flowg" {
		t.handleGossipStream(w, r)
	} else {
		t.handleGossipPacket(w, r)
	}
}

func (t *httpTransport) handleGossipStream(w http.ResponseWriter, r *http.Request) {
	if t.cookie != "" && r.Header.Get(COOKIE_HEADER_NAME) != t.cookie {
		http.Error(w, "invalid cluster key", http.StatusUnauthorized)
		return
	}

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "hijacking not supported", http.StatusNotImplemented)
		return
	}

	conn, _, err := hijacker.Hijack()
	if err != nil {
		message := fmt.Sprintf("failed to hijack connection: %v", err)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	data := []byte("HTTP/1.1 101 Switching Protocols\r\nUpgrade: flowg\r\n\r\n")
	if n, err := conn.Write(data); err != nil || n != len(data) {
		conn.Close()
		return
	}

	if err := t.connM.Send(r.Context(), conn); err != nil {
		t.delegate.logger.ErrorContext(
			r.Context(),
			"failed to accept connection",
			slog.String("error", err.Error()),
		)
		conn.Close()
		return
	}

	t.delegate.logger.InfoContext(
		r.Context(),
		"accepted connection",
		slog.String("remote", conn.RemoteAddr().String()),
	)
}

func (t *httpTransport) handleGossipPacket(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if t.cookie != "" && r.Header.Get(COOKIE_HEADER_NAME) != t.cookie {
		http.Error(w, "invalid cluster key", http.StatusUnauthorized)
		return
	}

	origin := r.Header.Get("Origin")
	if origin == "" {
		http.Error(w, "missing Origin header", http.StatusBadRequest)
		return
	}

	originUrl, err := url.Parse(origin)
	if err != nil {
		http.Error(w, "failed to parse Origin header", http.StatusBadRequest)
		return
	}

	originAddr, err := net.ResolveTCPAddr("tcp", originUrl.Host)
	if err != nil {
		http.Error(w, "failed to resolve Origin address", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}

	packet := &memberlist.Packet{
		Buf:       body,
		From:      originAddr,
		Timestamp: time.Now(),
	}

	if err := t.packetM.Send(r.Context(), packet); err != nil {
		http.Error(w, "failed to send packet", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
