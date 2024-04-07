package tests

import (
	"context"
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/service"
	"github.com/gilperopiola/grpc-gateway-impl/etc/tests/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

func TestServiceSignup(t *testing.T) {

	type getReturns mocks.RepoGetUserReturns
	type createReturns mocks.RepoCreateUserReturns

	setupMock := func(getReturns getReturns, createReturns createReturns) setupRepoMockFn {
		return func(repoMock *mocks.Storage) {
			repoMock.On("GetUser", mock.Anything).Return(copyUser(getReturns.User), getReturns.Err).Once()
			if getReturns.Err != gorm.ErrRecordNotFound {
				return
			}
			repoMock.On("CreateUser", mock.Anything, mock.Anything).Return(copyUser(createReturns.User), createReturns.Err).Once()
		}
	}

	type expected struct {
		response *pbs.SignupResponse
		err      error
	}

	tests := []struct {
		name      string
		setupMock setupRepoMockFn
		expected  expected
	}{
		{
			name:      "tc_service_signup_ok",
			setupMock: setupMock(getReturns{nil, gorm.ErrRecordNotFound}, createReturns{userWithID, nil}),
			expected:  expected{copyPB(signupResponse).(*pbs.SignupResponse), nil},
		},
		{
			name:      "tc_service_signup_already_exists",
			setupMock: setupMock(getReturns{userWithID, nil}, createReturns{nil, nil}),
			expected:  expected{nil, service.ErrAlreadyExists("user")},
		},
		{
			name:      "tc_service_signup_get_error",
			setupMock: setupMock(getReturns{nil, errGettingUser}, createReturns{nil, nil}),
			expected:  expected{nil, service.UsersDBError(context.Background(), errGettingUser)},
		},
		{
			name:      "tc_service_signup_create_error",
			setupMock: setupMock(getReturns{nil, gorm.ErrRecordNotFound}, createReturns{nil, errCreatingUser}),
			expected:  expected{nil, service.UsersDBError(context.Background(), errCreatingUser)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service, mock := newTestServiceQuick(test.setupMock) // Prepare

			response, err := service.Signup(context.Background(), signupRequest) // Act

			assert.True(t, proto.Equal(test.expected.response, response)) // Assert
			assertSvcError(t, test.expected.err, err)
			mock.AssertExpectations(t)
		})
	}
}
