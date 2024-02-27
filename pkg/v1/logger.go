package v1

import (
	"log"

	"go.uber.org/zap"
)

/* ----------------------------------- */
/*             - Logger -              */
/* ----------------------------------- */

// NewLogger returns a new instance of *zap.Logger.
func NewLogger(isProd bool, opts []zap.Option) *zap.Logger {
	newLoggerFn := zap.NewDevelopment
	if isProd {
		newLoggerFn = zap.NewProduction
	}

	logger, err := newLoggerFn(opts...)
	if err != nil {
		log.Fatalf(msgErrCreatingLogger_Fatal, err)
	}

	return logger
}

const (
	msgErrCreatingLogger_Fatal = "Failed to create zap logger: %v"
)
