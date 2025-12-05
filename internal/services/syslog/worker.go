package syslog

import (
	"log/slog"

	"sync"

	"net"

	gosyslog "gopkg.in/mcuadros/go-syslog.v2"

	"github.com/vladopajic/go-actor/actor"

	"link-society.com/flowg/internal/engines/pipelines"
	"link-society.com/flowg/internal/storage/config"
)

type worker struct {
	logger         *slog.Logger
	channel        gosyslog.LogPartsChannel
	configStorage  config.Storage
	pipelineRunner pipelines.Runner
}

var _ actor.Worker = (*worker)(nil)

func (w *worker) DoWork(ctx actor.Context) actor.WorkerStatus {
	select {
	case <-ctx.Done():
		return actor.WorkerEnd

	case logParts, ok := <-w.channel:
		if !ok {
			return actor.WorkerEnd
		}

		systemConfig, err := w.configStorage.ReadSystemConfig(ctx)
		if err != nil {
			w.logger.ErrorContext(
				ctx,
				"Failed to get allowed origins for syslog",
				slog.String("error", err.Error()),
			)
			return actor.WorkerContinue
		}

		if systemConfig.SyslogAllowedOrigins != nil {
			// no logging here to avoid potential performance issues

			client := logParts["client"].(string)
			clientIp, _, err := net.SplitHostPort(client)
			if err != nil {
				return actor.WorkerContinue
			}

			allowed := false

			for _, origin := range systemConfig.SyslogAllowedOrigins {
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
				return actor.WorkerContinue
			}
		}

		pipelineNames, err := w.configStorage.ListPipelines(ctx)
		if err != nil {
			w.logger.ErrorContext(
				ctx,
				"Failed to list pipelines",
				slog.String("error", err.Error()),
			)
			return actor.WorkerContinue
		}

		wg := sync.WaitGroup{}

		for _, pipelineName := range pipelineNames {
			wg.Add(1)
			go func(pipelineName string) {
				defer wg.Done()

				record := parseLogParts(logParts)

				err := w.pipelineRunner.Run(
					ctx,
					pipelineName,
					pipelines.SYSLOG_ENTRYPOINT,
					record,
				)
				if err != nil {
					w.logger.ErrorContext(
						ctx,
						"Failed to process log entry",
						slog.String("pipeline", pipelineName),
						slog.String("error", err.Error()),
					)
				}
			}(pipelineName)
		}

		wg.Wait()

		return actor.WorkerContinue
	}
}
