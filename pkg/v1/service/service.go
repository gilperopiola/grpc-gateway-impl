package service

import (
	"context"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/common"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository"

	"google.golang.org/grpc"
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

// getGRPCMethodFrom returns the gRPC method name from the context.
func getGRPCMethodFrom(ctx context.Context) string {
	if methodName, ok := grpc.Method(ctx); ok {
		return methodName
	}
	return ""
}
