package cfg

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

/* ----------------------------------- */
/*             - Config -              */
/* ----------------------------------- */

// Config holds the configuration of the entire API.
type Config struct {
	*MainConfig
	*DBConfig
	*JWTConfig
	*TLSConfig
	*LoggerConfig
	*RateLimiterConfig
}

// Load sets up the configuration from the environment variables.
func Load() *Config {

	// The project is either run from the root folder or the /cmd folder.
	// If it's run from /cmd, we add a '..' prefix to the filesystem paths to move them back to the root folder.
	// Otherwise, we just add a '.', staying on the root.
	projectName := envString("PROJECT_NAME", "grpc-gateway-impl")
	wdPrefix := getWdPrefix(projectName)

	developmentDBAdminPass := "n8zAyv96oAtfQoNof-_ulH4pS0Dqf61VThTZbbOLXCU=" // T0D0 remove.

	return &Config{
		MainConfig: &MainConfig{
			ProjectName: projectName,

			IsProd: envBulean("IS_PROD", false),
			IsDev:  !envBulean("IS_PROD", false),

			GRPCPort: envString("GRPC_PORT", ":50053"),
			HTTPPort: envString("HTTP_PORT", ":8083"),
			HashSalt: envString("HASH_SALT", "s0m3_s4l7"), // Used to hash passwords.
		},
		DBConfig: &DBConfig{
			Username: envString("DB_USERNAME", "root"),
			Password: envString("DB_PASSWORD", ""),
			Hostname: envString("DB_HOSTNAME", "localhost"),
			Port:     envString("DB_PORT", "3306"),
			Schema:   envString("DB_SCHEMA", "grpc-gateway-impl"),
			Params:   envString("DB_PARAMS", "?charset=utf8&parseTime=True&loc=Local"),

			Migrate:     envBulean("DB_MIGRATE", true),
			InsertAdmin: envBulean("DB_INSERT_ADMIN", true),
			AdminPwd:    envString("DB_ADMIN_PASSWORD", developmentDBAdminPass), // Should already be hashed with our salt.

			GormLogLevel: gormLevelMap[envString("DB_LOG_LEVEL", "error")],
		},
		JWTConfig: &JWTConfig{
			Secret:      envString("JWT_SECRET", "please_set_the_env_var"),
			SessionDays: envNumber("JWT_SESSION_DAYS", 7),
		},
		TLSConfig: &TLSConfig{
			Enabled:  envBulean("TLS_ENABLED", false),
			CertPath: envString("TLS_CERT_PATH", wdPrefix+"/server.crt"),
			KeyPath:  envString("TLS_KEY_PATH", wdPrefix+"/server.key"),
		},
		LoggerConfig: &LoggerConfig{
			Level:           zapLevelMap[envString("LOG_LEVEL", "info")],
			StacktraceLevel: zapLevelMap[envString("LOG_STACKTRACE_LEVEL", "dpanic")],
			LogCaller:       envBulean("LOG_CALLER", false),
		},
		RateLimiterConfig: &RateLimiterConfig{
			MaxTokens:       envNumber("RATE_LIMITER_MAX_TOKENS", 40),
			TokensPerSecond: envNumber("RATE_LIMITER_TOKENS_PER_SECOND", 10),
		},
	}
}

// Main configuration.
type MainConfig struct {
	ProjectName string

	IsProd bool
	IsDev  bool

	GRPCPort string
	HTTPPort string

	HashSalt string
}

// DB configuration.
type DBConfig struct {
	Username string
	Password string
	Hostname string
	Port     string
	Schema   string
	Params   string

	Migrate     bool
	InsertAdmin bool
	AdminPwd    string // Should already be hashed with our salt.

	GormLogLevel int // options are in the gormLevelMap var.
}

var gormLevelMap = map[string]int{
	"silent": int(logger.Silent),
	"error":  int(logger.Error),
	"warn":   int(logger.Warn),
	"info":   int(logger.Info),
}

// JWT Auth configuration.
type JWTConfig struct {
	Secret      string
	SessionDays int
}

// TLS configuration.
type TLSConfig struct {
	Enabled  bool // If enabled, use TLS between HTTP and gRPC.
	CertPath string
	KeyPath  string
}

// Logger configuration.
type LoggerConfig struct {
	Level           int // options are in the zapLevelMap var.
	StacktraceLevel int // options are in the zapLevelMap var.
	LogCaller       bool
}

var zapLevelMap = map[string]int{
	"debug":  int(zap.DebugLevel),
	"info":   int(zap.InfoLevel),
	"warn":   int(zap.WarnLevel),
	"error":  int(zap.ErrorLevel),
	"dpanic": int(zap.DPanicLevel),
	"panic":  int(zap.PanicLevel),
	"fatal":  int(zap.FatalLevel),
}

// Rate Limiter configuration.
type RateLimiterConfig struct {
	MaxTokens       int // Max tokens the bucket can hold.
	TokensPerSecond int // Tokens reloaded per second.
}

/* ----------------------------------- */
/*            - Helpers -              */
/* ----------------------------------- */

// envString returns the value of an env var or a fallback value if it doesn't exist.
func envString(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// envBulean returns the value of an env var as a boolean or a fallback value if it doesn't exist.
func envBulean(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		return value == "true" || value == "TRUE" || value == "1"
	}
	return fallback
}

// envNumber returns the value of an env var as an int or a fallback value if it doesn't exist.
func envNumber(key string, fallback int) int {
	if value, err := strconv.Atoi(envString(key, "")); err == nil {
		return value
	}
	return fallback
}

// getWdPrefix returns the prefix that needs to be added to the default paths to start at the root folder.
func getWdPrefix(projectName string) string {
	if isWorkingDirRootFolder(projectName) {
		return "."
	}
	return ".."
}

// isWorkingDirRootFolder returns true if the working directory is the root folder.
// The project name is used to determine where we are.
func isWorkingDirRootFolder(projectName string) bool {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf(errs.FatalErrMsgGettingWorkingDir, err) // Our logger is not initialized yet.
	}
	return strings.HasSuffix(workingDir, projectName)
}
