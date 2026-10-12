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
	maxWaitPer    []time.Duration
	eating        int
	maxEating     int
	cycles        []int
}

func New(n int) *Metrics {
	return &Metrics{
		cycles:     make([]int, n),
		maxWaitPer: make([]time.Duration, n),
	}
}

func (m *Metrics) RecordWait(id int, d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.totalWaitTime += d
	m.waitAmount++
	if d > m.maxWaitTime {
		m.maxWaitTime = d
	}

	if d > m.maxWaitPer[id] {
		m.maxWaitPer[id] = d
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
	fmt.Printf("Max Wait Per Philosopher: %v\n", m.maxWaitPer)
	fmt.Printf("Cycles per Philosopher: %v\n", m.cycles)
	fmt.Printf("Max eating simultaneously: %d\n", m.maxEating)
}

func (m *Metrics) StartEating() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.eating++
	if m.maxEating < m.eating {
		m.maxEating = m.eating
	}
}

func (m *Metrics) StopEating() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.eating--
}
