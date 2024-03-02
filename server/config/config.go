package config

import (
	"log"
	"os"
	"strings"

	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"
)

const (
	projectName = "grpc-gateway-impl" // Used to check if the working directory is the root folder.
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
	IsProd   bool
	GRPCPort string
	HTTPPort string
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
	workingDirIsRootFolder := isWorkingDirRootFolder()
	workingDirPathPrefix := getPathPrefix(workingDirIsRootFolder)

	return &Config{
		MainConfig: &MainConfig{
			IsProd:   getVarBool("IS_PROD", false),
			GRPCPort: getVar("GRPC_PORT", ":50053"),
			HTTPPort: getVar("HTTP_PORT", ":8083"),
		},
		TLS: TLSConfig{
			Enabled:  getVarBool("TLS_ENABLED", false),
			CertPath: getVar("TLS_CERT_PATH", workingDirPathPrefix+"/server.crt"),
			KeyPath:  getVar("TLS_KEY_PATH", workingDirPathPrefix+"/server.key"),
		},
		RateLimiter: RateLimiterConfig{
			MaxTokens:       20,
			TokensPerSecond: 5,
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

// isWorkingDirRootFolder returns true if the working directory is the /cmd folder.
func isWorkingDirRootFolder() bool {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf(v1.FatalErrMsgGettingWorkingDir, err)
	}
	return strings.HasSuffix(workingDir, projectName)
}

// getPathPrefix returns the prefix that needs to be added to the default paths
// so that we always start at the root folder.
func getPathPrefix(workingDirIsRootFolder bool) string {
	if workingDirIsRootFolder {
		return "."
	}
	return ".."
}
