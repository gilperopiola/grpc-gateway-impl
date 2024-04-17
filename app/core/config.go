package core

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
)

// Globals!!! Not sure about this but let's try it out.
// If zap does it then I can too right.
// These are the default values, LoadConfig() could override them.
var (
	AppName  = "grpc-gateway-impl"
	IsProd   = false
	GRPCPort = ":50053"
	HTTPPort = ":8083"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*             - Config -              */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Config holds the configuration values of our app.
type Config struct {
	LoggerCfg
	DatabaseCfg
	PwdHasherCfg
	RateLimiterCfg
	JWTCfg
	TLSCfg
}

// LoadConfig sets up the configuration from the environment variables.
func LoadConfig() *Config {
	AppName = envVar("APP_NAME", AppName)
	IsProd = envVar("IS_PROD", IsProd)
	GRPCPort = envVar("GRPC_PORT", GRPCPort)
	HTTPPort = envVar("HTTP_PORT", HTTPPort)

	return &Config{
		loadLoggerConfig(),
		loadDatabaseConfig(),
		loadPwdHasherConfig(),
		loadRateLimiterConfig(),
		loadJWTConfig(),
		loadTLSConfig(),
	}
}

func loadLoggerConfig() LoggerCfg {
	return LoggerCfg{
		Level:           LogLevels[envVar("LOG_LEVEL", "info")],
		LevelStackTrace: LogLevels[envVar("LOG_LEVEL_STACKTRACE", "dpanic")],
		LogCaller:       envVar("LOG_CALLER", false),
	}
}

func loadDatabaseConfig() DatabaseCfg {
	return DatabaseCfg{
		DatabaseConnCfg: loadDatabaseConnectionCfg(),
		MigrateModels:   envVar("DB_MIGRATE_MODELS", true),
		InsertAdmin:     envVar("DB_INSERT_ADMIN", true),
		AdminPwd:        envVar("DB_ADMIN_PWD", "n8zAyv96oAtfQoNof-_ulH4pS0Dqf61VThTZbbOLXCU="), // hashed, T0D0 change this unsafe!!!
		LogLevel:        LogLevels[envVar("DB_LOG_LEVEL", "error")],
	}
}

func loadDatabaseConnectionCfg() DatabaseConnCfg {
	return DatabaseConnCfg{
		envVar("DB_USERNAME", "root"),
		envVar("DB_PASSWORD", ""),
		envVar("DB_HOSTNAME", "localhost"),
		envVar("DB_PORT", "3306"),
		envVar("DB_SCHEMA", "grpc-gateway-impl"),
		envVar("DB_PARAMS", "?charset=utf8&parseTime=True&loc=Local"),
	}
}

func loadPwdHasherConfig() PwdHasherCfg {
	return PwdHasherCfg{Salt: envVar("HASH_SALT", "s0m3_s4l7")}
}

func loadRateLimiterConfig() RateLimiterCfg {
	return RateLimiterCfg{
		MaxTokens:       envVar("RATE_LIMITER_MAX_TOKENS", 40),
		TokensPerSecond: envVar("RATE_LIMITER_TOKENS_PER_SECOND", 10),
	}
}

func loadJWTConfig() JWTCfg {
	return JWTCfg{
		Secret:      envVar("JWT_SECRET", "s0m3_s3cr37"),
		SessionDays: envVar("JWT_SESSION_DAYS", 7),
	}
}

func loadTLSConfig() TLSCfg {
	rootPrefix := getRootPrefix(AppName) // "." or "../.." depending on where is the app being run from
	return TLSCfg{
		Enabled:  envVar("TLS_ENABLED", false),
		CertPath: envVar("TLS_CERT_PATH", rootPrefix+"/server.crt"),
		KeyPath:  envVar("TLS_KEY_PATH", rootPrefix+"/server.key"),
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*         - Config Structure -        */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type LoggerCfg struct {
	Level           int
	LevelStackTrace int
	LogCaller       bool
}

type DatabaseCfg struct {
	DatabaseConnCfg

	MigrateModels bool
	InsertAdmin   bool
	AdminPwd      string // hashed
	LogLevel      int
}

type DatabaseConnCfg struct {
	Username string
	Password string
	Hostname string
	Port     string
	Schema   string
	Params   string
}

type PwdHasherCfg struct {
	Salt string
}

type RateLimiterCfg struct {
	MaxTokens       int // Max tokens the bucket can hold.
	TokensPerSecond int // Tokens reloaded per second.
}

type JWTCfg struct {
	Secret      string
	SessionDays int
}

type TLSCfg struct {
	Enabled  bool
	CertPath string
	KeyPath  string
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*         - Config Helpers -          */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func envVar[T string | bool | int](key string, fallback T) T {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	switch any(fallback).(type) {
	case string:
		return any(value).(T)
	case bool:
		boolValue := value == "true" || value == "TRUE" || value == "1"
		return any(boolValue).(T)
	case int:
		if intValue, err := strconv.Atoi(value); err == nil {
			return any(intValue).(T)
		}
	}
	return fallback
}

// getRootPrefix returns the prefix that needs to be added to the default paths to point at the root folder.
func getRootPrefix(projectName string) string {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf(errs.FatalErrMsgGettingWorkingDir, err) // zap is not initialized yet.
	}
	if strings.HasSuffix(workingDir, projectName) {
		return "." // -> running from root folder
	}
	return "../.." // -> running from /etc/cmd
}
