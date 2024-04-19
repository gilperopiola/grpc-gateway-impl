package business

import (
	"context"
	"errors"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - v1 Service -            */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Service holds all of our particular services, here we just have 1.
// All business logic should be implemented here.
type Service interface {
	pbs.UsersServiceServer
}

// service is our concrete implementation of the BusinessLayer interface.
type service struct {
	Storage        core.Storage
	TokenGenerator core.TokenGenerator
	PwdHasher      core.PwdHasher

	*pbs.UnimplementedUsersServiceServer
}

// NewService returns a new instance of the business.
func NewService(storage core.Storage, tokenGen core.TokenGenerator, pwdHasher core.PwdHasher) *service {
	return &service{
		Storage:        storage,
		TokenGenerator: tokenGen,
		PwdHasher:      pwdHasher,
	}
}

var (
	ErrUnauthenticated = func() error { return errs.ErrSvcUnauthenticated() }
	ErrNotFound        = func(entity string) error { return errs.ErrSvcNotFound(entity) }
	ErrAlreadyExists   = func(entity string) error { return errs.ErrSvcAlreadyExists(entity) }
)

// getRoute returns the gRPC method name from the context.
func getRoute(ctx context.Context) string {
	if methodName, ok := grpc.Method(ctx); ok {
		return methodName
	}
	return ""
}

// errIsNotFound checks if the error is a gorm.ErrRecordNotFound.
func errIsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
