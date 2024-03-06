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

// initGRPCServer initializes the gRPC server and registers the API methods.
// The HTTP Gateway will point towards this server.
// This function also adds the gRPC interceptors to the server.
func initGRPCServer(api usersPB.UsersServiceServer, interceptors []grpc.ServerOption) *grpc.Server {
	grpcServer := grpc.NewServer(interceptors...)
	usersPB.RegisterUsersServiceServer(grpcServer, api)
	return grpcServer
}

// runGRPCServer runs the gRPC server on a given port.
// It listens for incoming gRPC requests and serves them.
func runGRPCServer(grpcServer *grpc.Server, grpcPort string) {
	log.Printf("Running gRPC on port %s!\n", grpcPort)

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf(v1.FatalErrMsgStartingGRPC, err)
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf(v1.FatalErrMsgServingGRPC, err)
		}
	}()
}

// shutdownGRPCServer gracefully shuts down the gRPC server.
func shutdownGRPCServer(grpcServer *grpc.Server) {
	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
}

/* ----------------------------------- */
/*        - gRPC Dial Options -        */
/* ----------------------------------- */

const (
	customUserAgent = "gRPC Gateway Implementation by @gilperopiola"
)

// getAllDialOptions returns the gRPC dial options.
func getAllDialOptions(clientTLSCredentials credentials.TransportCredentials) []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(clientTLSCredentials),
		grpc.WithUserAgent(customUserAgent),
	}
}
