package v1

import (
	"log"
	"os"
	"strings"
)

/* ----------------------------------- */
/*             - Config -              */
/* ----------------------------------- */

// APIConfig holds the configuration for the server.
type APIConfig struct {
	IsProd   bool
	GRPCPort string
	HTTPPort string

	// TLS holds the configuration for the SSL/TLS connection.
	TLS TLSConfig
}

type TLSConfig struct {
	// Enabled defines the use of SSL/TLS for the communication
	// between the HTTP Gateway and the gRPC server.
	Enabled bool

	// CertPath & KeyPath are the paths to the SSL/TLS certificate and key files.
	CertPath string
	KeyPath  string
}

// LoadConfig loads the configuration from the environment variables.
func LoadConfig() *APIConfig {

	// The project is either run from the root folder or the /cmd folder.
	// If it's run from the /cmd folder, we need to add a '..' prefix to the filesystem paths
	// to make them relative to the root folder.
	// Otherwise, we just add a '.', staying on the current directory.
	workingDirIsCmdFolder := isWorkingDirCmdFolder()
	workingDirPathPrefix := getPathPrefix(workingDirIsCmdFolder)

	return &APIConfig{
		IsProd:   getVarBool("IS_PROD", false),
		GRPCPort: getVar("GRPC_PORT", ":50053"),
		HTTPPort: getVar("HTTP_PORT", ":8083"),

		TLS: TLSConfig{
			Enabled:  getVarBool("TLS_ENABLED", true),
			CertPath: getVar("TLS_CERT_PATH", workingDirPathPrefix+"/server.crt"),
			KeyPath:  getVar("TLS_KEY_PATH", workingDirPathPrefix+"/server.key"),
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

// isWorkingDirCmdFolder returns true if the working directory is the /cmd folder.
func isWorkingDirCmdFolder() bool {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf(msgErrGettingWorkingDir, err)
	}
	return strings.Contains(workingDir, "cmd")
}

// getPathPrefix returns the prefix that needs to be added to the default paths
// so that we always start at the root folder.
func getPathPrefix(workingDirIsCmdFolder bool) string {
	if workingDirIsCmdFolder {
		return ".."
	}
	return "."
}

const (
	msgErrGettingWorkingDir = "Failed to get working directory: %v"
)
