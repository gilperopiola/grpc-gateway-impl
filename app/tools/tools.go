package tools

import (
	"crypto/x509"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"

	"golang.org/x/time/rate"
	"google.golang.org/grpc/credentials"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - App Tools -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type (
	// All tools used by the app. Tools do things and hold data.
	Tools struct {
		InputValidator core.InputValidator     // Validates GRPC and HTTP requests.
		Authenticator  core.TokenAuthenticator // Generate & Validate JWT Tokens.
		RateLimiter    *rate.Limiter           // Limit rate of requests.
		PwdHasher      core.PwdHasher          // Hash and compare passwords.
		*TLS                                   // Holds data for TLS communication.
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

func (tools *Tools) Setup(cfg *core.Config) {
	tools.SetupPwdHasher(cfg.PwdHasherCfg)
	tools.SetupRateLimiter(cfg.RateLimiterCfg)
	tools.SetupInputValidator()
	tools.SetupAuthenticator(cfg.JWTCfg)
	tools.SetupTLS(cfg.TLSCfg)
}

func (tools *Tools) SetupInputValidator() {
	tools.InputValidator = NewInputValidator()
}

func (tools *Tools) SetupAuthenticator(jwtCfg core.JWTCfg) {
	tools.Authenticator = NewJWTAuthenticator(jwtCfg.Secret, jwtCfg.SessionDays)
}

func (tools *Tools) SetupRateLimiter(rlCfg core.RateLimiterCfg) {
	tools.RateLimiter = rate.NewLimiter(rate.Limit(rlCfg.TokensPerSecond), rlCfg.MaxTokens)
}

func (tools *Tools) SetupPwdHasher(pwdHasherCfg core.PwdHasherCfg) {
	tools.PwdHasher = NewPwdHasher(pwdHasherCfg.Salt)
}

func (tools *Tools) SetupTLS(tlsCfg core.TLSCfg) {
	tools.TLS = &TLS{}
	if tlsCfg.Enabled {
		tools.TLS.ServerCert = NewTLSCertPool(tlsCfg.CertPath)                           // -> Server Certificate.
		tools.TLS.ServerCreds = NewServerTransportCreds(tlsCfg.CertPath, tlsCfg.KeyPath) // -> Server Credentials.
	}
	tools.TLS.ClientCreds = NewClientTransportCreds(tlsCfg.Enabled, tools.ServerCert) // -> Client Credentials.
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*     - Tools Concrete Methods -      */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Avoid depending on this package by using core.ToolsAccessor instead of this concrete struct.

func (t *Tools) GetInputValidator() core.InputValidator              { return t.InputValidator }
func (t *Tools) GetAuthenticator() core.TokenAuthenticator           { return t.Authenticator }
func (t *Tools) GetRateLimiter() *rate.Limiter                       { return t.RateLimiter }
func (t *Tools) GetPwdHasher() core.PwdHasher                        { return t.PwdHasher }
func (t *Tools) GetTLSServerCert() *x509.CertPool                    { return t.TLS.ServerCert }
func (t *Tools) GetTLSServerCreds() credentials.TransportCredentials { return t.TLS.ServerCreds }
func (t *Tools) GetTLSClientCreds() credentials.TransportCredentials { return t.TLS.ClientCreds }
