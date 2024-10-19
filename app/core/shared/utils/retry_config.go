package utils

// 🔽 Retry Config 🔽

const (
	defaultRetries = 5
)

// Used to configure the retrier's behavior.
type retryCfg struct {
	// Times to retry the operation.
	retries int

	// Skip logging the failures.
	skipLog bool

	// Call after the operation fails.
	fallbackFn func()
}

func NewRetryCfg(retries int, skipLog bool, fallbackFn func()) retryCfg {
	return retryCfg{retries, skipLog, fallbackFn}
}

// This is used unless overriden by a provided config.
func DefaultRetryCfg() retryCfg {
	return retryCfg{
		retries:    defaultRetries,
		skipLog:    false,
		fallbackFn: nil,
	}
}

// If the user provides a retry config, we override the default values.
func overrideRetryCfg(defaultCfg retryCfg, providedCfg retryCfg) retryCfg {
	if providedCfg.retries > 0 {
		defaultCfg.retries = providedCfg.retries
	}
	defaultCfg.skipLog = providedCfg.skipLog       // Defaults to false
	defaultCfg.fallbackFn = providedCfg.fallbackFn // Defaults to nil
	return defaultCfg
}
