package tools

import (
	"crypto/x509"
	"errors"
	"os"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/logs"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var _ core.TLSTool = &tlsTool{}

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
	tlsTool := tlsTool{}

	// Loads the server's certificate from a .crt file into a *x509.CertPool.
	// It holds 1 TLS cert, used to secure all communications with the GRPC Server.
	tlsTool.ServerCertPool = func(certPath string) *x509.CertPool {
		if !core.G.TLSEnabled {
			return nil
		}

		cert, err := os.ReadFile(certPath)
		logs.LogFatalIfErr(err, errs.FailedToReadTLSCert)

		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(cert) {
			logs.LogFatal(errors.New(errs.FailedToAppendTLSCert))
		}

		return certPool
	}(cfg.CertPath)

	// Returns the Server's transport credentials. Only called if TLS is enabled.
	tlsTool.ServerCreds = func(certPath, keyPath string) credentials.TransportCredentials {
		if !core.G.TLSEnabled {
			return nil
		}

		creds, err := credentials.NewServerTLSFromFile(certPath, keyPath)
		logs.LogFatalIfErr(err, errs.FailedToLoadTLSCreds)
		return creds
	}(cfg.CertPath, cfg.KeyPath)

	// Returns the client's transport credentials, either secure or insecure.
	tlsTool.ClientCreds = func(serverCertPool *x509.CertPool) credentials.TransportCredentials {
		if !core.G.TLSEnabled {
			return insecure.NewCredentials()
		}
		return credentials.NewClientTLSFromCert(serverCertPool, "")
	}(tlsTool.ServerCertPool)

	return &tlsTool
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (t *tlsTool) GetServerCertificate() *x509.CertPool {
	return t.ServerCertPool
}

func (t *tlsTool) GetServerCreds() credentials.TransportCredentials {
	return t.ServerCreds
}

func (t *tlsTool) GetClientCreds() credentials.TransportCredentials {
	return t.ClientCreds
}
