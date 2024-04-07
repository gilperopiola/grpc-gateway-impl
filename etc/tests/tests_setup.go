package tests

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/modules"
	"github.com/gilperopiola/grpc-gateway-impl/app/service"
	"github.com/gilperopiola/grpc-gateway-impl/app/storage"
	"github.com/gilperopiola/grpc-gateway-impl/etc/tests/mocks"
)

/* ----------------------------------- */
/*           - Tests Setup -           */
/* ----------------------------------- */

/* Service tests setup */

// newTestService returns a test service and a storage mock. Takes a function to setup the mock, a tokenGen and a pwdHasher.
func newTestService(setupMock setupRepoMockFn, tokenGen modules.TokenGenerator, pwdHasher modules.PwdHasher) (service.Service, *mocks.Storage) {
	repoMock := &mocks.Storage{}
	setupMock(repoMock)
	service := service.NewService(repoMock, tokenGen, pwdHasher)
	return service, repoMock
}

// newTestServiceQuick is like newTestService but with a default tokenGen and pwdHasher.
func newTestServiceQuick(setupMock setupRepoMockFn) (service.Service, *mocks.Storage) {
	return newTestService(setupMock, newTestTokenAuthenticator(), newTestPwdHasher())
}

type setupRepoMockFn func(repo *mocks.Storage)

/* Storage tests setup */

// newTestStorage returns a test storage and a gorm mock.
func newTestStorage(setupMock setupGormMockFn) (storage.Storage, *mocks.Gorm) {
	gormMock := &mocks.Gorm{}
	setupMock(gormMock)
	storage := storage.NewStorage(gormMock)
	return storage, gormMock
}

type setupGormMockFn func(*mocks.Gorm)

var setupGormMockEmpty = func(*mocks.Gorm) { /* Use this when a test case doesn't call any method on the gorm mock.*/ }

/* Etc */

func newTestServiceModules() (modules.TokenGenerator, modules.PwdHasher) {
	return newTestTokenAuthenticator(), newTestPwdHasher()
}

func newTestTokenAuthenticator() modules.TokenAuthenticator {
	return modules.NewJWTAuthenticator(jwtSecret, 10)
}

func newTestPwdHasher() modules.PwdHasher {
	return modules.NewPwdHasher(hashSalt)
}
