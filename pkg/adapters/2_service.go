package adapters

import (
	"github.com/gilperopiola/grpc-gateway-impl/pkg/entities"
	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
)

type Transport2ServiceAdapter interface {
	Signup(in *usersPB.SignupRequest) (entities.SignupRequest, error)
	Login(in *usersPB.LoginRequest) (entities.LoginRequest, error)
}

type transport2service struct{}

func NewTransport2ServiceAdapter() Transport2ServiceAdapter {
	return &transport2service{}
}

func (t *transport2service) Signup(in *usersPB.SignupRequest) (entities.SignupRequest, error) {
	return entities.SignupRequest{Username: in.Username, Password: in.Password}, nil
}

func (t *transport2service) Login(in *usersPB.LoginRequest) (entities.LoginRequest, error) {
	return entities.LoginRequest{Username: in.Username, Password: in.Password}, nil
}
