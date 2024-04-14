package sqldb

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// gormLoggerAdapter is an adapter for the Gorm Logger. It wraps our *zap.Logger and implements the Gorm Logger interface.
type gormLoggerAdapter struct {
	*zap.Logger
	logger.LogLevel
}

// newGormLoggerAdapter returns a new instance of *gormLoggerAdapter.
// We set the Log Level according to the configuration.
func newGormLoggerAdapter(l *zap.Logger, logLevel int) *gormLoggerAdapter {
	return &gormLoggerAdapter{l, logger.LogLevel(logLevel)}
}

// LogMode sets the Log Level.
func (gl *gormLoggerAdapter) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *gl
	newLogger.LogLevel = level
	return &newLogger
}

// Info logs info level logs.
func (gl *gormLoggerAdapter) Info(ctx context.Context, msg string, data ...interface{}) {
	if gl.LogLevel >= logger.Info {
		zap.S().Infof(msg, data...)
	}
}

// Warn logs warning level logs.
func (gl *gormLoggerAdapter) Warn(ctx context.Context, msg string, data ...interface{}) {
	if gl.LogLevel >= logger.Warn {
		zap.S().Warnf(msg, data...)
	}
}

// Error logs error level logs.
func (gl *gormLoggerAdapter) Error(ctx context.Context, msg string, data ...interface{}) {
	if gl.LogLevel >= logger.Error {
		zap.S().Errorf(msg, data...)
	}
}

// Trace logs trace level logs including the time taken for the operation, affected rows, and error if any.
func (gl *gormLoggerAdapter) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if gl.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	query, rows := fc()

	// If there's an error other than gorm.ErrRecordNotFound, log it.
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// For some errors, we don't want to log the query.
		var netError *net.OpError
		if errors.As(err, &netError) {
			query = "DB Network Error"
			return
		}

		zap.S().Errorf(getQueryInfo(elapsed.Nanoseconds(), rows, query), zap.Error(err))
		return
	}

	// Log the query if the log level is set to Info or if it took more than 1 second.
	if gl.LogLevel >= logger.Info || elapsed > 1000*time.Millisecond {
		zap.S().Infof(getQueryInfo(elapsed.Nanoseconds(), rows, query))
	}
}

// getQueryInfo returns the query info as a string. Example:
//
// [1 row in 25ms] -> INSERT INTO `table` (`field1`,`field2`) VALUES ('gorm','sucks')
//
// However, the standard format is actually:
// [25.14ms] [row:1] INSERT INTO `table` (`field1`,`field2`) VALUES ('gorm','sucks')
func getQueryInfo(nsElapsed, rowsAffected int64, query string) string {
	msElapsed := float64(nsElapsed) / 1e6
	rowOrRows := "row" + pluralPrefix(rowsAffected)

	return fmt.Sprintf("[%v %s in %0.fms] -> %s", rowsAffected, rowOrRows, msElapsed, query)
}

// pluralPrefix returns the 's' necessary for the plural form of a word, namely 'row', if the rows affected are not 1.
func pluralPrefix(rowsAffected int64) string {
	if rowsAffected == 1 {
		return ""
	}
	return "s"
}
