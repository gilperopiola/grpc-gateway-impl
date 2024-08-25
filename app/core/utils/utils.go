package utils

import (
	"errors"
	"io/fs"
	"strings"

	"github.com/gilperopiola/god"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Oh no, a utils package! We're all gonna die!

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*             - Utils -               */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

/* -~-~-~-~-~ General Utils -~-~-~-~-~- */

// Returns the first element of a string slice, or a fallback if the slice is empty.
func FirstOrDefault(slice []string, fallback string) string {
	if len(slice) > 0 {
		return slice[0]
	}
	return fallback
}

/* -~-~-~-~-~ Route Utils -~-~-~-~-~- */

// Our Routes are named by the last part of their GRPC Method.
// It's everything after the last slash.
//
//	Method = /pbs.Service/Signup
//	Route  = Signup
func RouteNameFromGRPC(method string) string {
	i := strings.LastIndex(method, "/")
	if i == -1 {
		return ""
	}
	return method[i+1:]
}

// Returns the route name from the context's data.
func RouteNameFromCtx(ctx god.Ctx) string {
	if method, ok := grpc.Method(ctx); ok {
		return RouteNameFromGRPC(method)
	}
	return ""
}

/* -~-~-~-~-~ Logger Utils -~-~-~-~-~- */

// On Windows, I'm getting a *fs.PathError when calling zap.L().Sync() to flush
// the Logger on shutdown.
// This just calls zap.L().Sync() and ignores that specific error.
// See https://github.com/uber-go/zap/issues/991
func SyncLogger() error {
	var pathErr *fs.PathError
	if err := zap.L().Sync(); err != nil && !errors.As(err, &pathErr) {
		return err
	}
	return nil
}
