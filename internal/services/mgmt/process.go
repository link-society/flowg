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

	"link-society.com/flowg/internal/app/logging"
	"link-society.com/flowg/internal/utils/proctree"
)

type procHandler struct {
	logger *slog.Logger

	bindAddress string
	tlsConfig   *tls.Config
	server      *http.Server
}

func (h *procHandler) Init(ctx actor.Context) proctree.ProcessResult {
	metricsHandler := promhttp.Handler()

	rootHandler := http.NewServeMux()
	rootHandler.HandleFunc(
		"/health",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK\r\n"))
		},
	)
	rootHandler.Handle("/metrics", metricsHandler)
	h.server = &http.Server{
		Addr:      h.bindAddress,
		Handler:   logging.NewMiddleware(rootHandler),
		TLSConfig: h.tlsConfig,
	}

	h.logger.InfoContext(ctx, "Starting Management HTTP server")

	l, err := net.Listen("tcp", h.bindAddress)
	if err != nil {
		h.logger.ErrorContext(
			ctx,
			"Failed to start Management HTTP server",
			slog.String("error", err.Error()),
		)

		return proctree.Terminate(err)
	}

	if h.tlsConfig != nil {
		go h.server.ServeTLS(l, "", "")
	} else {
		go h.server.Serve(l)
	}

	return proctree.Continue()
}

func (h *procHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	<-ctx.Done()
	return proctree.Terminate(ctx.Err())
}

func (h *procHandler) Terminate(ctx actor.Context, err error) error {
	h.logger.InfoContext(ctx, "Stopping Management HTTP server")

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	stopErr := h.server.Shutdown(ctx)

	if stopErr != nil {
		h.logger.ErrorContext(
			ctx,
			"Failed to shutdown Management HTTP server",
			slog.String("error", err.Error()),
		)

		err = errors.Join(err, stopErr)
	}

	return err
}
