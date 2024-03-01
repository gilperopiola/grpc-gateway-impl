package server

import (
	"log"

	"go.uber.org/zap"
)

const (
	errMsgCreatingLogger_Fatal = "Failed to create logger: %v" // Fatal error.
)

/* ----------------------------------- */
/*             - Logger -              */
/* ----------------------------------- */

// newLogger returns a new instance of *zap.Logger.
func newLogger(isProd bool, opts []zap.Option) *zap.Logger {
	newLoggerFn := zap.NewDevelopment
	if isProd {
		newLoggerFn = zap.NewProduction
	}

	logger, err := newLoggerFn(opts...)
	if err != nil {
		log.Fatalf(errMsgCreatingLogger_Fatal, err)
	}

	return logger
}

// newLoggerOptions returns the default options for the logger.
// For now it only adds a stack trace to panic logs.
func newLoggerOptions() []zap.Option {
	return []zap.Option{
		zap.AddStacktrace(zap.DPanicLevel),
	}
}
