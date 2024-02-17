package v1

import "os"

// Config holds the configuration for the server.
type Config struct {
	GRPCPort string
	HTTPPort string
}

// LoadConfig loads the configuration from the environment variables.
func LoadConfig() *Config {
	return &Config{
		GRPCPort: getVar("GRPC_PORT", ":50053"),
		HTTPPort: getVar("HTTP_PORT", ":8083"),
	}
}

// getVar returns the value of an env var or a fallback value if it doesn't exist.
func getVar(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
