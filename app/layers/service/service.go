package service

import (
	"context"
	"errors"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - v1 Service -            */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var _ core.ServiceLayer = (*serviceLayer)(nil)

// Concrete implementation of the interface above.
type serviceLayer struct {
	*pbs.UnimplementedUsersServiceServer

	External core.ExternalLayer // -> Holds a reference to the ExternalLayer.
	Toolbox  core.Toolbox       // -> Holds a reference to the Toolbox.
}

func SetupLayer(external core.ExternalLayer, toolbox core.Toolbox) *serviceLayer {
	return &serviceLayer{
		External: external,
		Toolbox:  toolbox,
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// getGRPCMethodFromCtx returns the GRPC method name from the context.
func getGRPCMethodFromCtx(ctx context.Context) string {
	if methodName, ok := grpc.Method(ctx); ok {
		return methodName
	}
	return ""
}

// errIsNotFound checks if the error is a gorm.ErrRecordNotFound or a mongo.ErrNoDocuments.
// T0D0 Move to storage?
func errIsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, mongo.ErrNoDocuments)
}

var (
	ErrUnauthenticated = func() error { return errs.ErrSvcUnauthenticated() }
	ErrNotFound        = func(entity string) error { return errs.ErrSvcNotFound(entity) }
	ErrAlreadyExists   = func(entity string) error { return errs.ErrSvcAlreadyExists(entity) }
)
