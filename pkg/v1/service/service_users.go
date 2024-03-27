package service

import (
	"context"
	"errors"
	"fmt"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/options"

	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

/* ----------------------------------- */
/*          - Users Service -          */
/* ----------------------------------- */

// Signup first tries to get the user with the given username.
// If the query succeeds, then that user already exists.
// If the query fails (without a gorm.ErrRecordNotFound), then we return an unknown error.
// If the query fails (with a gorm.ErrRecordNotFound), it means everything is OK, so we create the user and return its ID.
func (s *service) Signup(ctx context.Context, req *usersPB.SignupRequest) (*usersPB.SignupResponse, error) {
	user, err := s.Repo.GetUser(options.WithUsername(req.Username))
	if err == nil && user != nil {
		return nil, ErrAlreadyExists("user")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, UserErr(ctx, err)
	}

	if user, err = s.Repo.CreateUser(req.Username, s.PwdHasher.Hash(req.Password)); err != nil {
		return nil, UserErr(ctx, err)
	}

	return &usersPB.SignupResponse{Id: int32(user.ID)}, nil
}

// Login first tries to get the user with the given username.
// If the query fails (with a gorm.ErrRecordNotFound), then that user doesn't exist.
// If the query fails (for some other reason), then we return an unknown error.
// Then we compare both passwords. If they don't match, we return an unauthenticated error.
// If everything is OK, we generate a token and return it.
func (s *service) Login(ctx context.Context, req *usersPB.LoginRequest) (*usersPB.LoginResponse, error) {
	user, err := s.Repo.GetUser(options.WithUsername(req.Username))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("user")
	}
	if err != nil || user == nil {
		return nil, UserErr(ctx, err)
	}

	if !s.PwdHasher.Compare(req.Password, user.Password) {
		return nil, ErrUnauthenticated()
	}

	token, err := s.TokenGenerator.Generate(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, ErrOnTokenGeneration(err)
	}

	return &usersPB.LoginResponse{Token: token}, nil
}

// GetUser first tries to get the user with the given ID.
// If the query fails (with a gorm.ErrRecordNotFound), then that user doesn't exist.
// If the query fails (for some other reason), then it returns an unknown error.
// If everything is OK, it returns the user.
func (s *service) GetUser(ctx context.Context, req *usersPB.GetUserRequest) (*usersPB.GetUserResponse, error) {
	user, err := s.Repo.GetUser(options.WithUserID(int(req.UserId)))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound("user")
	}
	if err != nil {
		return nil, UserErr(ctx, err)
	}
	return &usersPB.GetUserResponse{User: user.ToUserInfo()}, nil
}

// GetUsers first gets the page, pageSize and filterQueryOptions from the request.
// With those values, it gets the users from the database. If there's an error, it returns unknown.
// If everything is OK, it returns the users and the pagination info.
func (s *service) GetUsers(ctx context.Context, req *usersPB.GetUsersRequest) (*usersPB.GetUsersResponse, error) {
	page, pageSize := getPaginationValues(req)
	filter := options.WithFilter("username", req.GetFilter())

	// While our page is 0-based, gorm offsets are 1-based. That's why we subtract 1.
	users, totalMatchingUsers, err := s.Repo.GetUsers(page-1, pageSize, filter)
	if err != nil {
		return nil, UserErr(ctx, err)
	}

	return &usersPB.GetUsersResponse{
		Users:      users.ToUserInfo(),
		Pagination: responsePagination(page, pageSize, totalMatchingUsers),
	}, nil
}

/* ----------------------------------- */
/*         - Service Errors -          */
/* ----------------------------------- */

var (
	UserErr = func(ctx context.Context, err error) error {
		return errs.NewGRPC(err, codes.Unknown, getGRPCMethodFrom(ctx))
	}

	ErrOnTokenGeneration = func(err error) error {
		return errs.NewGRPC(err, codes.Unknown)
	}

	ErrUnauthenticated = func() error {
		return errs.NewGRPC(nil, codes.Unauthenticated)
	}

	ErrNotFound = func(entity string) error {
		err := fmt.Errorf("%s not found", entity)
		return errs.NewGRPC(err, codes.NotFound)
	}

	ErrAlreadyExists = func(entity string) error {
		err := fmt.Errorf("%s already exists", entity)
		return errs.NewGRPC(err, codes.AlreadyExists)
	}
)
