package syslog

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/vladopajic/go-actor/actor"
	gosyslog "gopkg.in/mcuadros/go-syslog.v2"
	gosyslogformat "gopkg.in/mcuadros/go-syslog.v2/format"

	"link-society.com/flowg/internal/data/config"
	"link-society.com/flowg/internal/data/lognotify"
	"link-society.com/flowg/internal/data/logstorage"
	"link-society.com/flowg/internal/data/pipelines"
)

func NewServer(
	bindAddress string,
	configSorage *config.Storage,
	logStorage *logstorage.Storage,
	logNotifier *lognotify.LogNotifier,
) actor.Actor {
	worker := &worker{
		bindAddress: bindAddress,

		configStorage: configSorage,
		logStorage:    logStorage,
		logNotifier:   logNotifier,
	}

	return actor.New(
		worker,
		actor.OptOnStart(worker.onStart),
		actor.OptOnStop(worker.onStop),
	)
}

type worker struct {
	bindAddress string

	configStorage *config.Storage
	logStorage    *logstorage.Storage
	logNotifier   *lognotify.LogNotifier

	channel gosyslog.LogPartsChannel
	handler *gosyslog.ChannelHandler
	server  *gosyslog.Server
}

func (w *worker) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case logParts, ok := <-w.channel:
		if !ok {
			return actor.WorkerEnd
		}

		pipelineSys := config.NewPipelineSystem(w.configStorage)

		pipelineNames, err := pipelineSys.List()
		if err != nil {
			slog.ErrorContext(
				ctx,
				"Failed to list pipelines",
				"channel", "syslog",
				"error", err,
			)
			return actor.WorkerContinue
		}

		wg := sync.WaitGroup{}

		for _, pipelineName := range pipelineNames {
			wg.Add(1)
			go func(pipelineName string) {
				defer wg.Done()

				pipeline, err := pipelines.Build(pipelineSys, pipelineName)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to build pipeline",
						"channel", "syslog",
						"pipeline", pipelineName,
						"error", err.Error(),
					)
					return
				}

				entry := parseLogParts(logParts)
				runner := pipelines.NewRunner(
					ctx,
					w.configStorage,
					w.logStorage,
					w.logNotifier,
				)
				err = runner.Run(pipeline, pipelines.SYSLOG_ENTRYPOINT, entry)
				if err != nil {
					slog.ErrorContext(
						ctx,
						"Failed to process log entry",
						"channel", "syslog",
						"pipeline", pipelineName,
						"error", err.Error(),
					)
				}
			}(pipelineName)
		}

		wg.Wait()

		return actor.WorkerContinue
	}
}

func (w *worker) onStart(ctx actor.Context) {
	slog.InfoContext(
		ctx,
		"Starting syslog server",
		"channel", "syslog",
		"udp.bind", w.bindAddress,
	)

	w.channel = make(gosyslog.LogPartsChannel)
	w.handler = gosyslog.NewChannelHandler(w.channel)

	w.server = gosyslog.NewServer()
	w.server.SetFormat(gosyslog.Automatic)
	w.server.SetHandler(w.handler)
	w.server.ListenUDP(w.bindAddress)

	w.server.Boot()
}

func (w *worker) onStop() {
	w.server.Kill()
	w.server.Wait()
}

func parseLogParts(logParts gosyslogformat.LogParts) *logstorage.LogEntry {
	fields := make(map[string]string, len(logParts))

	for key, value := range logParts {
		switch value := value.(type) {
		case string:
			fields[key] = value

		case []byte:
			fields[key] = string(value)

		default:
			fields[key] = fmt.Sprintf("%v", value)
		}
	}

	return logstorage.NewLogEntry(fields)
}
