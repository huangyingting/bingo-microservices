package retry

import (
	"context"
	"time"
)

// Function signature of retryable function
type RetryableFunc func() error

func Do(retryableFunc RetryableFunc, opts ...Option) error {
	var n uint32

	// default
	config := newDefaultRetryConfig()

	// apply opts
	for _, opt := range opts {
		opt(config)
	}

	if err := config.context.Err(); err != nil {
		return err
	}

	var errorLog error

	// Setting attempts to 0 means we'll retry until we succeed
	if config.attempts == 0 {
		for err := retryableFunc(); err != nil; err = retryableFunc() {
			n++
			config.onRetry(n, err)

			select {
			case <-time.After(delay(config, n, err)):
			case <-config.context.Done():
				return config.context.Err()
			}
		}

		return nil
	}

	for n < config.attempts {
		err := retryableFunc()

		if err != nil {
			errorLog = unpackUnrecoverable(err)

			if !config.retryIf(err) {
				break
			}

			config.onRetry(n, err)

			// if this is last attempt - don't wait
			if n == config.attempts-1 {
				break
			}

			select {
			case <-time.After(delay(config, n, err)):
			case <-config.context.Done():
				return config.context.Err()
			}

		} else {
			return nil
		}

		n++
	}

	return errorLog
}

func newDefaultRetryConfig() *Config {
	return &Config{
		attempts:  uint32(10),
		delay:     100 * time.Millisecond,
		maxJitter: 100 * time.Millisecond,
		onRetry:   func(n uint32, err error) {},
		retryIf:   IsRecoverable,
		delayType: CombineDelay(BackOffDelay, RandomDelay),
		context:   context.Background(),
	}
}

type unrecoverableError struct {
	error
}

// Unrecoverable wraps an error in `unrecoverableError` struct
func Unrecoverable(err error) error {
	return unrecoverableError{err}
}

// IsRecoverable checks if error is an instance of `unrecoverableError`
func IsRecoverable(err error) bool {
	_, isUnrecoverable := err.(unrecoverableError)
	return !isUnrecoverable
}

func unpackUnrecoverable(err error) error {
	if unrecoverable, isUnrecoverable := err.(unrecoverableError); isUnrecoverable {
		return unrecoverable.error
	}

	return err
}

func delay(config *Config, n uint32, err error) time.Duration {
	delayTime := config.delayType(n, err, config)
	if config.maxDelay > 0 && delayTime > config.maxDelay {
		delayTime = config.maxDelay
	}

	return delayTime
}
