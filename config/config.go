package config

import "time"

const (
	ThinkTime        = 20 * time.Millisecond
	EatTime          = 30 * time.Millisecond
	Philosophers     = 5
	Rounds           = 20
	TimeBetweenForks = 1 * time.Millisecond // Set to > 1ms to force deadlock on v1
)
