package tools

import (
	"crypto/x509"
	"errors"
	"fmt"
	"os"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// NOTE: All TLS-related files must be in the root folder.

// To generate a self-signed TLS certificate, you can use
// -> openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes -subj '/CN=localhost'

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Loads the server's certificate from a file and returns a *x509.CertPool.
// It's only holds 1 TLS certficate, used to secure all communications with the gRPC Server.
// It must be in a .crt file.
func newTLSCertPool(certPath string) *x509.CertPool {
	certPool := x509.NewCertPool()

	cert, err := os.ReadFile(certPath)
	if err != nil {
		core.LogUnexpectedAndPanic(fmt.Errorf(errs.FatalErrMsgReadingTLSCert, err))
	}

	if !certPool.AppendCertsFromPEM(cert) {
		core.LogUnexpectedAndPanic(errors.New(errs.FatalErrMsgAppendingTLSCert))
	}

	return certPool
}

// Returns the Server's transport credentials. Only called if TLS is enabled.
func newServerTransportCreds(certPath, keyPath string) credentials.TransportCredentials {
	creds, err := credentials.NewServerTLSFromFile(certPath, keyPath)
	if err != nil {
		core.LogUnexpectedAndPanic(fmt.Errorf(errs.FatalErrMsgLoadingTLSCreds, err))
	}
	return creds
}

// Returns the client's transport credentials, either secure or insecure.
func newClientTransportCreds(tlsEnabled bool, serverCert *x509.CertPool) credentials.TransportCredentials {
	if tlsEnabled {
		return credentials.NewClientTLSFromCert(serverCert, "")
	}
	return insecure.NewCredentials()
}
