package sqldb

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/gilperopiola/god"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var _ gormLogger.Interface = &sqlDBLogger{}

// Adapter for the gormLogger.Interface.
// It wraps a *zap.Logger.
type sqlDBLogger struct {
	*zap.Logger
	gormLogger.LogLevel
}

// Returns a new instance of *sqlDBLogger with the given log level.
func newDBLogger(zapLogger *zap.Logger, level int) *sqlDBLogger {
	return &sqlDBLogger{
		zapLogger,
		gormLogger.LogLevel(level),
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (l *sqlDBLogger) Info(_ god.Ctx, msg string, data ...any) {
	l.InfoWarnOrError(l.LogLevel, gormLogger.Info, "ℹ️ "+msg, zap.S().Infof, data...)
}

func (l *sqlDBLogger) Warn(_ god.Ctx, msg string, data ...any) {
	l.InfoWarnOrError(l.LogLevel, gormLogger.Warn, "🚨 "+msg, zap.S().Warnf, data...)
}

func (l *sqlDBLogger) Error(_ god.Ctx, msg string, data ...any) {

	// Gorm tries to log this if it fails to connect to the DB,
	// however we handle it ourselves when we retry the connection, so skip.
	if strings.HasPrefix(msg, "failed to initialize database") {
		return
	}

	l.InfoWarnOrError(l.LogLevel, gormLogger.Error, "🛑 "+msg, zap.S().Errorf, data...)
}

// -> This gets called after every query. -> It logs the query, the time it took to execute, and the number of rows affected.
// -> If the log level is set to Silent, it doesn't log anything. -> If the log level is set to Info, it logs everything.
// -> If the log level is set to Warn, it logs only slow queries. -> If the log level is set to Error, it logs only errors.
// -> If the query returns an error, it logs the error (I'm not 100% sure about this :))
func (l *sqlDBLogger) Trace(_ god.Ctx, begin time.Time, fnCall func() (string, int64), err error) {
	query, rows := fnCall() // Execute query
	elapsed := time.Since(begin)

	if l.LogLevel <= gormLogger.Silent {
		return
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		l.LogSQLError(err, query, rows, elapsed)
		return
	}

	l.LogSQLQuery(query, rows, elapsed)
}

func (l *sqlDBLogger) LogSQLError(err error, query string, rows int64, elapsed time.Duration) {
	var netError *net.OpError
	if errors.As(err, &netError) {
		query = "DB Network Error" // Don't expose sensitive information.
	}
	zap.S().Errorf("🛑 "+newQueryInfoLog(elapsed.Nanoseconds(), rows, query), zap.Error(err))
}

func (l *sqlDBLogger) LogSQLQuery(query string, rows int64, elapsed time.Duration) {
	queryWasSlow := elapsed > time.Second // T0D0 move to config

	if l.LogLevel >= gormLogger.Info || queryWasSlow {
		zap.S().Infof("ℹ️ " + newQueryInfoLog(elapsed.Nanoseconds(), rows, query))
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (l *sqlDBLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	copiedLogger := *l
	copiedLogger.LogLevel = level
	return &copiedLogger
}

func (l *sqlDBLogger) InfoWarnOrError(currLogLevel gormLogger.LogLevel, logsAt gormLogger.LogLevel, msg string, fn func(string, ...interface{}), data ...interface{}) {
	if currLogLevel >= logsAt {
		fn(msg, data...)
	}
}

// newQueryInfoLog returns the query info formatted as a string:
// [1 row in 25ms] -> INSERT INTO `table` (`field1`,`field2`) VALUES ('gorm','sucks')
//
// The gorm default format isn't much different:
// [25.14ms] [row:1] INSERT INTO `table` (`field1`,`field2`) VALUES ('gorm','sucks')
func newQueryInfoLog(nsElapsed, rowsAffected int64, query string) string {
	msElapsed := float64(nsElapsed) / 1e6
	rowOrRows := "row" + sIfPlural(rowsAffected)

	return fmt.Sprintf("[%v %s in %0.fms] -> %s", rowsAffected, rowOrRows, msElapsed, query)
}

// Returns the 's' necessary for changing 'row' into 'rows' based on the number of rows affected.
func sIfPlural(rowsAffected int64) string {
	if rowsAffected == 1 {
		return ""
	}
	return "s"
}
