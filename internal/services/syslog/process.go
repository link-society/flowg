package syslog

import (
	"errors"
	"log/slog"

	"sync"

	"net"

	gosyslog "gopkg.in/mcuadros/go-syslog.v2"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/utils/proctree"

	"link-society.com/flowg/internal/engines/pipelines"
)

type procHandler struct {
	logger *slog.Logger

	opts *ServerOptions

	channel gosyslog.LogPartsChannel
	handler *gosyslog.ChannelHandler
	server  *gosyslog.Server
}

func (h *procHandler) Init(ctx actor.Context) proctree.ProcessResult {
	h.channel = make(gosyslog.LogPartsChannel)
	h.handler = gosyslog.NewChannelHandler(h.channel)

	h.server = gosyslog.NewServer()
	h.server.SetFormat(gosyslog.Automatic)
	h.server.SetHandler(h.handler)

	h.logger.InfoContext(ctx, "Starting Syslog server")

	switch {
	case h.opts.TcpMode && h.opts.TlsConfig != nil:
		if err := h.server.ListenTCPTLS(h.opts.BindAddress, h.opts.TlsConfig); err != nil {
			h.logger.ErrorContext(
				ctx,
				"Failed to listen on TCP+TLS",
				slog.String("error", err.Error()),
			)
			return proctree.Terminate(err)
		}

	case h.opts.TcpMode && h.opts.TlsConfig == nil:
		if err := h.server.ListenTCP(h.opts.BindAddress); err != nil {
			h.logger.ErrorContext(
				ctx,
				"Failed to listen on TCP",
				slog.String("error", err.Error()),
			)
			return proctree.Terminate(err)
		}

	case !h.opts.TcpMode:
		if err := h.server.ListenUDP(h.opts.BindAddress); err != nil {
			h.logger.ErrorContext(
				ctx,
				"Failed to listen on UDP",
				slog.String("error", err.Error()),
			)
			return proctree.Terminate(err)
		}
	}

	if err := h.server.Boot(); err != nil {
		h.logger.ErrorContext(
			ctx,
			"Failed to boot server",
			slog.String("error", err.Error()),
		)
		return proctree.Terminate(err)
	}

	return proctree.Continue()
}

func (h *procHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	select {
	case <-ctx.Done():
		return proctree.Terminate(ctx.Err())

	case logParts, ok := <-h.channel:
		if !ok {
			return proctree.Terminate(nil)
		}

		if h.opts.AllowOrigins != nil {
			// no logging here to avoid potential performance issues

			client := logParts["client"].(string)
			clientIp, _, err := net.SplitHostPort(client)
			if err != nil {
				return proctree.Continue()
			}

			allowed := false

			for _, origin := range h.opts.AllowOrigins {
				if origin == clientIp {
					allowed = true
					break
				}

				_, ipNet, err := net.ParseCIDR(origin)
				if err != nil {
					continue
				}

				if ipNet.Contains(net.ParseIP(clientIp)) {
					allowed = true
					break
				}
			}

			if !allowed {
				return proctree.Continue()
			}
		}

		pipelineNames, err := h.opts.ConfigStorage.ListPipelines(ctx)
		if err != nil {
			h.logger.ErrorContext(
				ctx,
				"Failed to list pipelines",
				slog.String("error", err.Error()),
			)
			return proctree.Continue()
		}

		wg := sync.WaitGroup{}

		for _, pipelineName := range pipelineNames {
			wg.Add(1)
			go func(pipelineName string) {
				defer wg.Done()

				record := parseLogParts(logParts)

				err := h.opts.PipelineRunner.Run(
					ctx,
					pipelineName,
					pipelines.SYSLOG_ENTRYPOINT,
					record,
				)
				if err != nil {
					h.logger.ErrorContext(
						ctx,
						"Failed to process log entry",
						slog.String("pipeline", pipelineName),
						slog.String("error", err.Error()),
					)
				}
			}(pipelineName)
		}

		wg.Wait()

		return proctree.Continue()
	}
}

func (h *procHandler) Terminate(ctx actor.Context, err error) error {
	h.logger.InfoContext(ctx, "Stopping Syslog server")

	if newErr := h.server.Kill(); newErr != nil {
		h.logger.ErrorContext(
			ctx,
			"Failed to kill server",
			slog.String("error", newErr.Error()),
		)
		return errors.Join(err, newErr)
	}

	h.server.Wait()
	return err
}
