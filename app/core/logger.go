package core

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	gormLogger "gorm.io/gorm/logger"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*             - Logger -              */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// We use zap as our Logger. It's fast and easy to use.
// We don't even need to wrap it in a struct, we just use it globally on the zap pkg.

const LogsTimeLayout = "02/01/06 15:04:05"

// Replaces the global Logger in the zap package with a new one.
// It uses a default zap.Config and allows for additional options to be passed.
func SetupLogger(cfg *LoggerCfg) *zap.Logger {
	zapOpts := newZapBuildOpts(cfg.LevelStackT)

	zapLogger, err := newZapConfig(cfg).Build(zapOpts...)
	if err != nil {
		log.Fatalf(errs.FailedToCreateLogger, err) // Don't use zap for this.
	}

	zap.ReplaceGlobals(zapLogger)

	return zapLogger
}

// This func is a GRPC Interceptor. Or technically a grpc.UnaryServerInterceptor.
func LogGRPCRequest(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)

	if err != nil {
		zap.S().Errorw("GRPC Error", ZapRouteFromGRPC(info.FullMethod), ZapDuration(duration), ZapError(err))
	} else {
		zap.S().Infow("GRPC Request", ZapRouteFromGRPC(info.FullMethod), ZapDuration(duration))
	}

	return resp, err
}

// Why log both HTTP and GRPC? Because we can. And it's cool.
func LogHTTPRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()
		handler.ServeHTTP(rw, req)
		duration := time.Since(start)

		zap.S().Infow("HTTP Request", ZapRouteFromHTTP(req), ZapDuration(duration))

		// Most HTTP logs come with a GRPC log before, as HTTP acts as a gateway to GRPC.
		// As such, we add a new line to separate the logs and easily identify different requests.
		// The only exception would be if there was an error before calling the GRPC handlers.
		zap.L().Info("\n") // T0D0 is this ok?
	})
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Helps keeping code clean and readable, lets you omit the error check on the caller when you just need to log it.
func LogIfErr(err error, optionalFmtMsg ...string) {
	if err != nil {
		fmtMsg := "untyped error: %v"
		if len(optionalFmtMsg) > 0 {
			fmtMsg = optionalFmtMsg[0]
		}
		zap.S().Errorf(fmtMsg, err)
	}
}

// Helps keeping code clean and readable, lets you omit the error check on the caller when you just need to log it.
func WarnIfErr(err error, optionalFmtMsg ...string) {
	if err != nil {
		fmtMsg := "untyped warning: %v"
		if len(optionalFmtMsg) > 0 {
			fmtMsg = optionalFmtMsg[0]
		}
		zap.S().Warnf(fmtMsg, err)
	}
}

// Helps keeping code clean and readable, lets you omit the error check on the caller.
func LogPanicIfErr(err error, optionalFmtMsg ...string) {
	if err != nil {
		fmtMsg := "untyped panic: %v"
		if len(optionalFmtMsg) > 0 {
			fmtMsg = optionalFmtMsg[0]
		}
		LogUnexpectedAndPanic(fmt.Errorf(fmtMsg, err))
	}
}

// Used to log unexpected errors, like panic recoveries or some connection errors.
func LogUnexpectedErr(err error) {
	zap.S().Error("Unexpected", ZapError(err), ZapStacktrace())
}

// Used to log unexpected errors that also should trigger a panic.
func LogUnexpectedAndPanic(err error) {
	zap.S().Fatal("Unexpected Fatal", ZapError(err), ZapStacktrace())
}

// Used to log things that shouldn't happen, like someone trying to access admin endpoints.
func LogPotentialThreat(msg string) {
	zap.S().Error("Threat", ZapMsg(msg))
}

// Used to log strange behaviour that isn't necessarily bad or an error.
func LogWeirdBehaviour(msg string, info ...any) {
	zap.S().Warn("Weird", ZapMsg(msg), ZapInfo(info...))
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Logs a simple message.
func ZapMsg(msg string) zap.Field {
	return zap.String("msg", msg)
}

// Logs any kind of info.
func ZapInfo(info ...any) zap.Field {
	if len(info) == 0 {
		return zap.Skip()
	}
	return zap.Any("info", info)
}

// Routes apply to both GRPC and HTTP.
//
//	-> In GRPC, it's the last part of the Method -> '/users.UsersService/GetUsers'.
func ZapRouteFromGRPC(method string) zap.Field {
	return zap.String("route", RouteNameFromGRPC(method))
}

// -> In HTTP, we join Method and Path -> 'GET /users'.
func ZapRouteFromHTTP(req *http.Request) zap.Field {
	return zap.String("route", RouteNameFromHTTP(req))
}

// Logs a duration.
func ZapDuration(duration time.Duration) zap.Field {
	return zap.Duration("duration", duration)
}

// Log error if not nil.
func ZapError(err error) zap.Field {
	if err == nil {
		return zap.Skip()
	}
	return zap.Error(err)
}

// Used to log where in the code a message comes from.
func ZapStacktrace() zap.Field {
	return zap.Stack("stack")
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Only log messages with a level equal or higher than the one we set in the config.
var LogLevels = map[string]int{
	"debug":  int(zap.DebugLevel),
	"info":   int(zap.InfoLevel),
	"warn":   int(zap.WarnLevel),
	"error":  int(zap.ErrorLevel),
	"dpanic": int(zap.DPanicLevel),
	"panic":  int(zap.PanicLevel),
	"fatal":  int(zap.FatalLevel),
}

// The selected DB Log Level will be used to log all SQL queries.
// 'silent' disables all logs, 'error' will only log errors, 'warn' logs errors and warnings, and 'info' logs everything.
var DBLogLevels = map[string]int{
	"silent": int(gormLogger.Silent),
	"error":  int(gormLogger.Error),
	"warn":   int(gormLogger.Warn),
	"info":   int(gormLogger.Info),
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Returns the default options for creating the zap Logger.
func newZapBuildOpts(levelStackT int) []zap.Option {
	return []zap.Option{
		zap.AddStacktrace(zapcore.Level(levelStackT)),
		zap.WithClock(zapcore.DefaultClock),
	}
}

// Returns a new zap.Config with the default options + *LoggerCfg settings.
func newZapConfig(cfg *LoggerCfg) zap.Config {
	zapCfg := newZapBaseConfig()

	zapCfg.Level = zap.NewAtomicLevelAt(zapcore.Level(cfg.Level))
	zapCfg.DisableCaller = !cfg.LogCaller

	zapCfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(LogsTimeLayout))
	}

	zapCfg.EncoderConfig.EncodeDuration = func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(d.Truncate(time.Millisecond).String())
	}

	return zapCfg
}

// Returns the default zap.Config for the current environment.
func newZapBaseConfig() zap.Config {
	if EnvIsProd {
		return zap.NewProductionConfig()
	}
	return zap.NewDevelopmentConfig()
}
