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

/* ----------------------------------- */
/*   - App Modules (Passive/Active) -  */
/* ----------------------------------- */

// Passive Modules are mainly used to hold objects that are loaded at runtime and are not used to perform any actions.
// They just hold stuff.
type Passive struct {
	*GRPC
	*HTTP
	*TLS
}

// Active Modules let us perform actions.
type Active struct {
	InputValidator interfaces.InputValidator     // Validates GRPC and HTTP requests.
	Authenticator  interfaces.TokenAuthenticator // Generate & Validate JWT Tokens.
	RateLimiter    *rate.Limiter                 // Limit rate of requests.
	PwdHasher      interfaces.PwdHasher          // Hash and compare passwords.
}

type (
	GRPC struct {
		ServerOptions []grpc.ServerOption
		DialOptions   []grpc.DialOption
	}
	HTTP struct {
		MuxOptionsMiddleware []runtime.ServeMuxOption
		MuxWrapperMiddleware func(http.Handler) http.Handler
	}
	TLS struct {
		ServerCert  *x509.CertPool
		ServerCreds credentials.TransportCredentials
		ClientCreds credentials.TransportCredentials
	}
)
