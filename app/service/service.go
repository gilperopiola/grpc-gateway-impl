package service

import (
	"context"
	"errors"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/modules"
	"github.com/gilperopiola/grpc-gateway-impl/app/storage"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

/* ----------------------------------- */
/*           - v1 Service -            */
/* ----------------------------------- */

// Service holds every gRPC method in pbs.UsersServiceServer. It handles all business logic.
type Service interface {
	pbs.UsersServiceServer
}

// service is our concrete implementation of the Service interface.
type service struct {
	Repo           storage.Storage
	TokenGenerator modules.TokenGenerator
	PwdHasher      modules.PwdHasher

	*pbs.UnimplementedUsersServiceServer
}

// NewService returns a new instance of the service.
func NewService(repo storage.Storage, tokenGen modules.TokenGenerator, pwdHasher modules.PwdHasher) *service {
	return &service{
		Repo:           repo,
		TokenGenerator: tokenGen,
		PwdHasher:      pwdHasher,
	}
}

var (
	ErrUnauthenticated = func() error { return errs.ErrSvcUnauthenticated() }
	ErrNotFound        = func(entity string) error { return errs.ErrSvcNotFound(entity) }
	ErrAlreadyExists   = func(entity string) error { return errs.ErrSvcAlreadyExists(entity) }
)

// getGRPCMethodName returns the gRPC method name from the context.
func getGRPCMethodName(ctx context.Context) string {
	if methodName, ok := grpc.Method(ctx); ok {
		return methodName
	}
	return ""
}

// errIsNotFound checks if the error is a gorm.ErrRecordNotFound.
func errIsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
