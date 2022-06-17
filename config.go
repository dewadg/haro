package haro

import "time"

type subscriberConfig struct {
	retry      int
	retryDelay time.Duration
	onError    func(error)
	onSuccess  func()
}

// SubscriberConfigFunc is signature of function for doing configuration
type SubscriberConfigFunc func(*subscriberConfig)

// Retry configures how many times should the subscriber retries before really failing
func Retry(count int) SubscriberConfigFunc {
	return func(cfg *subscriberConfig) {
		cfg.retry = count
	}
}

// DelayRetry configures delay between each retry
func DelayRetry(delay time.Duration) SubscriberConfigFunc {
	return func(cfg *subscriberConfig) {
		cfg.retryDelay = delay
	}
}

// OnError specify what to do when error occured
func OnError(onError func(error)) SubscriberConfigFunc {
	return func(cfg *subscriberConfig) {
		cfg.onError = onError
	}
}

// OnSuccess specify what to do when no error occured
func OnSuccess(onSuccess func()) SubscriberConfigFunc {
	return func(cfg *subscriberConfig) {
		cfg.onSuccess = onSuccess
	}
}
