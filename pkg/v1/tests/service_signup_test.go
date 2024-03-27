package tests

import (
	"context"
	"testing"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/mocks"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

func TestServiceSignup(t *testing.T) {

	type getExpectedFn func() (*usersPB.SignupResponse, error)

	expect := func(response *usersPB.SignupResponse, err error) getExpectedFn {
		return func() (*usersPB.SignupResponse, error) {
			return copyResponsePtr(response).(*usersPB.SignupResponse), err
		}
	}

	setupMock := func(getResult, createResult *models.User, getErr, createErr error) setupRepoMockFn {
		return func(repoMock *mocks.Repository) {
			repoMock.On("GetUser", mock.Anything).Return(copyUserPtr(getResult), getErr).Once()
			if getErr != gorm.ErrRecordNotFound {
				return
			}
			repoMock.On("CreateUser", mock.Anything, mock.Anything).Return(copyUserPtr(createResult), createErr).Once()
		}
	}

	tests := []struct {
		name        string
		setupMock   setupRepoMockFn
		getExpected getExpectedFn
	}{
		{
			name:        "tc_service_signup_ok",
			setupMock:   setupMock(user, userWithID, gorm.ErrRecordNotFound, nil),
			getExpected: expect(signupResponse, nil),
		},
		{
			name:        "tc_service_signup_already_exists",
			setupMock:   setupMock(userWithID, nil, nil, nil),
			getExpected: expect(nil, service.ErrAlreadyExists("user")),
		},
		{
			name:        "tc_service_signup_get_error",
			setupMock:   setupMock(nil, nil, errGettingUser, nil),
			getExpected: expect(nil, service.UserErr(context.Background(), errGettingUser)),
		},
		{
			name:        "tc_service_signup_create_error",
			setupMock:   setupMock(nil, nil, gorm.ErrRecordNotFound, errCreatingUser),
			getExpected: expect(nil, service.UserErr(context.Background(), errCreatingUser)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service, mock := newTestServiceQuick(test.setupMock) // Prepare
			expectedResponse, expectedErr := test.getExpected()

			response, err := service.Signup(context.Background(), signupRequest) // Act

			assert.True(t, proto.Equal(expectedResponse, response)) // Assert
			assertServiceError(t, expectedErr, err)
			mock.AssertExpectations(t)
		})
	}
}
