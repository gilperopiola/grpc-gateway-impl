package tools

import (
	"crypto/x509"
	"errors"
	"os"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var _ core.TLSTool = (*tlsTool)(nil)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - TLS Tool -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// NOTE: All TLS related files must be in the root folder.
// -> To generate a self-signed TLS certificate, you can use
// -> openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes -subj '/CN=localhost'

type tlsTool struct {
	ServerCertPool *x509.CertPool
	ServerCreds    credentials.TransportCredentials
	ClientCreds    credentials.TransportCredentials
}

func NewTLSTool(cfg *core.TLSCfg) core.TLSTool {
	tlsTool := &tlsTool{}
	tlsTool.ServerCertPool = newTLSCertPool(cfg.CertPath)
	tlsTool.ServerCreds = newServerTransportCreds(cfg.CertPath, cfg.KeyPath)
	tlsTool.ClientCreds = newClientTransportCreds(tlsTool.ServerCertPool)
	return tlsTool
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (t *tlsTool) GetServerCertificate() *x509.CertPool             { return t.ServerCertPool }
func (t *tlsTool) GetServerCreds() credentials.TransportCredentials { return t.ServerCreds }
func (t *tlsTool) GetClientCreds() credentials.TransportCredentials { return t.ClientCreds }

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Loads the server's certificate from a file and returns a *x509.CertPool.
// It's only holds 1 TLS certficate, used to secure all communications with the GRPC Server.
// It must be in a .crt file.
func newTLSCertPool(certPath string) *x509.CertPool {
	if !core.TLSEnabled {
		return nil
	}

	cert, err := os.ReadFile(certPath)
	core.LogFatalIfErr(err, errs.FailedToReadTLSCert)

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(cert) {
		core.LogFatal(errors.New(errs.FailedToAppendTLSCert))
	}

	return certPool
}

// Returns the Server's transport credentials. Only called if TLS is enabled.
func newServerTransportCreds(certPath, keyPath string) credentials.TransportCredentials {
	if !core.TLSEnabled {
		return nil
	}
	creds, err := credentials.NewServerTLSFromFile(certPath, keyPath)
	core.LogFatalIfErr(err, errs.FailedToLoadTLSCreds)
	return creds
}

// Returns the client's transport credentials, either secure or insecure.
func newClientTransportCreds(serverCertPool *x509.CertPool) credentials.TransportCredentials {
	if !core.TLSEnabled {
		core.LogImportant("TLS is not enabled! 🔒")
		return insecure.NewCredentials()
	}
	return credentials.NewClientTLSFromCert(serverCertPool, "")
}
