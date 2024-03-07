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

type GRPCServer struct {
	*grpc.Server

	port         string
	api          usersPB.UsersServiceServer
	interceptors []grpc.ServerOption
}

func NewGRPCServer(port string, api usersPB.UsersServiceServer, interceptors []grpc.ServerOption) *GRPCServer {
	return &GRPCServer{
		port:         port,
		api:          api,
		interceptors: interceptors,
	}
}

// Init initializes the gRPC server and registers the API methods.
// The HTTP Gateway will point towards this server.
// This function also adds the gRPC interceptors to the server.
func (g *GRPCServer) Init() {
	grpcServer := grpc.NewServer(g.interceptors...)
	usersPB.RegisterUsersServiceServer(grpcServer, g.api)
	g.Server = grpcServer
}

// Run runs the gRPC server on a given port.
// It listens for incoming gRPC requests and serves them.
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

// Shutdown gracefully shuts down the gRPC server.
func (g *GRPCServer) Shutdown() {
	log.Println("Shutting down gRPC server...")
	g.Server.GracefulStop()
}
