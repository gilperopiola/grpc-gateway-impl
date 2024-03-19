package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// gormLoggerAdapter is an adapter for the Gorm Logger. It wraps our zap Logger and implements the Gorm Logger interface.
type gormLoggerAdapter struct {
	*zap.Logger
	logger.LogLevel
}

// newGormLoggerAdapter returns a new instance of *gormLoggerAdapter. We set the Log Level to Warn to avoid logging failed queries.
func newGormLoggerAdapter(l *zap.Logger) *gormLoggerAdapter {
	return &gormLoggerAdapter{l, logger.Info}
}

// LogMode sets the log level for the logger.
func (g *gormLoggerAdapter) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *g
	newLogger.LogLevel = level
	return &newLogger
}

// Info logs info level logs.
func (g *gormLoggerAdapter) Info(ctx context.Context, msg string, data ...interface{}) {
	if g.LogLevel >= logger.Info {
		zap.S().Infof(msg, data...)
	}
}

// Warn logs warning level logs.
func (g *gormLoggerAdapter) Warn(ctx context.Context, msg string, data ...interface{}) {
	if g.LogLevel >= logger.Warn {
		zap.S().Warnf(msg, data...)
	}
}

// Error logs error level logs.
func (g *gormLoggerAdapter) Error(ctx context.Context, msg string, data ...interface{}) {
	if g.LogLevel >= logger.Error {
		zap.S().Errorf(msg, data...)
	}
}

// Trace logs trace level logs including the time taken for the operation, affected rows, and error if any.
func (g *gormLoggerAdapter) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if g.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// If there's an error other than gorm.ErrRecordNotFound, log it.
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		zap.S().Errorf(getQueryInfo(elapsed.Nanoseconds(), rows, sql), zap.Error(err))
		return
	}

	// Log the query info.
	if g.LogLevel >= logger.Info {
		zap.S().Infof(getQueryInfo(elapsed.Nanoseconds(), rows, sql))
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
	rowOrRows := "row" + getPluralPrefix(rowsAffected)

	return fmt.Sprintf("[%v %s in %0.fms] -> %s", rowsAffected, rowOrRows, msElapsed, query)
}

// getPluralPrefix returns the 's' necessary for the plural form of a word, namely 'row', if the rows affected are not 1.
func getPluralPrefix(rowsAffected int64) string {
	if rowsAffected == 1 {
		return ""
	}
	return "s"
}
