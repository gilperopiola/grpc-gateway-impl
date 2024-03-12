package v1

import (
	"context"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/dependencies"

	"gorm.io/gorm"
)

/* ----------------------------------- */
/*          - Users Service -          */
/* ----------------------------------- */

// Signup is the entrypoint of the Signup API method.
func (s *service) Signup(ctx context.Context, in *usersPB.SignupRequest) (*usersPB.SignupResponse, error) {
	user, err := s.DB.GetUser(0, in.Username)
	if (err != nil && err != gorm.ErrRecordNotFound) || user == nil {
		return nil, grpcUnknownErr("error getting user", err)
	}

	if user.Username != "" {
		return nil, grpcAlreadyExistsErr("user")
	}

	if user, err = s.DB.CreateUser(in.Username, s.PwdHasher.Hash(in.Password)); err != nil {
		return nil, grpcUnknownErr("error creating user", err)
	}

	return &usersPB.SignupResponse{Id: int32(user.ID)}, nil
}

// Login is the entrypoint of the Login API method.
func (s *service) Login(ctx context.Context, in *usersPB.LoginRequest) (*usersPB.LoginResponse, error) {
	user, err := s.DB.GetUser(0, in.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, grpcUnknownErr("error getting user", err)
	}

	if (user != nil && user.Username == "") || !s.PwdHasher.Compare(in.Password, user.Password) {
		return nil, grpcUnauthenticatedErr("wrong credentials.")
	}

	token, err := s.TokenGenerator.Generate(user.ID, user.Username, dependencies.DefaultRole)
	if err != nil {
		return nil, grpcUnknownErr("error generating token", err)
	}

	return &usersPB.LoginResponse{Token: token}, nil
}

// GetUser is the entrypoint of the GetUser API method.
func (s *service) GetUser(ctx context.Context, in *usersPB.GetUserRequest) (*usersPB.GetUserResponse, error) {
	user, err := s.DB.GetUser(int(in.UserId), "")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, grpcNotFoundErr("user")
		}
		return nil, grpcUnknownErr("error getting user", err)
	}

	return &usersPB.GetUserResponse{User: user.ToUserInfo()}, nil
}
