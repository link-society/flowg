package syslog

import (
	"log/slog"
	"sync"

	"github.com/vladopajic/go-actor/actor"

	gosyslog "gopkg.in/mcuadros/go-syslog.v2"

	"link-society.com/flowg/internal/engines/pipelines"
)

type workerState interface {
	DoWork(ctx actor.Context, worker *worker) workerState
}

type workerStarting struct {
	bindAddress string
}

type workerRunning struct {
	channel gosyslog.LogPartsChannel
	handler *gosyslog.ChannelHandler
	server  *gosyslog.Server
}

type workerStopping struct {
	server *gosyslog.Server
}

func (s *workerStarting) DoWork(ctx actor.Context, worker *worker) workerState {
	channel := make(gosyslog.LogPartsChannel)
	handler := gosyslog.NewChannelHandler(channel)

	server := gosyslog.NewServer()
	server.SetFormat(gosyslog.Automatic)
	server.SetHandler(handler)

	worker.logger.InfoContext(
		ctx,
		"Starting Syslog server",
		slog.Group("udp",
			slog.String("bind", s.bindAddress),
		),
	)

	if err := server.ListenUDP(s.bindAddress); err != nil {
		worker.logger.ErrorContext(
			ctx,
			"Failed to listen on UDP",
			slog.Group("udp",
				slog.String("bind", s.bindAddress),
			),
			slog.String("error", err.Error()),
		)
		worker.startCond.Broadcast(err)
		return nil
	}

	if err := server.Boot(); err != nil {
		worker.logger.ErrorContext(
			ctx,
			"Failed to boot server",
			slog.Group("udp",
				slog.String("bind", s.bindAddress),
			),
			slog.String("error", err.Error()),
		)
		worker.startCond.Broadcast(err)
		return nil
	}

	worker.startCond.Broadcast(nil)

	return &workerRunning{
		channel: channel,
		handler: handler,
		server:  server,
	}
}

func (s *workerRunning) DoWork(ctx actor.Context, worker *worker) workerState {
	select {
	case <-ctx.Done():
		return &workerStopping{server: s.server}

	case logParts, ok := <-s.channel:
		if !ok {
			return &workerStopping{server: s.server}
		}

		pipelineNames, err := worker.configStorage.ListPipelines(ctx)
		if err != nil {
			slog.ErrorContext(
				ctx,
				"Failed to list pipelines",
				slog.String("error", err.Error()),
			)
			return s
		}

		wg := sync.WaitGroup{}

		for _, pipelineName := range pipelineNames {
			wg.Add(1)
			go func(pipelineName string) {
				defer wg.Done()

				record := parseLogParts(logParts)

				err := worker.pipelineRunner.Run(
					ctx,
					pipelineName,
					pipelines.SYSLOG_ENTRYPOINT,
					record,
				)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to process log entry",
						slog.String("pipeline", pipelineName),
						slog.String("error", err.Error()),
					)
				}
			}(pipelineName)
		}

		wg.Wait()

		return s
	}
}

func (s *workerStopping) DoWork(ctx actor.Context, worker *worker) workerState {
	slog.InfoContext(ctx, "Stopping Syslog server")

	if err := s.server.Kill(); err != nil {
		slog.ErrorContext(
			ctx,
			"Failed to kill server",
			slog.String("error", err.Error()),
		)
		worker.stopCond.Broadcast(err)
		return nil
	}

	s.server.Wait()

	worker.stopCond.Broadcast(nil)
	return nil
}
