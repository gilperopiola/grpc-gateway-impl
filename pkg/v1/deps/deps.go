package deps

import (
	"crypto/x509"

	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc/credentials"
)

// Server is an interface that abstracts the gRPC & HTTP Servers.
// Both servers have the same methods, but they are implemented differently.
type Server interface {
	Init()
	Run()
	Shutdown()
}

// NewDeps returns a new empty Deps.
func NewDeps() *Deps {
	return &Deps{
		TLSDeps: &TLSDeps{},
	}
}

// Deps holds all the dependencies that are used in the API.
type Deps struct {
	// Logger is used to log every gRPC and HTTP request that comes in.
	// LoggerOptions are the options we can pass to the Logger.
	Logger        *Logger
	LoggerOptions []zap.Option

	// Validator is used to validate the incoming gRPC requests (also HTTP as they are converted to gRPC).
	Validator *Validator

	// RateLimiter is used to limit the number of requests that get processed.
	RateLimiter *rate.Limiter

	// Authenticator is used to generate and validate JWT Tokens.
	Authenticator *jwtAuthenticator

	// PwdHasher is used to hash and compare passwords.
	PwdHasher *pwdHasher

	// TLSDeps holds the server certificate and server & client credentials.
	*TLSDeps
}

// TLSDeps holds the TLS configuration for the gRPC Server and the HTTP Gateway.
type TLSDeps struct {
	// ServerCert is a pool of certificates to use for the Server's TLS configuration.
	ServerCert *x509.CertPool

	// ServerCreds and ClientCreds are used to secure the connection between the HTTP Gateway and the gRPC Server.
	ServerCreds credentials.TransportCredentials
	ClientCreds credentials.TransportCredentials
}
