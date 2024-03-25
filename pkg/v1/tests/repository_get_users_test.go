package tests

import (
	"errors"
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/mocks"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/options"

	"github.com/stretchr/testify/assert"
)

type GetUsersExpected struct {
	Result models.Users
	Total  int
	Error  error
}

func TestRepositoryGetUsers(t *testing.T) {
	testCases := []struct {
		name     string
		expected GetUsersExpected
	}{
		{
			name:     "tc_get_users_ok",
			expected: GetUsersExpected{models.Users{{ID: 1, Username: "username"}}, 20, nil},
		},
		{
			name:     "tc_get_users_no_results",
			expected: GetUsersExpected{nil, 0, nil},
		},
		{
			name:     "tc_get_users_error_in_count",
			expected: GetUsersExpected{nil, 0, errors.New("error counting users")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			setupFn := func(mock *mocks.Gorm) { setupGetUsersTC(mock, tc.expected) }
			repository, _ := NewTestRepository(setupFn)

			result, total, err := repository.GetUsers(0, 10, []options.QueryOpt{}...)

			assert.Equal(t, tc.expected.Result, result)
			assert.Equal(t, tc.expected.Total, total)
			assert.Equal(t, tc.expected.Error, err)
		})
	}
}

func setupGetUsersTC(m *mocks.Gorm, expected GetUsersExpected) {
	//m.WillCount(int64(expected.Total))
	//OnCallTo(m, "Model", &models.User{})
	//OnCallTo(m, "Count", mock.AnythingOfType(int64PtrTypeName))
	//OnGetError(m, expected.Error)
	//
	//if expected.Error != nil || expected.Total == 0 {
	//	return
	//}
	//
	//m.WillFind(expected.Result)
	//OnCallTo(m, "Offset", 0)
	//OnCallTo(m, "Limit", 10)
	//OnCallTo(m, "Find", new(models.Users), []interface{}(nil))
	//OnGetError(m, expected.Error)
}
