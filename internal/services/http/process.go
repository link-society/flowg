package http

import (
	"errors"
	"log/slog"

	"context"
	"time"

	"net"
	gohttp "net/http"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/api"
	"link-society.com/flowg/web"

	"link-society.com/flowg/internal/app/logging"
	"link-society.com/flowg/internal/utils/proctree"
)

type procHandler struct {
	logger *slog.Logger

	opts   *ServerOptions
	server *gohttp.Server
}

var _ proctree.ProcessHandler = (*procHandler)(nil)

func (h *procHandler) Init(ctx actor.Context) proctree.ProcessResult {
	apiHandler := api.NewHandler(&api.Dependencies{
		AuthStorage:   h.opts.AuthStorage,
		LogStorage:    h.opts.LogStorage,
		ConfigStorage: h.opts.ConfigStorage,

		LogNotifier:    h.opts.LogNotifier,
		PipelineRunner: h.opts.PipelineRunner,
	})
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
		Addr:      h.opts.BindAddress,
		Handler:   logging.NewMiddleware(rootHandler),
		TLSConfig: h.opts.TlsConfig,
	}

	h.logger.InfoContext(ctx, "Starting HTTP server")

	l, err := net.Listen("tcp", h.opts.BindAddress)
	if err != nil {
		h.logger.ErrorContext(
			ctx,
			"Failed to start HTTP server",
			slog.String("error", err.Error()),
		)

		return proctree.Terminate(err)
	}

	if h.opts.TlsConfig != nil {
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
