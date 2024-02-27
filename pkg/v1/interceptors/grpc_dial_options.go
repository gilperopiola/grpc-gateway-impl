package interceptors

import (
	"crypto/x509"

	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

/* ----------------------------------- */
/*        - gRPC Dial Options -        */
/* ----------------------------------- */

func GetAllDialOptions(tlsConfig v1.TLSConfig, serverCert *x509.CertPool) v1.GRPCDialOptionsI {
	return v1.GRPCDialOptionsI{
		newTLSSecurityDialOption(tlsConfig.Enabled, serverCert),
		grpc.WithUserAgent("gRPC Gateway Implementation by @gilperopiola"),
	}
}

// newTLSSecurityDialOption returns a gRPC dial option that enables the client to use TLS.
// If tlsEnabled is false, it returns an insecure dial option.
func newTLSSecurityDialOption(tlsEnabled bool, serverCert *x509.CertPool) grpc.DialOption {
	if !tlsEnabled {
		return grpc.WithTransportCredentials(insecure.NewCredentials())
	}
	return grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(serverCert, ""))
}
