package security

import (
	"crypto/x509"
	"log"
	"os"

	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// To generate a self-signed certificate, you can use the following command:
// openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes -subj '/CN=localhost'
// The certificate must be in the root directory of the project.

// NewServerTransportCredentials returns the server's transport credentials.
func NewServerTransportCredentials(certPath, keyPath string) credentials.TransportCredentials {
	creds, err := credentials.NewServerTLSFromFile(certPath, keyPath)
	if err != nil {
		log.Fatalf(v1.FatalErrMsgLoadingTLSCredentials, err)
	}
	return creds
}

// NewClientTransportCredentials returns the client's transport credentials.
func NewClientTransportCredentials(tlsEnabled bool, serverCert *x509.CertPool) credentials.TransportCredentials {
	if tlsEnabled {
		return credentials.NewClientTLSFromCert(serverCert, "")
	}
	return insecure.NewCredentials()
}

// NewTLSCertPool loads the server's certificate from a file and returns a certificate pool.
// It's a SSL/TLS certificate used to secure the communication between the HTTP Gateway and the gRPC server.
// It must be in a .crt format.
func NewTLSCertPool(tlsCertPath string) *x509.CertPool {

	// Create certificate pool.
	out := x509.NewCertPool()

	// Read certificate.
	cert, err := os.ReadFile(tlsCertPath)
	if err != nil {
		log.Fatalf(v1.FatalErrMsgReadingTLSCert, err)
	}

	// Append encoded certificate.
	if !out.AppendCertsFromPEM(cert) {
		log.Fatalf(v1.FatalErrMsgAppendingTLSCert)
	}

	return out
}
