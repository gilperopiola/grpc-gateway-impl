package service

import (
	"context"
	"strconv"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	sql "github.com/gilperopiola/grpc-gateway-impl/app/tools/db_tool/sqldb"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - Users Service -          */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Signup first tries to get the user with the given username.
// If the query succeeds, then that user already exists.
// If the query fails (without a gorm.ErrRecordNotFound), then we return an unknown error.
// If the query fails (with a gorm.ErrRecordNotFound), it means everything is OK, so we create the user and return its ID.
func (s *service) Signup(ctx context.Context, req *pbs.SignupRequest) (*pbs.SignupResponse, error) {
	user, err := s.Actions.GetUser(ctx, sql.WithUsername(req.Username))
	if err == nil && user != nil {
		return nil, errUserAlreadyExists()
	}
	if !s.Actions.IsNotFound(err) {
		return nil, errCallingUsersDB(ctx, err)
	}

	// If we're here, we should have gotten a not found in the function above.
	if user, err = s.Actions.CreateUser(ctx, req.Username, s.Actions.HashPassword(req.Password)); err != nil {
		return nil, errCallingUsersDB(ctx, err)
	}

	s.Actions.CreateFolders("data/users/user_" + strconv.Itoa(user.ID))

	s.Actions.CreateGroup(ctx, user.Username+"'s First Group", user.ID)

	return &pbs.SignupResponse{Id: int32(user.ID)}, nil
}

// Login first tries to get the user with the given username.
// If the query fails (with a gorm.ErrRecordNotFound), then that user doesn't exist.
// If the query fails (for some other reason), then we return an unknown error.
// Then we PasswordsMatch both passwords. If they don't match, we return an unauthenticated error.
// If everything is OK, we generate a token and return it.
func (s *service) Login(ctx context.Context, req *pbs.LoginRequest) (*pbs.LoginResponse, error) {
	user, err := s.Actions.GetUser(ctx, sql.WithUsername(req.Username))
	if s.Actions.IsNotFound(err) {
		return nil, errUserNotFound()
	}
	if err != nil || user == nil {
		return nil, errCallingUsersDB(ctx, err)
	}

	if !s.Actions.PasswordsMatch(req.Password, user.Password) {
		return nil, errUnauthenticated()
	}

	token, err := s.Actions.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errGeneratingToken(err)
	}

	return &pbs.LoginResponse{Token: token}, nil
}

// GetUser first tries to get the user with the given ID.
// If the query fails (with a gorm.ErrRecordNotFound), then that user doesn't exist.
// If the query fails (for some other reason), then it returns an unknown error.
// If everything is OK, it returns the user.
func (s *service) GetUser(ctx context.Context, req *pbs.GetUserRequest) (*pbs.GetUserResponse, error) {
	user, err := s.Actions.GetUser(ctx, sql.WithUserID(req.UserId))
	if s.Actions.IsNotFound(err) {
		return nil, errUserNotFound()
	}
	if err != nil || user == nil {
		return nil, errCallingUsersDB(ctx, err)
	}
	return &pbs.GetUserResponse{User: user.ToUserInfoPB()}, nil
}

// GetUsers first gets the page, pageSize and filterQueryOptions from the request.
// With those values, it gets the users from the database. If there's an error, it returns unknown.
// If everything is OK, it returns the users and the pagination info.
func (s *service) GetUsers(ctx context.Context, req *pbs.GetUsersRequest) (*pbs.GetUsersResponse, error) {
	page, pageSize := getPaginationFromRequest(req)
	usernameFilterOpt := sql.WithCondition(sql.Like, "username", req.GetFilter())

	// While our page is 0-based, gorm offsets are 1-based. That's why we subtract 1.
	users, totalMatches, err := s.Actions.GetUsers(ctx, page-1, pageSize, usernameFilterOpt)
	if err != nil {
		return nil, errCallingUsersDB(ctx, err)
	}

	return &pbs.GetUsersResponse{
		Users:      users.ToUsersInfoPB(),
		Pagination: newResponsePagination(page, pageSize, totalMatches),
	}, nil
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var (
	errGeneratingToken   = errs.GRPCGeneratingToken
	errUserNotFound      = func() error { return errNotFound("user") }
	errUserAlreadyExists = func() error { return errAlreadyExists("user") }
	errCallingUsersDB    = func(ctx context.Context, err error) error {
		return errs.GRPCUsersDBCall(err, core.RouteNameFromCtx(ctx), core.LogUnexpectedErr)
	}
)
