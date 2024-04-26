package servers

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - GRPC Server -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func newGRPCServer(usersService pbs.UsersServiceServer, serverOpts []grpc.ServerOption) *grpc.Server {
	grpcServer := grpc.NewServer(serverOpts...)
	pbs.RegisterUsersServiceServer(grpcServer, usersService)
	return grpcServer
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*       - GRPC Server Options -       */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// -> Server Options are used to configure the GRPC Server.
// -> Our interceptors are actually added here, chained together as a ServerOption.

// Returns the GRPC Server Options, interceptors included.
func defaultGRPCServerOpts(toolbox core.Toolbox, tlsEnabled bool) []grpc.ServerOption {
	serverOpts := []grpc.ServerOption{}

	if tlsEnabled {
		// Add TLS Option if enabled.
		serverOpts = append(serverOpts, grpc.Creds(toolbox.GetTLSTool().GetServerCreds()))
	}

	// Chain all Unary Interceptors into a single ServerOption and add it to the slice.
	defaultInterceptors := defaultGRPCInterceptors(toolbox)
	serverOpts = append(serverOpts, grpc.ChainUnaryInterceptor(defaultInterceptors...))

	return serverOpts
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*        - GRPC Dial Options -        */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Dial Options are used by the HTTP Gateway when connecting to the GRPC Server.
func defaultGRPCDialOpts(tlsClientCreds credentials.TransportCredentials) []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(tlsClientCreds),
		grpc.WithUserAgent(customUserAgent),
	}
}

const customUserAgent = "by @gilperopiola"
