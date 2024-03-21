package tests

import (
	"errors"
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"

	"github.com/stretchr/testify/assert"
)

/* ----------------------------------- */
/*      - Users Repository Tests -     */
/* ----------------------------------- */

/* - CreateUser
/*
/* --------------------------------------------------------------------------- */
/*   - repository.CreateUser(username, pwd string) -> (*models.User, error) -  */
/* --------------------------------------------------------------------------- */
//
// - GormAdapter.Create(userToCreate) 	-> (GormAdapter)
// - GormAdapter.Error() 					-> (error)

func TestRepositoryCreateUser(t *testing.T) {
	var (
		argUsername1 = "username"
		argPassword1 = "password"

		errOnCreate = errors.New("error creating user")
	)

	testRepositoryCreateUserOK(t, argUsername1, argPassword1)
	testRepositoryCreateUserError(t, argUsername1, argPassword1, errOnCreate)
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

func testRepositoryCreateUserError(t *testing.T, argUsername string, argPassword string, errOnCreate error) {
	repository := setupNewTest()
	gormMock := repository.GetGormMock()

	// - Mock Calls
	userToCreate := &models.User{Username: argUsername, Password: argPassword}
	gormMock.On("Create", userToCreate).Return(gormMock).Once()

	gormMock.On("Error").Return(errOnCreate).Once()

	// - Act & Assert
	user, err := repository.CreateUser(argUsername, argPassword)
	assert.Nil(t, user)
	assert.ErrorIs(t, err, errOnCreate)
}
