package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	WorkerQueueLength = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "sentinel",
			Subsystem: "worker",
			Name:      "queue_length",
			Help:      "Number of jobs waiting in the queue",
		},
	)

	WorkerQueueCapacity = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "sentinel",
			Subsystem: "worker",
			Name:      "queue_capacity",
			Help:      "Maximum capacity of the queue",
		},
	)

	WorkerJobsProcessedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "sentinel",
			Subsystem: "worker",
			Name:      "jobs_processed_total",
			Help:      "Total number of jobs processed",
		},
		[]string{"job_type", "status"},
	)

	WorkerJobDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "sentinel",
			Subsystem: "worker",
			Name:      "job_duration_seconds",
			Help:      "Time taken to process a job",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"job_type"},
	)

	WorkerJobsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "sentinel",
			Subsystem: "worker",
			Name:      "jobs_inflight",
			Help:      "Number of jobs currently being processed",
		},
	)

	WorkerErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "sentinel",
			Subsystem: "worker",
			Name:      "errors_total",
			Help:      "Total number of worker errors",
		},
		[]string{"error_type"},
	)
)

func RegisterWorkerMetrics() {
	prometheus.MustRegister(
		WorkerQueueLength,
		WorkerQueueCapacity,
		WorkerJobsProcessedTotal,
		WorkerJobDurationSeconds,
		WorkerJobsInFlight,
		WorkerErrorsTotal,
	)
}
