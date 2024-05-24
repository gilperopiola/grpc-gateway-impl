package service

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
)

var _ core.Service = (*Service)(nil)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*             - Service -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ v1 */

// -> Here lies... Our Service. It's the core of our application logic.
type Service struct {
	Toolbox core.Toolbox

	pbs.UnimplementedAuthServiceServer
	pbs.UnimplementedUsersServiceServer
	pbs.UnimplementedGroupsServiceServer
}

func Setup(toolbox core.Toolbox) core.Service {
	return &Service{
		Toolbox: toolbox,
	}
}

func (s *Service) RegisterGRPCServices(grpcServer core.GRPCServiceRegistrar) {
	pbs.RegisterAuthServiceServer(grpcServer, s)
	pbs.RegisterUsersServiceServer(grpcServer, s)
	pbs.RegisterGroupsServiceServer(grpcServer, s)
}

func (s *Service) RegisterHTTPServices(mux *core.HTTPMultiplexer, opts core.GRPCDialOptions) {
	ctx := core.NewCtx()
	port := core.GRPCPort

	core.LogPanicIfErr(pbs.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, port, opts))
	core.LogPanicIfErr(pbs.RegisterUsersServiceHandlerFromEndpoint(ctx, mux, port, opts))
	core.LogPanicIfErr(pbs.RegisterGroupsServiceHandlerFromEndpoint(ctx, mux, port, opts))
}