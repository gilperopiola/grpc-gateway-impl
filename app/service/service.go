package service

import (
	"context"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// IMPORTANT: If you add a new SubService, you'll need to include it in a few places in this file.

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*             - Service -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ v1 */

// Every service that we have defined on our protofiles is embedded here.
type Service struct {
	AuthSubService
	UsersSubService
	GroupsSubService
	HealthSubService
}

// I don't really like this way of calling each particular service a SubService,
// but I found no better way to differentiate our Service as the business layer
// and each .proto defined Service.

func Setup(tools core.Tools) *Service {
	// New services should be added here.
	return &Service{
		AuthSubService:   AuthSubService{Tools: tools},
		UsersSubService:  UsersSubService{Tools: tools},
		GroupsSubService: GroupsSubService{Tools: tools},
		HealthSubService: HealthSubService{Tools: tools},
	}
}

// Registers all of the GRPC services and their endpoints on the GRPC Server.
func (s *Service) RegisterGRPCEndpoints(grpcServer grpc.ServiceRegistrar) {

	// New services should be added here.
	servicesDescs := []grpc.ServiceDesc{
		pbs.AuthService_ServiceDesc,
		pbs.UsersService_ServiceDesc,
		pbs.GroupsService_ServiceDesc,
		pbs.HealthService_ServiceDesc,
	}

	for _, serviceDesc := range servicesDescs {
		grpcServer.RegisterService(&serviceDesc, s)
	}
}

// Registers all of the HTTP services and their endpoints on the HTTP Server.
func (s *Service) RegisterHTTPEndpoints(mux *runtime.ServeMux, opts ...grpc.DialOption) {

	// New services should be added here.
	httpServices := []func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error{
		pbs.RegisterAuthServiceHandlerFromEndpoint,
		pbs.RegisterUsersServiceHandlerFromEndpoint,
		pbs.RegisterGroupsServiceHandlerFromEndpoint,
		pbs.RegisterHealthServiceHandlerFromEndpoint,
	}

	for _, httpService := range httpServices {
		core.LogFatalIfErr(httpService(context.Background(), mux, core.GRPCPort, opts))
	}
}
