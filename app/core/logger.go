package core

import (
	"log"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gormLogger "gorm.io/gorm/logger"
)

/* ----------------------------------- */
/*             - Logger -              */
/* ----------------------------------- */

// We use zap as our Logger. It's fast and has a nice API.
// We don't even need to wrap it in a struct, we just use it globally on the zap pkg.

// SetupLogger replaces the global Logger in the zap package with a new one.
// It uses a default zap.Config and allows for additional options to be passed.
func SetupLogger(c *Config, opts ...zap.Option) *zap.Logger {
	zapLogger, err := newZapConfig(c).Build(opts...)
	if err != nil {
		log.Fatalf(errs.FatalErrMsgCreatingLogger, err)
	}

	zap.ReplaceGlobals(zapLogger)

	return zapLogger
}

// NewLoggerOptions returns the default options for the Logger.
func NewLoggerOptions(stackTraceLevel int) []zap.Option {
	return []zap.Option{
		zap.AddStacktrace(zapcore.Level(stackTraceLevel)),
		zap.WithClock(zapcore.DefaultClock),
	}
}

// newZapConfig returns a new zap.Config with the default options.
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

const LogTimeLayout = "02/01/06 15:04:05"

var LogLevels = map[string]int{
	"debug":  int(zap.DebugLevel),
	"info":   int(zap.InfoLevel),
	"warn":   int(zap.WarnLevel),
	"error":  int(zap.ErrorLevel),
	"dpanic": int(zap.DPanicLevel),
	"panic":  int(zap.PanicLevel),
	"fatal":  int(zap.FatalLevel),
}

var DBLogLevels = map[string]int{
	"silent": int(gormLogger.Silent),
	"error":  int(gormLogger.Error),
	"warn":   int(gormLogger.Warn),
	"info":   int(gormLogger.Info),
}

// ZapEndpoint unifies both HTTP and gRPC paths:
//
//	-> In HTTP, we join Method and Path -> 'GET /users'.
//	-> In gRPC, it's only the Method 	-> '/users.UsersService/GetUsers'.
func ZapEndpoint(value string) zap.Field {
	return zap.String("endpoint", value)
}

func ZapDuration(value time.Duration) zap.Field {
	return zap.Duration("duration", value)
}

func ZapError(err error) zap.Field {
	return zap.Error(err)
}
