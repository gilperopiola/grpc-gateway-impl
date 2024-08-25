package core

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/types"

	"github.com/joho/godotenv"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*                                   - Core: Configuration -                                 */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

/* - We have a Config struct, it's loaded on the App's Setup function and holds many
/*   different sub-configs: APIsCfg, DBCfg, LoggerCfg, etc.
/* - We also have some... 🌍 Global Configs! :)
/* - We use the .env file on the project's root to set the environment, some configs have
/*   defaults but some of them don't. For now we only support 1 file named '.env'.
/* - IMPORTANT: Add any new configs to the .env.example file.
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Standard Configs 🛠
type Config struct {
	APIsCfg      // -> APIs URLs, keys, etc.
	DBCfg        // -> DB Credentials and such
	JWTCfg       // -> JWT Secret
	TLSCfg       // -> TLS Certs paths
	LoggerCfg    // -> Logger settings
	PwdHasherCfg // -> Salt
	RetrierCfg   // -> N° Retries
	RLimiterCfg  // -> Rate settings
}

// Global Configs 🌍
// ~ Access these from A n y w h e r e !!! 🚀
//   with core.G.xxx.
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

var G = &Globals{
	AppName:     "grpc-gateway-impl",
	Version:     "v1.0",
	Env:         "local",
	IsProd:      false,
	IsStage:     false,
	IsDev:       true,
	GRPCPort:    ":50053",
	HTTPPort:    ":8083",
	TLSEnabled:  false,
	LogAPICalls: true,
}

// Loads and returns the Config from the .env file.
// Also sets globals in core.G.
func LoadConfig() *Config {
	if err := godotenv.Overload(); err != nil {
		log.Println("🚨 WARNING: No .env file found.")
	}

	G = loadGlobalConfig()

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

func loadGlobalConfig() *Globals {
	g := new(Globals)
	g.AppName = envVar("APP_NAME", G.AppName)
	g.Version = envVar("APP_VERSION", G.Version)
	g.Env = envVar("ENV", G.Env)
	g.IsProd = G.Env == "prod" || G.Env == "production" || G.Env == "live"
	g.IsStage = G.Env == "stage" || G.Env == "staging" || G.Env == "test"
	g.IsDev = !G.IsProd && !G.IsStage
	g.GRPCPort = envVar("GRPC_PORT", G.GRPCPort)
	g.HTTPPort = envVar("HTTP_PORT", G.HTTPPort)
	g.TLSEnabled = envVar("TLS_ENABLED", G.TLSEnabled)
	g.LogAPICalls = envVar("LOG_API_CALLS", G.LogAPICalls)
	return g
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
	Username string
	Password string
	Hostname string
	Port     string
	Database string
	Params   string
	Retries  int

	EraseAllData   bool
	MigrateModels  bool
	InsertAdmin    bool
	InsertAdminPwd string // hashed with our PwdHasher

	LogLevel int
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
		LogLevel:       types.DBLogLevels[envVar("DB_LOG_LEVEL", "error")],
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
		Level:       types.LogLevels[envVar("LOGGER_LEVEL", "info")],
		LevelStackT: types.LogLevels[envVar("LOGGER_LEVEL_STACKTRACE", "fatal")],
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
	return PwdHasherCfg{Salt: envVar("PWD_HASHER_SALT", "")}
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

// Used to connect to the DB.
func (c *DBCfg) GetSQLConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", c.Username, c.Password, c.Hostname, c.Port, c.Database, c.Params)
}

// Used if the DB we need is not yet created.
func (c *DBCfg) GetSQLConnectionStringNoDB() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.Username, c.Password, c.Hostname, c.Port, c.Params)
}
