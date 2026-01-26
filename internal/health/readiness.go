package health

import (
	"net/http"
	"sync/atomic"
)

var ready atomic.Bool

func SetReady(value bool) {
	ready.Store(value)
}

func Readiness(w http.ResponseWriter, _ *http.Request) {
	if ready.Load() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ready"))
		return
	}

	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte("not ready"))
}
