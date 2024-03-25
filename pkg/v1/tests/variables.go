package tests

import (
	"errors"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"
)

var (
	id       = 1
	username = "username"
	password = "password"

	user       = &models.User{Username: username, Password: password}
	userEmpty  = &models.User{}
	userWithID = &models.User{ID: 1, Username: username, Password: password}

	errCreatingUser = errors.New("error creating user")
	errGettingUser  = errors.New("error getting user")
)
