package core

import (
	"log"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gormLogger "gorm.io/gorm/logger"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*             - Logger -              */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// We use zap as our Logger. It's fast and easy to use.
// We don't even need to wrap it in a struct, we just use it globally on the zap pkg.

const LogTimeLayout = "02/01/06 15:04:05"

// Replaces the global Logger in the zap package with a new one.
// It uses a default zap.Config and allows for additional options to be passed.
func SetupLogger(c *Config, opts ...zap.Option) *zap.Logger {
	zapLogger, err := newZapConfig(c).Build(opts...)
	if err != nil {
		log.Fatalf(errs.FatalErrMsgCreatingLogger, err)
	}

	zap.ReplaceGlobals(zapLogger)

	return zapLogger
}

// Returns the default options for the Logger.
func SetupLoggerOptions(stackTraceLevel int) []zap.Option {
	return []zap.Option{
		zap.AddStacktrace(zapcore.Level(stackTraceLevel)),
		zap.WithClock(zapcore.DefaultClock),
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Used to log things that shouldn't happen, like someone trying to access admin endpoints.
func LogPotentialThreat(msg string) {
	zap.S().Error("Potential Threat", ZapMsg(msg), ZapStacktrace())
}

// Used to log unexpected errors, like panic recoveries or some connection errors.
func LogUnexpected(err error) {
	zap.S().Warn("Unexpected Error", ZapError(err), ZapStacktrace())
}

// Used to log unexpected errors that also should trigger a panic.
func LogUnexpectedAndPanic(err error) {
	zap.S().Fatal("Unexpected Error: Fatal", ZapError(err), ZapStacktrace())
}

// Used to log strange behaviour that isn't necessarily bad or an error.
func LogWeirdBehaviour(msg string, info ...any) {
	zap.S().Error("Weird Behaviour", ZapMsg(msg), ZapStacktrace(), ZapInfo(info...))
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// This unifies both HTTP and gRPC paths:
//
//	-> In gRPC, it's only the Method 	-> '/users.UsersService/GetUsers'.
//	-> In HTTP, we join Method and Path -> 'GET /users'.
func ZapRoute(route string) zap.Field {
	return zap.String("route", route)
}

// Logs a simple message.
func ZapMsg(info string) zap.Field {
	return zap.String("info", info)
}

// Logs a duration.
func ZapDuration(duration time.Duration) zap.Field {
	return zap.Duration("duration", duration)
}

// Logs any kind of info.
func ZapInfo(info ...any) zap.Field {
	if len(info) == 0 {
		return zap.Skip()
	}
	return zap.Any("info", info)
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

// Returns a new zap.Config with the default options.
func newZapConfig(cfg *Config) zap.Config {
	zapCfg := getDefaultZapConfig(IsProd)

	zapCfg.Level = zap.NewAtomicLevelAt(zapcore.Level(cfg.LoggerCfg.Level))
	zapCfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(LogTimeLayout))
	}
	zapCfg.DisableCaller = !cfg.LoggerCfg.LogCaller

	return zapCfg
}

// Returns the default zap.Config for the current environment.
func getDefaultZapConfig(isProd bool) zap.Config {
	if isProd {
		return zap.NewProductionConfig()
	}
	return zap.NewDevelopmentConfig()
}
