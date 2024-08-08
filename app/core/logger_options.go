package core

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - Log Options -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type logOpt func(*zap.Logger)

// Prepares a new child Logger, with fields defined by the given logOpts.
func prepareLog(opts ...logOpt) *zap.Logger {
	childLogger := *zap.L()
	for _, opt := range opts {
		opt(&childLogger)
	}
	return &childLogger
}

// Shorthand versions of each logOpt.
var (
	wMsg   = withMsg
	wErr   = withError
	wDurat = withDuration
	wStack = withStacktrace
	wGRPC  = withGRPC
	wHTTP  = withHTTP
)

// Logs a simple message.
var withMsg = func(msg string) logOpt {
	return func(logger *zap.Logger) {
		*logger = *logger.With(zap.String("msg", msg))
	}
}

// Logs any kind of info.
var withData = func(data ...any) logOpt {
	return func(logger *zap.Logger) {
		if len(data) == 0 {
			return
		}
		*logger = *logger.With(zap.Any("data", data))
	}
}

// Logs a duration.
var withDuration = func(duration time.Duration) logOpt {
	return func(logger *zap.Logger) {
		*logger = *logger.With(zap.Duration("duration", duration))
	}
}

// Log error if not nil.
var withError = func(err error) logOpt {
	return func(logger *zap.Logger) {
		if err == nil {
			return
		}
		*logger = *logger.With(zap.Error(err))
	}
}

// Used to log where in the code a message comes from.
var withStacktrace = func() logOpt {
	return func(logger *zap.Logger) {
		*logger = *logger.With(zap.Stack("trace"))
	}
}

// Routes apply to both GRPC and HTTP.
//
//	-> In GRPC, it's the last part of the Method -> '/users.UsersService/GetUsers'.
//
// See routes.go for more info.
var withGRPC = func(method string) logOpt {
	return func(logger *zap.Logger) {
		*logger = *logger.With(zap.String("route", RouteNameFromGRPC(method)))
	}
}

// -> In HTTP, for the route we join Method and Path -> 'GET /users'.
var withHTTP = func(req *http.Request) logOpt {
	return func(logger *zap.Logger) {
		*logger = *logger.With(zap.String("route", req.Method+" "+req.URL.Path))
	}
}
