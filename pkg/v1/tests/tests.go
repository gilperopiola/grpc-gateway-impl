package tests

import (
	"errors"
	"testing"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/common"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/mocks"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"
	"google.golang.org/protobuf/proto"

	"github.com/stretchr/testify/assert"
)

/* ----------------------------------- */
/*         - Tests Variables -         */
/* ----------------------------------- */

var (

	// JWT Auth
	jwtSecret     = "jwt_secret"
	generateToken = func(g common.TokenGenerator, id int, username, role string) string {
		token, _ := g.Generate(id, username, models.Role(role))
		return token
	}

	// Pwd Hasher
	hashSalt = "hash_salt"
	hashUser = func(h common.PwdHasher, user *models.User) *models.User {
		user.Password = h.Hash(user.Password)
		return user
	}

	// User Model
	userID   = 1
	username = "username"
	password = "password"

	user       = &models.User{Username: username, Password: password}
	userWithID = &models.User{ID: 1, Username: username, Password: password}
	userEmpty  = &models.User{}
	users      = models.Users{user}
	usersNil   = models.Users(nil)

	// Repository Errors
	errCreatingUser  = errors.New(repository.CreateUserErr)
	errGettingUser   = errors.New(repository.GetUserErr)
	errGettingUsers  = errors.New(repository.GetUsersErr)
	errCountingUsers = errors.New(repository.CountUsersErr)
	errNoOpts        = errors.New(repository.ErrNoOptions)

	// PB Requests & Responses
	signupRequest  = &usersPB.SignupRequest{Username: username, Password: password}
	signupResponse = &usersPB.SignupResponse{Id: int32(userID)}
	loginRequest   = &usersPB.LoginRequest{Username: username, Password: password}
)

/* ----------------------------------- */
/*           - Tests Setup -           */
/* ----------------------------------- */

// newTestRepository returns a new testing repository with a gorm mock inside.
func newTestRepository(setupMock setupGormMockFn) (repository.Repository, *mocks.Gorm) {
	gormMock := &mocks.Gorm{}
	repository := repository.NewRepository(gormMock)
	setupMock(gormMock)
	return repository, gormMock
}

func newTestService(setupMock setupRepoMockFn, tokenGen common.TokenGenerator, pwdHasher common.PwdHasher) (service.Service, *mocks.Repository) {
	repoMock := &mocks.Repository{}
	service := service.NewService(repoMock, tokenGen, pwdHasher)
	setupMock(repoMock)
	return service, repoMock
}

func newTestServiceQuick(setupMock setupRepoMockFn) (service.Service, *mocks.Repository) {
	return newTestService(setupMock, newTestTokenGenerator(), newTestPwdHasher())
}

func newTestCommonComponents() (common.TokenGenerator, common.PwdHasher) {
	return newTestTokenGenerator(), newTestPwdHasher()
}

func newTestTokenGenerator() common.TokenGenerator {
	jwtSessionDays := 10
	return common.NewJWTAuthenticator(jwtSecret, jwtSessionDays)
}

func newTestPwdHasher() common.PwdHasher {
	return common.NewPwdHasher(hashSalt)
}

type setupGormMockFn func(*mocks.Gorm)

var emptyGormMockFn = func(*mocks.Gorm) { /* Use this when a test case doesn't ever call any method on the mock.*/ }

type setupRepoMockFn func(repo *mocks.Repository)

/* ----------------------------------- */
/*          - Tests Helpers -          */
/* ----------------------------------- */

func assertDBError(t *testing.T, expectedErr, err error) {
	if expectedErr == nil {
		assert.NoError(t, err)
		return
	}

	var dbErr *errs.DBError
	if assert.ErrorAs(t, err, &dbErr) {
		assert.Equal(t, expectedErr, dbErr.Unwrap())
	}
}

func assertServiceError(t *testing.T, expectedErr, err error) {
	if expectedErr == nil {
		assert.NoError(t, err)
		return
	}

	assert.ErrorIs(t, err, expectedErr)
}

func copyUserPtr(user *models.User) *models.User {
	if user == nil {
		return nil
	}
	copy := *user
	return &copy
}

func copyUsers(users models.Users) models.Users {
	if users == nil {
		return nil
	}
	copy := make(models.Users, len(users))
	for i, user := range users {
		copy[i] = copyUserPtr(user)
	}
	return copy
}

func copyResponsePtr(response models.PBResponse) models.PBResponse {
	if response == nil {
		return nil
	}
	return proto.Clone(response).(models.PBResponse)
}
