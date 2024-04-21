package tools

import (
	"crypto/x509"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"

	"golang.org/x/time/rate"
	"google.golang.org/grpc/credentials"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*              - Tools -              */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// -> Tools do things and hold data.

var _ core.ToolsAccessor = (*Tools)(nil)

type (
	Tools struct {
		core.RequestsValidator  // Validates gRPC requests.
		core.TokenAuthenticator // Generates & Validates JWT Tokens.
		*rate.Limiter           // Limits rate of requests.
		core.PwdHasher          // Hashes and compares passwords.
		*TLS                    // Holds data for TLS communication.
	}

	TLS struct {
		ServerCert  *x509.CertPool
		ServerCreds credentials.TransportCredentials
		ClientCreds credentials.TransportCredentials
	}
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - Setup Tools -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (t *Tools) Setup(cfg *core.Config) {
	t.PwdHasher = NewPwdHasher(cfg.PwdHasherCfg.Salt)
	t.Limiter = rate.NewLimiter(rate.Limit(cfg.TokensPerSecond), cfg.MaxTokens)
	t.RequestsValidator = NewProtoValidator()
	t.TokenAuthenticator = NewJWTAuthenticator(cfg.JWTCfg.Secret, cfg.JWTCfg.SessionDays)
	t.SetupTLS(cfg.TLSCfg)
}

func (t *Tools) SetupTLS(cfg core.TLSCfg) {
	t.TLS = &TLS{}
	if cfg.Enabled {
		t.TLS.ServerCert = newTLSCertPool(cfg.CertPath)
		t.TLS.ServerCreds = newServerTransportCreds(cfg.CertPath, cfg.KeyPath)
	}
	t.TLS.ClientCreds = newClientTransportCreds(cfg.Enabled, t.ServerCert)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*        - Get Tools Methods -        */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Avoid depending on this package by using core.ToolsAccessor instead of this concrete struct.
func (t *Tools) GetRequestsValidator() core.RequestsValidator        { return t.RequestsValidator }
func (t *Tools) GetAuthenticator() core.TokenAuthenticator           { return t.TokenAuthenticator }
func (t *Tools) GetRateLimiter() *rate.Limiter                       { return t.Limiter }
func (t *Tools) GetPwdHasher() core.PwdHasher                        { return t.PwdHasher }
func (t *Tools) GetTLSServerCert() *x509.CertPool                    { return t.TLS.ServerCert }
func (t *Tools) GetTLSServerCreds() credentials.TransportCredentials { return t.TLS.ServerCreds }
func (t *Tools) GetTLSClientCreds() credentials.TransportCredentials { return t.TLS.ClientCreds }
