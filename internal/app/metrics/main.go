package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	streamLogCounter   *prometheus.CounterVec
	pipelineLogCounter *prometheus.CounterVec
)

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

func IncStreamLogCounter(stream string) {
	streamLogCounter.WithLabelValues(stream).Inc()
}

func IncPipelineLogCounter(pipeline string, success bool) {
	var status string
	if success {
		status = "success"
	} else {
		status = "error"
	}

	pipelineLogCounter.WithLabelValues(pipeline, status).Inc()
}
