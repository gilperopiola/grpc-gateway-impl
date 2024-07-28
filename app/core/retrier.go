package core

import (
	"fmt"
	"math"
	"time"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - Global Retrier -         */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// This function just calls FallbackAndRetryOnFailure with a nil fallbackFn.
// As it has a different signature for the 1st parameter, we need an adapter.
func Retry(fn func() error, nTries int) error {
	adapterFn := func() (any, error) {
		return nil, fn()
	}

	_, err := FallbackAndRetry(adapterFn, func() {}, nTries)
	return err
}

// Retries up to nTries, with exponential backoff.
func FallbackAndRetry(fn func() (any, error), fallbackFn func(), nTries int) (any, error) {
	var result any
	var err error

	for nTry := 1; nTry <= nTries; nTry++ {
		if result, err = fn(); err == nil {
			return result, nil
		}

		LogUnexpected(fmt.Errorf("try %d (of %d) failed: %v", nTry, nTries, err))

		if nTry == nTries { // Don't fallback or sleep on the last try
			break
		}

		fallbackFn()
		sleepFor := math.Pow(2, float64(nTry)) // 2, 4, 8, 16, 32...
		time.Sleep(time.Second * time.Duration(sleepFor))
	}

	return result, err
}
