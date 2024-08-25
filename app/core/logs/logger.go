package logs

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// We use zap. It's fast and easy.
// Set it up and then just use it with zap.L() or zap.S().

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*             - Logger -              */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var (
	L = zap.L
	S = zap.S
)

const LogsTimeLayout = "02/01/06 15:04:05"

// Replaces the global Logger in the zap package with a new one.
// It uses a default zap.Config and allows for additional options to be passed.
func SetupLogger(cfg *core.LoggerCfg, opts ...zap.Option) *zap.Logger {

	// Default options: Add stacktrace and use the default clock.
	opts = append([]zap.Option{
		zap.AddStacktrace(zapcore.Level(cfg.LevelStackT)),
		zap.WithClock(zapcore.DefaultClock),
	}, opts...)

	logger, err := newZapConfig(cfg).Build(opts...)
	if err != nil {
		log.Fatalf(errs.FailedToCreateLogger, err) // Std log, don't use zap.
	}

	zap.ReplaceGlobals(logger)

	return logger
}

func LogGRPC(route string, duration time.Duration, err error) {
	l := prepareLog(withGRPC(route), withDuration(duration))
	if err == nil {
		l.Info("GRPC Request")
	} else {
		l.Error("GRPC Error", zap.Error(err))
	}
}

// We log both GRPC and HTTP requests, just because.
func LogHTTPRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()
		handler.ServeHTTP(rw, req)
		duration := time.Since(start)

		l := prepareLog(withHTTP(req), withDuration(duration))

		customRW := rw.(core.HTTPRespWriter)
		if customRW.GetWrittenStatus() < 400 {
			l.Info("HTTP Request")
		} else {
			err := errors.New(string(customRW.GetWrittenBody()))
			l.Error("HTTP Error", zap.Error(err))
		}
	})
}

// Prefix used when Infof or Infoln are called.
var ServerLogPrefix = "Server | "

func ServerLog(s string) {
	zap.L().Info(ServerLogPrefix + s)
}

func ServerLogf(s string, args ...any) {
	zap.S().Infof(ServerLogPrefix+s, args...)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - Shorthand -            */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func LogDebug(msg string) {
	prepareLog(withMsg(msg)).Debug("🐞 Debug")
}

// Used to log unexpected errors, like panic recoveries or some connection errors.
// Returns the error so the caller can -> return LogUnexpected(err).
func LogUnexpected(err error) error {
	prepareLog(withError(err)).Error("🛑 Unexpected Error")
	return err
}

// Helps keeping code clean and readable, lets you omit the error check
// on the caller when you can get away with just logging it.
// Use this is for errors that are somewhat expected.
func LogIfErr(err error, optionalFmt ...string) {
	if err == nil {
		return
	}
	format := utils.FirstOrDefault(optionalFmt, "untyped error: %v")
	zap.S().Errorf(format, err)
}

// Used to log unexpected errors that also should trigger a panic.
func LogFatal(err error) {
	prepareLog(wErr(err), wStack()).Fatal("🛑 Fatal Error")
}

// Helps keeping code clean and readable, lets you omit the error check on the caller.
func LogFatalIfErr(err error, optionalFmt ...string) {
	if err == nil {
		return
	}

	format := utils.FirstOrDefault(optionalFmt, "untyped fatal: %v")
	LogFatal(fmt.Errorf(format, err))
}

// Helps keeping code clean and readable, lets you omit the error check
// on the caller when you just need to log-warn the err.
func WarnIfErr(err error, optionalFmt ...string) {
	if err != nil {
		format := "untyped warning: %v"
		if len(optionalFmt) > 0 {
			format = optionalFmt[0]
		}
		zap.S().Warnf(format, err)
	}
}

func LogImportant(msg string) {
	prepareLog(withMsg(msg)).Info("⭐-⭐-⭐")
}

// Used to log strange behaviour that isn't necessarily bad or an error.
func LogStrange(msg string, info ...any) {
	prepareLog(withMsg(msg), withData(info...)).Warn("🤔 Strange")
}

// Used to log security-related things that shouldn't happen,
// like a non-admin trying to access admin endpoints.
func LogThreat(msg string) {
	prepareLog(withMsg(msg)).Warn("🚨 Threat")
}

func LogResult(ofWhat string, err error) {
	if err == nil {
		prepareLog().Info("✅ " + ofWhat + " succeeded!")
	} else {
		prepareLog(withError(err)).Error("❌ " + ofWhat + " failed!")
	}
}

// If Debug is true, logs external API calls.
func LogAPICall(url string, status int, body []byte) {
	if core.G.LogAPICalls {
		zap.L().Info("External API Call",
			zap.String("url", url),
			zap.Int("status", status),
			zap.ByteString("body", body),
		)
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*               - Etc -               */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Returns a new zap.Config based on our LoggerCfg.
func newZapConfig(cfg *core.LoggerCfg) *zap.Config {

	// Start off with a default dev/prod config.
	zapCfg := zap.NewDevelopmentConfig()
	if core.G.IsProd {
		zapCfg = zap.NewProductionConfig()
		zapCfg.Sampling = nil
	}

	// Sets the log level - Shows or hides the function caller.
	zapCfg.Level = zap.NewAtomicLevelAt(zapcore.Level(cfg.Level))
	zapCfg.DisableCaller = !cfg.LogCaller

	// Format dates. Default is -> "02/01/06 15:04:05"
	zapCfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(LogsTimeLayout))
	}

	// Format durations. Default is ms.
	zapCfg.EncoderConfig.EncodeDuration = func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(d.Truncate(time.Millisecond).String())
	}

	return &zapCfg
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - Log Options -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type logOpt func(*zap.Logger)

// Shorthand versions of each logOpt.
var (
	wMsg   = withMsg
	wErr   = withError
	wDurat = withDuration
	wStack = withStacktrace
	wGRPC  = withGRPC
	wHTTP  = withHTTP
)

// Prepares a new child Logger, with fields defined by the given logOpts.
func prepareLog(opts ...logOpt) *zap.Logger {
	childLogger := *zap.L()
	for _, opt := range opts {
		opt(&childLogger)
	}
	return &childLogger
}

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
		*logger = *logger.With(zap.String("route", utils.RouteNameFromGRPC(method)))
	}
}

// -> In HTTP, for the route we join Method and Path -> 'GET /users'.
var withHTTP = func(req *http.Request) logOpt {
	return func(logger *zap.Logger) {
		*logger = *logger.With(zap.String("route", req.Method+" "+req.URL.Path))
	}
}
