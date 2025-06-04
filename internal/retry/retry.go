package retry

import (
	"context"
	"time"
)

type retryOption func(*retryConfig)

type retryConfig struct {
	maxRetries int
	backOff    time.Duration
}

func defaultConfig() *retryConfig {
	return &retryConfig{
		maxRetries: 3,
		backOff:    time.Millisecond * 500,
	}
}

func withRetries(n int) retryOption {
	return func(cfg *retryConfig) {
		cfg.maxRetries = n
	}
}

func withBackOff(interval time.Duration) retryOption {
	return func(cfg *retryConfig) {
		cfg.backOff = interval
	}
}

func Retry(ctx context.Context, operation func() error, opts ...retryOption) error {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	var err error
	for attempt := 0; attempt <= cfg.maxRetries; attempt++ {
		err = operation()
		if err == nil {
			return nil
		}
		if attempt <= cfg.maxRetries {
			select {
			case <-time.After(cfg.backOff):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return err
}
