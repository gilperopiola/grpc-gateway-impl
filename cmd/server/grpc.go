package server

import (
	"log"
	"net"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"
	v1Service "github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"

	"google.golang.org/grpc"
)

/* ----------------------------------- */
/*             - gRPC -                */
/* ----------------------------------- */

// InitGRPCServer initializes the gRPC server and registers the API methods.
// The HTTP Gateway will point towards this server.
// This function also adds the gRPC interceptors.
func InitGRPCServer(service v1Service.ServiceLayer) *grpc.Server {
	grpcServer := grpc.NewServer(v1.GetInterceptorsAsServerOptions())
	usersPB.RegisterUsersServiceServer(grpcServer, v1.NewAPI(service))
	return grpcServer
}

// RunGRPCServer runs the gRPC server on a given port.
// It listens for incoming gRPC requests and serves them.
func RunGRPCServer(grpcServer *grpc.Server, grpcPort string) {
	log.Println("Running gRPC!")

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf(errMsgListenGRPC, err)
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf(errMsgServeGRPC, err)
		}
	}()
}

// ShutdownGRPCServer gracefully shuts down the gRPC server.
func ShutdownGRPCServer(grpcServer *grpc.Server) {
	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
}

const (
	errMsgListenGRPC = "Failed to listen gRPC: %v"
	errMsgServeGRPC  = "Failed to serve gRPC: %v"
)
