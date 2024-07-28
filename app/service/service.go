package service

import (
	"context"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*             - Service -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ v1 */

// This holds every service that we have defined on our protofiles.
type Services struct {
	AuthService
	UsersService
	GroupsService
}

func Setup(tools core.Tools) *Services {
	return &Services{
		AuthService:   AuthService{Tools: tools},
		UsersService:  UsersService{Tools: tools},
		GroupsService: GroupsService{Tools: tools},
	}
}

type AuthService struct {
	pbs.UnimplementedAuthServiceServer
	Tools core.Tools
}

type UsersService struct {
	pbs.UnimplementedUsersServiceServer
	Tools core.Tools
}

type GroupsService struct {
	pbs.UnimplementedGroupsServiceServer
	Tools core.Tools
}

// Registers all of the GRPC services and their endpoints on the GRPC Server.
func (s *Services) RegisterGRPCEndpoints(grpcServer grpc.ServiceRegistrar) {

	// New services should be added here.
	grpcServices := []grpc.ServiceDesc{
		pbs.AuthService_ServiceDesc,
		pbs.UsersService_ServiceDesc,
		pbs.GroupsService_ServiceDesc,
	}

	for _, grpcService := range grpcServices {
		grpcServer.RegisterService(&grpcService, s)
	}
}

// Registers all of the HTTP services and their endpoints on the HTTP Server.
func (s *Services) RegisterHTTPEndpoints(mux *runtime.ServeMux, opts ...grpc.DialOption) {

	// New services should be added here.
	httpServices := []func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error){
		pbs.RegisterAuthServiceHandlerFromEndpoint,
		pbs.RegisterUsersServiceHandlerFromEndpoint,
		pbs.RegisterGroupsServiceHandlerFromEndpoint,
	}

	ctx := context.Background()
	for _, httpService := range httpServices {
		core.LogFatalIfErr(httpService(ctx, mux, core.GRPCPort, opts))
	}
}
