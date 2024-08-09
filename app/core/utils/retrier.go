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

// Executes a function.
// On failure, it calls a fallback (if set), sleeps some time, then retries.
func Retry(fn func() (any, error), nTries int, opts ...RetryOpt) (any, error) {
	var result any
	var err error

	cfg := newRetryCallCfg()

	// Apply options.
	for _, opt := range opts {
		opt(cfg)
	}

	for nTry := 1; nTry <= nTries; nTry++ {

		// Perform the operation.
		// If it succeeds, return directly.
		if result, err = fn(); err == nil {
			return result, nil
		}

		// Operation failed.
		if !cfg.dontLogFailures {
			zap.L().Error(fmt.Sprintf("try %d (of %d) failed: %v", nTry, nTries, err))
		}

		// Don't fallback or sleep on the last try.
		if nTry == nTries {
			break
		}

		cfg.fallbackFn()
		sleepSeconds := math.Pow(2, float64(nTry)) // 2, 4, 8, 16, 32...
		time.Sleep(time.Second * time.Duration(sleepSeconds))
	}

	return result, err
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Each call to Retry() creates an instance of this.
// It's used to configure the retry behavior via the RetryOpts.
type retryCallCfg struct {

	// Function to call if the operation fails.
	fallbackFn func()

	// The function used to log failures.
	logFn func(error)

	// Logs failures by default.
	// Set to true to disable logging.
	dontLogFailures bool
}

func newRetryCallCfg() *retryCallCfg {
	return &retryCallCfg{
		fallbackFn: func() {
			// Don't do anything.
		},
		logFn: func(err error) {
			zap.L().Error(err.Error())
		},
		dontLogFailures: false,
	}
}

type RetryOpt func(*retryCallCfg)

// DontLog disables logs for failed tries.
func DontLog() RetryOpt {
	return func(cfg *retryCallCfg) {
		cfg.dontLogFailures = true
	}
}

// Fallback sets the fallback function.
func Fallback(fallbackFn func()) RetryOpt {
	return func(cfg *retryCallCfg) {
		cfg.fallbackFn = fallbackFn
	}
}
