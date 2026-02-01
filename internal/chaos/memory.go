package chaos

import "sync"

var (
	memoryMu   sync.Mutex
	memoryHold []byte
)

func AllocateMemory(sizeMB int) {
	memoryMu.Lock()
	defer memoryMu.Unlock()

	memoryHold = nil

	if sizeMB <= 0 {
		return
	}

	size := sizeMB * 1024 * 1024
	memoryHold = make([]byte, size)
	for i := 0; i < size; i += 4096 {
		memoryHold[i] = 0
	}
}

func ReleaseMemory() {
	memoryMu.Lock()
	defer memoryMu.Unlock()
	memoryHold = nil
}

func MemoryBytesHeld() int {
	memoryMu.Lock()
	defer memoryMu.Unlock()
	return len(memoryHold)
}
