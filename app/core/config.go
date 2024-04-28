package core

import (
	"fmt"
	"os"
	"strconv"
)

// Globals! Because why not?
// If Zap uses globals then I can too :) -> **gets decapitated by the Go community**
var AppName = "grpc-gateway-impl"
var AppAlias = "GGWI"
var Env = "local"
var EnvIsProd = false
var GRPCPort = ":50053"
var HTTPPort = ":8083"

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*             - Config -              */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Holds the configurable settings of our app.
type Config struct {
	DBCfg
	JWTCfg
	TLSCfg
	LoggerCfg
	PwdHasherCfg
	RLimiterCfg
}

// Loads all settings from environment variables.
func SetupConfig() *Config {

	// Globals
	AppName = envVar("APP_NAME", AppName)
	AppAlias = envVar("APP_ACRONYM", AppAlias)
	Env = envVar("ENV_NAME", Env)
	EnvIsProd = Env == "prod" || Env == "production" || Env == "live"
	GRPCPort = envVar("GRPC_PORT", GRPCPort)
	HTTPPort = envVar("HTTP_PORT", HTTPPort)

	return &Config{
		DBCfg:        setupDBConfig(),
		JWTCfg:       setupJWTConfig(),
		TLSCfg:       setupTLSConfig(),
		LoggerCfg:    setupLoggerConfig(),
		PwdHasherCfg: setupPwdHasherConfig(),
		RLimiterCfg:  setupRateLimiterConfig(),
	}
}

func setupDBConfig() DBCfg {
	return DBCfg{
		Username:       envVar("DB_USERNAME", "root"),
		Password:       envVar("DB_PASSWORD", ""),
		Hostname:       envVar("DB_HOSTNAME", "localhost"),
		Port:           envVar("DB_PORT", "3306"),
		Schema:         envVar("DB_SCHEMA", "grpc-gateway-impl"),
		Params:         envVar("DB_PARAMS", "?charset=utf8&parseTime=True&loc=Local"),
		MigrateModels:  envVar("DB_MIGRATE_MODELS", true),
		InsertAdmin:    envVar("DB_INSERT_ADMIN", true),
		InsertAdminPwd: envVar("DB_INSERT_ADMIN_PWD", "n8zAyv96oAtfQoNof-_ulH4pS0Dqf61VThTZbbOLXCU="), // T0D0 unsafe!!!! But it's local so...
		LogLevel:       LogLevels[envVar("DB_LOG_LEVEL", "error")],
	}
}

func setupJWTConfig() JWTCfg {
	return JWTCfg{
		Secret:      envVar("JWT_SECRET", "s0m3_s3cr37"),
		SessionDays: envVar("JWT_SESSION_DAYS", 7),
	}
}

func setupTLSConfig() TLSCfg {
	return TLSCfg{
		Enabled:  envVar("TLS_ENABLED", false),
		CertPath: envVar("TLS_CERT_PATH", "./server.crt"),
		KeyPath:  envVar("TLS_KEY_PATH", "./server.key"),
	}
}

func setupLoggerConfig() LoggerCfg {
	return LoggerCfg{
		Level:       LogLevels[envVar("LOGGER_LEVEL", "info")],
		LevelStackT: LogLevels[envVar("LOGGER_LEVEL_STACKTRACE", "dpanic")],
		LogCaller:   envVar("LOGGER_LOG_CALLER", false),
	}
}

func setupPwdHasherConfig() PwdHasherCfg {
	return PwdHasherCfg{Salt: envVar("PWD_HASHER_SALT", "s0m3_s4l7")}
}

func setupRateLimiterConfig() RLimiterCfg {
	return RLimiterCfg{
		MaxTokens:       envVar("RLIMITER_MAX_TOKENS", 40),
		TokensPerSecond: envVar("RLIMITER_TOKENS_PER_SECOND", 10),
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type (
	DBCfg struct {
		Username string
		Password string
		Hostname string
		Port     string
		Schema   string
		Params   string

		MigrateModels  bool
		InsertAdmin    bool
		InsertAdminPwd string // hashed with our PwdHasher

		LogLevel int
	}
	JWTCfg struct {
		Secret      string
		SessionDays int
	}
	TLSCfg struct {
		Enabled  bool
		CertPath string
		KeyPath  string
	}
	LoggerCfg struct {
		Level       int
		LevelStackT int
		LogCaller   bool
	}
	PwdHasherCfg struct {
		Salt string
	}
	RLimiterCfg struct {
		MaxTokens       int // Max tokens the bucket can hold.
		TokensPerSecond int // Tokens reloaded per second.
	}
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func envVar[T string | bool | int](key string, fallback T) T {
	val, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	switch any(fallback).(type) {
	case string:
		return any(val).(T)
	case bool:
		return any(val == "true" || val == "TRUE" || val == "1").(T)
	case int:
		if intVal, err := strconv.Atoi(val); err == nil {
			return any(intVal).(T)
		}
	}
	return fallback
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (c *DBCfg) GetSQLConnString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", c.Username, c.Password, c.Hostname, c.Port, c.Schema, c.Params)
}
