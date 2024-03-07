package v1

import (
	"context"
	"fmt"
	"math/rand"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
)

/* ----------------------------------- */
/*           - v1 Service -            */
/* ----------------------------------- */

// Service is the interface that defines the methods of the API.
type Service interface {
	usersPB.UsersServiceServer
}

// Service is our concrete implementation of the gRPC API defined in the .proto files.
type service struct {
	*usersPB.UnimplementedUsersServiceServer
}

// NewService returns a new instance of the Service.
func NewService() *service {
	return &service{}
}

/* ----------------------------------- */
/*        - Service Handlers -         */
/* ----------------------------------- */

// Signup is the handler / entrypoint for the Signup API method. Both gRPC and HTTP.
func (s *service) Signup(ctx context.Context, in *usersPB.SignupRequest) (*usersPB.SignupResponse, error) {

	// ... check username is available, hash password, create user in DB, etc.

	// if err := something(); err != nil {
	// 		return nil, fmt.Errorf("error in something(): %w", err)
	// }

	// Simulate a random error sometimes.
	if err := simulateRandomErr(); err != nil {
		return nil, fmt.Errorf("Service.Signup: %w", err)
	}

	return signupOKResponse(rand.Intn(1000))
}

func signupOKResponse(id int) (*usersPB.SignupResponse, error) {
	return &usersPB.SignupResponse{Id: int32(id)}, nil
}

// Login is the handler / entrypoint for the Login API method. Both gRPC and HTTP.
func (s *service) Login(ctx context.Context, in *usersPB.LoginRequest) (*usersPB.LoginResponse, error) {

	// ... get user from DB, hash password, compare passwords, etc.

	// if err := something(); err != nil {
	// 		return nil, fmt.Errorf("error in something(): %w", err)
	// }

	// Simulate a random error sometimes.
	if err := simulateRandomErr(); err != nil {
		return nil, fmt.Errorf("Service.Login: %w", err)
	}

	return loginOKResponse("some.jwt.token")
}

func loginOKResponse(token string) (*usersPB.LoginResponse, error) {
	return &usersPB.LoginResponse{Token: token}, nil
}

// simulateRandomErr returns an error 1 out of 5 times.
func simulateRandomErr() error {
	if rand.Intn(5) == 1 {
		return fmt.Errorf("random error")
	}
	return nil
}
