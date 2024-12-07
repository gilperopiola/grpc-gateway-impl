package utils

// 🔽 Retry Config 🔽

const (
	defaultRetries = 5
)

// Used to configure the retrier's behavior.
type RetryCfg struct {
	// Times to retry the operation.
	Times int

	// Skip logging the failures.
	SkipLog bool

	// Call after the operation fails.
	OnFailure func()
}

func NewRetryCfg(retries int, skipLog bool, fallbackFn func()) RetryCfg {
	return RetryCfg{retries, skipLog, fallbackFn}
}

// This is used unless overriden by a provided config.
func DefaultRetryCfg() RetryCfg {
	return RetryCfg{
		Times:     defaultRetries,
		SkipLog:   false,
		OnFailure: nil,
	}
}

// If the user provides a retry config, we override the default values.
func overrideRetryCfg(defaultCfg RetryCfg, providedCfg RetryCfg) RetryCfg {
	if providedCfg.Times > 0 {
		defaultCfg.Times = providedCfg.Times
	}
	defaultCfg.SkipLog = providedCfg.SkipLog     // Defaults to false
	defaultCfg.OnFailure = providedCfg.OnFailure // Defaults to nil
	return defaultCfg
}
