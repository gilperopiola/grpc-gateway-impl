package transport

import (
	"context"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/adapters"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/service"
	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"

	"github.com/bufbuild/protovalidate-go"
)

type TransportLayer interface {
	Signup(ctx context.Context, pbRequest *usersPB.SignupRequest) (*usersPB.SignupResponse, error)
	Login(ctx context.Context, pbRequest *usersPB.LoginRequest) (*usersPB.LoginResponse, error)
}

type transport struct {
	protoValidator *protovalidate.Validator
	service        service.ServiceLayer
	toService      adapters.Transport2ServiceAdapter
	toTransport    adapters.Service2TransportAdapter
}

func NewTransport(protoValidator *protovalidate.Validator, serviceLayer service.ServiceLayer) *transport {
	return &transport{
		protoValidator: protoValidator,
		service:        serviceLayer,
		toService:      adapters.NewTransport2ServiceAdapter(),
		toTransport:    adapters.NewService2TransportAdapter(),
	}
}
