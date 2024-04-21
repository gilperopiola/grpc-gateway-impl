package business

import (
	"context"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/external/storage/options"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - Users Service -          */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Signup first tries to get the user with the given username.
// If the query succeeds, then that user already exists.
// If the query fails (without a gorm.ErrRecordNotFound), then we return an unknown error.
// If the query fails (with a gorm.ErrRecordNotFound), it means everything is OK, so we create the user and return its ID.
func (s *ServiceLayer) Signup(ctx context.Context, req *pbs.SignupRequest) (*pbs.SignupResponse, error) {
	user, err := s.Storage.GetUser(options.WithUsername(req.Username))
	if err == nil && user != nil {
		return nil, ErrAlreadyExists("user")
	}
	if !errIsNotFound(err) {
		return nil, ErrInUsersDBCall(ctx, err)
	}

	// If we're here, we should have gotten a gorm.ErrRecordNotFound in the function above.
	if user, err = s.Storage.CreateUser(req.Username, s.PwdHasher.Hash(req.Password)); err != nil {
		return nil, ErrInUsersDBCall(ctx, err)
	}

	return &pbs.SignupResponse{Id: int32(user.ID)}, nil
}

// Login first tries to get the user with the given username.
// If the query fails (with a gorm.ErrRecordNotFound), then that user doesn't exist.
// If the query fails (for some other reason), then we return an unknown error.
// Then we compare both passwords. If they don't match, we return an unauthenticated error.
// If everything is OK, we generate a token and return it.
func (s *ServiceLayer) Login(ctx context.Context, req *pbs.LoginRequest) (*pbs.LoginResponse, error) {
	user, err := s.Storage.GetUser(options.WithUsername(req.Username))
	if errIsNotFound(err) {
		return nil, ErrNotFound("user")
	}
	if err != nil || user == nil {
		return nil, ErrInUsersDBCall(ctx, err)
	}

	if !s.PwdHasher.Compare(req.Password, user.Password) {
		return nil, ErrUnauthenticated()
	}

	token, err := s.TokenGenerator.Generate(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, ErrGeneratingToken(err)
	}

	return &pbs.LoginResponse{Token: token}, nil
}

// GetUser first tries to get the user with the given ID.
// If the query fails (with a gorm.ErrRecordNotFound), then that user doesn't exist.
// If the query fails (for some other reason), then it returns an unknown error.
// If everything is OK, it returns the user.
func (s *ServiceLayer) GetUser(ctx context.Context, req *pbs.GetUserRequest) (*pbs.GetUserResponse, error) {
	user, err := s.Storage.GetUser(options.WithUserID(int(req.UserId)))
	if errIsNotFound(err) {
		return nil, ErrNotFound("user")
	}
	if err != nil || user == nil {
		return nil, ErrInUsersDBCall(ctx, err)
	}
	return &pbs.GetUserResponse{User: user.ToUserInfo()}, nil
}

// GetUsers first gets the page, pageSize and filterQueryOptions from the request.
// With those values, it gets the users from the database. If there's an error, it returns unknown.
// If everything is OK, it returns the users and the pagination info.
func (s *ServiceLayer) GetUsers(ctx context.Context, req *pbs.GetUsersRequest) (*pbs.GetUsersResponse, error) {
	page, pageSize := getPaginationValues(req)
	filter := options.WithFilter("username", req.GetFilter())

	// While our page is 0-based, gorm offsets are 1-based. That's why we subtract 1.
	users, totalMatches, err := s.Storage.GetUsers(page-1, pageSize, filter)
	if err != nil {
		return nil, ErrInUsersDBCall(ctx, err)
	}

	return &pbs.GetUsersResponse{
		Users:      users.ToUserInfo(),
		Pagination: responsePagination(page, pageSize, totalMatches),
	}, nil
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*       - Users Service Errors -      */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var (
	ErrInUsersDBCall = func(ctx context.Context, err error) error {
		route := getRouteFromCtx(ctx)
		go core.LogUnexpected(err)
		return errs.ErrSvcUserRelated(err, route)
	}
	ErrGeneratingToken = func(err error) error {
		return errs.ErrSvcOnTokenGeneration(err)
	}
)
