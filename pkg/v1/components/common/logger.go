package common

import (
	"log"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/* ----------------------------------- */
/*             - Logger -              */
/* ----------------------------------- */

// We use zap as our logger. It's fast and has a nice API.
// We don't even need to wrap it in a struct, we just use it globally on the zap pkg.

// InitGlobalLogger replaces the global logger in the zap package with a new one.
// It uses a default zap.Config and allows for additional options to be passed.
func InitGlobalLogger(isProd bool, opts ...zap.Option) {
	zapLogger, err := getZapConfig(isProd).Build(opts...)
	if err != nil {
		log.Fatalf(errs.FatalErrMsgCreatingLogger, err)
	}

	zap.ReplaceGlobals(zapLogger)
}

// NewLoggerOptions returns the default options for the Logger.
func NewLoggerOptions() []zap.Option {
	return []zap.Option{
		zap.AddStacktrace(zap.DPanicLevel),  // T0D0 -> Config var.
		zap.WithClock(zapcore.DefaultClock), // T0D0 -> Config var.
	}
}

// getZapConfig returns a new zap.Config with the default options.
func getZapConfig(isProd bool) zap.Config {
	zapConfigFunc := zap.NewDevelopmentConfig
	if isProd {
		zapConfigFunc = zap.NewProductionConfig
	}

	zapConfig := zapConfigFunc()
	zapConfig.DisableCaller = true                         // T0D0 -> Config var.
	zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel) // T0D0 -> Config var.
	zapConfig.EncoderConfig.EncodeTime = myTimeEncoder     // T0D0 -> Config var.

	return zapConfig
}

// ZapEndpoint unifies both HTTP and gRPC paths:
//
//   - In HTTP, we join Method and Path -> 'GET /users'.
//   - In gRPC, it's only the Method 		-> '/users.UsersService/GetUsers'.
//
// -
func ZapEndpoint(value string) zap.Field {
	return zap.String("endpoint", value)
}

func ZapDuration(value time.Duration) zap.Field {
	return zap.Duration("duration", value)
}

func ZapError(err error) zap.Field {
	return zap.Error(err)
}

func encodeTimeLayout(t time.Time, layout string, enc zapcore.PrimitiveArrayEncoder) {
	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}

	if enc, ok := enc.(appendTimeEncoder); ok {
		enc.AppendTimeLayout(t, layout)
		return
	}

	enc.AppendString(t.Format(layout))
}

func myTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	encodeTimeLayout(t, "02/01/06 15:04:05", enc)
}
