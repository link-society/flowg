package rafthttp

import (
	"fmt"
	"log/slog"

	"time"

	"bytes"
	"encoding/json"

	"net"
	"net/http"
	"net/url"

	"github.com/vladopajic/go-actor/actor"

	"github.com/hashicorp/raft"

	"link-society.com/flowg/internal/cluster/raftmembership"
)

func NewHandler(
	connM actor.MailboxSender[net.Conn],
	membershipServer *raftmembership.Server,
	membershipRequestTimeout time.Duration,
	joinAddr string,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /cluster/consensus",
		handleConsensus(connM),
	)
	mux.HandleFunc(
		"GET /cluster/nodes",
		handleStatus(joinAddr, membershipServer),
	)
	mux.HandleFunc(
		"PUT /cluster/nodes/{id}",
		handleJoinCluster(joinAddr, membershipServer, membershipRequestTimeout),
	)
	mux.HandleFunc(
		"DELETE /cluster/nodes/{id}",
		handleLeaveCluster(joinAddr, membershipServer, membershipRequestTimeout),
	)

	return mux
}

func handleConsensus(connM actor.MailboxSender[net.Conn]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Upgrade") != "raft" {
			http.Error(w, "expected `Upgrade: raft` header", http.StatusBadRequest)
			return
		}

		hijacker, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "hijacking not supported", http.StatusInternalServerError)
			return
		}

		conn, _, err := hijacker.Hijack()
		if err != nil {
			message := fmt.Errorf("failed to hijack connection: %w", err).Error()
			http.Error(w, message, http.StatusInternalServerError)
			return
		}

		data := []byte("HTTP/1.1 101 Switching Protocols\r\nUpgrade: raft\r\n\r\n")
		if n, err := conn.Write(data); err != nil || n != len(data) {
			conn.Close()
			return
		}

		if err := connM.Send(r.Context(), conn); err != nil {
			slog.ErrorContext(
				r.Context(),
				"failed to accept connection",
				slog.String("channel", "cluster.consensus"),
				slog.String("error", err.Error()),
			)
			conn.Close()
		}
	}
}

func handleStatus(
	joinAddr string,
	membershipServer *raftmembership.Server,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reply := make(chan []raftmembership.ServerInfo, 1)
		req := raftmembership.NewStatusRequest(reply)
		err := membershipServer.SendRequest(r.Context(), req)
		if err != nil {
			if err == raft.ErrNotLeader {
				u := getRedirectUrl(joinAddr, r)
				slog.InfoContext(
					r.Context(),
					"redirecting to leader",
					slog.String("channel", "cluster.status"),
					slog.String("url", u),
				)
				http.Redirect(w, r, u, http.StatusTemporaryRedirect)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		select {
		case <-r.Context().Done():
			http.Error(w, r.Context().Err().Error(), http.StatusRequestTimeout)

		case nodes := <-reply:
			body := struct {
				Nodes []raftmembership.ServerInfo `json:"nodes"`
			}{Nodes: nodes}

			buf := new(bytes.Buffer)
			if err := json.NewEncoder(buf).Encode(body); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(buf.Bytes())
		}
	}
}

func handleJoinCluster(
	joinAddr string,
	membershipServer *raftmembership.Server,
	membershipRequestTimeout time.Duration,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := raft.ServerID(r.PathValue("id"))

		var body struct {
			Address string `json:"address"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		addr := raft.ServerAddress(body.Address)
		req := raftmembership.NewJoinRequest(id, addr, membershipRequestTimeout)
		err := membershipServer.SendRequest(r.Context(), req)
		if err == raft.ErrNotLeader {
			u := getRedirectUrl(joinAddr, r)
			slog.InfoContext(
				r.Context(),
				"redirecting to leader",
				slog.String("channel", "cluster.status"),
				slog.String("url", u),
			)
			http.Redirect(w, r, u, http.StatusTemporaryRedirect)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func handleLeaveCluster(
	joinAddr string,
	membershipServer *raftmembership.Server,
	membershipRequestTimeout time.Duration,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := raft.ServerID(r.PathValue("id"))

		req := raftmembership.NewLeaveRequest(id, membershipRequestTimeout)
		err := membershipServer.SendRequest(r.Context(), req)
		if err != nil {
			if err == raft.ErrNotLeader {
				u := getRedirectUrl(joinAddr, r)
				slog.InfoContext(
					r.Context(),
					"redirecting to leader",
					slog.String("channel", "cluster.status"),
					slog.String("url", u),
				)
				http.Redirect(w, r, u, http.StatusTemporaryRedirect)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func getRedirectUrl(joinAddr string, r *http.Request) string {
	u, err := url.JoinPath(joinAddr, r.URL.Path)
	if err != nil {
		panic(err)
	}

	return u
}
