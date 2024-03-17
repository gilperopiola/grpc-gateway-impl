package grpc

import (
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

/* ----------------------------------- */
/*       - gRPC Server Options -       */
/* ----------------------------------- */

/* Server Options are used to configure the gRPC Server.
/* Our interceptors are actually added here, chained together as a ServerOption. */

// AllServerOptions returns the gRPC Server Options.
func AllServerOptions(components *components.Wrapper, tlsEnabled bool) []grpc.ServerOption {
	serverOptions := []grpc.ServerOption{}

	// Add TLS Option if enabled.
	if tlsEnabled {
		tlsServerOption := grpc.Creds(components.ServerCreds)
		serverOptions = append(serverOptions, tlsServerOption)
	}

	// Chain all Unary Interceptors into a single ServerOption and add it to the slice.
	unaryInterceptorsServerOption := grpc.ChainUnaryInterceptor(getUnaryInterceptors(components)...)
	serverOptions = append(serverOptions, unaryInterceptorsServerOption)

	return serverOptions
}

/* ----------------------------------- */
/*        - gRPC Dial Options -        */
/* ----------------------------------- */

/* Dial Options are used by the HTTP Gateway when connecting to the gRPC server. */

// AllDialOptions returns the gRPC Dial Options.
func AllDialOptions(tlsClientCreds credentials.TransportCredentials) []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(tlsClientCreds),
		grpc.WithUserAgent(customUserAgent),
	}
}

const (
	customUserAgent = "by @gilperopiola"
)
