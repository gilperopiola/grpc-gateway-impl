package service

import (
	"context"
	"math/rand"

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

// Signup should be the implementation of the Signup service method.
func (s *service) Signup(ctx context.Context, in *usersPB.SignupRequest) (*usersPB.SignupResponse, error) {

	// ... check username is available, hash password, create user in DB, etc.

	// if err := something(); err != nil {
	// 		return entities.SignupResponse{}, fmt.Errorf("error in something(): %w", err)
	// }

	return &usersPB.SignupResponse{
		Id: int32(rand.Intn(1000)),
	}, nil
}

// Login should be the implementation of the Login service method.
func (s *service) Login(ctx context.Context, in *usersPB.LoginRequest) (*usersPB.LoginResponse, error) {

	// ... get user from DB, hash password, compare passwords, etc.

	// if err := something(); err != nil {
	// 		return entities.LoginResponse{}, fmt.Errorf("error in something(): %w", err)
	// }

	return &usersPB.LoginResponse{
		Token: "some.jwt.token",
	}, nil
}
