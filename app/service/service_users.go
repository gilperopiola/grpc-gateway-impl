package service

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	sql "github.com/gilperopiola/grpc-gateway-impl/app/toolbox/db_tool/sqldb"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - Users Service -          */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// GetUser first tries to get the user with the given ID.
// If the query fails (with a gorm.ErrRecordNotFound), then that user doesn't exist.
// If the query fails (for some other reason), then it returns an unknown error.
// If everything is OK, it returns the user.
func (s *Service) GetUser(ctx core.Ctx, req *pbs.GetUserRequest) (*pbs.GetUserResponse, error) {
	user, err := s.Toolbox.GetUser(ctx, sql.WithID(req.UserId))
	if s.Toolbox.IsNotFound(err) {
		return nil, errUserNotFound()
	}

	return &pbs.GetUserResponse{User: user.ToUserInfoPB()}, nil
}

// GetUsers first gets the page, pageSize and filterQueryOptions from the request.
// With those values, it gets the users from the database. If there's an error, it returns unknown.
// If everything is OK, it returns the users and the pagination info.
func (s *Service) GetUsers(ctx core.Ctx, req *pbs.GetUsersRequest) (*pbs.GetUsersResponse, error) {
	page, pageSize := getPaginationFromRequest(req)
	usernameFilterOpt := sql.WithCondition(sql.Like, "username", req.GetFilter())

	// While our page is 0-based, gorm offsets are 1-based. That's why we subtract 1.
	users, totalMatches, err := s.Toolbox.GetUsers(ctx, page-1, pageSize, usernameFilterOpt)
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
	errUserNotFound      = func() error { return errs.GRPCNotFound("user") }
	errUserAlreadyExists = func() error { return errs.GRPCAlreadyExists("user") }
	errCallingUsersDB    = func(ctx core.Ctx, err error) error {
		return errs.GRPCUsersDBCall(err, core.RouteNameFromCtx(ctx), core.LogUnexpectedErr)
	}
)
