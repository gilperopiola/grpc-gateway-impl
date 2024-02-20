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
// It implements a handler for each API method, connecting it with the Service.
type API struct {
	Service v1Service.ServiceLayer
	usersPB.UnimplementedUsersServiceServer
}

// NewAPI returns a new instance of the API.
func NewAPI(service v1Service.ServiceLayer) *API {
	return &API{Service: service}
}

/* ----------------------------------- */
/*             - Config -              */
/* ----------------------------------- */

// APIConfig holds the configuration for the server.
type APIConfig struct {
	GRPCPort string
	HTTPPort string
}

// LoadConfig loads the configuration from the environment variables.
func LoadConfig() *APIConfig {
	return &APIConfig{
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
