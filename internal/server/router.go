package server

import (
	"net/http"

	"github.com/rjsej12/sentinel-go/internal/health"
	"github.com/rjsej12/sentinel-go/internal/metrics"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/metrics", metrics.Handler())

	mux.HandleFunc("/api/ping", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("pong"))
	})

	mux.HandleFunc("/healthz", health.Liveness)
	mux.HandleFunc("/readyz", health.Readiness)

	return mux
}
