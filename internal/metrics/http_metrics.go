package metrics

import (
	"net/http"
	"strconv"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func HTTPMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		recorder := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		elapsed := time.Since(start).Seconds()

		HTTPRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(elapsed)
		HTTPRequestsTotal.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(recorder.status)).Inc()
	})
}
