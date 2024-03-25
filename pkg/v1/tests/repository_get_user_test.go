package tests

import (
	"fmt"
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/mocks"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/options"

	"github.com/stretchr/testify/assert"
)

func TestRepositoryGetUser(t *testing.T) {

	type getExpectedFn func() (*models.User, error)

	expect := func(result *models.User, err error) getExpectedFn {
		return func() (*models.User, error) {
			return copyUserPtr(result), err
		}
	}

	setupMock := func(argUser *models.User, argUserID int, argUsername string, result *models.User, err error) setupGormMockFn {
		return func(mock *mocks.Gorm) {
			mock.OnModel(&models.User{})
			mock.OnWhereIDOrUsername(argUserID, argUsername)
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
			options:     queryOptions(options.WithUserID(id)),
			setupMock:   setupMock(userEmpty, id, "", userWithID, nil),
			getExpected: expect(userWithID, nil),
		},
		{
			name:        "tc_repository_get_user_with_username_ok",
			options:     queryOptions(options.WithUsername(username)),
			setupMock:   setupMock(userEmpty, 0, username, userWithID, nil),
			getExpected: expect(userWithID, nil),
		},
		{
			name:        "tc_repository_get_user_no_options_error",
			options:     queryOptions(nil),
			setupMock:   emptyGormMockFn,
			getExpected: expect(nil, fmt.Errorf(repository.ErrNoOpts)),
		},
		{
			name:        "tc_repository_get_user_error",
			options:     queryOptions(options.WithUserID(id)),
			setupMock:   setupMock(userEmpty, id, "", nil, errGettingUser),
			getExpected: expect(nil, errGettingUser),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Prepare ⬇️
			repository, mock := NewTestRepository(test.setupMock)

			// Act ⬇️
			user, err := repository.GetUser(test.options...)

			// Assert ⬇️
			expectedUser, expectedErr := test.getExpected()

			assert.Equal(t, expectedUser, user)
			assert.ErrorIs(t, err, expectedErr)
			mock.AssertExpectations(t)
		})
	}
}
