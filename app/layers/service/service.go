package service

import (
	"context"
	"errors"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/external"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - v1 Service -            */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var _ core.ServiceAPI = (*ServiceLayer)(nil)

// Concrete implementation of the interface above.
type ServiceLayer struct {
	*pbs.UnimplementedUsersServiceServer

	*external.ExternalLayer // -> Holds a reference to the ExternalLayer. This probably should be handled differently T0D0.

	core.TokenGenerator // Tool
	core.PwdHasher      // Tool
}

func SetupLayer(external *external.ExternalLayer, tokenGen core.TokenGenerator, pwdHasher core.PwdHasher) *ServiceLayer {
	return &ServiceLayer{
		ExternalLayer:  external,
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
