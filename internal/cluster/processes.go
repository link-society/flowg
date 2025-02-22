package cluster

import (
	"errors"
	"fmt"
	"log/slog"

	"bytes"
	"encoding/json"

	"io"
	"net"
	"net/http"
	"net/url"

	"time"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/utils/proctree"

	"github.com/hashicorp/raft"

	"link-society.com/flowg/internal/cluster/raftfsm"
	"link-society.com/flowg/internal/cluster/rafthttp"
	"link-society.com/flowg/internal/cluster/raftmembership"
	"link-society.com/flowg/internal/cluster/raftstore"

	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"
)

type state struct {
	logger *slog.Logger

	listener net.Listener
	dial     rafthttp.Dial
	timeout  time.Duration
	connM    actor.Mailbox[net.Conn]

	nodeID   string
	joinAddr string

	authStorage   *auth.Storage
	configStorage *config.Storage
	logStorage    *log.Storage

	store            *raftstore.Store
	snapshots        raft.SnapshotStore
	raft             *raft.Raft
	membershipServer *raftmembership.Server
	httpHandler      http.Handler
}

type storeHandler struct {
	state *state
	dir   string
}

func (h *storeHandler) Init(ctx actor.Context) proctree.ProcessResult {
	h.state.logger.InfoContext(ctx, "Initialize consensus store")

	store, err := raftstore.New(h.dir)
	if err != nil {
		h.state.logger.ErrorContext(
			ctx,
			"Failed to initialize consensus store",
			slog.String("error", err.Error()),
		)
		return proctree.Terminate(err)
	}

	h.state.store = store

	return proctree.Continue()
}

func (h *storeHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	<-ctx.Done()
	return proctree.Terminate(ctx.Err())
}

func (h *storeHandler) Terminate(ctx actor.Context, parentErr error) error {
	h.state.logger.InfoContext(ctx, "Teardown consensus store")

	err := h.state.store.Close()
	if err != nil {
		h.state.logger.ErrorContext(
			ctx,
			"Failed to teardown consensus store",
			slog.String("error", err.Error()),
		)
		return errors.Join(parentErr, err)
	}

	return parentErr
}

type snapshotsHandler struct {
	state *state
	dir   string
}

func (h *snapshotsHandler) Init(ctx actor.Context) proctree.ProcessResult {
	h.state.logger.InfoContext(ctx, "Initialize consensus snapshots store")

	snapshots, err := raft.NewFileSnapshotStore(
		h.dir,
		1,
		&raftLogger{logger: h.state.logger},
	)
	if err != nil {
		h.state.logger.ErrorContext(
			ctx,
			"Failed to initialize consensus snapshots store",
			slog.String("error", err.Error()),
		)
		return proctree.Terminate(err)
	}

	h.state.snapshots = snapshots

	return proctree.Continue()
}

func (h *snapshotsHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	<-ctx.Done()
	return proctree.Terminate(ctx.Err())
}

func (h *snapshotsHandler) Terminate(ctx actor.Context, err error) error {
	return err
}

type consensusHandler struct {
	state *state
}

func (h *consensusHandler) Init(ctx actor.Context) proctree.ProcessResult {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(h.state.nodeID)
	config.Logger = newRaftLogger("raft", h.state.logger)

	transport := raft.NewNetworkTransportWithLogger(
		rafthttp.NewLayer(
			h.state.connM,
			"/cluster/consensus",
			h.state.listener.Addr(),
			h.state.dial,
		),
		2,
		h.state.timeout,
		newRaftLogger("raft.transport", h.state.logger),
	)

	h.state.logger.InfoContext(ctx, "Initialize consensus server")

	raftCluster, err := raft.NewRaft(
		config,
		raftfsm.New(h.state.authStorage, h.state.configStorage, h.state.logStorage),
		h.state.store,
		h.state.store,
		h.state.snapshots,
		transport,
	)
	if err != nil {
		h.state.logger.ErrorContext(
			ctx,
			"Failed to initialize consensus server",
			slog.String("error", err.Error()),
		)
		return proctree.Terminate(err)
	}

	if h.state.joinAddr == "" {
		bootstrapConfig := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		raftCluster.BootstrapCluster(bootstrapConfig)
	}

	h.state.raft = raftCluster

	return proctree.Continue()
}

func (h *consensusHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	<-ctx.Done()
	return proctree.Terminate(ctx.Err())
}

func (h *consensusHandler) Terminate(ctx actor.Context, parentErr error) error {
	h.state.logger.InfoContext(ctx, "Teardown consensus server")

	future := h.state.raft.Shutdown()

	if err := future.Error(); err != nil {
		h.state.logger.ErrorContext(
			ctx,
			"Failed to teardown consensus server",
			slog.String("error", err.Error()),
		)
		return errors.Join(parentErr, err)
	}

	return parentErr
}

type membershipHandler struct {
	state *state
}

func (h *membershipHandler) Init(ctx actor.Context) proctree.ProcessResult {
	h.state.logger.InfoContext(ctx, "Initialize membership server")

	h.state.membershipServer = raftmembership.NewServer(h.state.raft)
	h.state.membershipServer.Start()

	err := h.state.membershipServer.WaitReady(ctx)
	if err != nil {
		h.state.logger.ErrorContext(
			ctx,
			"Failed to initialize membership server",
			slog.String("error", err.Error()),
		)
		return proctree.Terminate(err)
	}

	return proctree.Continue()
}

func (h *membershipHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	<-ctx.Done()
	return proctree.Terminate(ctx.Err())
}

func (h *membershipHandler) Terminate(ctx actor.Context, parentErr error) error {
	h.state.logger.InfoContext(ctx, "Teardown membership server")
	h.state.membershipServer.Stop()

	err := h.state.membershipServer.Join(ctx)
	if err != nil {
		h.state.logger.ErrorContext(
			ctx,
			"Failed to teardown membership server",
			slog.String("error", err.Error()),
		)
		return errors.Join(parentErr, err)
	}

	return parentErr
}

type httpHandler struct {
	state *state
}

func (h *httpHandler) Init(ctx actor.Context) proctree.ProcessResult {
	h.state.httpHandler = rafthttp.NewHandler(
		h.state.connM,
		h.state.membershipServer,
		h.state.timeout,
		h.state.joinAddr,
	)

	return proctree.Continue()
}

func (h *httpHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	<-ctx.Done()
	return proctree.Terminate(ctx.Err())
}

func (h *httpHandler) Terminate(ctx actor.Context, err error) error {
	return err
}

type joinHandler struct {
	state *state
}

func (h *joinHandler) Init(ctx actor.Context) proctree.ProcessResult {
	if h.state.joinAddr == "" {
		return proctree.Continue()
	}

	addr := h.state.listener.Addr().String()

	h.state.logger.InfoContext(ctx, "Join cluster")

	nodesUrl, err := url.JoinPath(h.state.joinAddr, "/cluster/nodes/", h.state.nodeID)
	if err != nil {
		h.state.logger.ErrorContext(
			ctx,
			"invalid join address",
			slog.String("error", err.Error()),
		)
		return proctree.Terminate(err)
	}

	payload, err := json.Marshal(map[string]string{"address": addr})
	if err != nil {
		h.state.logger.ErrorContext(
			ctx,
			"could not encode join request payload",
			slog.String("error", err.Error()),
		)
		return proctree.Terminate(err)
	}

	req, err := http.NewRequest(http.MethodPut, nodesUrl, bytes.NewReader(payload))
	if err != nil {
		h.state.logger.ErrorContext(
			ctx,
			"could not send join request",
			slog.String("error", err.Error()),
		)
		return proctree.Terminate(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		h.state.logger.ErrorContext(
			ctx,
			"invalid join response",
			slog.String("error", err.Error()),
		)
		return proctree.Terminate(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		h.state.logger.ErrorContext(
			ctx,
			"invalid join response status code",
			slog.Group("error",
				slog.String("status", resp.Status),
				slog.String("body", string(body)),
			),
		)
		return proctree.Terminate(err)
	}

	return proctree.Continue()
}

func (h *joinHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	<-ctx.Done()
	return proctree.Terminate(ctx.Err())
}

func (h *joinHandler) Terminate(ctx actor.Context, parentErr error) error {
	if h.state.joinAddr == "" {
		return parentErr
	}

	h.state.logger.InfoContext(ctx, "Leave cluster")

	nodesUrl, err := url.JoinPath(h.state.joinAddr, "/cluster/nodes/", h.state.nodeID)
	if err != nil {
		h.state.logger.ErrorContext(
			ctx,
			"invalid leave address",
			slog.String("error", err.Error()),
		)
		return errors.Join(parentErr, err)
	}

	req, err := http.NewRequest(http.MethodDelete, nodesUrl, nil)
	if err != nil {
		h.state.logger.ErrorContext(
			ctx,
			"coud not send leave request",
			slog.String("error", err.Error()),
		)
		return errors.Join(parentErr, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		h.state.logger.ErrorContext(
			ctx,
			"invalid leave response",
			slog.String("error", err.Error()),
		)
		return errors.Join(parentErr, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		h.state.logger.ErrorContext(
			ctx,
			"invalid leave response status code",
			slog.Group("error",
				slog.String("status", resp.Status),
				slog.String("body", string(body)),
			),
		)
		return errors.Join(parentErr, fmt.Errorf("invalid lieave response"))
	}

	return parentErr
}
