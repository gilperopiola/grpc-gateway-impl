package service

import (
	"strconv"

	"github.com/gilperopiola/god"
	sql "github.com/gilperopiola/grpc-gateway-impl/app/clients/dbs/sqldb"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/utils"
)

type AuthSvc struct {
	pbs.UnimplementedAuthServiceServer
	Clients core.Clients
	Tools   core.Tools
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - Auth Service -          */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

//  1. We try to get a user with the username that we want to create.
//     a. If we don't get any error, that means the user already exists.
//     b. If we get an error but it's not a DBNotFound, we return an unknown error.
func (s *AuthSvc) Signup(ctx god.Ctx, req *pbs.SignupRequest) (*pbs.SignupResponse, error) {
	user, err := s.Clients.DBGetUser(ctx, sql.WithUsername(req.Username))
	if err == nil || user != nil {
		return nil, errUserAlreadyExists()
	}

	if !utils.IsNotFound(err) {
		return nil, errCallingUsersDB(ctx, err)
	}

	// If we're here, we should have gotten a not found in the function above.
	if user, err = s.Clients.DBCreateUser(ctx, req.Username, s.Tools.HashPassword(req.Password)); err != nil {
		return nil, errCallingUsersDB(ctx, err)
	}

	defer func() {
		go s.doAfterSignup(ctx, user)
	}()

	return &pbs.SignupResponse{Id: int32(user.ID)}, nil
}

func (s *AuthSvc) doAfterSignup(ctx god.Ctx, user *models.User) {
	s.Tools.CreateFolder("users/user_" + strconv.Itoa(user.ID))
	s.Clients.DBCreateGroup(ctx, user.Username+"'s First Group", user.ID, []int{})
}

// Login first tries to get the user with the given username.
// If the query fails (with a gorm.ErrRecordNotFound), then that user doesn't exist.
// If the query fails (for some other reason), then we return an unknown error.
// Then we PasswordsMatch both passwords. If they don't match, we return an unauthenticated error.
// If everything is OK, we generate a token and return it.
func (s *AuthSvc) Login(ctx god.Ctx, req *pbs.LoginRequest) (*pbs.LoginResponse, error) {
	user, err := s.Clients.DBGetUser(ctx, sql.WithUsername(req.Username))
	if utils.IsNotFound(err) {
		return nil, errs.GRPCNotFound("user", req.Username)
	}
	if err != nil || user == nil {
		return nil, errs.GRPCFromDB(err, shared.GetRouteFromCtx(ctx).Name)
	}

	if !s.Tools.PasswordsMatch(req.Password, user.Password) {
		return nil, errs.GRPCWrongLoginInfo()
	}

	token, err := s.Tools.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errs.GRPCGeneratingToken(err)
	}

	return &pbs.LoginResponse{Token: token}, nil
}
