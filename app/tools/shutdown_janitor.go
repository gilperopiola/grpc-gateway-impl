package tools

import "github.com/gilperopiola/grpc-gateway-impl/app/core"

// This was previously named CleanupTool, but
// ShutdownJanitor has a bit more style.

// Used to release all resources on server shutdown.
// Basically, on init we add here all the de-init methods
// of everything we create.
type shutdownJanitor struct {
	cleanupFns []func()
}

func NewShutdownJanitor(cleanupFns ...func()) core.ShutdownJanitor {
	return &shutdownJanitor{
		cleanupFns,
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (sj shutdownJanitor) Cleanup() {
	for _, cleanupFn := range sj.cleanupFns {
		cleanupFn()
	}
}

func (sj *shutdownJanitor) AddCleanupFunc(fn func()) {
	sj.cleanupFns = append(sj.cleanupFns, fn)
}

func (sj *shutdownJanitor) AddCleanupFuncWithErr(fn func() error) {
	// We use an adapter here that acts like a closure: it captures the original function and once called
	// just calls it and logs the error if there was one.
	adapter := func() {
		core.LogIfErr(fn(), "ungraceful shutdown: %v")
	}

	sj.cleanupFns = append(sj.cleanupFns, adapter)
}

var _ core.ShutdownJanitor = &shutdownJanitor{}
