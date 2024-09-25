package service

import (
	"context"
	"log/slog"
	"net"

	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"link-society.com/flowg/internal/app/logging"
	"link-society.com/flowg/internal/data/auth"
	"link-society.com/flowg/internal/data/config"
	"link-society.com/flowg/internal/data/lognotify"
	"link-society.com/flowg/internal/data/logstorage"

	"link-society.com/flowg/api"
	"link-society.com/flowg/web"
)

type httpA struct {
	bindAddress string
	srv         *http.Server

	startErrC chan struct{}
	stopErrC  chan struct{}
}

func newHttpA(
	bindAddress string,
	authDb *auth.Database,
	logDb *logstorage.Storage,
	configStorage *config.Storage,
	logNotifier *lognotify.LogNotifier,
) *httpA {
	apiHandler := api.NewHandler(authDb, logDb, configStorage, logNotifier)
	webHandler := web.NewHandler()
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
	rootHandler.Handle("/api/", apiHandler)
	rootHandler.Handle("/web/", webHandler)

	rootHandler.HandleFunc(
		"GET /{$}",
		func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/web/", http.StatusPermanentRedirect)
		},
	)

	server := &http.Server{
		Addr:    bindAddress,
		Handler: logging.NewMiddleware(rootHandler),
	}

	return &httpA{
		bindAddress: bindAddress,
		srv:         server,

		startErrC: make(chan struct{}, 1),
		stopErrC:  make(chan struct{}, 1),
	}
}

func (a *httpA) Start() {
	defer close(a.startErrC)

	slog.Info(
		"Starting HTTP server",
		"channel", "http",
		"http.bind", a.bindAddress,
	)

	l, err := net.Listen("tcp", a.bindAddress)
	if err != nil {
		slog.Error(
			"Failed to start HTTP server",
			"channel", "http",
			"http.bind", a.bindAddress,
			"error", err,
		)
		a.startErrC <- struct{}{}
		return
	}

	go a.srv.Serve(l)
}

func (a *httpA) Stop() {
	defer close(a.stopErrC)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.srv.Shutdown(ctx); err != nil {
		slog.Error(
			"Failed to shutdown HTTP server",
			"channel", "main",
			"http.bind", a.bindAddress,
			"error", err,
		)
		a.stopErrC <- struct{}{}
		return
	}
}

func (a *httpA) StartErrC() <-chan struct{} {
	return a.startErrC
}

func (a *httpA) StopErrC() <-chan struct{} {
	return a.stopErrC
}
