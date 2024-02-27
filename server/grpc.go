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

// getAllDialOptions returns the gRPC dial options.
func getAllDialOptions(tlsEnabled bool, serverCert *x509.CertPool) []grpc.DialOption {
	return []grpc.DialOption{
		newTLSSecurityDialOption(tlsEnabled, serverCert),
		grpc.WithUserAgent("gRPC Gateway Implementation by @gilperopiola"),
	}
}

/* ----------------------------------- */
/*           - gRPC Server -           */
/* ----------------------------------- */

// InitGRPCServer initializes the gRPC server and registers the API methods.
// The HTTP Gateway will point towards this server.
// This function also adds the gRPC interceptors to the server.
func InitGRPCServer(api usersPB.UsersServiceServer, tlsConfig TLSConfig, interceptors grpc.ServerOption) *grpc.Server {
	serverOptions := addTLSToInterceptors(tlsConfig, interceptors)
	grpcServer := grpc.NewServer(serverOptions...)
	usersPB.RegisterUsersServiceServer(grpcServer, api)
	return grpcServer
}

// runGRPCServer runs the gRPC server on a given port.
// It listens for incoming gRPC requests and serves them.
func runGRPCServer(grpcServer *grpc.Server, grpcPort string) {
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
func addTLSToInterceptors(tlsConfig TLSConfig, interceptors grpc.ServerOption) []grpc.ServerOption {
	options := []grpc.ServerOption{}
	if tlsConfig.Enabled {
		options = append(options, newTLSSecurityServerOption(tlsConfig))
	}
	options = append(options, interceptors)
	return options
}

// shutdownGRPCServer gracefully shuts down the gRPC server.
func shutdownGRPCServer(grpcServer *grpc.Server) {
	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
}

const (
	msgErrListeningGRPC_Fatal = "Failed to listen gRPC: %v"
	msgErrServingGRPC_Fatal   = "Failed to serve gRPC: %v"
)

/* ----------------------------------- */
/*        - gRPC Dial Options -        */
/* ----------------------------------- */

// newTLSSecurityDialOption returns a gRPC dial option that enables the client to use TLS.
// If tlsEnabled is false, it returns an insecure dial option.
func newTLSSecurityDialOption(tlsEnabled bool, serverCert *x509.CertPool) grpc.DialOption {
	if !tlsEnabled {
		return grpc.WithTransportCredentials(insecure.NewCredentials())
	}
	return grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(serverCert, ""))
}
