package v1

import (
	"context"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"
)

// API is our concrete implementation of the gRPC API defined in the .proto files.
// It implements a handler for each API method, connecting it with the Service.
type API struct {
	Service service.ServiceLayer
	usersPB.UnimplementedUsersServiceServer
}

// Signup is the handler for the Signup API method. Both gRPC and HTTP calls will trigger this method.
func (api *API) Signup(ctx context.Context, in *usersPB.SignupRequest) (*usersPB.SignupResponse, error) {
	return api.Service.Signup(ctx, in)
}

// Login is the handler for the Login API method. Both gRPC and HTTP calls will trigger this method.
func (api *API) Login(ctx context.Context, in *usersPB.LoginRequest) (*usersPB.LoginResponse, error) {
	return api.Service.Login(ctx, in)
}
