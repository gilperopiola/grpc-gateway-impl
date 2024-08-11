package servers

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - GRPC Stuff -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Returns the GRPC Server Options:
// For now it's only TLS + Interceptors.
func getGRPCInterceptors(tools core.Tools, tls bool) []grpc.ServerOption {
	serverOpts := []grpc.ServerOption{}

	if tls {
		serverOpts = append(serverOpts, grpc.Creds(tools.GetServerCreds()))
	}

	// Chain all Interceptors together under a single Server Option.
	serverOpts = append(serverOpts, grpc.ChainUnaryInterceptor(getInterceptors(tools)...))

	return serverOpts
}

// Dial Options are used by the HTTP Gateway when connecting to the GRPC Server.
func getGRPCDialOpts(tlsClientCreds credentials.TransportCredentials) []grpc.DialOption {
	const userAgent = "@gilperopiola"
	return []grpc.DialOption{
		grpc.WithTransportCredentials(tlsClientCreds),
		grpc.WithUserAgent(userAgent),
	}
}
