package service

import (
	"context"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
)

/* ----------------------------------- */
/*            - Service -              */
/* ----------------------------------- */

// ServiceLayer is the interface that wraps the service methods. All the business logic should be implemented here.
type ServiceLayer interface {
	Signup(ctx context.Context, in *usersPB.SignupRequest) (*usersPB.SignupResponse, error)
	Login(ctx context.Context, in *usersPB.LoginRequest) (*usersPB.LoginResponse, error)
}

// service is our concrete implementation of the ServiceLayer interface.
// This struct should usually contain a RepositoryLayer interface.
type service struct{}

// NewService returns a new instance of the service.
func NewService() *service {
	return &service{}
}
