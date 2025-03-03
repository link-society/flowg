package mgmt

import (
	"errors"
	"log/slog"

	"context"
	"time"

	"net"
	"net/http"
	"net/url"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/utils/proctree"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"link-society.com/flowg/internal/cluster"
)

type state struct {
	logger *slog.Logger

	opts *ServerOptions

	listener       net.Listener
	clusterManager *cluster.Manager
	server         *http.Server
}

type listenerHandler struct {
	state *state
}

func (h *listenerHandler) Init(ctx actor.Context) proctree.ProcessResult {
	h.state.logger.InfoContext(ctx, "Start Management HTTP server")

	l, err := net.Listen("tcp", h.state.opts.BindAddress)
	if err != nil {
		h.state.logger.ErrorContext(
			ctx,
			"Failed to start Management HTTP server",
			slog.String("error", err.Error()),
		)

		return proctree.Terminate(err)
	}

	h.state.listener = l

	return proctree.Continue()
}

func (h *listenerHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	<-ctx.Done()
	return proctree.Terminate(ctx.Err())
}

func (h *listenerHandler) Terminate(ctx actor.Context, err error) error {
	return err
}

type clusterHandler struct {
	state *state
}

func (h *clusterHandler) Init(ctx actor.Context) proctree.ProcessResult {
	h.state.logger.InfoContext(ctx, "Start cluster manager")

	localEndpoint := &url.URL{
		Scheme: "http",
		Host:   h.state.listener.Addr().String(),
	}

	if h.state.opts.TlsConfig != nil {
		localEndpoint.Scheme = "https"
	}

	h.state.clusterManager = cluster.NewManager(
		h.state.opts.ClusterNodeID,
		localEndpoint,
		h.state.opts.ClusterJoinNodeID,
		h.state.opts.ClusterJoinEndpoint,
	)

	h.state.clusterManager.Start()
	err := h.state.clusterManager.WaitReady(context.Background())
	if err != nil {
		h.state.logger.ErrorContext(
			ctx,
			"Failed to start cluster manager",
			slog.String("error", err.Error()),
		)
		return proctree.Terminate(err)
	}

	return proctree.Continue()
}

func (h *clusterHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	<-ctx.Done()
	return proctree.Terminate(ctx.Err())
}

func (h *clusterHandler) Terminate(ctx actor.Context, parentErr error) error {
	h.state.logger.InfoContext(ctx, "Stop cluster manager")

	h.state.clusterManager.Stop()
	err := h.state.clusterManager.Join(context.Background())
	if err != nil {
		h.state.logger.ErrorContext(
			ctx,
			"Failed to stop cluster manager",
			slog.String("error", err.Error()),
		)
		return errors.Join(parentErr, err)
	}

	return parentErr
}

type serverHandler struct {
	state *state
}

func (h *serverHandler) Init(ctx actor.Context) proctree.ProcessResult {
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

	h.state.server = &http.Server{
		Addr:      h.state.opts.BindAddress,
		Handler:   rootHandler,
		TLSConfig: h.state.opts.TlsConfig,
	}

	if h.state.opts.TlsConfig != nil {
		go h.state.server.ServeTLS(h.state.listener, "", "")
	} else {
		go h.state.server.Serve(h.state.listener)
	}

	return proctree.Continue()
}

func (h *serverHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	<-ctx.Done()
	return proctree.Terminate(ctx.Err())
}

func (h *serverHandler) Terminate(ctx actor.Context, parentErr error) error {
	h.state.logger.InfoContext(ctx, "Stopping Management HTTP server")

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if err := h.state.server.Shutdown(ctx); err != nil {
		h.state.logger.ErrorContext(
			ctx,
			"Failed to shutdown Management HTTP server",
			slog.String("error", err.Error()),
		)

		return errors.Join(parentErr, err)
	}

	return parentErr
}
