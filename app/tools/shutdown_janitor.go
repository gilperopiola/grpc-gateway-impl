package tools

import "github.com/gilperopiola/grpc-gateway-impl/app/core"

type shutdownJanitor struct {
	cleanupFns        []func()
	cleanupFnsWithErr []func() error
}

func NewShutdownJanitor(cleanupFns ...func()) core.ShutdownJanitor {
	cleanupFnsWithErr := []func() error{}

	return &shutdownJanitor{
		cleanupFns,
		cleanupFnsWithErr,
	}
}

var _ core.ShutdownJanitor = (*shutdownJanitor)(nil)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (sj *shutdownJanitor) AddCleanupFunc(fn func()) {
	sj.cleanupFns = append(sj.cleanupFns, fn)
}

func (sj *shutdownJanitor) AddCleanupFuncWithErr(fnWithErr func() error) {
	sj.cleanupFnsWithErr = append(sj.cleanupFnsWithErr, fnWithErr)
}

func (sj shutdownJanitor) Cleanup() {
	for _, cleanupFn := range sj.cleanupFns {
		cleanupFn()
	}
	for _, cleanupFnWithErr := range sj.cleanupFnsWithErr {
		core.LogIfErr(cleanupFnWithErr(), "ungraceful shutdown :(")
	}
}
