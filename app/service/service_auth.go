package service

import (
	"strconv"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	sql "github.com/gilperopiola/grpc-gateway-impl/app/toolbox/db_tool/sqldb"

	"go.uber.org/zap"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - Auth Service -          */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (s *Service) Signup(ctx god.Ctx, req *pbs.SignupRequest) (*pbs.SignupResponse, error) {
	user, err := s.Toolbox.GetUser(ctx, sql.WithUsername(req.Username))
	if err == nil && user != nil {
		return nil, errUserAlreadyExists()
	}

	if !s.Toolbox.IsNotFound(err) {
		return nil, errCallingUsersDB(ctx, err)
	}

	// If we're here, we should have gotten a not found in the function above.
	if user, err = s.Toolbox.CreateUser(ctx, req.Username, s.Toolbox.HashPassword(req.Password)); err != nil {
		return nil, errCallingUsersDB(ctx, err)
	}

	go s.doAfterSignup(ctx, user)

	return &pbs.SignupResponse{Id: int32(user.ID)}, nil
}

func (s *Service) doAfterSignup(ctx god.Ctx, user *models.User) {
	s.Toolbox.CreateFolder("users/user_" + strconv.Itoa(user.ID))
	s.Toolbox.CreateGroup(ctx, user.Username+"'s First Group", user.ID, []int{})

	w, err := s.Toolbox.GetCurrentWeather(ctx, 44.34, 10.99)
	core.LogIfErr(err)
	zap.S().Info("Weather", w)
}

// Login first tries to get the user with the given username.
// If the query fails (with a gorm.ErrRecordNotFound), then that user doesn't exist.
// If the query fails (for some other reason), then we return an unknown error.
// Then we PasswordsMatch both passwords. If they don't match, we return an unauthenticated error.
// If everything is OK, we generate a token and return it.
func (s *Service) Login(ctx god.Ctx, req *pbs.LoginRequest) (*pbs.LoginResponse, error) {
	user, err := s.Toolbox.GetUser(ctx, sql.WithUsername(req.Username))
	if s.Toolbox.IsNotFound(err) {
		return nil, errUserNotFound()
	}
	if err != nil || user == nil {
		return nil, errCallingUsersDB(ctx, err)
	}

	if !s.Toolbox.PasswordsMatch(req.Password, user.Password) {
		return nil, errs.GRPCUnauthenticated()
	}

	token, err := s.Toolbox.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errGeneratingToken(err)
	}

	return &pbs.LoginResponse{Token: token}, nil
}
