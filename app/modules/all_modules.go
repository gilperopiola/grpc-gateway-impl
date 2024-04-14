package modules

import (
	"crypto/x509"
	"net/http"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/interfaces"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type All struct {
	*GRPC
	*HTTP
	*TLS

	InputValidator interfaces.InputValidator     // Validates GRPC and HTTP requests.
	Authenticator  interfaces.TokenAuthenticator // Generate & Validate JWT Tokens.
	RateLimiter    *rate.Limiter                 // Limit rate of requests.
	PwdHasher      interfaces.PwdHasher          // Hash and compare passwords.
}

type GRPC struct {
	ServerOptions []grpc.ServerOption
	DialOptions   []grpc.DialOption
}

type HTTP struct {
	MuxOptionsMiddleware []runtime.ServeMuxOption
	MuxWrapperMiddleware func(http.Handler) http.Handler
}

type TLS struct {
	ServerCert  *x509.CertPool
	ServerCreds credentials.TransportCredentials
	ClientCreds credentials.TransportCredentials
}
