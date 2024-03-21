package repository

import (
	"errors"
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/mocks"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/db"
	"github.com/stretchr/testify/assert"
)

/* ----------------------------------- */
/*      - Users Repository Tests -     */
/* ----------------------------------- */

type TestRepository interface {
	Repository
	GetGormMock() *mocks.GormMock
}

type testRepository struct {
	Repository
	gormMock *mocks.GormMock
}

func (r *testRepository) GetGormMock() *mocks.GormMock {
	return r.gormMock
}

func setupNewTest() TestRepository {
	gormMock := &mocks.GormMock{}
	dbWrapper := db.NewDatabaseWrapper(nil, gormMock)
	repository := NewRepository(dbWrapper)
	return &testRepository{repository, gormMock}
}

/* -------------------------------------------------------------------------- */
/*   - repository.CreateUser(argUsername, argPwd) -> (*models.User, error) -  */
/* -------------------------------------------------------------------------- */
//
// - GormAdapter.Create(userToCreate) 	-> (GormAdapter)
// - GormAdapter.Error() 					-> (error)

func TestRepositoryCreateUser(t *testing.T) {
	var (
		argUsername1 = "username"
		argPassword1 = "password"

		expectedError = errors.New("error creating user")
	)

	testRepositoryCreateUserOK(t, argUsername1, argPassword1)
	testRepositoryCreateUserError(t, argUsername1, argPassword1, expectedError)
}

func testRepositoryCreateUserOK(t *testing.T, argUsername string, argPassword string) {
	repository := setupNewTest()
	gormMock := repository.GetGormMock()

	// - Mock Calls
	userToCreate := &models.User{Username: argUsername, Password: argPassword}
	gormMock.On("Create", userToCreate).Return(gormMock).Once()

	gormMock.On("Error").Return(nil).Once()

	// - Act & Assert
	user, err := repository.CreateUser(argUsername, argPassword)
	assert.NotNil(t, user)
	assert.Nil(t, err)
}

func testRepositoryCreateUserError(t *testing.T, argUsername string, argPassword string, expectedErr error) {
	repository := setupNewTest()
	gormMock := repository.GetGormMock()

	// - Mock Calls
	userToCreate := &models.User{Username: argUsername, Password: argPassword}
	gormMock.On("Create", userToCreate).Return(gormMock).Once()

	gormMock.On("Error").Return(expectedErr).Once()

	// - Act & Assert
	user, err := repository.CreateUser(argUsername, argPassword)
	assert.Nil(t, user)
	assert.ErrorIs(t, err, expectedErr)
}
