package components

import (
	"crypto/x509"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/common"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Wrapper holds all the components that are used in the API.
type Wrapper struct {

	// GRPC and HTTP hold the gRPC and HTTP Servers, their dependencies and configuration.
	*GRPC
	*HTTP

	// TLS holds the server certificate and server & client credentials.
	*TLS

	// InputValidator is used to validate the incoming gRPC requests (also HTTP as they are converted to gRPC).
	InputValidator common.InputValidator

	// Authenticator is used to generate and validate JWT Tokens.
	Authenticator common.Authenticator

	// RateLimiter is used to limit the number of requests that get processed.
	RateLimiter *rate.Limiter

	// PwdHasher is used to hash and compare passwords.
	PwdHasher common.PwdHasher
}

// NewWrapper returns a new empty Wrapper to hold all components.
func NewWrapper() *Wrapper {
	return &Wrapper{
		GRPC: &GRPC{},
		HTTP: &HTTP{},
		TLS:  &TLS{},
	}
}

// GRPC holds the gRPC Server, interceptors and dial options.
type GRPC struct {
	Server        Server              // Server is, surprisingly, our gRPC Server.
	ServerOptions []grpc.ServerOption // ServerOptions configure the gRPC Server.
	DialOptions   []grpc.DialOption   // DialOptions configure the communication between HTTP and gRPC.
}

// HTTP holds the HTTP Gateway, middleware and Mux wrapper.
type HTTP struct {
	Gateway    Server                   // Gateway is our HTTP Gateway (it's also a Server).
	Middleware []runtime.ServeMuxOption // Middleware configure the HTTP Gateway.
	MuxWrapper http.MuxWrapperFunc      // MiddlewareWrapper are the middleware that wrap around the HTTP Gateway.
}

// TLS holds the Transport Layer Security components, such as certificates and credentials.
type TLS struct {
	// ServerCert is a pool of certificates.
	ServerCert *x509.CertPool

	// ServerCreds and ClientCreds are used to secure the connection between the HTTP Gateway and the gRPC Server.
	ServerCreds credentials.TransportCredentials
	ClientCreds credentials.TransportCredentials
}

// Server is an interface that abstracts the gRPC & HTTP Servers.
// Both servers have the same methods, but they are implemented differently.
type Server interface {
	Init()
	Run()
	Shutdown()
}
