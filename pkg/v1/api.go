package v1

import (
	"os"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	v1Service "github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"
)

/* ----------------------------------- */
/*             - v1 API -              */
/* ----------------------------------- */

// API is our concrete implementation of the gRPC API defined in the .proto files.
// It has a handler for each API method, connecting it with the Service.
// It implements the usersPB.UsersServiceServer interface.
type API struct {
	Service v1Service.ServiceLayer
	usersPB.UnimplementedUsersServiceServer
}

// NewAPI returns a new instance of the API.
func NewAPI(service v1Service.ServiceLayer) *API {
	return &API{Service: service}
}

// NewService returns a new instance of the Service.
func NewService() v1Service.ServiceLayer {
	return v1Service.NewService()
}

/* ----------------------------------- */
/*             - Config -              */
/* ----------------------------------- */

// APIConfig holds the configuration for the server.
type APIConfig struct {
	// GRPCPort is the port where the gRPC server will listen.
	GRPCPort string

	// HTTPPort is the port where the HTTP Gateway will listen.
	HTTPPort string

	// TLSEnabled defines the use of SSL/TLS for the communication
	// between the HTTP Gateway and the gRPC server.
	TLSEnabled bool
}

// LoadConfig loads the configuration from the environment variables.
func LoadConfig() *APIConfig {
	return &APIConfig{
		GRPCPort:   getVar("GRPC_PORT", ":50053"),
		HTTPPort:   getVar("HTTP_PORT", ":8083"),
		TLSEnabled: getVarBool("TLS_ENABLED", false),
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
