package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	sql "github.com/gilperopiola/grpc-gateway-impl/app/tools/db_tool/sqldb"

	"go.uber.org/zap"
)

type AuthSubService struct {
	pbs.UnimplementedAuthServiceServer
	Tools core.Tools
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - Auth Service -          */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (s *AuthSubService) Signup(ctx god.Ctx, req *pbs.SignupRequest) (*pbs.SignupResponse, error) {
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

func (s *AuthSubService) doAfterSignup(ctx god.Ctx, user *models.User) {
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
func (s *AuthSubService) Login(ctx god.Ctx, req *pbs.LoginRequest) (*pbs.LoginResponse, error) {
	user, err := s.Tools.GetUser(ctx, sql.WithUsername(req.Username))
	if s.Tools.IsNotFound(err) {
		return nil, errs.GRPCNotFound("user", req.Username)
	}
	if err != nil || user == nil {
		return nil, errs.GRPCFromDB(err, core.RouteNameFromCtx(ctx))
	}

	if !s.Tools.PasswordsMatch(req.Password, user.Password) {
		return nil, errs.GRPCWrongLoginInfo()
	}

	token, err := s.Tools.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errs.GRPCGeneratingToken(err)
	}

	zap.L().Info(req.String())

	w, err := s.Tools.GetCurrentWeather(ctx, 44.34, 10.99)
	core.LogIfErr(err)
	zap.S().Info("Weather", w)

	got, err := s.Tools.NewCompletion(ctx, "Gimme just the latitude and longitude of buenos aires, response must follow format: -35.03, -54.33", "")
	core.LogIfErr(err)
	zap.S().Info(fmt.Sprintf("Got %s", got))

	splat := strings.Split(got, ",")
	if len(splat) == 2 {
		lat, err := strconv.ParseFloat(strings.TrimSpace(splat[0]), 64)
		core.LogIfErr(err)
		lon, err := strconv.ParseFloat(strings.TrimSpace(splat[1]), 64)
		core.LogIfErr(err)

		w, err := s.Tools.GetCurrentWeather(ctx, lat, lon)
		core.LogIfErr(err)
		zap.S().Info("Weather in bsas", w)
	}

	return &pbs.LoginResponse{Token: token}, nil
}
