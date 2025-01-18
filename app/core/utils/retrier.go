package utils

import (
	"errors"
	"math"
	"time"

	"go.uber.org/zap"
)

type Retryable func() (any, error)
type RetryableNoErr func() any

const maxAllowedRetries = 7

// Executes a given function. If it fails, it logs the error, optionally calls a fallback,
// waits with exponential backoff, and retries the initial operation.
func TryAndRetry(fn Retryable, maxRetries int, skipLog bool, onFailure func()) (any, error) {

	// Try
	out, err := fn()
	if err == nil {
		return out, nil
	}

	// It failed
	if !skipLog {
		zap.S().Errorf("retryable operation failed: %v\n", err)
		zap.S().Errorf("retrier kicking in...\n")
	}

	if maxRetries <= 0 {
		maxRetries = 1
	}

	if maxRetries > maxAllowedRetries {
		maxRetries = maxAllowedRetries
	}

	for nRetry := 0; nRetry <= maxRetries; nRetry++ {
		if onFailure != nil {
			onFailure()
		}

		sleepFor := math.Pow(2, float64(nRetry+2)) // 4, 8, 16, 32...
		time.Sleep(time.Second * time.Duration(sleepFor))

		if out, err = fn(); err == nil {
			zap.S().Errorf("retry [%d of %d] succeeded\n", nRetry+1, maxRetries)
			return out, nil
		}

		if !skipLog {
			zap.S().Errorf("retry [%d of %d] failed: %v\n", nRetry+1, maxRetries, err)
		}
	}

	return out, err
}

// Same as RetryFunc but for functions that return a non-error result; nil result is treated as error.
func TryAndRetryNoErr(fn RetryableNoErr, maxRetries int, skipLog bool, onFailure func()) (any, error) {
	return TryAndRetry(func() (any, error) {
		if got := fn(); got != nil {
			return got, nil
		}
		return nil, errors.New("function call returned nil")
	}, maxRetries, skipLog, onFailure)
}
