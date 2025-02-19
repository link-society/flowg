package mgmt

import (
	"log/slog"

	"context"
	"time"

	"crypto/tls"
	"net"
	"net/http"

	"github.com/vladopajic/go-actor/actor"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"link-society.com/flowg/internal/app/logging"
)

type workerState interface {
	DoWork(ctx actor.Context, worker *worker) workerState
}

type workerStarting struct {
	bindAddress string
	tlsConfig   *tls.Config
}

type workerRunning struct {
	server *http.Server
}

type workerStopping struct {
	server *http.Server
}

func (s *workerStarting) DoWork(ctx actor.Context, worker *worker) workerState {
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
	server := &http.Server{
		Addr:      s.bindAddress,
		Handler:   logging.NewMiddleware(rootHandler),
		TLSConfig: s.tlsConfig,
	}

	worker.logger.InfoContext(
		ctx,
		"Starting Management HTTP server",
		slog.Group("mgmt",
			slog.String("bind", s.bindAddress),
		),
	)

	l, err := net.Listen("tcp", s.bindAddress)
	if err != nil {
		worker.logger.ErrorContext(
			ctx,
			"Failed to start Management HTTP server",
			slog.Group("mgmt",
				slog.String("bind", s.bindAddress),
			),
			slog.String("error", err.Error()),
		)

		worker.startCond.Broadcast(err)
		return nil
	}

	if s.tlsConfig != nil {
		go server.ServeTLS(l, "", "")
	} else {
		go server.Serve(l)
	}

	worker.startCond.Broadcast(nil)
	return &workerRunning{server: server}
}

func (s *workerRunning) DoWork(ctx actor.Context, worker *worker) workerState {
	<-ctx.Done()
	return &workerStopping{server: s.server}
}

func (s *workerStopping) DoWork(ctx actor.Context, worker *worker) workerState {
	worker.logger.InfoContext(
		ctx,
		"Stopping Management HTTP server",
		slog.Group("mgmt",
			slog.String("bind", s.server.Addr),
		),
	)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := s.server.Shutdown(ctx)
	worker.stopCond.Broadcast(err)

	if err != nil {
		worker.logger.ErrorContext(
			ctx,
			"Failed to shutdown Management HTTP server",
			slog.Group("mgmt",
				slog.String("bind", s.server.Addr),
			),
			slog.String("error", err.Error()),
		)
		return nil
	}

	return nil
}
