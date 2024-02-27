package server

import (
	"crypto/x509"
	"log"
	"os"

	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// LoadTLSCertPool loads the server's certificate from a file and returns a certificate pool.
// It's a SSL/TLS certificate used to secure the communication between the HTTP Gateway and the gRPC server.
// It must be in a .crt format.
//
// To generate a self-signed certificate, you can use the following command:
// openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes -subj '/CN=localhost'
// The certificate must be in the root directory of the project.
func LoadTLSCertPool(tlsCertPath string) *x509.CertPool {
	// Read certificate.
	cert, err := os.ReadFile(tlsCertPath)
	if err != nil {
		log.Fatalf(msgErrReadingTLSCert_Fatal, err)
	}

	// Create certificate pool.
	if out := x509.NewCertPool(); out.AppendCertsFromPEM(cert) {
		return out
	}

	// Error appending certificate.
	log.Fatalf(msgErrAppendingTLSCert_Fatal)
	return nil
}

// newTLSSecurityServerOption returns a grpc.ServerOption that enables TLS communication.
// It loads the server's certificate and key from a file.
func newTLSSecurityServerOption(tlsConfig v1.TLSConfig) grpc.ServerOption {
	creds, err := credentials.NewServerTLSFromFile(tlsConfig.CertPath, tlsConfig.KeyPath)
	if err != nil {
		log.Fatalf(msgErrLoadingTLSCredentials_Fatal, err)
	}
	return grpc.Creds(creds)
}

const (
	msgErrLoadingTLSCredentials_Fatal = "Failed to load server TLS credentials: %v"
	msgErrReadingTLSCert_Fatal        = "Failed to read TLS certificate: %v"
	msgErrAppendingTLSCert_Fatal      = "Failed to append TLS certificate"
)
