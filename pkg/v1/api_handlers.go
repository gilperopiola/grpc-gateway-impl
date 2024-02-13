package v1

import (
	"context"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
)

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
