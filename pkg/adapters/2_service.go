package adapters

import (
	"github.com/gilperopiola/grpc-gateway-impl/pkg/entities"
	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
)

// Transport2ServiceAdapter is the interface that wraps the methods to convert data from the transport format to the service format.
type Transport2ServiceAdapter interface {
	Signup(in *usersPB.SignupRequest) (entities.SignupRequest, error)
	Login(in *usersPB.LoginRequest) (entities.LoginRequest, error)
}

// transport2service is our concrete implementation of the Transport2ServiceAdapter interface.
type transport2service struct{}

// NewTransport2ServiceAdapter returns a new instance of the transport2service.
func NewTransport2ServiceAdapter() Transport2ServiceAdapter {
	return &transport2service{}
}

// Signup takes a PB Signup Request and returns a Service Signup Request.
func (t *transport2service) Signup(in *usersPB.SignupRequest) (entities.SignupRequest, error) {
	return entities.SignupRequest{Username: in.Username, Password: in.Password}, nil
}

// Login takes a PB Login Request and returns a Service Login Request.
func (t *transport2service) Login(in *usersPB.LoginRequest) (entities.LoginRequest, error) {
	return entities.LoginRequest{Username: in.Username, Password: in.Password}, nil
}
