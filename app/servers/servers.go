package servers

import (
	"net"
	"net/http"
	"time"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Servers struct {
	GRPC *grpc.Server
	HTTP *http.Server

	checkHealth func() error
}

func Setup(service core.Service, toolbox core.Toolbox) *Servers {
	core.Infof(" 🚀 Starting Servers")

	var (
		grpcServerOpts   = getGRPCServerOpts(toolbox, core.TLSEnabled)
		grpcDialOpts     = getGRPCDialOpts(toolbox.GetClientCreds())
		httpMiddleware   = getHTTPMiddlewareChain()
		httpServeMuxOpts = getHTTPMuxOpts()
	)

	return &Servers{
		GRPC: setupGRPCServer(service, grpcServerOpts),
		HTTP: setupHTTPGateway(service, httpServeMuxOpts, httpMiddleware, grpcDialOpts),

		checkHealth: toolbox.CheckHealth,
	}
}

func setupGRPCServer(service core.Service, serverOpts god.GRPCServerOpts) *grpc.Server {
	grpcServer := grpc.NewServer(serverOpts...)
	service.RegisterGRPCServices(grpcServer)
	return grpcServer
}

func setupHTTPGateway(svc core.Service, serveMuxOpts []runtime.ServeMuxOption, mw middlewareFunc, grpcDialOpts god.GRPCDialOpts) *http.Server {
	mux := runtime.NewServeMux(serveMuxOpts...)
	svc.RegisterHTTPServices(mux, grpcDialOpts)
	return &http.Server{
		Addr:    core.HTTPPort,
		Handler: mw(mux),
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (s *Servers) Run() {
	runGRPC(s.GRPC)
	runHTTP(s.HTTP)

	go func() {
		// Check service health after 1 second
		time.Sleep(time.Second)
		err := s.checkHealth()
		core.LogFatalIfErr(err)
		core.Infoln(" 🚀 All OK!")
	}()
}

func runGRPC(grpcServer *grpc.Server) {
	core.Infof(" 🚀 GRPC OK on Port %s", core.GRPCPort)

	lis, err := net.Listen("tcp", core.GRPCPort)
	core.LogFatalIfErr(err)

	logServiceInfo(grpcServer)

	go core.LogFatalIfErr(grpcServer.Serve(lis))
}

func runHTTP(httpGateway *http.Server) {
	core.Infof(" 🚀 HTTP OK on Port %s", core.HTTPPort)
	go core.LogFatalIfErr(httpGateway.ListenAndServe())
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func logServiceInfo(grpcServer *grpc.Server) {
	for _, svcInfo := range grpcServer.GetServiceInfo() {
		core.Infof(" 🐸 Service Loaded: %s", svcInfo.Metadata)
		for _, svcMethod := range svcInfo.Methods {
			core.Infof(" \t - Endpoint Loaded: %s", svcMethod.Name)
		}
	}
}

func (s *Servers) Shutdown() {
	core.Infof(" 🛑 - Shutting down GRPC")
	s.GRPC.GracefulStop()

	core.Infof(" 🛑 - Shutting down HTTP")
	ctx, cancelCtx := god.NewCtxWithTimeout(5 * time.Second)
	defer cancelCtx()
	core.LogFatalIfErr(s.HTTP.Shutdown(ctx))
}
