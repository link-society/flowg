package pipelines

import "link-society.com/flowg/internal/models"

const (
	DIRECT_ENTRYPOINT = "direct"
	SYSLOG_ENTRYPOINT = "syslog"
	OTLP_ENTRYPOINT   = "otlp"
)

type message struct {
	replyTo chan<- error

	pipelineName string
	entrypoint   string
	record       *models.LogRecord
}
