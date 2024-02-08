package adapters

import (
	"github.com/gilperopiola/grpc-gateway-impl/pkg/entities"
	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
)

type Service2TransportAdapter interface {
	Signup(in entities.SignupResponse) (*usersPB.SignupResponse, error)
	Login(in entities.LoginResponse) (*usersPB.LoginResponse, error)
}

type service2transport struct{}

func NewService2TransportAdapter() Service2TransportAdapter {
	return &service2transport{}
}

func (s *service2transport) Signup(in entities.SignupResponse) (*usersPB.SignupResponse, error) {
	return &usersPB.SignupResponse{Message: &usersPB.ResponseMessage{Data: "ok"}}, nil
}

func (s *service2transport) Login(in entities.LoginResponse) (*usersPB.LoginResponse, error) {
	return &usersPB.LoginResponse{Message: &usersPB.ResponseMessage{Data: "ok"}}, nil
}
