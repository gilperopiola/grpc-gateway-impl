package servers

import (
	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"

	"google.golang.org/grpc"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - GRPC Stuff -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// GRPC Server Options are used to configure the Server.
// GRPC Dial Options are used to configure the connection between HTTP and GRPC.

// Returns the GRPC Server Options:
// For now it's only TLS + Interceptors.
func getGRPCServerOpts(tools core.Tools, tls bool) []grpc.ServerOption {
	serverOpts := []grpc.ServerOption{}

	if tls {
		serverOpts = append(serverOpts, grpc.Creds(tools.GetServerCreds()))
	}

	// Chain all Interceptors together under a single Server Option.
	serverOpts = append(serverOpts, grpc.ChainUnaryInterceptor(getInterceptors(tools)...))

	return serverOpts
}

// Dial Options are used by the HTTP Gateway when connecting to the GRPC Server.
func getGRPCDialOpts(tlsClientCreds god.TLSCreds) []grpc.DialOption {
	const userAgent = "by @gilperopiola"
	return []grpc.DialOption{
		grpc.WithTransportCredentials(tlsClientCreds),
		grpc.WithUserAgent(userAgent),
	}
}
