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
	*TLSConfig
	*RateLimiterConfig
}

// MainConfig holds the main configuration settings of the API.
type MainConfig struct {
	ProjectName string

	IsProd bool

	GRPCPort string
	HTTPPort string
}

// TLS configuration.
type TLSConfig struct {
	// Enabled defines the use of SSL/TLS for the communication between the servers.
	Enabled bool

	// CertPath & KeyPath are the paths to the SSL/TLS certificate and key files.
	CertPath string
	KeyPath  string
}

// RateLimiter configuration.
type RateLimiterConfig struct {
	MaxTokens       int // MaxTokens is the maximum number of tokens the bucket can hold.
	TokensPerSecond int // TokensPerSecond is the number of tokens reloaded per second.
}

// Init loads the configuration from the environment variables.
func Init() *Config {
	projectName := getVar("PROJECT_NAME", "grpc-gateway-impl")

	// The project is either run from the root folder or the /cmd folder.
	// If it's run from /cmd, we add a '..' prefix to the filesystem paths to move them back to the root folder.
	// Otherwise, we just add a '.', staying on the root.
	filePathPrefix := getPathPrefix(projectName)

	return &Config{
		MainConfig: &MainConfig{
			ProjectName: projectName,
			IsProd:      getVarBool("IS_PROD", false),
			GRPCPort:    getVar("GRPC_PORT", ":50053"),
			HTTPPort:    getVar("HTTP_PORT", ":8083"),
		},
		TLSConfig: &TLSConfig{
			Enabled:  getVarBool("TLS_ENABLED", false),
			CertPath: getVar("TLS_CERT_PATH", filePathPrefix+"/server.crt"),
			KeyPath:  getVar("TLS_KEY_PATH", filePathPrefix+"/server.key"),
		},
		RateLimiterConfig: &RateLimiterConfig{
			MaxTokens:       getVarInt("RATE_LIMITER_MAX_TOKENS", 20),
			TokensPerSecond: getVarInt("RATE_LIMITER_TOKENS_PER_SECOND", 4),
		},
	}
}

// getVar returns the value of an env var or a fallback value if it doesn't exist.
func getVar(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// getVarBool returns the value of an env var as a boolean or a fallback value if it doesn't exist.
func getVarBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		return value == "true" || value == "TRUE" || value == "1"
	}
	return fallback
}

// getVarBool returns the value of an env var as an int or a fallback value if it doesn't exist.
func getVarInt(key string, fallback int) int {
	valueStr := getVar(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
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
func isWorkingDirRootFolder(projectName string) bool {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf(errs.FatalErrMsgGettingWorkingDir, err)
	}

	// The project name is the last part of the root folder path.
	return strings.HasSuffix(workingDir, projectName)
}
