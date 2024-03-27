package tests

import (
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/mocks"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/options"

	"github.com/stretchr/testify/assert"
)

func TestRepositoryGetUser(t *testing.T) {

	type getExpectedFn func() (*models.User, error)

	expect := func(user *models.User, err error) getExpectedFn {
		return func() (*models.User, error) {
			return copyUserPtr(user), err
		}
	}

	setupMock := func(argUser *models.User, argUserID int, argUsername string, result *models.User, err error) setupGormMockFn {
		return func(mock *mocks.Gorm) {
			mock.OnModel(&models.User{})
			mock.OnWhereUser(argUserID, argUsername)
			mock.OnFirstUser(copyUserPtr(argUser), copyUserPtr(result)).ErrorWillBe(err)
		}
	}

	tests := []struct {
		name        string
		options     []options.QueryOpt
		setupMock   setupGormMockFn
		getExpected getExpectedFn
	}{
		{
			name:        "tc_repository_get_user_with_id_ok",
			options:     options.Slice(options.WithUserID(userID)),
			setupMock:   setupMock(userEmpty, userID, "", userWithID, nil),
			getExpected: expect(userWithID, nil),
		},
		{
			name:        "tc_repository_get_user_with_username_ok",
			options:     options.Slice(options.WithUsername(username)),
			setupMock:   setupMock(userEmpty, 0, username, userWithID, nil),
			getExpected: expect(userWithID, nil),
		},
		{
			name:        "tc_repository_get_user_no_options_error",
			options:     options.Slice(nil),
			setupMock:   emptyGormMockFn,
			getExpected: expect(nil, errNoOpts),
		},
		{
			name:        "tc_repository_get_user_error",
			options:     options.Slice(options.WithUserID(userID)),
			setupMock:   setupMock(userEmpty, userID, "", nil, errGettingUser),
			getExpected: expect(nil, errGettingUser),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository, mock := newTestRepository(test.setupMock) // Prepare
			expectedUser, expectedErr := test.getExpected()

			user, err := repository.GetUser(test.options...) // Act

			assert.Equal(t, expectedUser, user) // Assert
			assertDBError(t, expectedErr, err)
			mock.AssertExpectations(t)
		})
	}
}
