package modules

import (
	"crypto/x509"
	"os"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

	"go.uber.org/zap"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// To generate a self-signed certificate, you can use the following command:
// openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes -subj '/CN=localhost'
// The certificate must be in the root directory of the project.

// NewTLSCertPool loads the server's certificate from a file and returns a certificate pool.
// It's a SSL/TLS certificate used to secure the communication between the HTTP Gateway and the gRPC Server.
// It must be in a .crt format.
func NewTLSCertPool(certPath string) *x509.CertPool {
	certPool := x509.NewCertPool() // Create certificate pool.

	cert, err := os.ReadFile(certPath) // Read certificate.
	if err != nil {
		zap.S().Fatalf(errs.FatalErrMsgReadingTLSCert, err)
	}

	if !certPool.AppendCertsFromPEM(cert) { // Append encoded certificate.
		zap.S().Fatalf(errs.FatalErrMsgAppendingTLSCert)
	}

	return certPool
}

// NewServerTransportCreds returns the Server's transport credentials.
func NewServerTransportCreds(certPath, keyPath string) credentials.TransportCredentials {
	creds, err := credentials.NewServerTLSFromFile(certPath, keyPath)
	if err != nil {
		zap.S().Fatalf(errs.FatalErrMsgLoadingTLSCreds, err)
	}
	return creds
}

// NewClientTransportCreds returns the client's transport credentials, either secure or insecure.
func NewClientTransportCreds(tlsEnabled bool, svCert *x509.CertPool) credentials.TransportCredentials {
	if tlsEnabled {
		return credentials.NewClientTLSFromCert(svCert, "")
	}
	return insecure.NewCredentials()
}
