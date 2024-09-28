package http

import (
	"log/slog"

	"context"
	"time"

	"net"
	gohttp "net/http"

	"github.com/vladopajic/go-actor/actor"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"link-society.com/flowg/internal/app/logging"

	"link-society.com/flowg/api"
	"link-society.com/flowg/web"
)

type workerState interface {
	DoWork(ctx actor.Context, worker *worker) workerState
}

type workerStarting struct {
	bindAddress string
}

type workerRunning struct {
	server *gohttp.Server
}

type workerStopping struct {
	server *gohttp.Server
}

func (s *workerStarting) DoWork(ctx actor.Context, worker *worker) workerState {
	apiHandler := api.NewHandler(
		worker.authStorage,
		worker.logStorage,
		worker.configStorage,
		worker.logNotifier,
		worker.pipelineRunner,
	)
	webHandler := web.NewHandler()
	metricsHandler := promhttp.Handler()

	rootHandler := gohttp.NewServeMux()
	rootHandler.HandleFunc(
		"/health",
		func(w gohttp.ResponseWriter, r *gohttp.Request) {
			w.WriteHeader(gohttp.StatusOK)
			w.Write([]byte("OK\r\n"))
		},
	)
	rootHandler.Handle("/metrics", metricsHandler)
	rootHandler.Handle("/api/", apiHandler)
	rootHandler.Handle("/web/", webHandler)

	rootHandler.HandleFunc(
		"GET /{$}",
		func(w gohttp.ResponseWriter, r *gohttp.Request) {
			gohttp.Redirect(w, r, "/web/", gohttp.StatusPermanentRedirect)
		},
	)

	server := &gohttp.Server{
		Addr:    s.bindAddress,
		Handler: logging.NewMiddleware(rootHandler),
	}

	worker.logger.InfoContext(
		ctx,
		"Starting HTTP server",
		slog.Group("http",
			slog.String("bind", s.bindAddress),
		),
	)

	l, err := net.Listen("tcp", s.bindAddress)
	if err != nil {
		worker.logger.ErrorContext(
			ctx,
			"Failed to start HTTP server",
			slog.Group("http",
				slog.String("bind", s.bindAddress),
			),
			slog.String("error", err.Error()),
		)

		worker.startCond.Broadcast(err)
		return nil
	}

	go server.Serve(l)

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
		"Stopping HTTP server",
		slog.Group("http",
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
			"Failed to shutdown HTTP server",
			slog.Group("http",
				slog.String("bind", s.server.Addr),
			),
			slog.String("error", err.Error()),
		)
		return nil
	}

	return nil
}
