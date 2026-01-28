package server

import (
	"fmt"
	"net/http"

	"github.com/rjsej12/sentinel-go/internal/health"
	"github.com/rjsej12/sentinel-go/internal/metrics"
	"github.com/rjsej12/sentinel-go/internal/worker"
)

func NewRouter(queue *worker.Queue, processor *worker.Processor) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/metrics", metrics.Handler())

	mux.HandleFunc("/healthz", health.Liveness)
	mux.HandleFunc("/readyz", health.Readiness)

	mux.HandleFunc("/api/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("pong"))
	})

	mux.HandleFunc("/api/stats", func(w http.ResponseWriter, _ *http.Request) {
		stats := processor.Stats()
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{
			"queue_length": %d,
			"queue_capacity": %d,
			"worker_pool": %d
		}`, stats.QueueLength, stats.QueueCapacity, stats.WorkerPool)))
	})

	mux.HandleFunc("/api/jobs/process", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		jobType := r.URL.Query().Get("type")
		if jobType == "" {
			jobType = "process"
		}

		var jobTypeEnum worker.JobType
		switch jobType {
		case "process":
			jobTypeEnum = worker.JobTypeProcess
		case "delay":
			jobTypeEnum = worker.JobTypeDelay
		default:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid job type"))
			return
		}

		job := worker.NewJob(jobTypeEnum, []byte("test data"))
		if err := queue.Enqueue(job); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to enqueue job"))
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(fmt.Sprintf(`{"job_id": "%s", "status": "queued"}`, job.ID)))
	})

	return mux
}
