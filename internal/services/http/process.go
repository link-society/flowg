package http

import (
	"errors"
	"log/slog"

	"context"
	"time"

	"crypto/tls"
	"net"
	gohttp "net/http"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/api"
	"link-society.com/flowg/web"

	"link-society.com/flowg/internal/app/logging"
	"link-society.com/flowg/internal/utils/proctree"

	"link-society.com/flowg/internal/storage/auth"
	"link-society.com/flowg/internal/storage/config"
	"link-society.com/flowg/internal/storage/log"

	"link-society.com/flowg/internal/engines/lognotify"
	"link-society.com/flowg/internal/engines/pipelines"
)

type procHandler struct {
	logger *slog.Logger

	bindAddress string
	tlsConfig   *tls.Config
	server      *gohttp.Server

	authStorage   *auth.Storage
	configStorage *config.Storage
	logStorage    *log.Storage

	logNotifier    *lognotify.LogNotifier
	pipelineRunner *pipelines.Runner
}

func (h *procHandler) Init(ctx actor.Context) proctree.ProcessResult {
	apiHandler := api.NewHandler(
		h.authStorage,
		h.logStorage,
		h.configStorage,
		h.logNotifier,
		h.pipelineRunner,
	)
	webHandler := web.NewHandler()

	rootHandler := gohttp.NewServeMux()
	rootHandler.Handle("/api/", apiHandler)
	rootHandler.Handle("/web/", webHandler)

	rootHandler.HandleFunc(
		"GET /{$}",
		func(w gohttp.ResponseWriter, r *gohttp.Request) {
			gohttp.Redirect(w, r, "/web/", gohttp.StatusPermanentRedirect)
		},
	)

	h.server = &gohttp.Server{
		Addr:      h.bindAddress,
		Handler:   logging.NewMiddleware(rootHandler),
		TLSConfig: h.tlsConfig,
	}

	h.logger.InfoContext(ctx, "Starting HTTP server")

	l, err := net.Listen("tcp", h.bindAddress)
	if err != nil {
		h.logger.ErrorContext(
			ctx,
			"Failed to start HTTP server",
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
	h.logger.InfoContext(ctx, "Stopping HTTP server")

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	stopErr := h.server.Shutdown(ctx)

	if stopErr != nil {
		h.logger.ErrorContext(
			ctx,
			"Failed to shutdown HTTP server",
			slog.String("error", err.Error()),
		)

		err = errors.Join(err, stopErr)
	}

	return err
}
