package v1

import (
	"context"
	"math/rand"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// shouldReturnInternalError helps simulate a random internal error. Returns true 1 out of 4 times.
func shouldReturnInternalError() bool {
	return rand.Intn(4) == 1
}

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
	if shouldReturnInternalError() {
		return nil, status.Error(codes.Internal, "error creating user.")
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
	if shouldReturnInternalError() {
		return nil, status.Error(codes.Internal, "error logging in user.")
	}

	return loginOKResponse("some.jwt.token")
}

func loginOKResponse(token string) (*usersPB.LoginResponse, error) {
	return &usersPB.LoginResponse{Token: token}, nil
}
