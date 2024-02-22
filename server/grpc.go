package server

import (
	"log"
	"net"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"

	"google.golang.org/grpc"
)

/* ----------------------------------- */
/*           - gRPC Server -           */
/* ----------------------------------- */

// InitGRPCServer initializes the gRPC server and registers the API methods.
// The HTTP Gateway will point towards this server.
// This function also adds the gRPC interceptors to the server.
func InitGRPCServer(api usersPB.UsersServiceServer, interceptors grpc.ServerOption) *grpc.Server {
	options := []grpc.ServerOption{NewGRPCServerCredentials(), interceptors}
	grpcServer := grpc.NewServer(options...)
	usersPB.RegisterUsersServiceServer(grpcServer, api)
	return grpcServer
}

// RunGRPCServer runs the gRPC server on a given port.
// It listens for incoming gRPC requests and serves them.
func RunGRPCServer(grpcServer *grpc.Server, grpcPort string) {
	log.Println("Running gRPC!")

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf(msgErrListeningGRPC_Fatal, err)
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf(msgErrServingGRPC_Fatal, err)
		}
	}()
}

// ShutdownGRPCServer gracefully shuts down the gRPC server.
func ShutdownGRPCServer(grpcServer *grpc.Server) {
	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
}

const (
	msgErrListeningGRPC_Fatal = "Failed to listen gRPC: %v"
	msgErrServingGRPC_Fatal   = "Failed to serve gRPC: %v"
)
