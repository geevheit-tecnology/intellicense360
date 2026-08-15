package events

import "time"

type RetryConfig struct {
	MaxAttempts int
	Delay       time.Duration
}

func DefaultRetryConfig() RetryConfig {
	return RetryConfig{MaxAttempts: 3, Delay: time.Minute}
}
