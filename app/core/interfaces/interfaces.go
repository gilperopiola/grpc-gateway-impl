package interfaces

import (
	"context"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/external/storage/options"

	"google.golang.org/grpc"
)

// BusinessLayer holds every gRPC method in pbs.UsersServiceServer. It handles all business logic.
type BusinessLayer interface {
	pbs.UsersServiceServer
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

// InputValidator is the interface that wraps the ValidateInput method.
// It is used to validate incoming gRPC requests. The rules are defined in the .proto files.
type InputValidator interface {
	ValidateInput(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
}

// PwdHasher is the interface that wraps the Hash and Compare methods.
type PwdHasher interface {
	Hash(pwd string) string
	Compare(plainPwd, hashedPwd string) bool
}
