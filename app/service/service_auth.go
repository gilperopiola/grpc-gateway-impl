package service

import (
	"strconv"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	sql "github.com/gilperopiola/grpc-gateway-impl/app/tools/db_tool/sqldb"

	"go.uber.org/zap"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - Auth Service -          */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (s *AuthService) Signup(ctx god.Ctx, req *pbs.SignupRequest) (*pbs.SignupResponse, error) {
	user, err := s.Tools.GetUser(ctx, sql.WithUsername(req.Username))
	if err == nil && user != nil {
		return nil, errUserAlreadyExists()
	}

	if !s.Tools.IsNotFound(err) {
		return nil, errCallingUsersDB(ctx, err)
	}

	// If we're here, we should have gotten a not found in the function above.
	if user, err = s.Tools.CreateUser(ctx, req.Username, s.Tools.HashPassword(req.Password)); err != nil {
		return nil, errCallingUsersDB(ctx, err)
	}

	go s.doAfterSignup(ctx, user)

	return &pbs.SignupResponse{Id: int32(user.ID)}, nil
}

func (s *AuthService) doAfterSignup(ctx god.Ctx, user *models.User) {
	s.Tools.CreateFolder("users/user_" + strconv.Itoa(user.ID))
	s.Tools.CreateGroup(ctx, user.Username+"'s First Group", user.ID, []int{})

	w, err := s.Tools.GetCurrentWeather(ctx, 44.34, 10.99)
	core.LogIfErr(err)
	zap.S().Info("Weather", w)
}

// Login first tries to get the user with the given username.
// If the query fails (with a gorm.ErrRecordNotFound), then that user doesn't exist.
// If the query fails (for some other reason), then we return an unknown error.
// Then we PasswordsMatch both passwords. If they don't match, we return an unauthenticated error.
// If everything is OK, we generate a token and return it.
func (s *AuthService) Login(ctx god.Ctx, req *pbs.LoginRequest) (*pbs.LoginResponse, error) {
	user, err := s.Tools.GetUser(ctx, sql.WithUsername(req.Username))
	if s.Tools.IsNotFound(err) {
		return nil, errUserNotFound()
	}
	if err != nil || user == nil {
		return nil, errCallingUsersDB(ctx, err)
	}

	if !s.Tools.PasswordsMatch(req.Password, user.Password) {
		return nil, errs.GRPCUnauthenticated()
	}

	token, err := s.Tools.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errGeneratingToken(err)
	}

	return &pbs.LoginResponse{Token: token}, nil
}
