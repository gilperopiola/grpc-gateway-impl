package service

import (
	"context"
	"encoding/json"
	"math/rand"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
)

/* ----------------------------------- */
/*            - Service -              */
/* ----------------------------------- */

// ServiceLayer is the interface that wraps the service methods. All the business logic should be implemented here.
type ServiceLayer interface {
	Signup(ctx context.Context, in *usersPB.SignupRequest) (*usersPB.SignupResponse, error)
	Login(ctx context.Context, in *usersPB.LoginRequest) (*usersPB.LoginResponse, error)
}

// service is our concrete implementation of the ServiceLayer interface.
type service struct{}

// NewService returns a new instance of the service.
func NewService() *service {
	return &service{}
}

func successfulResponse(data string) *usersPB.ResponseMessage {
	return &usersPB.ResponseMessage{
		Success: usersPB.Success_TRUE,
		Data:    data,
		Error:   "",
	}
}

func structToJSONString(s interface{}) string {
	str, _ := json.Marshal(s)
	return string(str)
}

// Signup should be the implementation of the Signup service method.
func (s *service) Signup(ctx context.Context, in *usersPB.SignupRequest) (*usersPB.SignupResponse, error) {

	// ... check username is available, hash password, create user in DB, etc.

	// if err := something(); err != nil {
	// 		return entities.SignupResponse{}, fmt.Errorf("error in something(): %w", err)
	// }

	out := struct {
		ID int
	}{
		ID: rand.Intn(1000),
	}

	return &usersPB.SignupResponse{Message: successfulResponse(structToJSONString(out))}, nil
}

// Login should be the implementation of the Login service method.
func (s *service) Login(ctx context.Context, in *usersPB.LoginRequest) (*usersPB.LoginResponse, error) {

	// ... get user from DB, hash password, compare passwords, etc.

	// if err := something(); err != nil {
	// 		return entities.LoginResponse{}, fmt.Errorf("error in something(): %w", err)
	// }

	out := struct {
		Token string
	}{
		Token: "some.jwt.token",
	}

	return &usersPB.LoginResponse{Message: successfulResponse(structToJSONString(out))}, nil
}
