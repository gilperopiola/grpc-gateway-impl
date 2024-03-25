package tests

import (
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/mocks"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"

	"github.com/stretchr/testify/assert"
)

func TestRepositoryCreateUser(t *testing.T) {

	type getExpectedFn func() (*models.User, error)

	expect := func(user *models.User, err error) getExpectedFn {
		return func() (*models.User, error) {
			return copyUserPtr(user), err
		}
	}

	setupMock := func(userBeforeCreate *models.User, userAfterCreate *models.User, err error) setupGormMockFn {
		return func(mock *mocks.Gorm) {
			mock.OnCreateUser(copyUserPtr(userBeforeCreate), copyUserPtr(userAfterCreate)).ErrorWillBe(err)
		}
	}

	tests := []struct {
		name        string
		setupMock   setupGormMockFn
		getExpected getExpectedFn
	}{
		{
			name:        "tc_repository_create_user_ok",
			setupMock:   setupMock(user, userWithID, nil),
			getExpected: expect(userWithID, nil),
		},
		{
			name:        "tc_repository_create_user_error",
			setupMock:   setupMock(user, nil, errCreatingUser),
			getExpected: expect(nil, errCreatingUser),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Prepare ⬇️
			repository, mock := NewTestRepository(tc.setupMock)

			// Act ⬇️
			result, err := repository.CreateUser(username, password)

			// Assert ⬇️
			expResult, expErr := tc.getExpected()

			assert.Equal(t, expResult, result)
			assert.ErrorIs(t, err, expErr)
			mock.AssertExpectations(t)
		})
	}
}
