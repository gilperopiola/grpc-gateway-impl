package tools

import "github.com/gilperopiola/grpc-gateway-impl/app/core"

var _ core.ShutdownJanitor = (*shutdownJanitor)(nil)

type shutdownJanitor struct {
	cleanupFns        []func()
	cleanupFnsWithErr []func() error
}

func NewShutdownJanitor(cleanupFns ...func()) core.ShutdownJanitor {
	return &shutdownJanitor{cleanupFns, []func() error{}}
}

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
