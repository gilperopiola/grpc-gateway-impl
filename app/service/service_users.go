package service

import (
	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	sql "github.com/gilperopiola/grpc-gateway-impl/app/tools/db_tool/sqldb"
)

type UsersSubService struct {
	pbs.UnimplementedUsersServiceServer
	Tools core.Tools
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - Users Service -          */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// GetUser first tries to get the user with the given ID.
// If the query fails (with a gorm.ErrRecordNotFound), then that user doesn't exist.
// If the query fails (for some other reason), then it returns an unknown error.
// If everything is OK, it returns the user info.
func (s *UsersSubService) GetUser(ctx god.Ctx, req *pbs.GetUserRequest) (*pbs.GetUserResponse, error) {
	user, err := s.Tools.GetUser(ctx, sql.WithID(req.UserId))
	if s.Tools.IsNotFound(err) {
		return nil, errUserNotFound(int(req.UserId))
	}

	return &pbs.GetUserResponse{User: s.Tools.UserToUserInfoPB(user)}, nil
}

// GetUsers first gets the page, pageSize and filterQueryOptions from the request.
// With those values, it gets the users from the database. If there's an error, it returns unknown.
// If everything is OK, it returns the users and the pagination info.
func (s *UsersSubService) GetUsers(ctx god.Ctx, req *pbs.GetUsersRequest) (*pbs.GetUsersResponse, error) {
	page, pageSize := s.Tools.PaginatedRequest(req)
	usernameFilterOpt := sql.WithCondition(sql.Like, "username", req.GetFilter())

	// While our page is 0-based, gorm offsets are 1-based. We subtract 1.
	users, totalMatches, err := s.Tools.GetUsers(ctx, page-1, pageSize, usernameFilterOpt)
	if err != nil {
		return nil, errCallingUsersDB(ctx, err)
	}

	return &pbs.GetUsersResponse{
		Users:      s.Tools.UsersToUsersInfoPB(users),
		Pagination: s.Tools.PaginatedResponse(page, pageSize, totalMatches),
	}, nil
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var (
	errUserNotFound      = func(id int) error { return errs.GRPCNotFound("user", id) }
	errUserAlreadyExists = func() error { return errs.GRPCAlreadyExists("user") }
	errCallingUsersDB    = func(ctx god.Ctx, err error) error {
		route := core.RouteNameFromCtx(ctx)
		core.LogUnexpected(err)
		return errs.GRPCUsersDBCall(err, route)
	}
)
