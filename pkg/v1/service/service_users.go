package service

import (
	"context"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"

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
	user, err := s.Repo.GetUser(signupQueryOptions(in.Username)...)
	if (err != nil && err != gorm.ErrRecordNotFound) || user == nil {
		return nil, grpcUnknownErr(errGettingUser, err)
	}

	if user.Username != "" {
		return nil, grpcAlreadyExistsErr("user")
	}

	if user, err = s.Repo.CreateUser(in.Username, s.PwdHasher.Hash(in.Password)); err != nil {
		return nil, grpcUnknownErr(errCreatingUser, err)
	}

	return &usersPB.SignupResponse{Id: int32(user.ID)}, nil
}

// Login is the entrypoint of the Login API method.
func (s *service) Login(ctx context.Context, in *usersPB.LoginRequest) (*usersPB.LoginResponse, error) {
	user, err := s.Repo.GetUser(loginQueryOptions(in.Username)...)
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
	user, err := s.Repo.GetUser(getUserQueryOptions(in.UserId)...)
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

	// We get page and pageSize from the request.
	page, pageSize := getPaginationValues(in)

	// DB is 0-indexed, our API is 1-indexed.
	// We filter by username field.
	users, totalPages, err := s.Repo.GetUsers(page-1, pageSize, getUsersQueryOptions(in, "username")...)
	if err != nil {
		return nil, grpcUnknownErr("error getting users", err)
	}

	// If the page is greater than the total pages, return an Invalid Argument error.
	if page > 1 && page > totalPages {
		return nil, grpcInvalidArgumentErr("page")
	}

	return &usersPB.GetUsersResponse{
		Users:      users.ToUserInfo(),
		Pagination: &usersPB.PaginationInfo{Current: int32(page), Total: int32(totalPages)},
	}, nil
}
