package chaos

import (
	"sync"
	"time"
)

var (
	latencyMu     sync.RWMutex
	latencyMs     int
	latencyActive bool
)

func SetLatency(ms int) {
	latencyMu.Lock()
	defer latencyMu.Unlock()

	if ms <= 0 {
		latencyActive = false
		latencyMs = 0
		return
	}

	latencyMs = ms
	latencyActive = true
}

func Latency() (ms int, active bool) {
	latencyMu.RLock()
	defer latencyMu.RUnlock()
	return latencyMs, latencyActive
}

func ApplyLatency() {
	latencyMu.RLock()
	ms := latencyMs
	active := latencyActive
	latencyMu.RUnlock()

	if !active || ms <= 0 {
		return
	}
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
