package server

import (
	"crypto/x509"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	tlsCertPath = "./server.crt"
	tlsKeyPath  = "./server.key"
)

// NewGRPCServerCredentials loads the server's certificate and key files and returns a gRPC server option
// that enables the server to use TLS.
func NewGRPCServerCredentials() grpc.ServerOption {
	creds, err := credentials.NewServerTLSFromFile(tlsCertPath, tlsKeyPath)
	if err != nil {
		log.Fatalf(msgErrLoadingTLSCredentials_Fatal, err)
	}
	return grpc.Creds(creds)
}

// LoadTLSCertPool loads the server's certificate from a file and returns a certificate pool.
// It's a SSL/TLS certificate used to secure the communication between the HTTP Gateway and the gRPC server.
// It must be in a .crt format.
//
// To generate a self-signed certificate, you can use the following command:
// openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes -subj '/CN=localhost'
// The certificate must be in the root directory of the project.
func LoadTLSCertPool() (out *x509.CertPool) {

	// Read certificate.
	cert, err := os.ReadFile(tlsCertPath)
	if err != nil {
		log.Fatalf(msgErrReadingTLSCert_Fatal, err)
	}

	// Create certificate pool.
	if out = x509.NewCertPool(); !out.AppendCertsFromPEM(cert) {
		log.Fatalf(msgErrAppendingTLSCert_Fatal)
	}

	return out
}

const (
	msgErrLoadingTLSCredentials_Fatal = "Failed to load server TLS credentials: %v"
	msgErrReadingTLSCert_Fatal        = "Failed to read TLS certificate: %v"
	msgErrAppendingTLSCert_Fatal      = "Failed to append TLS certificate"
)
