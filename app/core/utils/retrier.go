package utils

import (
	"fmt"
	"math"
	"time"

	"go.uber.org/zap"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - Global Retrier -         */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Executes a function, can be kinda configured.
func Retry(fn func() (any, error), cfg retryCfg) (any, error) {

	out, err := fn()
	if err == nil {
		return out, nil
	}

	// Operation failed.
	for nRetry := 0; nRetry <= cfg.maxRetries; nRetry++ {

		if cfg.logFailures {
			if nRetry == 0 {
				zap.L().Error(fmt.Sprintf("operation failed, will retry up to %d times: %v", cfg.maxRetries, err))
			} else {
				zap.L().Error(fmt.Sprintf("retry %d of %d failed: %v", nRetry, cfg.maxRetries, err))
			}
		}

		// Fallback.
		if cfg.fallbackFn != nil {
			cfg.fallbackFn()
		}

		// Wait.
		sleepSeconds := math.Pow(2, float64(nRetry+2)) // 4, 8, 16, 32...
		time.Sleep(time.Second * time.Duration(sleepSeconds))

		// Retry.
		// If it succeeds, return directly.
		if out, err = fn(); err == nil {
			return out, nil
		}
	}

	return out, err
}

func RetryV2(fn func() any, cfg retryCfg) (any, error) {
	return Retry(func() (any, error) {
		return fn(), nil
	}, cfg)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Used to configure the retrier's behavior.
type retryCfg struct {

	// The number of times to retry the operation.
	maxRetries int

	// Whether to log the failures.
	logFailures bool

	// Function to call if the operation fails.
	fallbackFn func()
}

func BasicRetryCfg(maxRetries int, fallbackFn func()) retryCfg {
	return retryCfg{
		maxRetries:  maxRetries,
		fallbackFn:  fallbackFn,
		logFailures: true,
	}
}
