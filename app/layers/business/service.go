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

var _ core.Services = (*ServiceLayer)(nil)

// This is our concrete implementation of the core.Service interface.
// Holds a Storage layer to interact with the db and some tools.
type ServiceLayer struct {
	*pbs.UnimplementedUsersServiceServer

	Storage core.StorageLayer
	core.TokenGenerator
	core.PwdHasher
}

func NewService(storage core.StorageLayer, tokenGen core.TokenGenerator, pwdHasher core.PwdHasher) *ServiceLayer {
	return &ServiceLayer{
		Storage:        storage,
		TokenGenerator: tokenGen,
		PwdHasher:      pwdHasher,
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// getRouteFromCtx returns the gRPC method name from the context.
func getRouteFromCtx(ctx context.Context) string {
	if methodName, ok := grpc.Method(ctx); ok {
		return methodName
	}
	return ""
}

// errIsNotFound checks if the error is a gorm.ErrRecordNotFound.
func errIsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

var (
	ErrUnauthenticated = func() error { return errs.ErrSvcUnauthenticated() }
	ErrNotFound        = func(entity string) error { return errs.ErrSvcNotFound(entity) }
	ErrAlreadyExists   = func(entity string) error { return errs.ErrSvcAlreadyExists(entity) }
)
