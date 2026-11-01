package syslog

import (
	"log/slog"

	"net"

	gosyslog "gopkg.in/mcuadros/go-syslog.v2"

	"github.com/vladopajic/go-actor/actor"

	"go.link-society.com/flowg/internal/engines/pipelines"
	storage "go.link-society.com/flowg/internal/storage/interfaces"
)

// worker is the body of the syslog service's actors. Several of them drain the
// message channel, each driving one message at a time through the pipeline
// engine.
type worker struct {
	logger         *slog.Logger
	channel        gosyslog.LogPartsChannel
	configStorage  storage.ConfigStorage
	pipelineRunner pipelines.Runner
}

var _ actor.Worker = (*worker)(nil)

// DoWork handles one received syslog message: it enforces the configured
// allowed-origin list (exact IPs or CIDR ranges), then runs the message through
// every running pipeline's syslog entrypoint, one after the other, logging but
// not aborting on per-pipeline errors.
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

		if len(systemConfig.SyslogAllowedOrigins) > 0 {
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

		for _, pipelineName := range w.pipelineRunner.RunningPipelines() {
			record := parseLogParts(logParts)

			err := w.pipelineRunner.Process(
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
		}

		return actor.WorkerContinue
	}
}
