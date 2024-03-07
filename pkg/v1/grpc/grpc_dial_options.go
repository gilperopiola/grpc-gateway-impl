package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

/* ----------------------------------- */
/*        - gRPC Dial Options -        */
/* ----------------------------------- */

const (
	customUserAgent = "gRPC Gateway Implementation by @gilperopiola"
)

// AllDialOptions returns the gRPC dial options.
func AllDialOptions(clientTLSCreds credentials.TransportCredentials) []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(clientTLSCreds),
		grpc.WithUserAgent(customUserAgent),
	}
}
