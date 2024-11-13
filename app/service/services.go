package service

import (
	"context"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/logs"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// ➤ Service: The entire Business Layer of our application.
// ➤ Svc: A specific domain of our Service, each defined in its own .proto file.

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*             - Service -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ v1 */

// This holds every Svc and all the info needed to register them on our Servers.
type Service struct {

	// ➤ Each Svc needs to be registered in GRPC and HTTP. When we create a new Svc, we must append its
	// GRPC Descriptor and HTTP Register Function to the corresponding slice inside this and voilá.
	// That is done in [Setup].
	RegistrationInfo

	// ➤ New Svcs also go here.
	AuthSvc
	UserSvc
	GroupSvc
	GPTSvc
	HealthSvc
	// ...
}

func Setup(clients core.Clients, tools core.Tools) *Service {
	service := Service{
		HealthSvc: HealthSvc{Clients: clients, Tools: tools},
		AuthSvc:   AuthSvc{Clients: clients, Tools: tools},
		UserSvc:   UserSvc{Clients: clients, Tools: tools},
		GroupSvc:  GroupSvc{Clients: clients, Tools: tools},
		GPTSvc:    GPTSvc{Clients: clients, Tools: tools},
		// ...
		RegistrationInfo: RegistrationInfo{
			GRPCServiceDescs: []*grpc.ServiceDesc{
				&pbs.AuthService_ServiceDesc,
				&pbs.UsersSvc_ServiceDesc,
				&pbs.GroupsService_ServiceDesc,
				&pbs.GPTService_ServiceDesc,
				&pbs.HealthService_ServiceDesc,
				// ...
			},
			HTTPRegisterFns: []func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error{
				pbs.RegisterAuthServiceHandlerFromEndpoint,
				pbs.RegisterUsersSvcHandlerFromEndpoint,
				pbs.RegisterGroupsServiceHandlerFromEndpoint,
				pbs.RegisterGPTServiceHandlerFromEndpoint,
				pbs.RegisterHealthServiceHandlerFromEndpoint,
				// ...
			},
		},
	}
	logs.InitModuleOK("Service", "⚡")
	return &service
}

// Has all the information needed to register the Svcs both on the GRPC and HTTP servers.
// Add any new Svcs to these slices on the [Setup].
type RegistrationInfo struct {
	GRPCServiceDescs []*grpc.ServiceDesc
	HTTPRegisterFns  []func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
}

// Registers all Svcs and endpoints on the GRPC Server.
func (s *Service) RegisterInGRPC(grpcServer grpc.ServiceRegistrar) {
	for _, serviceDesc := range s.GRPCServiceDescs {
		grpcServer.RegisterService(serviceDesc, s)
	}
}

// Registers all Svcs and endpoints on the HTTP Server.
func (s *Service) RegisterInHTTP(mux *runtime.ServeMux, opts ...grpc.DialOption) {
	for _, registerServiceFn := range s.HTTPRegisterFns {
		ctx := context.Background()
		logs.LogFatalIfErr(registerServiceFn(ctx, mux, core.G.GRPCPort, opts))
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

/* -> Scraped Ideas:

type SvcBase[T pbs.UnimplementedAuthServiceServer | pbs.UnimplementedUsersSvcServer | pbs.UnimplementedGroupsServiceServer | pbs.UnimplementedGPTServiceServer | pbs.UnimplementedHealthServiceServer] struct {
	T       T
	Clients core.Clients
	Tools   core.Tools
}

func newBaseSvc[T pbs.UnimplementedAuthServiceServer | pbs.UnimplementedUsersSvcServer | pbs.UnimplementedGroupsServiceServer | pbs.UnimplementedGPTServiceServer | pbs.UnimplementedHealthServiceServer](t T, c core.Clients, tools core.Tools) SvcBase[T] {
	return SvcBase[T]{T: t, Clients: c, Tools: tools}
}

// --------------------------------------------------

type EmbeddedSubService struct {
	Name       string
	Code       string
	InstanceID string
}

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
		GRPCServiceDesc: &pbs.UsersSvc_ServiceDesc,
		HTTPRegisterFn:  pbs.RegisterUsersSvcHandlerFromEndpoint},
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
