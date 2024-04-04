package tests

import (
	"context"
	"testing"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/common"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/mocks"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
)

func TestServiceLogin(t *testing.T) {
	expect := func(response *usersPB.LoginResponse, err error) (*usersPB.LoginResponse, error) {
		return copyPB(response).(*usersPB.LoginResponse), err
	}

	setupMock := func(h common.PwdHasher, getResult *models.User, getErr error) setupRepoMockFn {
		getResult.Password = hashPwd(h, getResult.Password)
		return func(repoMock *mocks.Repository) {
			repoMock.On("GetUser", mock.Anything, mock.Anything).Return(copyUser(getResult), getErr).Once()
		}
	}

	tests := []struct {
		name        string
		request     *usersPB.LoginRequest
		setupMock   func(common.PwdHasher) setupRepoMockFn
		getExpected func(common.TokenGenerator) (*usersPB.LoginResponse, error)
	}{
		{
			name:    "tc_service_login_ok",
			request: loginRequest,
			setupMock: func(h common.PwdHasher) setupRepoMockFn {
				return setupMock(h, userWithID, nil)
			},
			getExpected: func(g common.TokenGenerator) (*usersPB.LoginResponse, error) {
				return expect(&usersPB.LoginResponse{Token: generateToken(g, userID, username, "")}, nil)
			},
		},
		//{
		//	name:        "tc_service_login_not_found",
		//	setupMock:   setupMock(nil, gorm.ErrRecordNotFound),
		//	getExpected: expect(nil, service.ErrNotFound("user")),
		//},
		//{
		//	name:        "tc_service_login_get_error",
		//	setupMock:   setupMock(nil, errors.New("unknown error")),
		//	getExpected: expect(nil, service.UserErr(context.Background(), errors.New("unknown error"))),
		//},
		//{
		//	name:        "tc_service_login_password_mismatch",
		//	setupMock:   setupMock(&models.User{ID: "1", Username: "user1", Password: "hashedpassword"}, nil),
		//	getExpected: expect(nil, service.ErrUnauthenticated()),
		//},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokenGen, pwdHasher := newTestServiceComponents() // Prepare
			service, repoMock := newTestService(test.setupMock(pwdHasher), tokenGen, pwdHasher)
			expectedResponse, expectedErr := test.getExpected(tokenGen)

			response, err := service.Login(context.Background(), test.request) // Act

			assert.True(t, proto.Equal(expectedResponse, response)) // Assert
			assert.Equal(t, expectedErr, err)
			repoMock.AssertExpectations(t)
		})
	}
}
