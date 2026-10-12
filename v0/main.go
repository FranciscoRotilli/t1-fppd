package main

// Example code - Sequential Dinning Philosophers
// By Cesar De Rose
//
// Adaptado para T1 FPPD, tempo em milissegundos, constantes globais.

import (
	"fmt"
	"t1-fppd/config"
	"t1-fppd/metrics"
	"time"
)

func Philosopher(id, iteration int, m *metrics.Metrics) {

	// think
	for s := 1; s <= id*6; s++ {
		fmt.Printf(" ")
	}
	fmt.Printf(" T%d\n", iteration) // Thinking
	time.Sleep(config.ThinkTime)

	t0 := time.Now()

	// acquire forks (empty on v0)

	m.RecordWait(time.Since(t0))

	// eat

	for s := 1; s <= id*6; s++ {
		fmt.Printf(" ")
	}
	fmt.Printf(" E%d\n", iteration) // Eating
	time.Sleep(config.EatTime)

	m.CycleDone(id) // Record cycle completion
}

func main() {
	fmt.Println("\n[P1]  [P2]  [P3]  [P4]  [P5]")

	// run philosophers

	start := time.Now()

	m := metrics.New(config.Philosophers)

	for i := 1; i <= config.Rounds; i++ {
		for j := 0; j < config.Philosophers; j++ {
			Philosopher(j, i, m)
		}
	}

	m.Print()

	elapsed := time.Since(start)

	fmt.Printf("\nSequential Dinner took %s\n\n", elapsed)
}
