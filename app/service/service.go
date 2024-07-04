package service

import (
	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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

func Setup(toolbox core.Toolbox) *Service {
	return &Service{
		Toolbox: toolbox,
	}
}

func (s *Service) RegisterGRPCServices(grpcServer god.GRPCSvcRegistrar) {
	pbs.RegisterAuthServiceServer(grpcServer, s)
	pbs.RegisterUsersServiceServer(grpcServer, s)
	pbs.RegisterGroupsServiceServer(grpcServer, s)
}

func (s *Service) RegisterHTTPServices(mux *runtime.ServeMux, opts god.GRPCDialOpts) {
	ctx := god.NewCtx()
	port := core.GRPCPort

	core.LogFatalIfErr(pbs.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, port, opts))
	core.LogFatalIfErr(pbs.RegisterUsersServiceHandlerFromEndpoint(ctx, mux, port, opts))
	core.LogFatalIfErr(pbs.RegisterGroupsServiceHandlerFromEndpoint(ctx, mux, port, opts))
}
