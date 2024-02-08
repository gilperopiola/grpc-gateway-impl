package service

import (
	"context"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/entities"
)

/* ----------------------------------- */
/*            - Service -              */
/* ----------------------------------- */

// ServiceLayer is the interface that wraps the service methods. All the business logic should be implemented here.
// Note that the methods receive and return entities, not PBs.
type ServiceLayer interface {
	Signup(ctx context.Context, in entities.SignupRequest) (entities.SignupResponse, error)
	Login(ctx context.Context, in entities.LoginRequest) (entities.LoginResponse, error)
}

// service is our concrete implementation of the ServiceLayer interface.
type service struct{}

// NewService returns a new instance of the service.
func NewService() *service {
	return &service{}
}

// Signup should be the implementation of the Signup service method.
func (s *service) Signup(ctx context.Context, in entities.SignupRequest) (entities.SignupResponse, error) {

	// ... check if user exists, hash password, create user, etc.

	// if err := something(); err != nil {
	// 		return entities.SignupResponse{}, fmt.Errorf("some error: %w", err)
	// }

	generatedUserID := 5

	return entities.SignupResponse{ID: generatedUserID}, nil
}

// Login should be the implementation of the Login service method.
func (s *service) Login(ctx context.Context, in entities.LoginRequest) (entities.LoginResponse, error) {

	// ... get user from DB, hash password, compare passwords, etc.

	// if err := something(); err != nil {
	// 		return entities.LoginResponse{}, fmt.Errorf("some error: %w", err)
	// }

	generatedJWTToken := "some.jwt.token"

	return entities.LoginResponse{Token: generatedJWTToken}, nil
}
