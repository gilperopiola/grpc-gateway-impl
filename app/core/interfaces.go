package core

import (
	"context"
	"crypto/x509"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/external/storage/options"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ToolsAccessor interface {
	GetInputValidator() InputValidator
	GetAuthenticator() TokenAuthenticator
	GetRateLimiter() *rate.Limiter
	GetPwdHasher() PwdHasher
	GetTLSServerCert() *x509.CertPool
	GetTLSServerCreds() credentials.TransportCredentials
	GetTLSClientCreds() credentials.TransportCredentials
}

// Storage is the interface that wraps the basic methods to interact with the Database.
// App -> Service -> Storage -> DB.
type Storage interface {
	CreateUser(username, hashedPwd string) (*models.User, error)
	GetUser(opts ...options.QueryOpt) (*models.User, error)
	GetUsers(page, pageSize int, opts ...options.QueryOpt) (models.Users, int, error)
}

type TokenAuthenticator interface {
	TokenGenerator
	TokenValidator
}

type TokenGenerator interface {
	Generate(id int, username string, role models.Role) (string, error)
}

type TokenValidator interface {
	Validate(ctx context.Context, req interface{}, svInfo *grpc.UnaryServerInfo, h grpc.UnaryHandler) (any, error)
}

// Used to validate incoming gRPC requests. Rules are defined on the protofiles.
type InputValidator interface {
	ValidateInput(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
}

// PwdHasher is the interface that wraps the Hash and Compare methods.
type PwdHasher interface {
	Hash(pwd string) string
	Compare(plainPwd, hashedPwd string) bool
}
