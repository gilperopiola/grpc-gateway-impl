package tests

import (
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
	"github.com/gilperopiola/grpc-gateway-impl/etc/tests/mocks"

	"github.com/stretchr/testify/assert"
)

func TestStorageCreateUser(t *testing.T) {

	type getExpectedFn func() (*models.User, error)

	expect := func(user *models.User, err error) getExpectedFn {
		return func() (*models.User, error) {
			return copyUser(user), err
		}
	}

	setupMock := func(userBeforeCreate *models.User, userAfterCreate *models.User, err error) setupGormMockFn {
		return func(mock *mocks.Gorm) {
			mock.OnCreateUser(copyUser(userBeforeCreate), copyUser(userAfterCreate)).ErrorWillBe(err)
		}
	}

	tests := []struct {
		name        string
		setupMock   setupGormMockFn
		getExpected getExpectedFn
	}{
		{
			name:        "tc_storage_create_user_ok",
			setupMock:   setupMock(user, userWithID, nil),
			getExpected: expect(userWithID, nil),
		},
		{
			name:        "tc_storage_create_user_error",
			setupMock:   setupMock(user, nil, errCreatingUser),
			getExpected: expect(nil, errCreatingUser),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			external, mock := newTestStorage(tc.setupMock) // Prepare
			expectedUser, expectedErr := tc.getExpected()

			user, err := external.CreateUser(username, password) // Act

			assert.Equal(t, expectedUser, user) // Assert
			assertDBError(t, expectedErr, err)
			mock.AssertExpectations(t)
		})
	}
}
