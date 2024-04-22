package core

import (
	"net/http"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Servers struct {
	GRPC *grpc.Server
	HTTP *http.Server
}

// Sets up both gRPC and HTTP servers.
func SetupServers(usersSvc pbs.UsersServiceServer, tools Toolbox, tlsEnabled bool) *Servers {
	var (
		grpcServerOpts        = defaultGRPCServerOpts(tools, tlsEnabled)
		grpcDialOpts          = defaultGRPCDialOpts(tools.GetTLSClientCreds())
		httpGatewayOpts       = defaultHTTPServeOpts()
		httpGatewayMiddleware = defaultHTTPMiddleware()
	)

	zap.S().Info("GRPC Gateway Implementation | Starting up 🚀")

	return &Servers{
		newGRPCServer(usersSvc, grpcServerOpts),
		newHTTPGateway(httpGatewayOpts, httpGatewayMiddleware, grpcDialOpts),
	}
}

func (s *Servers) Run() {
	runGRPC(s.GRPC)
	runHTTP(s.HTTP)
}

func (s *Servers) Shutdown() {
	shutdownGRPC(s.GRPC)
	shutdownHTTP(s.HTTP)
}
