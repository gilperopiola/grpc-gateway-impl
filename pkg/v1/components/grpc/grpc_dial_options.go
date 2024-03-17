package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

/* ----------------------------------- */
/*        - gRPC Dial Options -        */
/* ----------------------------------- */

/* Dial Options are used by the HTTP Gateway when connecting to the gRPC server. */

const (
	customUserAgent = "by @gilperopiola"
)

// AllDialOptions returns the gRPC dial options.
func AllDialOptions(tlsClientCreds credentials.TransportCredentials) []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(tlsClientCreds),
		grpc.WithUserAgent(customUserAgent),
	}
}
