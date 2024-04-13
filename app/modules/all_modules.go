package modules

import (
	"crypto/x509"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type All struct {
	*GRPC
	*HTTP
	*TLS

	InputValidator InputValidator     // Validates GRPC and HTTP requests.
	Authenticator  TokenAuthenticator // Generate & Validate JWT Tokens.
	RateLimiter    *rate.Limiter      // Limit rate of requests.
	PwdHasher      PwdHasher          // Hash and compare passwords.
}

type GRPC struct {
	Server        Server
	ServerOptions []grpc.ServerOption
	DialOptions   []grpc.DialOption
}

type HTTP struct {
	Gateway              Server
	MuxOptionsMiddleware []runtime.ServeMuxOption
	MuxWrapperMiddleware func(http.Handler) http.Handler
}

type TLS struct {
	ServerCert  *x509.CertPool
	ServerCreds credentials.TransportCredentials
	ClientCreds credentials.TransportCredentials
}

// Server abstracts the gRPC Server & HTTP Gateway.
// We use it to avoid import cycles between this pkg and modules/grpc or modules/http.
type Server interface {
	Init()
	Run()
	Shutdown()
}
