package core

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
func SetupLogger(cfg *LoggerCfg, opts ...zap.Option) *zap.Logger {

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
		prepareLog(withMsg(msg)).Info("🐞 Debug")
	}
}

// Used to log unexpected errors, like panic recoveries or some connection errors.
func LogUnexpected(err error) {
	prepareLog(withError(err)).Error("🛑 Unexpected Error")
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

func LogImportant(msg string) {
	prepareLog(withMsg(msg)).Info("⭐-⭐-⭐")
}

// Used to log strange behaviour that isn't necessarily bad or an error.
func LogStrange(msg string, info ...any) {
	prepareLog(withMsg(msg), withData(info...)).Warn("🤔 Hmm... Strange")
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

// Returns a new zap.Config based on our LoggerCfg.
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

// Only logs with a level equal or higher than the one we set in the config will be logged.
// For example -> if the config is 1, warn, then no info or debug logs will be logged.
var LogLevels = map[string]int{
	"debug": int(zap.DebugLevel), "info": int(zap.InfoLevel),
	"warn": int(zap.WarnLevel), "error": int(zap.ErrorLevel),
	"fatal": int(zap.FatalLevel),
}

// Gorm has its own set of log levels, the logging of SQL queries
// depends on the config.
//
// silent 	-> disables all logs 		| info -> 	logs everything
// warn 	-> logs errors and warnings | error -> 	only logs errors.
var DBLogLevels = map[string]int{
	"silent": int(gormLogger.Silent), "info": int(gormLogger.Info),
	"warn": int(gormLogger.Warn), "error": int(gormLogger.Error),
}
