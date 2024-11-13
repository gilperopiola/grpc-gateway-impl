package tools

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/logs"
)

// This was previously named CleanupTool, but
// ShutdownJanitor has a bit more style.
type shutdownJanitor struct {
	cleanupFns []func()
}

// Used to release all used resources on server shutdown.
// Basically, on init we add here all the de-init methods
// of everything we create.
func NewShutdownJanitor(cleanupFns ...func()) core.ShutdownJanitor {
	return &shutdownJanitor{cleanupFns}
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

	// This allows us to handle the error and still
	// be able to add the function to the cleanup list.
	adapter := func() {
		logs.LogIfErr(fn(), "ungraceful shutdown: %v")
	}

	sj.cleanupFns = append(sj.cleanupFns, adapter)
}

var _ core.ShutdownJanitor = &shutdownJanitor{}
