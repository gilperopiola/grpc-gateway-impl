package server

import (
	"log"
	"net"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

func newGRPCServer(port string, api usersPB.UsersServiceServer, interceptors []grpc.ServerOption) *GRPCServer {
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
		log.Fatalf(v1.FatalErrMsgStartingGRPC, err)
	}
	go func() {
		if err := g.Server.Serve(lis); err != nil {
			log.Fatalf(v1.FatalErrMsgServingGRPC, err)
		}
	}()
}

// Shutdown gracefully shuts down the gRPC server.
func (g *GRPCServer) Shutdown() {
	log.Println("Shutting down gRPC server...")
	g.Server.GracefulStop()
}

/* ----------------------------------- */
/*        - gRPC Dial Options -        */
/* ----------------------------------- */

const (
	customUserAgent = "gRPC Gateway Implementation by @gilperopiola"
)

// AllDialOptions returns the gRPC dial options.
func AllDialOptions(clientTLSCreds credentials.TransportCredentials) []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(clientTLSCreds),
		grpc.WithUserAgent(customUserAgent),
	}
}
