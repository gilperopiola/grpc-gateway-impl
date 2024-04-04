package tests

import (
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/common"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/mocks"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"
)

/* ----------------------------------- */
/*           - Tests Setup -           */
/* ----------------------------------- */

/* Service tests setup */

// newTestService returns a test service and a repository mock. Takes a function to setup the mock, a tokenGen and a pwdHasher.
func newTestService(setupMock setupRepoMockFn, tokenGen common.TokenGenerator, pwdHasher common.PwdHasher) (service.Service, *mocks.Repository) {
	repoMock := &mocks.Repository{}
	setupMock(repoMock)
	service := service.NewService(repoMock, tokenGen, pwdHasher)
	return service, repoMock
}

// newTestServiceQuick is like newTestService but with a default tokenGen and pwdHasher.
func newTestServiceQuick(setupMock setupRepoMockFn) (service.Service, *mocks.Repository) {
	return newTestService(setupMock, newTestTokenAuthenticator(), newTestPwdHasher())
}

type setupRepoMockFn func(repo *mocks.Repository)

/* Repository tests setup */

// newTestRepository returns a test repository and a gorm mock.
func newTestRepository(setupMock setupGormMockFn) (repository.Repository, *mocks.Gorm) {
	gormMock := &mocks.Gorm{}
	setupMock(gormMock)
	repository := repository.NewRepository(gormMock)
	return repository, gormMock
}

type setupGormMockFn func(*mocks.Gorm)

var setupGormMockEmpty = func(*mocks.Gorm) { /* Use this when a test case doesn't call any method on the gorm mock.*/ }

/* Etc */

func newTestServiceComponents() (common.TokenGenerator, common.PwdHasher) {
	return newTestTokenAuthenticator(), newTestPwdHasher()
}

func newTestTokenAuthenticator() common.TokenAuthenticator {
	return common.NewJWTAuthenticator(jwtSecret, 10)
}

func newTestPwdHasher() common.PwdHasher {
	return common.NewPwdHasher(hashSalt)
}
