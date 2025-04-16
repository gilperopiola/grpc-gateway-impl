package logs

import (
	"go.uber.org/zap"
)

// GetZapLogger returns the global zap logger instance
// for use by other packages
func GetZapLogger() *zap.Logger {
	return zap.L()
}
