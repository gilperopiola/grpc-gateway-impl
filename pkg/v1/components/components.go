package components

import (
	"crypto/x509"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/common"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/http"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
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

	// Logger is used to log every gRPC and HTTP request that comes in.
	// LoggerOptions are the options we pass to the Logger.
	Logger        *common.Logger
	LoggerOptions []zap.Option

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
	Server       Server              // GRPCServer is, surprisingly, our gRPC Server.
	Interceptors []grpc.ServerOption // GRPCInterceptors run before or after gRPC calls.
	DialOptions  []grpc.DialOption   // GRPCDialOptions configure the communication between the HTTP and gRPC.
}

// HTTP holds the HTTP Gateway, middleware and Mux wrapper.
type HTTP struct {
	Gateway           Server                   // HTTPGateway is our HTTP Gateway (it's also a Server).
	Middleware        []runtime.ServeMuxOption // HTTPMiddleware are the ServeMuxOptions that run before or after HTTP calls.
	MiddlewareWrapper http.MuxWrapperFunc      // HTTPMiddlewareWrapper are the middleware that wrap around the HTTP Gateway.
}

// TLS holds the TLS components to guarantee secure communication.
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