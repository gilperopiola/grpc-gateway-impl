package server

import (
	"crypto/x509"
	"log"
	"net"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	customUserAgent = "gRPC Gateway Implementation by @gilperopiola"

	errMsgListeningGRPC_Fatal = "Failed to listen gRPC: %v" // Fatal error.
	errMsgServingGRPC_Fatal   = "Failed to serve gRPC: %v"  // Fatal error.
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
	log.Println("Running gRPC!")

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf(errMsgListeningGRPC_Fatal, err)
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf(errMsgServingGRPC_Fatal, err)
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

// getAllDialOptions returns the gRPC dial options.
func getAllDialOptions(tlsEnabled bool, serverCert *x509.CertPool) []grpc.DialOption {
	return []grpc.DialOption{
		newTLSDialOption(tlsEnabled, serverCert),
		grpc.WithUserAgent(customUserAgent),
	}
}

// newTLSDialOption returns a gRPC dial option that enables the client to use TLS.
// If tlsEnabled is false, it returns an insecure dial option.
func newTLSDialOption(tlsEnabled bool, serverCert *x509.CertPool) grpc.DialOption {
	if !tlsEnabled {
		return grpc.WithTransportCredentials(insecure.NewCredentials())
	}
	return grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(serverCert, ""))
}
