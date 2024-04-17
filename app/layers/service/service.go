package service

import (
	"context"
	"errors"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/interfaces"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - v1 Service -            */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// service is our concrete implementation of the BusinessLayer interface.
type service struct {
	Storage        interfaces.Storage
	TokenGenerator interfaces.TokenGenerator
	PwdHasher      interfaces.PwdHasher

	*pbs.UnimplementedUsersServiceServer
}

// NewService returns a new instance of the service.
func NewService(storage interfaces.Storage, tokenGen interfaces.TokenGenerator, pwdHasher interfaces.PwdHasher) *service {
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
