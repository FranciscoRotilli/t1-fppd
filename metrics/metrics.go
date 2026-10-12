package metrics

import (
	"fmt"
	"sync"
	"time"
)

type Metrics struct {
	mu            sync.Mutex
	totalWaitTime time.Duration
	waitAmount    int
	maxWaitTime   time.Duration
	maxEating     int
	cycles        []int
}

func New(n int) *Metrics {
	return &Metrics{
		cycles: make([]int, n),
	}
}

func (m *Metrics) RecordWait(d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.totalWaitTime += d
	m.waitAmount++
	if d > m.maxWaitTime {
		m.maxWaitTime = d
	}
}

func (m *Metrics) CycleDone(id int) {
	m.cycles[id]++
}

func (m *Metrics) Print() {
	m.mu.Lock()
	defer m.mu.Unlock()
	meanWaitTime := time.Duration(0)
	if m.waitAmount > 0 {
		meanWaitTime = m.totalWaitTime / time.Duration(m.waitAmount)
	}
	fmt.Printf("\nMetrics:\n")
	fmt.Printf("Total Wait Time: %s\n", m.totalWaitTime)
	fmt.Printf("Mean Wait Time: %s\n", meanWaitTime)
	fmt.Printf("Max Wait Time: %s\n", m.maxWaitTime)
	fmt.Printf("Cycles per Philosopher: %v\n", m.cycles)
}
