package core

import (
	"context"
	"crypto/x509"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/external/storage/options"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - Interfaces -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Services holds all of our particular services, here we just have 1.
type Services interface {
	pbs.UsersServiceServer
}

// StorageLayer is the interface that wraps the basic methods to interact with the Database.
// App -> Service -> StorageLayer -> DB.
type StorageLayer interface {
	CreateUser(username, hashedPwd string) (*models.User, error)
	GetUser(opts ...options.QueryOpt) (*models.User, error)
	GetUsers(page, pageSize int, opts ...options.QueryOpt) (models.Users, int, error)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type ToolsAccessor interface {
	GetRequestsValidator() RequestsValidator
	GetAuthenticator() TokenAuthenticator
	GetRateLimiter() *rate.Limiter
	GetPwdHasher() PwdHasher
	GetTLSServerCert() *x509.CertPool
	GetTLSServerCreds() credentials.TransportCredentials
	GetTLSClientCreds() credentials.TransportCredentials
}

type RequestsValidator interface {
	ValidateRequest(ctx context.Context, req any, svInfo *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)
}

type TokenAuthenticator interface {
	TokenGenerator
	TokenValidator
}

type TokenGenerator interface {
	Generate(id int, username string, role models.Role) (string, error)
}

type TokenValidator interface {
	Validate(ctx context.Context, req any, svInfo *grpc.UnaryServerInfo, h grpc.UnaryHandler) (any, error)
}

type PwdHasher interface {
	Hash(pwd string) string
	Compare(plainPwd, hashedPwd string) bool
}
