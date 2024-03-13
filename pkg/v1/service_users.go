package v1

import (
	"context"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/db"

	"gorm.io/gorm"
)

const (
	errGettingUser  = "error getting user"
	errCreatingUser = "error creating user"
)

/* ----------------------------------- */
/*          - Users Service -          */
/* ----------------------------------- */

// Signup is the entrypoint of the Signup API method.
func (s *service) Signup(ctx context.Context, in *usersPB.SignupRequest) (*usersPB.SignupResponse, error) {
	user, err := s.DB.GetUser(0, in.Username)
	if (err != nil && err != gorm.ErrRecordNotFound) || user == nil {
		return nil, grpcUnknownErr(errGettingUser, err)
	}

	if user.Username != "" {
		return nil, grpcAlreadyExistsErr("user")
	}

	if user, err = s.DB.CreateUser(in.Username, s.PwdHasher.Hash(in.Password)); err != nil {
		return nil, grpcUnknownErr(errCreatingUser, err)
	}

	return &usersPB.SignupResponse{Id: int32(user.ID)}, nil
}

// Login is the entrypoint of the Login API method.
func (s *service) Login(ctx context.Context, in *usersPB.LoginRequest) (*usersPB.LoginResponse, error) {
	user, err := s.DB.GetUser(0, in.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, grpcUnknownErr(errGettingUser, err)
	}

	if (user != nil && user.Username == "") || !s.PwdHasher.Compare(in.Password, user.Password) {
		return nil, grpcUnauthenticatedErr("wrong credentials.")
	}

	token, err := s.TokenGenerator.Generate(user.ID, user.Username, user.Role)
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
		return nil, grpcUnknownErr(errGettingUser, err)
	}

	return &usersPB.GetUserResponse{User: user.ToUserInfo()}, nil
}

// GetUsers is the entrypoint of the GetUsers API method.
func (s *service) GetUsers(ctx context.Context, in *usersPB.GetUsersRequest) (*usersPB.GetUsersResponse, error) {

	// Page starts at 1 and pageSize at 10.
	page := max(in.GetPage(), 1)
	pageSize := max(in.GetPageSize(), 10)
	filter := in.GetFilter()

	// DB is 0-indexed.
	users, totalPages, err := s.DB.GetUsers(int(page-1), int(pageSize), filter)
	if err != nil {
		return nil, grpcUnknownErr("error getting users", err)
	}

	// If the page is greater than the total pages, return an invalid argument error.
	if page > 1 && page > int32(totalPages) {
		return nil, grpcInvalidArgumentErr("page")
	}

	return &usersPB.GetUsersResponse{
		Users:      db.Users(users).ToUserInfo(),
		Pagination: &usersPB.PaginationInfo{Current: int32(page), Total: int32(totalPages)},
	}, nil
}

func max(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
