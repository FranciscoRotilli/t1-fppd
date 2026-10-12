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
		for s := 1; s <= id*6; s++ {
			fmt.Printf(" ")
		}
		fmt.Printf(" T%d\n", i) // Thinking
		time.Sleep(config.ThinkTime)

		t0 := time.Now()

		leftFork := id
		rightFork := (id + 1) % config.Philosophers

		forks[leftFork].Lock()
		time.Sleep(config.TimeBetweenForks)
		forks[rightFork].Lock()

		m.RecordWait(time.Since(t0))

		// eat

		for s := 1; s <= id*6; s++ {
			fmt.Printf(" ")
		}

		m.StartEating()

		fmt.Printf(" E%d\n", i) // Eating
		time.Sleep(config.EatTime)

		m.StopEating()

		forks[rightFork].Unlock()
		forks[leftFork].Unlock()

		m.CycleDone(id) // Record cycle completion
	}
}

func main() {
	forks := make([]sync.Mutex, config.Philosophers)

	fmt.Println("\n[P1]  [P2]  [P3]  [P4]  [P5]")

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
