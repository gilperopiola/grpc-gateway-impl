package service

import (
	"context"
	"math/rand"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/* ----------------------------------- */
/*          - Users Service -          */
/* ----------------------------------- */

// Signup should be the implementation of the Signup service method.
func (s *service) Signup(ctx context.Context, in *usersPB.SignupRequest) (*usersPB.SignupResponse, error) {

	// ... check username is available, hash password, create user in DB, etc.

	// if err := something(); err != nil {
	// 		return entities.SignupResponse{}, fmt.Errorf("error in something(): %w", err)
	// }

	// Simulate a random error sometimes.
	if rand.Intn(4) == 1 {
		return nil, status.Error(codes.Internal, "error creating user.")
	}

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

	// Simulate a random error sometimes.
	if rand.Intn(4) == 1 {
		return nil, status.Error(codes.Internal, "error logging in user.")
	}

	return &usersPB.LoginResponse{
		Token: "some.jwt.token",
	}, nil
}
