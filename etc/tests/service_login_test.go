package tests

//func TestServiceLogin(t *testing.T) {
//	expect := func(response *pbs.LoginResponse, err error) (*pbs.LoginResponse, error) {
//		return copyPB(response).(*pbs.LoginResponse), err
//	}
//
//	setupMock := func(h modules.PwdHasher, getResult *models.User, getErr error) setupRepoMockFn {
//		getResult.Password = hashPwd(h, getResult.Password)
//		return func(repoMock *mocks.Storage) {
//			repoMock.On("GetUser", mock.Anything, mock.Anything).Return(copyUser(getResult), getErr).Once()
//		}
//	}
//
//	tests := []struct {
//		name        string
//		request     *pbs.LoginRequest
//		setupMock   func(modules.PwdHasher) setupRepoMockFn
//		getExpected func(modules.TokenGenerator) (*pbs.LoginResponse, error)
//	}{
//		{
//			name:    "tc_service_login_ok",
//			request: loginRequest,
//			setupMock: func(h modules.PwdHasher) setupRepoMockFn {
//				return setupMock(h, userWithID, nil)
//			},
//			getExpected: func(g modules.TokenGenerator) (*pbs.LoginResponse, error) {
//				return expect(&pbs.LoginResponse{Token: generateToken(g, userID, username, "")}, nil)
//			},
//		},
//		//{
//		//	name:        "tc_service_login_not_found",
//		//	setupMock:   setupMock(nil, gorm.ErrRecordNotFound),
//		//	getExpected: expect(nil, service.ErrNotFound("user")),
//		//},
//		//{
//		//	name:        "tc_service_login_get_error",
//		//	setupMock:   setupMock(nil, errors.New("unknown error")),
//		//	getExpected: expect(nil, service.UserErr(context.Background(), errors.New("unknown error"))),
//		//},
//		//{
//		//	name:        "tc_service_login_password_mismatch",
//		//	setupMock:   setupMock(&models.User{ID: "1", Username: "user1", Password: "hashedpassword"}, nil),
//		//	getExpected: expect(nil, service.ErrUnauthenticated()),
//		//},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			tokenGen, pwdHasher := newTestServiceModules() // Prepare
//			service, repoMock := newTestService(test.setupMock(pwdHasher), tokenGen, pwdHasher)
//			expectedResponse, expectedErr := test.getExpected(tokenGen)
//
//			response, err := service.Login(context.Background(), test.request) // Act
//
//			assert.True(t, proto.Equal(expectedResponse, response)) // Assert
//			assert.Equal(t, expectedErr, err)
//			repoMock.AssertExpectations(t)
//		})
//	}
//}
//
