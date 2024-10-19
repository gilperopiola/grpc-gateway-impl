package utils

import (
	"errors"
	"math"
	"time"

	"go.uber.org/zap"
)

// 🔻 Global Retrier 🔻

type RetryableFn func() (any, error)
type RetryableFnNoErr func() any

// Executes a given function. If it fails, it logs the error, falls back to another function,
// waits with exponential backoff and then retries the initial operation again.
func RetryFunc(fn RetryableFn, optionalCfg ...retryCfg) (any, error) {

	// Try.
	out, err := fn()
	if err == nil {
		return out, nil
	}

	// Operation failed.
	// Create a config to handle the retry behavior.
	cfg := DefaultRetryCfg()
	if len(optionalCfg) > 0 {
		cfg = overrideRetryCfg(cfg, optionalCfg[0])
	}

	if !cfg.skipLog {
		zap.S().Errorf("operation failed, will retry up to %d times: %v", cfg.retries, err)
	}

	for nRetry := 0; nRetry <= cfg.retries; nRetry++ {

		// Fallback.
		if cfg.fallbackFn != nil {
			cfg.fallbackFn()
		}

		// Wait.
		sleepFor := math.Pow(2, float64(nRetry+2)) // 4, 8, 16, 32...
		time.Sleep(time.Second * time.Duration(sleepFor))

		// Retry.
		if out, err = fn(); err == nil {
			// If it succeeds, return directly.
			return out, nil
		}

		if !cfg.skipLog {
			zap.S().Errorf("retry %d of %d failed: %v", nRetry+1, cfg.retries, err)
		}
	}

	return out, err
}

// This version of Retry is for functions that don't return an error with the result.
// If the result gotten is nil, it will be treated as an error and will be retried.
//
// This is just a wrapper around RetryFunc.
func RetryFuncNoError(fn RetryableFnNoErr, optionalCfg ...retryCfg) (any, error) {
	return RetryFunc(func() (any, error) {
		if got := fn(); got != nil {
			return got, nil
		}
		return nil, errors.New("function call returned nil")
	}, optionalCfg...)
}
