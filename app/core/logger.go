package core

import (
	"context"
	"errors"
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

// We use zap. It's fast and easy.
// Set it up and then just use it with zap.L() or zap.S().

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*             - Logger -              */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

const LogsTimeLayout = "02/01/06 15:04:05"

// Replaces the global Logger in the zap pkg with a new one.
// It uses a default zap.Config and allows for additional options to be passed.
func SetupLogger(cfg *LoggerCfg) *zap.Logger {

	// Default options: Add stacktrace and use the default clock.
	opts := []zap.Option{
		zap.AddStacktrace(zapcore.Level(cfg.LevelStackT)),
		zap.WithClock(zapcore.DefaultClock),
	}

	logger, err := newZapConfig(cfg).Build(opts...)
	if err != nil {
		log.Fatalf(errs.FailedToCreateLogger, err) // Std log, don't use zap.
	}

	zap.ReplaceGlobals(logger)

	return logger
}

// This func is a GRPC Interceptor. Or technically a grpc.UnaryServerInterceptor.
func LogGRPCRequest(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)

	l := newLog(withGRPC(info.FullMethod), withDuration(duration))

	if err == nil {
		l.Info("GRPC Request")
	} else {
		l.Error("GRPC Error", zap.Error(err))
	}

	return resp, err
}

// We log both GRPC and HTTP requests, just because.
func LogHTTPRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()
		handler.ServeHTTP(rw, req)
		duration := time.Since(start)

		l := newLog(withHTTP(req), withDuration(duration))

		customRW := rw.(HTTPRespWriter)
		if customRW.GetWrittenStatus() < 400 {
			l.Info("HTTP Request")
		} else {
			err := errors.New(string(customRW.GetWrittenBody()))
			l.Error("HTTP Error", zap.Error(err))
		}
	})
}

// Prefix used when Infof or Infoln are called.
var ServerLogPrefix = AppEmoji + " " + AppAlias + " | "

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
	if Debug {
		newLog(withMsg(msg)).Info("🐞 Debug")
	}
}

// Used to log unexpected errors, like panic recoveries or some connection errors.
func LogUnexpected(err error) {
	newLog(withError(err)).Error("🛑 Unexpected Error")
}

// Helps keeping code clean and readable, lets you omit the error check
// on the caller when you just need to log the err.
// Use this is for errors that are expected.
func LogIfErr(err error, optionalFmt ...string) {
	if err != nil {
		format := "untyped error: %v"
		if len(optionalFmt) > 0 {
			format = optionalFmt[0]
		}
		zap.S().Errorf(format, err)
	}
}

// Used to log unexpected errors that also should trigger a panic.
func LogFatal(err error) {
	newLog(withError(err), withStacktrace()).Fatal("🛑 Fatal Error")
}

// Helps keeping code clean and readable, lets you omit the error check on the caller.
func LogFatalIfErr(err error, optionalFmt ...string) {
	if err == nil {
		return
	}

	format := "untyped fatal: %v"
	if len(optionalFmt) > 0 {
		format = optionalFmt[0]
	}

	LogFatal(fmt.Errorf(format, err))
}

func LogImportant(msg string) {
	newLog(withMsg(msg)).Info("⭐ Important!")
}

// Used to log strange behaviour that isn't necessarily bad or an error.
func LogWeirdBehaviour(msg string, info ...any) {
	newLog(withMsg(msg), withData(info...)).Warn("🤔 Weird")
}

// Used to log things that shouldn't happen, like someone trying to access admin endpoints.
func LogPotentialThreat(msg string) {
	newLog(withMsg(msg)).Warn("🚨 Potential Threat")
}

func LogResult(operation string, err error) {
	if err == nil {
		newLog().Info("✅ " + operation + " succeeded!")
	} else {
		newLog(withError(err)).Error("❌ " + operation + " failed!")
	}
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

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*               - Etc -               */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Returns a new zap.Config.
// You can pass it custom options like the log level.
func newZapConfig(cfg *LoggerCfg) zap.Config {

	// Start off with a default dev/prod config.
	zapCfg := zap.NewDevelopmentConfig()
	if EnvIsProd {
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

	return zapCfg
}

// Only logs with a level equal or higher than the one we set in the config
// will be logged.
// For example, if the config is 'warn' no info or debug logs will be logged.
var LogLevels = map[string]int{
	"debug": int(zap.DebugLevel), "info": int(zap.InfoLevel), "warn": int(zap.WarnLevel),
	"error": int(zap.ErrorLevel), "dpanic": int(zap.DPanicLevel), "panic": int(zap.PanicLevel),
	"fatal": int(zap.FatalLevel),
}

// The selected DB Log Level will be used to log all SQL queries.
// 'silent' disables all logs, 'info' logs everything,
// 'warn' logs errors and warnings, and 'error' will only log errors.
var DBLogLevels = map[string]int{
	"silent": int(gormLogger.Silent), "info": int(gormLogger.Info),
	"warn": int(gormLogger.Warn), "error": int(gormLogger.Error),
}
