package mgmt

import (
	"errors"
	"log/slog"

	"context"
	"time"

	"crypto/tls"
	"net/http"
	"net/http/pprof"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/utils/proctree"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"link-society.com/flowg/internal/cluster"
)

type serverHandler struct {
	logger *slog.Logger

	bindAddress string
	tlsConfig   *tls.Config

	listenerH      *listenerHandler
	clusterManager *cluster.Manager

	server *http.Server
}

var _ proctree.ProcessHandler = (*serverHandler)(nil)

func (h *serverHandler) Init(ctx actor.Context) proctree.ProcessResult {
	h.logger.InfoContext(ctx, "Start Management HTTP server")

	metricsHandler := promhttp.Handler()
	clusterHandler := h.clusterManager.HttpHandler()

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

	rootHandler.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	rootHandler.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	rootHandler.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	rootHandler.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	rootHandler.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	h.server = &http.Server{
		Addr:      h.bindAddress,
		Handler:   rootHandler,
		TLSConfig: h.tlsConfig,
	}

	if h.tlsConfig != nil {
		go h.server.ServeTLS(h.listenerH.listener, "", "")
	} else {
		go h.server.Serve(h.listenerH.listener)
	}

	return proctree.Continue()
}

func (h *serverHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	<-ctx.Done()
	return proctree.Terminate(ctx.Err())
}

func (h *serverHandler) Terminate(ctx actor.Context, parentErr error) error {
	h.logger.InfoContext(ctx, "Stopping Management HTTP server")

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if err := h.server.Shutdown(ctx); err != nil {
		h.logger.ErrorContext(
			ctx,
			"Failed to shutdown Management HTTP server",
			slog.String("error", err.Error()),
		)

		return errors.Join(parentErr, err)
	}

	return parentErr
}
