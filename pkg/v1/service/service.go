package service

import (
	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/common"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/* ----------------------------------- */
/*           - v1 Service -            */
/* ----------------------------------- */

// Service holds every gRPC method in usersPB.UsersServiceServer. It handles all business logic.
type Service interface {
	usersPB.UsersServiceServer
}

// service is our concrete implementation of the Service interface.
type service struct {
	Repo           repository.Repository
	TokenGenerator common.TokenGenerator
	PwdHasher      common.PwdHasher

	*usersPB.UnimplementedUsersServiceServer
}

// NewService returns a new instance of the service.
func NewService(repo repository.Repository, tokenGen common.TokenGenerator, pwdHasher common.PwdHasher) *service {
	return &service{
		Repo:           repo,
		TokenGenerator: tokenGen,
		PwdHasher:      pwdHasher,
	}
}

/* ----------------------------------- */
/*         - Service Errors -          */
/* ----------------------------------- */

var (
	grpcUnknownErr         = func(msg string, err error) error { return status.Errorf(codes.Unknown, "%s: %v", msg, err) }
	grpcNotFoundErr        = func(entity string) error { return status.Errorf(codes.NotFound, "%s not found.", entity) }
	grpcAlreadyExistsErr   = func(entity string) error { return status.Errorf(codes.AlreadyExists, "%s already exists.", entity) }
	grpcUnauthenticatedErr = func(reason string) error { return status.Errorf(codes.Unauthenticated, reason) }
)
