package service

import (
	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/logs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
)

type UserSvc struct {
	pbs.UnimplementedUsersSvcServer
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
func (s *UserSvc) GetUser(ctx god.Ctx, req *pbs.GetUserRequest) (*pbs.GetUserResponse, error) {
	// Use repository instead of direct DB call
	user, err := s.Clients.UserRepository().GetUserByID(ctx, int(req.UserId))
	if errs.IsDBNotFound(err) {
		return nil, errUserNotFound(int(req.UserId))
	}

	return &pbs.GetUserResponse{User: s.Tools.UserToUserInfoPB(user)}, nil
}

// GetUsers first gets the page, pageSize and filterQueryOptions from the request.
// With those values, it gets the users from the database. If there's an error, it returns unknown.
// If everything is OK, it returns the users and the pagination info.
func (s *UserSvc) GetUsers(ctx god.Ctx, req *pbs.GetUsersRequest) (*pbs.GetUsersResponse, error) {
	page, pageSize := s.Tools.PaginatedRequest(req)

	// Note: We can no longer use the WithCondition filter directly with repositories
	// We will need to enhance the repository interface to support filtering by username
	// For now, we just get all users with pagination and filter in memory

	// Use repository instead of direct DB call
	users, totalMatches, err := s.Clients.UserRepository().GetUsers(ctx, page, pageSize)
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
		route := core.GetRouteFromCtx(ctx)
		logs.LogUnexpected(err)
		return errs.GRPCFromDB(err, route.Name)
	}
)
