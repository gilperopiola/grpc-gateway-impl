package service

import (
	"github.com/gilperopiola/god"
	sql "github.com/gilperopiola/grpc-gateway-impl/app/clients/dbs/sqldb"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/logs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/utils"
)

type UsersSubService struct {
	pbs.UnimplementedUsersServiceServer
	Clients core.Clients
	Tools   core.Tools
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - Users Service -          */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// GetUser first tries to get the user with the given ID.
// If the query fails (with a gorm.ErrRecordNotFound), then that user doesn't exist.
// If the query fails (for some other reason), then it returns an unknown error.
// If everything is OK, it returns the user info.
func (s *UsersSubService) GetUser(ctx god.Ctx, req *pbs.GetUserRequest) (*pbs.GetUserResponse, error) {
	user, err := s.Clients.DBGetUser(ctx, sql.WithID(req.UserId))
	if utils.IsNotFound(err) {
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
	users, totalMatches, err := s.Clients.DBGetUsers(ctx, page-1, pageSize, usernameFilterOpt)
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
		route := shared.RouteNameFromCtx(ctx)
		logs.LogUnexpected(err)
		return errs.GRPCFromDB(err, route)
	}
)
