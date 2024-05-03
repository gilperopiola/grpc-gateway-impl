package service

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
)

var _ core.Service = (*service)(nil)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*             - Service -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ v1 */

// -> Here lies... Our Service. It's the core of our application logic.
// -> It holds all the methods that the GRPC and HTTP Servers will call.
type service struct {
	pbs.UnimplementedUsersServiceServer
	pbs.UnimplementedGroupsServiceServer

	// -> DB and other stuff are here.
	core.Actions
}

func Setup(actions core.Actions) core.Service {
	return &service{
		Actions: actions,
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var errNotFound = errs.GRPCNotFound
var errAlreadyExists = errs.GRPCAlreadyExists
var errUnauthenticated = errs.GRPCUnauthenticated
