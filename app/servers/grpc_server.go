package servers

import (
	"net"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

/* ----------------------------------- */
/*           - gRPC Server -           */
/* ----------------------------------- */

// GRPCServer is a wrapper around the actual gRPC Server.
type GRPCServer struct {
	*grpc.Server
}

// NewGRPCServer returns a new instance of GRPCServer.
func NewGRPCServer(service pbs.UsersServiceServer, options []grpc.ServerOption) *GRPCServer {
	grpcServer := &GRPCServer{Server: grpc.NewServer(options...)}
	pbs.RegisterUsersServiceServer(grpcServer.Server, service)
	return grpcServer
}

// Run makes the gRPC Server listen for incoming gRPC requests and serves them.
func (g *GRPCServer) Run() {
	zap.S().Infof("Running gRPC on port %s!\n", core.GRPCPort)

	lis, err := net.Listen("tcp", core.GRPCPort)
	if err != nil {
		zap.S().Fatalf(errs.FatalErrMsgStartingGRPC, err)
	}

	go func() {
		if err := g.Server.Serve(lis); err != nil {
			zap.S().Fatalf(errs.FatalErrMsgServingGRPC, err)
		}
	}()
}

// Shutdown gracefully shuts down the gRPC Server.
func (g *GRPCServer) Shutdown() {
	zap.S().Info("Shutting down gRPC server...")
	g.Server.GracefulStop()
}
