package cfg

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"
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
	*RateLimiterConfig
}

// Main configuration.
type MainConfig struct {
	ProjectName string
	IsProd      bool
	GRPCPort    string
	HTTPPort    string
	HashSalt    string
}

// DB configuration.
type DBConfig struct {
	Username string
	Password string
	Hostname string
	Port     string
	Schema   string
	Params   string

	Migrate       bool
	InsertAdmin   bool
	AdminPassword string // T0D0 - This doesn't feel right. Hash it beforehand.
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

// Rate Limiter configuration.
type RateLimiterConfig struct {
	MaxTokens       int // Max tokens the bucket can hold.
	TokensPerSecond int // Tokens reloaded per second.
}

// Load sets up the configuration from the environment variables.
func Load() *Config {

	// The project is either run from the root folder or the /cmd folder.
	// If it's run from /cmd, we add a '..' prefix to the filesystem paths to move them back to the root folder.
	// Otherwise, we just add a '.', staying on the root.
	projectName := envStr("PROJECT_NAME", "grpc-gateway-impl")
	filesPrefix := getPathPrefix(projectName)

	return &Config{
		MainConfig: &MainConfig{
			ProjectName: projectName,
			IsProd:      envBool("IS_PROD", false),
			GRPCPort:    envStr("GRPC_PORT", ":50053"),
			HTTPPort:    envStr("HTTP_PORT", ":8083"),
			HashSalt:    envStr("HASH_SALT", "s0m3_s4l7"), // This is used to hash passwords.
		},
		DBConfig: &DBConfig{
			Username:      envStr("DB_USERNAME", "root"),
			Password:      envStr("DB_PASSWORD", ""),
			Hostname:      envStr("DB_HOSTNAME", "localhost"),
			Port:          envStr("DB_PORT", "3306"),
			Schema:        envStr("DB_SCHEMA", "grpc-gateway-impl"),
			Params:        envStr("DB_PARAMS", "?charset=utf8&parseTime=True&loc=Local"),
			Migrate:       envBool("DB_MIGRATE", true),
			InsertAdmin:   envBool("DB_INSERT_ADMIN", true),
			AdminPassword: envStr("DB_ADMIN_PASSWORD", "please_set_the_env_var"), // This gets hashed before being used.
		},
		JWTConfig: &JWTConfig{
			Secret:      envStr("JWT_SECRET", "please_set_the_env_var"),
			SessionDays: envInt("JWT_SESSION_DAYS", 7),
		},
		TLSConfig: &TLSConfig{
			Enabled:  envBool("TLS_ENABLED", false),
			CertPath: envStr("TLS_CERT_PATH", filesPrefix+"/server.crt"),
			KeyPath:  envStr("TLS_KEY_PATH", filesPrefix+"/server.key"),
		},
		RateLimiterConfig: &RateLimiterConfig{
			MaxTokens:       envInt("RATE_LIMITER_MAX_TOKENS", 40),
			TokensPerSecond: envInt("RATE_LIMITER_TOKENS_PER_SECOND", 10),
		},
	}
}

/* ----------------------------------- */
/*            - Helpers -              */
/* ----------------------------------- */

// envStr returns the value of an env var or a fallback value if it doesn't exist.
func envStr(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// envBool returns the value of an env var as a boolean or a fallback value if it doesn't exist.
func envBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		return value == "true" || value == "TRUE" || value == "1"
	}
	return fallback
}

// envInt returns the value of an env var as an int or a fallback value if it doesn't exist.
func envInt(key string, fallback int) int {
	if value, err := strconv.Atoi(envStr(key, "")); err == nil {
		return value
	}
	return fallback
}

// getPathPrefix returns the prefix that needs to be added to the default paths to start at the root folder.
func getPathPrefix(projectName string) string {
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
