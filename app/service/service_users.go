package service

import (
	"context"

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
	user, err := s.GetDBTool().GetUser(ctx, sql.WithUsername(req.Username))
	if err == nil && user != nil {
		return nil, errs.GRPCAlreadyExists("user")
	}
	if !errIsNotFound(err) {
		return nil, errInUsersDB(ctx, err)
	}

	// If we're here, we should have gotten a gorm.ErrRecordNotFound in the function above.
	if user, err = s.GetDBTool().CreateUser(ctx, req.Username, s.Toolbox.HashPassword(req.Password)); err != nil {
		return nil, errInUsersDB(ctx, err)
	}

	return &pbs.SignupResponse{Id: int32(user.ID)}, nil
}

// Login first tries to get the user with the given username.
// If the query fails (with a gorm.ErrRecordNotFound), then that user doesn't exist.
// If the query fails (for some other reason), then we return an unknown error.
// Then we PasswordsMatch both passwords. If they don't match, we return an unauthenticated error.
// If everything is OK, we generate a token and return it.
func (s *service) Login(ctx context.Context, req *pbs.LoginRequest) (*pbs.LoginResponse, error) {
	user, err := s.GetDBTool().GetUser(ctx, sql.WithUsername(req.Username))
	if errIsNotFound(err) {
		return nil, errs.GRPCNotFound("user")
	}
	if err != nil || user == nil {
		return nil, errInUsersDB(ctx, err)
	}

	if !s.Toolbox.PasswordsMatch(req.Password, user.Password) {
		return nil, errs.GRPCUnauthenticated()
	}

	token, err := s.Toolbox.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errs.GRPCGeneratingToken(err)
	}

	return &pbs.LoginResponse{Token: token}, nil
}

// GetUser first tries to get the user with the given ID.
// If the query fails (with a gorm.ErrRecordNotFound), then that user doesn't exist.
// If the query fails (for some other reason), then it returns an unknown error.
// If everything is OK, it returns the user.
func (s *service) GetUser(ctx context.Context, req *pbs.GetUserRequest) (*pbs.GetUserResponse, error) {
	user, err := s.GetDBTool().GetUser(ctx, sql.WithUserID(req.UserId))
	if errIsNotFound(err) {
		return nil, errs.GRPCNotFound("user")
	}
	if err != nil || user == nil {
		return nil, errInUsersDB(ctx, err)
	}
	return &pbs.GetUserResponse{User: user.ToUserInfoPB()}, nil
}

// GetUsers first gets the page, pageSize and filterQueryOptions from the request.
// With those values, it gets the users from the database. If there's an error, it returns unknown.
// If everything is OK, it returns the users and the pagination info.
func (s *service) GetUsers(ctx context.Context, req *pbs.GetUsersRequest) (*pbs.GetUsersResponse, error) {
	page, pageSize := getPaginationFromRequest(req)
	usernameFilter := sql.WithCondition(sql.Like, "username", req.GetFilter())

	// While our page is 0-based, gorm offsets are 1-based. That's why we subtract 1.
	users, totalMatches, err := s.GetDBTool().GetUsers(ctx, page-1, pageSize, usernameFilter)
	if err != nil {
		return nil, errInUsersDB(ctx, err)
	}

	return &pbs.GetUsersResponse{
		Users:      users.ToUsersInfoPB(),
		Pagination: newResponsePagination(page, pageSize, totalMatches),
	}, nil
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var errInUsersDB = func(ctx context.Context, err error) error {
	core.LogUnexpectedErr(err)
	return errs.GRPCUsersDBCall(err, core.RouteNameFromCtx(ctx))
}
