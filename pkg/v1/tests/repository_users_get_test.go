package tests

/* ----------------------------------- */
/*      - Users Repository Tests -     */
/* ----------------------------------- */

/* - GetUser
/*
/* --------------------------------------------------------------------------- */
/*     - repository.GetUser(...options.QueryOption) (*models.User, error) -    */
/* --------------------------------------------------------------------------- */
//
// - GormAdapter.Model(&models.User{}) 	-> (GormAdapter)
// - GormAdapter.First(&user) 				-> (GormAdapter)
// - GormAdapter.Error() 						-> (error)

//func TestRepositoryGetUser(t *testing.T) {
//	var (
//		id       = 1
//		username = "username"
//
//		argOptWithID = options.WithUserID(id)
//		// argOptUsername = options.WithUsername("username")
//
//		// errOnFirst = errors.New("error getting user")
//	)
//
//	testRepositoryGetUserWithIDOK(t, argOptWithID, id, username)
//	// testRepositoryGetUserWithUsernameOK(t, argOptUsername)
//}
//
//func testRepositoryGetUserWithIDOK(t *testing.T, argOptID options.QueryOption, id int, username string) {
//	repository := setupNewTest()
//	gormMock := repository.GetGormMock()
//
//	// - Mock Calls
//	userToGet := &models.User{}
//	gormMock.On("Model", userToGet).Return(gormMock).Once()
//
//	gormMock.On("Where", "id = ?", []interface{}{fmt.Sprintf("%d", id)}).Return(gormMock).Once()
//
//	gormMock.On("First", userToGet).Return(gormMock).Once()
//
//	gormMock.On("Error").Return(nil).Once()
//
//	// - Act & Assert
//	user, err := repository.GetUser(argOptID)
//	assert.NotNil(t, user)
//	assert.Nil(t, err)
//}
