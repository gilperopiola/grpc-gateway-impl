package core

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// —► A .env file on the root of the project is required to run the app.
// —► —► —► PLEASE add any new configs to the .env.example file!!!

/* ———————————————————————————————— — — — CORE: CONFIGURATION — — — ———————————————————————————————— */

// init ⚡ reads the .env file located at the root folder.
// Sets environment vars and globals.
func init() {
	if err := godotenv.Overload(); err != nil {
		log.Printf("🚨 WARNING: No .env file found: %v", err)
	}
	setGlobals()
}

// ⭐️ Our App holds a reference to one of this, which contains all the config values
// to be passed from the App to the different services and tools.
type Config struct {
	APIsCfg      // —► API URLs, keys, etc
	DBCfg        // —► DB Credentials and such
	JWTCfg       // —► JWT Secret
	TLSCfg       // —► TLS Certs paths
	LoggerCfg    // —► Logger settings
	PwdHasherCfg // —► Salt
	RetrierCfg   // —► N° Retries
	RLimiterCfg  // —► Rate settings
}

// As on the init func we load the .env file, in here we already
// have available all of the env vars.
func LoadConfig() *Config {
	defer func() {
		log.Println(" \t🎈 APP = " + G.AppName + " " + G.Version)
		log.Println(" \t🌐 ENV = " + G.Env)
		if G.TLSEnabled {
			log.Println(" \t✅ TLS = on")
		} else {
			log.Println(" \t⚠️  TLS = off")
		}
	}()

	return &Config{
		APIsCfg:      loadAPIsConfig(),
		DBCfg:        loadDBConfig(),
		JWTCfg:       loadJWTConfig(),
		TLSCfg:       loadTLSConfig(),
		LoggerCfg:    loadLoggerConfig(),
		PwdHasherCfg: loadPwdHasherConfig(),
		RetrierCfg:   loadRetrierConfig(),
		RLimiterCfg:  loadRateLimiterConfig(),
	}
}

/* -~-~-~-~ Global Config ~-~-~-~- */

var G = new(Globals)

// 🌍 Global Configs
type Globals struct {
	AppName     string
	Version     string // -> TODO - Makefile should set this value.
	Env         string
	IsProd      bool
	IsStage     bool
	IsDev       bool
	GRPCPort    string
	HTTPPort    string
	TLSEnabled  bool
	LogAPICalls bool
}

func setGlobals() {
	var env = envVar("ENV", "local")
	G = &Globals{
		Env:         env,
		IsProd:      env == "prod" || env == "production" || env == "live",
		IsStage:     env == "stage" || env == "staging" || env == "test",
		IsDev:       env == "local" || env == "dev" || env == "development",
		AppName:     envVar("APP_NAME", "grpc-gateway-impl"),
		Version:     envVar("APP_VERSION", "v1.0"),
		GRPCPort:    envVar("GRPC_PORT", ":50053"),
		HTTPPort:    envVar("HTTP_PORT", ":8083"),
		TLSEnabled:  envVar("TLS_ENABLED", false),
		LogAPICalls: envVar("LOG_API_CALLS", true),
	}
}

/* -~-~-~-~ APIs Config ~-~-~-~- */

type (
	APIsCfg struct {
		Weather OpenWeatherMapAPICfg
		GPT     ChatGptAPICfg
	}
	OpenWeatherMapAPICfg struct {
		BaseURL string
		AppID   string
	}
	ChatGptAPICfg struct {
		BaseURL string
		APIKey  string
	}
)

func loadAPIsConfig() APIsCfg {
	return APIsCfg{
		OpenWeatherMapAPICfg{
			BaseURL: envVar("API_OPENWEATHERMAP_BASE_URL", "https://api.weathermap.org/data/2.5/weather"),
			AppID:   envVar("API_OPENWEATHERMAP_APP_ID", ""),
		},
		ChatGptAPICfg{
			BaseURL: envVar("API_CHATGPT_BASE_URL", "https://api.openai.com/v1"),
			APIKey:  envVar("API_CHATGPT_API_KEY", ""),
		},
	}
}

/* -~-~-~-~ DB Config ~-~-~-~- */

type DBCfg struct {
	Username       string
	Password       string
	Hostname       string
	Port           string
	Database       string
	Params         string
	Retries        int
	EraseAllData   bool
	MigrateModels  bool
	InsertAdmin    bool
	InsertAdminPwd string // hashed with our PwdHasher
	LogLevel       int
}

func (c *DBCfg) IsPostgres() bool {
	return true
}

func loadDBConfig() DBCfg {
	return DBCfg{
		Username:       envVar("DB_USERNAME", "root"),
		Password:       envVar("DB_PASSWORD", ""),
		Hostname:       envVar("DB_HOSTNAME", "localhost"),
		Port:           envVar("DB_PORT", "3306"),
		Database:       envVar("DB_DATABASE", "grpc-gateway-impl"),
		Params:         envVar("DB_PARAMS", "?charset=utf8&parseTime=True&loc=Local"),
		Retries:        envVar("DB_RETRIES", 7),
		EraseAllData:   envVar("DB_ERASE_ALL_DATA", false),
		MigrateModels:  envVar("DB_MIGRATE_MODELS", true),
		InsertAdmin:    envVar("DB_INSERT_ADMIN", true),
		InsertAdminPwd: envVar("DB_INSERT_ADMIN_PWD", ""),
		LogLevel:       int(DBLogLevels[envVar("DB_LOG_LEVEL", "error")]),
	}
}

/* -~-~-~-~ Logger Config ~-~-~-~- */

type LoggerCfg struct {
	Level       int
	LevelStackT int
	LogCaller   bool
}

func loadLoggerConfig() LoggerCfg {
	return LoggerCfg{
		Level:       LogLevels[envVar("LOGGER_LEVEL", "info")],
		LevelStackT: LogLevels[envVar("LOGGER_LEVEL_STACKTRACE", "fatal")],
		LogCaller:   envVar("LOGGER_LOG_CALLER", false),
	}
}

/* -~-~-~-~ JWT Config ~-~-~-~- */

type JWTCfg struct {
	Secret      string
	SessionDays int
}

func loadJWTConfig() JWTCfg {
	return JWTCfg{
		Secret:      envVar("JWT_SECRET", ""),
		SessionDays: envVar("JWT_SESSION_DAYS", 7),
	}
}

/* -~-~-~-~ TLS Config ~-~-~-~- */

// For TLSEnabled, see [Globals].
type TLSCfg struct {
	CertPath string
	KeyPath  string
}

func loadTLSConfig() TLSCfg {
	return TLSCfg{
		CertPath: envVar("TLS_CERT_PATH", "./server.crt"),
		KeyPath:  envVar("TLS_KEY_PATH", "./server.key"),
	}
}

/* -~-~-~-~ Pwd Hasher Config ~-~-~-~- */

type PwdHasherCfg struct {
	Salt string
}

func loadPwdHasherConfig() PwdHasherCfg {
	return PwdHasherCfg{
		Salt: envVar("PWD_HASHER_SALT", ""),
	}
}

/* -~-~-~-~ Retrier Config ~-~-~-~- */

type RetrierCfg struct{}

func loadRetrierConfig() RetrierCfg {
	return RetrierCfg{}
}

/* -~-~-~-~ Rate Limiter Config ~-~-~-~- */

type RLimiterCfg struct {
	MaxTokens       int
	TokensPerSecond int
}

func loadRateLimiterConfig() RLimiterCfg {
	return RLimiterCfg{
		MaxTokens:       envVar("RLIMITER_MAX_TOKENS", 40),
		TokensPerSecond: envVar("RLIMITER_TOKENS_PER_SECOND", 10),
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func envVar[T string | bool | int](key string, fallback T) T {
	val, isSet := os.LookupEnv(key)
	if !isSet {
		return fallback
	}

	switch any(fallback).(type) {
	case string:
		return any(val).(T)
	case bool:
		return any(strings.ToLower(val) == "true" || val == "1").(T)
	case int:
		if intVal, err := strconv.Atoi(val); err == nil {
			return any(intVal).(T)
		}
	}
	log.Printf("🚨 WARNING: Env var %s set to an unsupported value: %s", key, val)
	return fallback
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Used to connect to the DB.
func (c *DBCfg) GetSQLConnString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", c.Username, c.Password, c.Hostname, c.Port, c.Database, c.Params)
}

// Used if the DB we need is not yet created.
func (c *DBCfg) GetSQLConnStringNoDB() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.Username, c.Password, c.Hostname, c.Port, c.Params)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

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
	"debug":  4, // Info
	"info":   4, // Info
	"warn":   3, // Warn
	"error":  2, // Error
	"fatal":  2, // Error
	"silent": 1, // Silent
}

// TODO -> Be able to change some of these config values at runtime.
