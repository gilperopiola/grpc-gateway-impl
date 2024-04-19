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

// Origin is actually WHERE the log message is coming from. i.e. 'Auth' or 'GetUsers'.
// T0D0 make an enum?
func LogUnexpected(err error, origin string) {
	zap.S().Warn("Unexpected Error", ZapError(err), ZapOrigin(origin))
}

func LogPotentialThreat(msg string) {
	zap.S().Error("Potential Threat", ZapInfo(msg))
}

// Used to log strange behaviours that aren't necessarily bad. i.e. 'Route not found'.
func LogWeirdBehaviour(msg string) {
	zap.S().Error("Weird Behaviour", ZapInfo(msg))
}

func ZapError(err error) zap.Field {
	if err == nil {
		return zap.Skip()
	}
	return zap.Error(err)
}

// ZapEndpoint unifies both HTTP and gRPC paths:
//
//	-> In gRPC, it's only the Method 	-> '/users.UsersService/GetUsers'.
//	-> In HTTP, we join Method and Path -> 'GET /users'.
func ZapEndpoint(value string) zap.Field {
	return zap.String("endpoint", value)
}

// Used to log where the log message is coming from. i.e. 'Auth' or 'GetUsers'.
func ZapOrigin(value string) zap.Field {
	return zap.String("origin", value)
}

func ZapInfo(value string) zap.Field {
	return zap.String("message", value)
}

func ZapDuration(value time.Duration) zap.Field {
	return zap.Duration("duration", value)
}

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

// Returns a new zap.Config with the default options.
func newZapConfig(cfg *Config) zap.Config {
	newZapConfigFn := zap.NewDevelopmentConfig
	if IsProd {
		newZapConfigFn = zap.NewProductionConfig
	}

	zapConfig := newZapConfigFn()
	zapConfig.DisableCaller = !cfg.LoggerCfg.LogCaller
	zapConfig.Level = zap.NewAtomicLevelAt(zapcore.Level(cfg.LoggerCfg.Level))
	zapConfig.EncoderConfig.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(t.Format(LogTimeLayout))
	}

	return zapConfig
}
