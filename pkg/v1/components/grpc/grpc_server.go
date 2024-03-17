package grpc

import (
	"log"
	"net"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"

	"google.golang.org/grpc"
)

/* ----------------------------------- */
/*           - gRPC Server -           */
/* ----------------------------------- */

// GRPCServer is a wrapper around the actual gRPC Server.
type GRPCServer struct {
	*grpc.Server

	service usersPB.UsersServiceServer
	port    string
	options []grpc.ServerOption
}

// NewGRPCServer returns a new instance of GRPCServer.
func NewGRPCServer(port string, api usersPB.UsersServiceServer, options []grpc.ServerOption) *GRPCServer {
	return &GRPCServer{
		port:    port,
		service: api,
		options: options,
	}
}

// Init initializes the gRPC Server, adds the interceptors and registers the API methods.
func (g *GRPCServer) Init() {
	g.Server = grpc.NewServer(g.options...)
	usersPB.RegisterUsersServiceServer(g.Server, g.service)
}

// Run makes the gRPC Server listen for incoming gRPC requests and serves them.
func (g *GRPCServer) Run() {
	log.Printf("Running gRPC on port %s!\n", g.port)

	lis, err := net.Listen("tcp", g.port)
	if err != nil {
		log.Fatalf(errs.FatalErrMsgStartingGRPC, err)
	}

	go func() {
		if err := g.Server.Serve(lis); err != nil {
			log.Fatalf(errs.FatalErrMsgServingGRPC, err)
		}
	}()
}

// Shutdown gracefully shuts down the gRPC Server.
func (g *GRPCServer) Shutdown() {
	log.Println("Shutting down gRPC server...")
	g.Server.GracefulStop()
}
