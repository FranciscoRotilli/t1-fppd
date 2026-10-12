package config

import "time"

const (
	ThinkTime        = 20 * time.Millisecond
	EatTime          = 30 * time.Millisecond
	Philosophers     = 5
	Rounds           = 20
	TimeBetweenForks = 1 * time.Millisecond
)
