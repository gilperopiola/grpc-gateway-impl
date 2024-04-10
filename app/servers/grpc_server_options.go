package servers

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/modules"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

/* ----------------------------------- */
/*       - gRPC Server Options -       */
/* ----------------------------------- */

// Server Options are used to configure the gRPC Server.
// Our interceptors are actually added here, chained together as a ServerOption.

// AllServerOptions returns the gRPC Server Options.
func AllServerOptions(allModules *modules.All, tlsEnabled bool) []grpc.ServerOption {
	serverOptions := []grpc.ServerOption{}

	// Add TLS Option if enabled.
	if tlsEnabled {
		serverOptions = append(serverOptions, grpc.Creds(allModules.ServerCreds))
	}

	// Chain all Unary Interceptors into a single ServerOption and add it to the slice.
	unaryInterceptorSvOpt := grpc.ChainUnaryInterceptor(getUnaryInterceptors(allModules)...)
	serverOptions = append(serverOptions, unaryInterceptorSvOpt)

	return serverOptions
}

/* ----------------------------------- */
/*        - gRPC Dial Options -        */
/* ----------------------------------- */

// Dial Options are used by the HTTP Gateway when connecting to the gRPC Server.

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
