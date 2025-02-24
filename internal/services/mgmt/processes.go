package mgmt

import (
	"errors"
	"log/slog"

	"context"
	"time"

	"crypto/tls"
	"net"
	"net/http"

	"github.com/vladopajic/go-actor/actor"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"link-society.com/flowg/internal/utils/proctree"

	"link-society.com/flowg/internal/cluster"

	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"
)

type state struct {
	logger *slog.Logger

	bindAddress string
	tlsConfig   *tls.Config

	clusterStateDir string
	clusterNodeID   string
	clusterJoinAddr string
	clusterTimeout  time.Duration

	listener       net.Listener
	clusterManager *cluster.Manager
	httpServer     *http.Server

	authStorage   *auth.Storage
	configStorage *config.Storage
	logStorage    *log.Storage
}

type listenerHandler struct {
	state *state
}

func (h *listenerHandler) Init(ctx actor.Context) proctree.ProcessResult {
	h.state.logger.InfoContext(ctx, "Starting Management HTTP server")

	listener, err := net.Listen("tcp", h.state.bindAddress)
	if err != nil {
		h.state.logger.ErrorContext(
			ctx,
			"Failed to start Management HTTP server",
			slog.String("error", err.Error()),
		)

		return proctree.Terminate(err)
	}

	h.state.listener = listener

	return proctree.Continue()
}

func (h *listenerHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	<-ctx.Done()
	return proctree.Terminate(ctx.Err())
}

func (h *listenerHandler) Terminate(ctx actor.Context, err error) error {
	return err
}

type clusterManagerHandler struct {
	state *state
}

func (h *clusterManagerHandler) Init(ctx actor.Context) proctree.ProcessResult {
	h.state.clusterManager = cluster.NewManager(
		h.state.clusterStateDir,
		h.state.clusterNodeID,
		h.state.clusterJoinAddr,

		h.state.listener,
		h.state.tlsConfig,
		h.state.clusterTimeout,

		h.state.authStorage,
		h.state.configStorage,
		h.state.logStorage,
	)
	h.state.clusterManager.Start()

	err := h.state.clusterManager.WaitReady(ctx)
	if err != nil {
		h.state.logger.ErrorContext(
			ctx,
			"Failed to initialize cluster manager",
			slog.String("error", err.Error()),
		)
		return proctree.Terminate(err)
	}

	return proctree.Continue()
}

func (h *clusterManagerHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	<-ctx.Done()
	return proctree.Terminate(ctx.Err())
}

func (h *clusterManagerHandler) Terminate(ctx actor.Context, parentErr error) error {
	h.state.clusterManager.Stop()

	err := h.state.clusterManager.Join(ctx)
	if err != nil {
		h.state.logger.ErrorContext(
			ctx,
			"Failed to teardown cluster manager",
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
	metricsHandler := promhttp.Handler()
	clusterHandler := h.state.clusterManager.HttpHandler()

	rootHandler := http.NewServeMux()

	rootHandler.HandleFunc(
		"/health",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK\r\n"))
		},
	)

	rootHandler.Handle("/metrics", metricsHandler)
	rootHandler.Handle("/cluster/", clusterHandler)

	h.state.httpServer = &http.Server{
		Addr:      h.state.listener.Addr().String(),
		Handler:   rootHandler,
		TLSConfig: h.state.tlsConfig,
	}

	if h.state.tlsConfig != nil {
		go h.state.httpServer.ServeTLS(h.state.listener, "", "")
	} else {
		go h.state.httpServer.Serve(h.state.listener)
	}

	return proctree.Continue()
}

func (h *httpHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	<-ctx.Done()
	return proctree.Terminate(ctx.Err())
}

func (h *httpHandler) Terminate(ctx actor.Context, parentErr error) error {
	h.state.logger.InfoContext(ctx, "Stopping Management HTTP server")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := h.state.httpServer.Shutdown(ctx)
	if err != nil {
		h.state.logger.ErrorContext(
			ctx,
			"Failed to shutdown Management HTTP server",
			slog.String("error", err.Error()),
		)

		return errors.Join(parentErr, err)
	}

	return parentErr
}
