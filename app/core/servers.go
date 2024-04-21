package core

import (
	"net/http"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	"google.golang.org/grpc"
)

type Servers struct {
	GRPC *grpc.Server
	HTTP *http.Server
}

func SetupServers(tools ToolsAccessor, usersService pbs.UsersServiceServer, tlsEnabled bool) *Servers {
	return &Servers{
		NewGRPCServer(usersService, AllServerOptions(tools, tlsEnabled)),
		NewHTTPGateway(MiddlewareServeOpts(), MiddlewareWrapper(), AllDialOptions(tools.GetTLSClientCreds())),
	}
}

func (s *Servers) Run() {
	RunGRPCServer(s.GRPC)
	RunHTTPGateway(s.HTTP)
}

func (s *Servers) Shutdown() {
	ShutdownGRPCServer(s.GRPC)
	ShutdownHTTPGateway(s.HTTP)
}
