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
// It implements a handler for each API method, connecting it with the Service.
type API struct {
	// Service is the business logic layer of our API. It should have 1 method for each API endpoint.
	Service service.ServiceLayer

	// usersPB.UnimplementedUsersServiceServer is an empty struct that we embed to satisfy the usersPB.UsersServiceServer interface.
	usersPB.UnimplementedUsersServiceServer
}

/* ----------------------------------- */
/*            - Handlers -             */
/* ----------------------------------- */

// Signup is the handler / entrypoint for the Signup API method. Both gRPC and HTTP.
func (api *API) Signup(ctx context.Context, in *usersPB.SignupRequest) (*usersPB.SignupResponse, error) {
	return api.Service.Signup(ctx, in)
}

// Login is the handler / entrypoint for the Login API method. Both gRPC and HTTP.
func (api *API) Login(ctx context.Context, in *usersPB.LoginRequest) (*usersPB.LoginResponse, error) {
	return api.Service.Login(ctx, in)
}
