package tests

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core/interfaces"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/external"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/service"
	"github.com/gilperopiola/grpc-gateway-impl/app/modules"
	"github.com/gilperopiola/grpc-gateway-impl/etc/tests/mocks"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - Tests Setup -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

/* Service tests setup */

// newTestService returns a test service and a external mock. Takes a function to setup the mock, a tokenGen and a pwdHasher.
func newTestService(setupMock setupRepoMockFn, tokenGen interfaces.TokenGenerator, pwdHasher interfaces.PwdHasher) (interfaces.BusinessLayer, *mocks.Storage) {
	repoMock := &mocks.Storage{}
	setupMock(repoMock)
	service := service.NewService(repoMock, tokenGen, pwdHasher)
	return service, repoMock
}

// newTestServiceQuick is like newTestService but with a default tokenGen and pwdHasher.
func newTestServiceQuick(setupMock setupRepoMockFn) (interfaces.BusinessLayer, *mocks.Storage) {
	return newTestService(setupMock, newTestTokenAuthenticator(), newTestPwdHasher())
}

type setupRepoMockFn func(external *mocks.Storage)

/* Storage tests setup */

// newTestStorage returns a test external and a gorm mock.
func newTestStorage(setupMock setupGormMockFn) (interfaces.Storage, *mocks.Gorm) {
	gormMock := &mocks.Gorm{}
	setupMock(gormMock)
	external := external.NewExternalLayer(gormMock)
	return external.GetStorage(), gormMock
}

type setupGormMockFn func(*mocks.Gorm)

var setupGormMockEmpty = func(*mocks.Gorm) { /* Use this when a test case doesn't call any method on the gorm mock.*/ }

/* Etc */

func newTestServiceinterfaces() (interfaces.TokenGenerator, interfaces.PwdHasher) {
	return newTestTokenAuthenticator(), newTestPwdHasher()
}

func newTestTokenAuthenticator() interfaces.TokenAuthenticator {
	return modules.NewJWTAuthenticator(jwtSecret, 10)
}

func newTestPwdHasher() interfaces.PwdHasher {
	return modules.NewPwdHasher(hashSalt)
}
