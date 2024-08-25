package types

import (
	"go.uber.org/zap"
	gormLogger "gorm.io/gorm/logger"
)

// Unsure about the location of these 2 maps.

// Levels go from debug (lowest) to fatal (highest).
//
// We only log messages with a Level equal or higher
// than the one in the [LoggerCfg].
var LogLevels = map[string]int{
	"debug": int(zap.DebugLevel), // Lowest.
	"info":  int(zap.InfoLevel),
	"warn":  int(zap.WarnLevel),
	"error": int(zap.ErrorLevel),
	"fatal": int(zap.FatalLevel), // Highest.
}

// Gorm Levels go from silent (lowest) to info (highest),
// as opposed to our Logger Levels.
//
//	info 	-> 	logs everything
//	warn 	-> 	logs warnings + errors
//	error 	-> 	logs errors
//	silent 	-> 	logs nothing
//
// We also added debug and fatal, but they just map to info and error.
var DBLogLevels = map[string]int{
	"debug":  int(gormLogger.Info),
	"info":   int(gormLogger.Info),
	"warn":   int(gormLogger.Warn),
	"error":  int(gormLogger.Error),
	"fatal":  int(gormLogger.Error),
	"silent": int(gormLogger.Silent),
}
