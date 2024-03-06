package v1

import (
	"context"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"
)

/* ----------------------------------- */
/*             - v1 API -              */
/* ----------------------------------- */

// API is our concrete implementation of the gRPC API defined in the .proto files.
// It has a handler for each API method, connecting it with the Service.
// It implements the usersPB.UsersServiceServer interface.
type API struct {
	Service service.Service
	usersPB.UnimplementedUsersServiceServer
}

// NewAPI returns a new instance of the API.
func NewAPI(service service.Service) *API {
	return &API{Service: service}
}

/* ----------------------------------- */
/*          - API Handlers -           */
/* ----------------------------------- */

// Signup is the handler / entrypoint for the Signup API method. Both gRPC and HTTP.
func (api *API) Signup(ctx context.Context, in *usersPB.SignupRequest) (*usersPB.SignupResponse, error) {
	return api.Service.Signup(ctx, in)
}

// Login is the handler / entrypoint for the Login API method. Both gRPC and HTTP.
func (api *API) Login(ctx context.Context, in *usersPB.LoginRequest) (*usersPB.LoginResponse, error) {
	return api.Service.Login(ctx, in)
}
