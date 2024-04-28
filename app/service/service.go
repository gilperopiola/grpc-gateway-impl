package service

import (
	"errors"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var _ core.Service = (*service)(nil)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - Service (v1) -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Main App Service.
// Holds all the methods that the GRPC server will use.
type service struct {
	core.Toolbox // -> Our Service has a handy Toolbox.

	*pbs.UnimplementedUsersServiceServer
}

func Setup(toolbox core.Toolbox) *service {
	return &service{
		Toolbox: toolbox,
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// errIsNotFound checks if the error is a gorm.ErrRecordNotFound or a mongo.ErrNoDocuments.
// T0D0 Move to storage?
func errIsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, mongo.ErrNoDocuments)
}
