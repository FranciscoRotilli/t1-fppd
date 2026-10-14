package main

import (
	"fmt"
	"sync"
	"t1-fppd/config"
	"t1-fppd/metrics"
	"time"
)

func Philosopher(id int, m *metrics.Metrics, forks []sync.Mutex) {
	for i := 1; i <= config.Rounds; i++ {
		// think
		fmt.Printf("%*s T%d\n", id*6, "", i) // Thinking
		time.Sleep(config.ThinkTime)

		t0 := time.Now()

		first, second := id, (id+1)%config.Philosophers
		if id%2 == 1 {
			first, second = second, first
		}
		forks[first].Lock()
		forks[second].Lock()

		m.RecordWait(id, time.Since(t0))

		// eat

		m.StartEating()

		fmt.Printf("%*s E%d\n", id*6, "", i) // Eating
		time.Sleep(config.EatTime)

		m.StopEating()

		forks[second].Unlock()
		forks[first].Unlock()

		m.CycleDone(id) // Record cycle completion
	}
}

func main() {
	forks := make([]sync.Mutex, config.Philosophers)

	fmt.Println("\n[P0]  [P1]  [P2]  [P3]  [P4]")

	// run philosophers

	start := time.Now()

	m := metrics.New(config.Philosophers)

	var wg sync.WaitGroup

	for id := 0; id < config.Philosophers; id++ {
		wg.Go(func() {
			Philosopher(id, m, forks)
		})
	}

	wg.Wait()

	elapsed := time.Since(start)

	fmt.Printf("\nConcurrent Dinner took %s\n\n", elapsed)
	m.Print()
}
