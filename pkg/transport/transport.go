package transport

import (
	"context"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/adapters"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/service"
	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"

	"github.com/bufbuild/protovalidate-go"
)

// TransportLayer is the interface that wraps the transport methods.
type TransportLayer interface {
	Signup(ctx context.Context, pbRequest *usersPB.SignupRequest) (*usersPB.SignupResponse, error)
	Login(ctx context.Context, pbRequest *usersPB.LoginRequest) (*usersPB.LoginResponse, error)
}

// transport is our concrete implementation of the TransportLayer interface.
type transport struct {
	protoValidator *protovalidate.Validator
	service        service.ServiceLayer
	toService      adapters.Transport2ServiceAdapter
	toTransport    adapters.Service2TransportAdapter
}

// NewTransport returns a new instance of the transport.
func NewTransport(protoValidator *protovalidate.Validator, serviceLayer service.ServiceLayer) *transport {
	return &transport{
		protoValidator: protoValidator,
		service:        serviceLayer,
		toService:      adapters.NewTransport2ServiceAdapter(),
		toTransport:    adapters.NewService2TransportAdapter(),
	}
}
