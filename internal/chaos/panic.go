package chaos

import (
	"sync"
	"sync/atomic"
)

var (
	panicTrigger atomic.Bool
	panicMu      sync.Mutex
	panicCount   int
)

func SetPanicTrigger(trigger bool) {
	panicTrigger.Store(trigger)
}

func ShouldPanic() bool {
	return panicTrigger.Swap(false)
}

func TriggerPanicIfSet() {
	if ShouldPanic() {
		incPanicCount()
		panic("chaos: simulated panic")
	}
}

func incPanicCount() {
	panicMu.Lock()
	panicCount++
	panicMu.Unlock()
}

func PanicCount() int {
	panicMu.Lock()
	defer panicMu.Unlock()
	return panicCount
}
