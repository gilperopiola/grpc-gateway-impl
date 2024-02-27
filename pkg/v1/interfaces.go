package v1

import (
	"context"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// These are our custom types for every dependency we need.
// They are used to make the code more readable and to avoid having to import the actual types in this file.
type MiddlewareI []runtime.ServeMuxOption
type InterceptorsI []grpc.ServerOption
type GRPCDialOptionsI []grpc.DialOption

type ServiceLayer interface {
	Signup(ctx context.Context, in *usersPB.SignupRequest) (*usersPB.SignupResponse, error)
	Login(ctx context.Context, in *usersPB.LoginRequest) (*usersPB.LoginResponse, error)
}
