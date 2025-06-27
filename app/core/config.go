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

		MockCalls bool
		MockData  map[string]string
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
	mockAPICalls := envVar("API_MOCK_CALLS", false)

	return APIsCfg{
		Weather: OpenWeatherMapAPICfg{
			BaseURL: envVar("API_OPENWEATHERMAP_BASE_URL", "https://api.weathermap.org/data/2.5/weather"),
			AppID:   envVar("API_OPENWEATHERMAP_APP_ID", ""),
		},
		GPT: ChatGptAPICfg{
			BaseURL: envVar("API_CHATGPT_BASE_URL", "https://api.openai.com/v1"),
			APIKey:  envVar("API_CHATGPT_API_KEY", ""),
		},
		MockCalls: mockAPICalls,
		MockData: func() map[string]string {
			if !mockAPICalls {
				return nil
			}

			// Keys must match a substring of the actual API URL to trigger the mock response.
			mockData := make(map[string]string)
			mockData["https://api.openai.com/v1/chat/completions"] = `
{
    "id": "chatcmpl-...",
    "object": "chat.completion",
    "created": 1250814434,
    "model": "gpt-4o-2024-08-06",
    "choices": [
		{"index": 0, "message": {"role": "assistant","content": "Mock response gotten","refusal": null, "annotations": []}, "logprobs": null, "finish_reason": "stop"}
    ],
    "usage": {
        "prompt_tokens": 31,
        "completion_tokens": 604,
        "total_tokens": 635,
        "prompt_tokens_details": {"cached_tokens": 0,"audio_tokens": 0},
        "completion_tokens_details": {"reasoning_tokens": 0,"audio_tokens": 0,"accepted_prediction_tokens": 0,"rejected_prediction_tokens": 0}
    },
    "service_tier": "default",
    "system_fingerprint": "fp_0..."
}`

			mockData["https://api.openai.com/v1/images/generations"] = `
{
    "created": 175033,
    "data": [
        {
            "revised_prompt": "Create a widescreen, photorealistic image in ultra HD quality with rich colors, mimicking the style of the kitsch era. The image should depict a dark starry sky that stretches across the entire breadth of the image. The sky should be teeming with fiery meteorites, each leaving behind a trail of sparkling dust, accentuated with flamboyant lens flares. The colors should be vibrant and intense to capture the essence of the kitsch aesthetic.",
            "url": "https://oaidalleapiprodscus.blob.core.windows.net/private/org-QSC7lVI62AyfPgP7NyIh4OTS/user-T2w0qXeRmFDqwUgIrWWRSiZv/img-VeoC7pt5ZGHUgVQsCxg3daD8.png?st=2025-06-25T00%3A27%3A50Z&se=2025-06-25T02%3A27%3A50Z&sp=r&sv=2024-08-04&sr=b&rscd=inline&rsct=image/png&skoid=cc612491-d948-4d2e-9821-2683df3719f5&sktid=a48cca56-e6da-484e-a814-9c849652bcb3&skt=2025-06-24T18%3A52%3A08Z&ske=2025-06-25T18%3A52%3A08Z&sks=b&skv=2024-08-04&sig=9X6F7IJgLO%2BghFPa/0kZ4FiE55y0j0/mXTNU%2BHxUbak%3D"
        }
    ]
}`
			return mockData
		}(),
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
