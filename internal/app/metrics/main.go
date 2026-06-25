package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	streamLogCounter   *prometheus.CounterVec
	pipelineLogCounter *prometheus.CounterVec
)

// Setup creates the FlowG metrics and registers them with the default
// Prometheus registry. It must be called once during startup before any of the
// Inc* helpers are used.
func Setup() {
	streamLogCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "flowg_stream_log_total",
			Help: "Total number of log messages ingested in a stream",
		},
		[]string{"stream"},
	)
	pipelineLogCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "flowg_pipeline_log_total",
			Help: "Total number of log messages processed in a pipeline",
		},
		[]string{"pipeline", "status"},
	)

	prometheus.MustRegister(
		streamLogCounter,
		pipelineLogCounter,
	)
}

// IncStreamLogCounter records that one log record was ingested into the given
// stream.
func IncStreamLogCounter(stream string) {
	streamLogCounter.WithLabelValues(stream).Inc()
}

// IncPipelineLogCounter records that one log record was processed by the given
// pipeline, labelling the sample as a success or an error.
func IncPipelineLogCounter(pipeline string, success bool) {
	var status string
	if success {
		status = "success"
	} else {
		status = "error"
	}

	pipelineLogCounter.WithLabelValues(pipeline, status).Inc()
}
