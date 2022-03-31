package haro

import "time"

type config struct {
	retry      int
	retryDelay time.Duration
	onError    func(error)
	onSuccess  func()
}

// ConfigFunc is signature of function for doing configuration
type ConfigFunc func(*config)

// Retry configures how many times should the subscriber retries before really failing
func Retry(count int) ConfigFunc {
	return func(cfg *config) {
		cfg.retry = count
	}
}

// DelayRetry configures delay between each retry
func DelayRetry(delay time.Duration) ConfigFunc {
	return func(cfg *config) {
		cfg.retryDelay = delay
	}
}

// OnError specify what to do when error occured
func OnError(onError func(error)) ConfigFunc {
	return func(cfg *config) {
		cfg.onError = onError
	}
}

// OnSuccess specify what to do when no error occured
func OnSuccess(onSuccess func()) ConfigFunc {
	return func(cfg *config) {
		cfg.onSuccess = onSuccess
	}
}
