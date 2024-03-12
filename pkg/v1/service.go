package v1

import (
	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/dependencies"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/* ----------------------------------- */
/*           - v1 Service -            */
/* ----------------------------------- */

// Service holds every gRPC method of the usersPB.UsersServiceServer.
// It handles the business logic of the API.
type Service interface {
	usersPB.UsersServiceServer
}

// service is our concrete implementation of the Service interface.
// It has an embedded Repository to interact with the database.
type service struct {
	DB             Repository
	TokenGenerator dependencies.TokenGenerator
	PwdHasher      dependencies.PwdHasher

	*usersPB.UnimplementedUsersServiceServer
}

// NewService returns a new instance of the service.
func NewService(db Repository, tokenGen dependencies.TokenGenerator, pwdHasher dependencies.PwdHasher) *service {
	return &service{
		DB:             db,
		TokenGenerator: tokenGen,
		PwdHasher:      pwdHasher,
	}
}

/* ----------------------------------- */
/*         - Service Errors -          */
/* ----------------------------------- */

var grpcUnknownErr = func(str string, err error) error {
	return status.Errorf(codes.Unknown, "%s: %v", str, err)
}

var grpcNotFoundErr = func(entity string) error {
	return status.Errorf(codes.NotFound, "%s not found.", entity)
}

var grpcAlreadyExistsErr = func(entity string) error {
	return status.Errorf(codes.AlreadyExists, "%s already exists.", entity)
}

var grpcUnauthenticatedErr = func(reason string) error {
	return status.Errorf(codes.Unauthenticated, reason)
}
