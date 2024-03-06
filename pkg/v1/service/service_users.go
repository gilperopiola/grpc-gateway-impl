package service

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
/*          - Users Service -          */
/* ----------------------------------- */

// - Signup

// Signup should be the implementation of the Signup service method.
func (s *service) Signup(ctx context.Context, in *usersPB.SignupRequest) (*usersPB.SignupResponse, error) {

	// ... check username is available, hash password, create user in DB, etc.

	// if err := something(); err != nil {
	// 		return nil, fmt.Errorf("error in something(): %w", err)
	// }

	// Simulate a random error sometimes.
	if shouldReturnInternalError() {
		return nil, status.Error(codes.Internal, "error creating user.")
	}

	return newSignupOKResponse(rand.Intn(1000))
}

func newSignupOKResponse(id int) (*usersPB.SignupResponse, error) {
	return &usersPB.SignupResponse{Id: int32(id)}, nil
}

// - Login

// Login should be the implementation of the Login service method.
func (s *service) Login(ctx context.Context, in *usersPB.LoginRequest) (*usersPB.LoginResponse, error) {

	// ... get user from DB, hash password, compare passwords, etc.

	// if err := something(); err != nil {
	// 		return nil, fmt.Errorf("error in something(): %w", err)
	// }

	// Simulate a random error sometimes.
	if shouldReturnInternalError() {
		return nil, status.Error(codes.Internal, "error logging in user.")
	}

	return newLoginOKResponse("some.jwt.token")
}

func newLoginOKResponse(token string) (*usersPB.LoginResponse, error) {
	return &usersPB.LoginResponse{Token: token}, nil
}
