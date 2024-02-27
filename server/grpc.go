package server

import (
	"log"
	"net"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"

	"google.golang.org/grpc"
)

/* ----------------------------------- */
/*           - gRPC Server -           */
/* ----------------------------------- */

// InitGRPCServer initializes the gRPC server and registers the API methods.
// The HTTP Gateway will point towards this server.
// This function also adds the gRPC interceptors to the server.
func InitGRPCServer(api usersPB.UsersServiceServer, tlsConfig v1.TLSConfig, interceptors grpc.ServerOption) *grpc.Server {
	serverOptions := addTLSToInterceptors(tlsConfig, interceptors)
	grpcServer := grpc.NewServer(serverOptions...)
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

// addTLSToInterceptors returns the gRPC server options.
// It returns a TLS server option if tlsEnabled is true and the given interceptors.
func addTLSToInterceptors(tlsConfig v1.TLSConfig, interceptors grpc.ServerOption) v1.InterceptorsI {
	options := v1.InterceptorsI{}
	if tlsConfig.Enabled {
		options = append(options, newTLSSecurityServerOption(tlsConfig))
	}
	options = append(options, interceptors)
	return options
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
