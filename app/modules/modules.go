package modules

import (
	"crypto/x509"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// All holds all the modules that are used in the App.
type All struct {
	GRPC // GRPC Module.
	HTTP // HTTP Module.

	InputValidator InputValidator     // Validates GRPC and HTTP requests.
	Authenticator  TokenAuthenticator // Generate & Validate JWT Tokens.
	RateLimiter    *rate.Limiter      // Limit rate of requests.
	PwdHasher      PwdHasher          // Hash and compare passwords.
	TLS                               // TLS Module.
}

// Server abstracts the gRPC Server & HTTP Gateway.
// We use it to avoid import cycles between this pkg and modules/grpc or modules/http.
type Server interface {
	Init()
	Run()
	Shutdown()
}

// GRPC Module.
type GRPC struct {
	Server        Server
	ServerOptions []grpc.ServerOption
	DialOptions   []grpc.DialOption
}

// HTTP Module.
type HTTP struct {
	Gateway              Server
	MuxOptionsMiddleware []runtime.ServeMuxOption
	MuxWrapperMiddleware func(http.Handler) http.Handler
}

// TLS Module.
type TLS struct {
	ServerCert  *x509.CertPool
	ServerCreds credentials.TransportCredentials
	ClientCreds credentials.TransportCredentials
}
