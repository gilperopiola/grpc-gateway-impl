package sql

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var _ gormLogger.Interface = (&sqlLogger{})

// sqlLogger is an adapter for gormLogger.Interface. It wraps our *zap.Logger.
type sqlLogger struct {
	*zap.Logger
	gormLogger.LogLevel
}

// Returns a new instance of *sqlLogger.
// We set the Log Level according to the configuration.
func newSQLLogger(zapLogger *zap.Logger, logLevel int) *sqlLogger {
	return &sqlLogger{zapLogger, gormLogger.LogLevel(logLevel)}
}

// LogMode sets the Log Level.
func (l *sqlLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info logs info level logs.
func (l *sqlLogger) Info(_ context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Info {
		zap.S().Infof(msg, data...)
	}
}

// Warn logs warning level logs.
func (l *sqlLogger) Warn(_ context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Warn {
		zap.S().Warnf(msg, data...)
	}
}

// Error logs error level logs.
func (l *sqlLogger) Error(_ context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Error {
		zap.S().Errorf(msg, data...)
	}
}

// Trace logs trace level logs including the time taken for the operation, affected rows, and error if any.
func (l *sqlLogger) Trace(_ context.Context, begin time.Time, fnCall func() (string, int64), err error) {
	if l.LogLevel <= gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	query, rows := fnCall()

	// If there's an error other than gorm.ErrRecordNotFound, log it.
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		var netError *net.OpError
		if errors.As(err, &netError) {
			// For some errors, we don't want to log the query.
			query = "DB Network Error"
		}

		zap.S().Errorf(newQueryInfoLog(elapsed.Nanoseconds(), rows, query), zap.Error(err))
		return
	}

	// Log the query if the log level is set to Info or if it took more than 1 second.
	queryWasSlow := elapsed > time.Second
	if l.LogLevel >= gormLogger.Info || queryWasSlow {
		zap.S().Infof(newQueryInfoLog(elapsed.Nanoseconds(), rows, query))
	}
}

// newQueryInfoLog returns the query info formatted as a string:
// [1 row in 25ms] -> INSERT INTO `table` (`field1`,`field2`) VALUES ('gorm','sucks')
//
// The gorm default format isn't much different:
// [25.14ms] [row:1] INSERT INTO `table` (`field1`,`field2`) VALUES ('gorm','sucks')
func newQueryInfoLog(nsElapsed, rowsAffected int64, query string) string {
	msElapsed := float64(nsElapsed) / 1e6
	rowOrRows := "row" + plural(rowsAffected)

	return fmt.Sprintf("[%v %s in %0.fms] -> %s", rowsAffected, rowOrRows, msElapsed, query)
}

// Returns the 's' necessary for changing 'row' into 'rows' based on the number of rows affected.
func plural(rowsAffected int64) string {
	if rowsAffected == 1 {
		return ""
	}
	return "s"
}
