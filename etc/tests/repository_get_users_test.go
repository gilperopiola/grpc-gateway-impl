package tests

import (
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/external/storage/options"
	"github.com/gilperopiola/grpc-gateway-impl/etc/tests/mocks"

	"github.com/stretchr/testify/assert"
)

func TestStorageGetUsers(t *testing.T) {

	type getExpectedFn func() (models.Users, int, error)

	expect := func(expectUsers models.Users, expectTotal int, expectErr error) getExpectedFn {
		return func() (models.Users, int, error) {
			return copyUsers(expectUsers), expectTotal, expectErr
		}
	}

	setupMock := func(whereUserID int, whereUsername string, countResult int, countErr error, findResultUsers *models.Users, findErr error) setupGormMockFn {
		return func(mock *mocks.Gorm) {
			mock.OnModel(&models.User{})
			mock.OnWhereUser(whereUserID, whereUsername)
			mock.OnCount(countResult).ErrorWillBe(countErr)

			if countErr != nil || countResult == 0 {
				return
			}

			mock.OnOffset()
			mock.OnLimit()
			var findUsersIn models.Users
			var findUsersOut = copyUsers(*findResultUsers)
			mock.OnFindUsers(&findUsersIn, &findUsersOut).ErrorWillBe(findErr)
		}
	}

	tests := []struct {
		name        string
		setupMock   setupGormMockFn
		getExpected getExpectedFn
	}{
		{
			name:        "tc_storage_get_users_ok",
			setupMock:   setupMock(0, "", 20, nil, &users, nil),
			getExpected: expect(users, 20, nil),
		},
		{
			name:        "tc_storage_get_users_no_results",
			setupMock:   setupMock(0, "", 0, nil, &usersNil, nil),
			getExpected: expect(usersNil, 0, nil),
		},
		{
			name:        "tc_storage_get_users_error_in_count",
			setupMock:   setupMock(0, "", 0, errCountingUsers, nil, nil),
			getExpected: expect(nil, 0, errCountingUsers),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			external, mock := newTestStorage(test.setupMock)
			expectedUsers, expectedTotal, expectedErr := test.getExpected()

			users, total, err := external.GetUsers(0, 10, []options.QueryOpt{}...)

			assert.Equal(t, expectedUsers, users)
			assert.Equal(t, expectedTotal, total)
			assertDBError(t, expectedErr, err)
			mock.AssertExpectations(t)
		})
	}
}
