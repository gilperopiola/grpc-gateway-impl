package core

import (
	"fmt"
	"os"
	"strconv"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*             - Config -              */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Our most used Configs are... -> 🌍 Globals!~
// -> Just call them like core.Whatever from anywhere and you're good to go.

var Env = "local"
var EnvIsProd = false

var AppName = "grpc-gateway-impl"
var AppAlias = "GrpcG8Way"
var AppEmoji = "📱"
var AppVersion = "v1.0" // TODO - Makefile should pass this as an env var.

var GRPCPort = ":50053"
var HTTPPort = ":8083"
var TLSEnabled = false
var Debug = false
var LogAPICalls = true

// These are our non-global Configs 🌍❌
// -> The App loads an instance of this on startup and passes it around.
type Config struct {
	DBCfg        // -> DB Credentials and such
	JWTCfg       // -> JWT Secret
	TLSCfg       // -> TLS Certs paths
	LoggerCfg    // -> Logger settings
	PwdHasherCfg // -> Salt
	RetrierCfg   // -> N° Retries
	RLimiterCfg  // -> Rate settings
}

func LoadConfig() *Config {

	// -> 🌍 Globals
	Env = envVar("ENV", Env)
	EnvIsProd = Env == "prod" || Env == "production" || Env == "live"

	AppName = envVar("APP_NAME", AppName)
	AppAlias = envVar("APP_ALIAS", AppAlias)
	AppEmoji = envVar("APP_EMOJI", AppEmoji)
	AppVersion = envVar("APP_VERSION", AppVersion)

	GRPCPort = envVar("GRPC_PORT", GRPCPort)
	HTTPPort = envVar("HTTP_PORT", HTTPPort)
	TLSEnabled = envVar("TLS_ENABLED", TLSEnabled)
	Debug = envVar("DEBUG", Debug)

	return &Config{
		DBCfg:        loadDBConfig(),
		JWTCfg:       loadJWTConfig(),
		TLSCfg:       loadTLSConfig(),
		LoggerCfg:    loadLoggerConfig(),
		PwdHasherCfg: loadPwdHasherConfig(),
		RetrierCfg:   loadRetrierConfig(),
		RLimiterCfg:  loadRateLimiterConfig(),
	}
}

func loadDBConfig() DBCfg {
	return DBCfg{
		Username:       envVar("DB_USERNAME", "root"),
		Password:       envVar("DB_PASSWORD", ""),
		Hostname:       envVar("DB_HOSTNAME", "localhost"),
		Port:           envVar("DB_PORT", "3306"),
		Schema:         envVar("DB_SCHEMA", "grpc-gateway-impl"),
		Params:         envVar("DB_PARAMS", "?charset=utf8&parseTime=True&loc=Local"),
		Retries:        envVar("DB_RETRIES", 7),
		EraseAllData:   envVar("DB_ERASE_ALL_DATA", false),
		MigrateModels:  envVar("DB_MIGRATE_MODELS", true),
		InsertAdmin:    envVar("DB_INSERT_ADMIN", true),
		InsertAdminPwd: envVar("DB_INSERT_ADMIN_PWD", "n8zAyv96oAtfQoNof-_ulH4pS0Dqf61VThTZbbOLXCU="), // T0D0 unsafe!!!! But it's local so...
		LogLevel:       LogLevels[envVar("DB_LOG_LEVEL", "error")],
	}
}

func loadJWTConfig() JWTCfg {
	return JWTCfg{
		Secret:      envVar("JWT_SECRET", "s0m3_s3cr37"),
		SessionDays: envVar("JWT_SESSION_DAYS", 7),
	}
}

func loadTLSConfig() TLSCfg {
	return TLSCfg{
		CertPath: envVar("TLS_CERT_PATH", "./server.crt"),
		KeyPath:  envVar("TLS_KEY_PATH", "./server.key"),
	}
}

func loadLoggerConfig() LoggerCfg {
	return LoggerCfg{
		Level:       LogLevels[envVar("LOGGER_LEVEL", "info")],
		LevelStackT: LogLevels[envVar("LOGGER_LEVEL_STACKTRACE", "fatal")],
		LogCaller:   envVar("LOGGER_LOG_CALLER", false),
	}
}

func loadPwdHasherConfig() PwdHasherCfg {
	return PwdHasherCfg{Salt: envVar("PWD_HASHER_SALT", "s0m3_s4l7")}
}

func loadRetrierConfig() RetrierCfg {
	return RetrierCfg{}
}

func loadRateLimiterConfig() RLimiterCfg {
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
		Retries  int

		EraseAllData   bool
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
	RetrierCfg struct {
	}
	RLimiterCfg struct {
		MaxTokens       int // Max tokens the bucket can hold
		TokensPerSecond int // Tokens reloaded per second
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

// Used on init if the DB we need is not yet created
func (c *DBCfg) GetSQLConnStringNoSchema() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.Username, c.Password, c.Hostname, c.Port, c.Params)
}
