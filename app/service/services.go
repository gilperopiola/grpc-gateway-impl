package service

import (
	"context"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/logs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// ➤ If you add a new SubService, you'll need to include it in a few places in this file.
// TODO Automate.
//
// ➤ Service: The entire Business Layer of our application.
// ➤ SubService: A specific part of our Service, each one defined in a .proto file.

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*             - Service -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ v1 */

// Every SubService that we have defined on our protofiles is embedded here.
type Service struct {
	AuthSubService
	UsersSubService
	GroupsSubService
	GPTSubService
	HealthSubService
}

// Add new SubServices here (1/3).
func Setup(clients core.Clients, tools core.Tools) *Service {
	return &Service{
		AuthSubService:   AuthSubService{Clients: clients, Tools: tools},
		UsersSubService:  UsersSubService{Clients: clients, Tools: tools},
		GroupsSubService: GroupsSubService{Clients: clients, Tools: tools},
		GPTSubService:    GPTSubService{Clients: clients, Tools: tools},
		HealthSubService: HealthSubService{Clients: clients, Tools: tools},
	}
}

// Registers all of the GRPC services and endpoints on the GRPC Server.
func (s *Service) RegisterGRPCEndpoints(grpcServer grpc.ServiceRegistrar) {
	for _, serviceDesc := range grpcServicesDescriptors {
		grpcServer.RegisterService(serviceDesc, s)
	}
}

// Registers all of the HTTP services and endpoints on the HTTP Server.
func (s *Service) RegisterHTTPEndpoints(mux *runtime.ServeMux, opts ...grpc.DialOption) {
	for _, registerService := range httpRegisterServicesFns {
		ctx := context.Background()
		logs.LogFatalIfErr(registerService(ctx, mux, core.G.GRPCPort, opts))
	}
}

// And here (2/3).
var grpcServicesDescriptors = []*grpc.ServiceDesc{
	&pbs.AuthService_ServiceDesc,
	&pbs.UsersService_ServiceDesc,
	&pbs.GroupsService_ServiceDesc,
	&pbs.GPTService_ServiceDesc,
	&pbs.HealthService_ServiceDesc,
}

// And here (3/3).
var httpRegisterServicesFns = []func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error{
	pbs.RegisterAuthServiceHandlerFromEndpoint,
	pbs.RegisterUsersServiceHandlerFromEndpoint,
	pbs.RegisterGroupsServiceHandlerFromEndpoint,
	pbs.RegisterGPTServiceHandlerFromEndpoint,
	pbs.RegisterHealthServiceHandlerFromEndpoint,
}

/*
type SubServ[T SubService] struct {
	T
	RegisterGRPC *grpc.ServiceDesc
	RegisterHTTP  func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
}

type SubServ struct {
	GRPCServiceDesc *grpc.ServiceDesc
	HTTPRegisterFn  func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
}

type SubService interface {
	New(c core.Clients, t core.Tools) SubService
}

func NewSubService[T SubService](c core.Clients, t core.Tools) *T {
	var zeroValSubService T
	subService := zeroValSubService.New(c, t).(T)
	return &subService
}

var (
	AuthSubServiceCreator   = NewSubService[AuthSubService]
	UsersSubServiceCreator  = NewSubService[UsersSubService]
	GroupsSubServiceCreator = NewSubService[GroupsSubService]
	GPTSubServiceCreator    = NewSubService[GPTSubService]
	HealthSubServiceCreator = NewSubService[HealthSubService]
)

var subServices = []struct {
	NewFn           func(core.Clients, core.Tools) SubService
	GRPCServiceDesc *grpc.ServiceDesc
	HTTPRegisterFn  func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
}{
	{
		NewFn:           func(c core.Clients, t core.Tools) SubService { return AuthSubServiceCreator(c, t) },
		GRPCServiceDesc: &pbs.AuthService_ServiceDesc,
		HTTPRegisterFn:  pbs.RegisterAuthServiceHandlerFromEndpoint},
	{
		NewFn:           func(c core.Clients, t core.Tools) SubService { return UsersSubServiceCreator(c, t) },
		GRPCServiceDesc: &pbs.UsersService_ServiceDesc,
		HTTPRegisterFn:  pbs.RegisterUsersServiceHandlerFromEndpoint},
	{
		NewFn:           func(c core.Clients, t core.Tools) SubService { return GroupsSubServiceCreator(c, t) },
		GRPCServiceDesc: &pbs.GroupsService_ServiceDesc,
		HTTPRegisterFn:  pbs.RegisterGroupsServiceHandlerFromEndpoint},
	{
		NewFn:           func(c core.Clients, t core.Tools) SubService { return GPTSubServiceCreator(c, t) },
		GRPCServiceDesc: &pbs.GPTService_ServiceDesc,
		HTTPRegisterFn:  pbs.RegisterGPTServiceHandlerFromEndpoint},
	{
		NewFn:           func(c core.Clients, t core.Tools) SubService { return HealthSubServiceCreator(c, t) },
		GRPCServiceDesc: &pbs.HealthService_ServiceDesc,
		HTTPRegisterFn:  pbs.RegisterHealthServiceHandlerFromEndpoint},
}
*/
