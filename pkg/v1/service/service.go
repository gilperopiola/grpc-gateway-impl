package service

import (
	"context"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
)

/* ----------------------------------- */
/*            - Service -              */
/* ----------------------------------- */

// Service is the interface that wraps the service methods. All the business logic should be implemented here.
type Service interface {
	Signup(ctx context.Context, in *usersPB.SignupRequest) (*usersPB.SignupResponse, error)
	Login(ctx context.Context, in *usersPB.LoginRequest) (*usersPB.LoginResponse, error)
}

// service is our concrete implementation of the Service interface.
// This struct should usually contain a Repository interface inside.
type service struct{}

// NewService returns a new instance of the service.
func NewService() *service {
	return &service{}
}
