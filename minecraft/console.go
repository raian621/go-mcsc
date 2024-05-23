package minecraft

import (
	"io"
	"sync"
)

type ConsoleMonitor struct {
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser
	mutex  sync.Mutex
}

func (cm *ConsoleMonitor) Lock()         { cm.mutex.Lock() }
func (cm *ConsoleMonitor) Unlock()       { cm.mutex.Unlock() }
func (cm *ConsoleMonitor) TryLock() bool { return cm.mutex.TryLock() }
