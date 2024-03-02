package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"
)

/* ----------------------------------- */
/*             - Config -              */
/* ----------------------------------- */

// Config holds the configuration for the entire server.
type Config struct {
	*MainConfig // Main holds the main configuration settings of the server.

	TLS         TLSConfig         // TLS holds the configuration of the SSL/TLS connection.
	RateLimiter RateLimiterConfig // RateLimiter holds the configuration of the rate limiter.
}

// MainConfig is the main configuration settings of the server.
type MainConfig struct {
	ProjectName string
	IsProd      bool
	GRPCPort    string
	HTTPPort    string
}

// TLS configuration.
type TLSConfig struct {
	// Enabled defines the use of SSL/TLS for the communication
	// between the HTTP Gateway and the gRPC server.
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

// LoadConfig loads the configuration from the environment variables.
func LoadConfig() *Config {

	// The project is either run from the root folder or the /cmd folder.
	// If it's run from the /cmd folder, we need to add a '..' prefix to the filesystem paths
	// to make them relative to the root folder.
	// Otherwise, we just add a '.', staying on the current directory.
	projectName := getVar("PROJECT_NAME", "grpc-gateway-impl")
	workingDirPathPrefix := getPathPrefix(projectName)

	return &Config{
		MainConfig: &MainConfig{
			ProjectName: projectName,
			IsProd:      getVarBool("IS_PROD", false),
			GRPCPort:    getVar("GRPC_PORT", ":50053"),
			HTTPPort:    getVar("HTTP_PORT", ":8083"),
		},
		TLS: TLSConfig{
			Enabled:  getVarBool("TLS_ENABLED", false),
			CertPath: getVar("TLS_CERT_PATH", workingDirPathPrefix+"/server.crt"),
			KeyPath:  getVar("TLS_KEY_PATH", workingDirPathPrefix+"/server.key"),
		},
		RateLimiter: RateLimiterConfig{
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
		return value == "true" || value == "1" || value == "TRUE"
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

// isWorkingDirRootFolder returns true if the working directory is the /cmd folder.
func isWorkingDirRootFolder(projectName string) bool {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf(v1.FatalErrMsgGettingWorkingDir, err)
	}
	return strings.HasSuffix(workingDir, projectName)
}

// getPathPrefix returns the prefix that needs to be added to the default paths
// so that we always start at the root folder.
func getPathPrefix(projectName string) string {
	if isWorkingDirRootFolder(projectName) {
		return "."
	}
	return ".."
}