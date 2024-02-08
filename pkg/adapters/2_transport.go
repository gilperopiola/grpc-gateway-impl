package adapters

import (
	"github.com/gilperopiola/grpc-gateway-impl/pkg/entities"
	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
)

// Service2TransportAdapter is the interface that wraps the methods to convert data from the service format to the transport format.
type Service2TransportAdapter interface {
	Signup(in entities.SignupResponse) (*usersPB.SignupResponse, error)
	Login(in entities.LoginResponse) (*usersPB.LoginResponse, error)
}

// service2transport is our concrete implementation of the Service2TransportAdapter interface.
type service2transport struct{}

// NewService2TransportAdapter returns a new instance of the service2transport.
func NewService2TransportAdapter() Service2TransportAdapter {
	return &service2transport{}
}

// Signup takes a Service Signup Response and returns a PB Signup Response.
func (s *service2transport) Signup(in entities.SignupResponse) (*usersPB.SignupResponse, error) {
	return &usersPB.SignupResponse{Message: &usersPB.ResponseMessage{Data: "ok"}}, nil
}

// Login takes a Service Login Response and returns a PB Login Response.
func (s *service2transport) Login(in entities.LoginResponse) (*usersPB.LoginResponse, error) {
	return &usersPB.LoginResponse{Message: &usersPB.ResponseMessage{Data: "ok"}}, nil
}
