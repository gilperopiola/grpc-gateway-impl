package servers

import (
	"net"
	"net/http"
	"time"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/service"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// I just asked GPT and he said technically this project runs 1 server with 2 protocols.
// Still naming this 'Servers'.
type Servers struct {
	GRPC *grpc.Server
	HTTP *http.Server

	checkHealth func() error
}

func Setup(services *service.Services, tools core.Tools) *Servers {
	core.ServerLogf(" 🚀 Starting Servers")

	var (
		grpcServerOpts   = getGRPCServerOpts(tools, core.TLSEnabled)
		grpcDialOpts     = getGRPCDialOpts(tools.GetClientCreds())
		httpMiddleware   = getHTTPMiddlewareChain()
		httpServeMuxOpts = getHTTPMuxOpts()
	)

	return &Servers{
		GRPC:        setupGRPC(services, grpcServerOpts),
		HTTP:        setupHTTP(services, httpServeMuxOpts, httpMiddleware, grpcDialOpts...),
		checkHealth: tools.CheckHealth,
	}
}

func setupGRPC(services *service.Services, serverOpts god.GRPCServerOpts) *grpc.Server {
	grpcServer := grpc.NewServer(serverOpts...)
	services.RegisterGRPCEndpoints(grpcServer)
	return grpcServer
}

func setupHTTP(services *service.Services, serveMuxOpts []runtime.ServeMuxOption, mw middlewareFunc, grpcDialOpts ...grpc.DialOption) *http.Server {
	mux := runtime.NewServeMux(serveMuxOpts...)
	services.RegisterHTTPEndpoints(mux, grpcDialOpts...)
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
		time.Sleep(time.Second)
		core.LogFatalIfErr(core.Retry(s.checkHealth, 5))
		core.ServerLog("Services Healthy! 🌈\n")
	}()
}

func runGRPC(grpcServer *grpc.Server) {
	core.ServerLogf(" 🚀 GRPC OK | Port %s", core.GRPCPort)

	listenOnTCP := func() (any, error) { return net.Listen("tcp", core.GRPCPort) }
	result, err := core.FallbackAndRetry(listenOnTCP, func() {}, 5)
	core.LogFatalIfErr(err)

	lis := result.(net.Listener)
	go func(lis net.Listener) {
		logServicesAndEndpoints(grpcServer)
		core.LogFatalIfErr(grpcServer.Serve(lis))
	}(lis)
}

func runHTTP(httpGateway *http.Server) {
	core.ServerLogf(" 🚀 HTTP OK | Port %s", core.HTTPPort)
	go core.LogFatalIfErr(core.Retry(httpGateway.ListenAndServe, 5))
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func logServicesAndEndpoints(grpcServer *grpc.Server) {
	for _, svcInfo := range grpcServer.GetServiceInfo() {
		core.ServerLogf(" 🐸 Service Loaded: %s", svcInfo.Metadata)
		for _, svcMethod := range svcInfo.Methods {
			core.ServerLogf(" \t - Endpoint Loaded: %s", svcMethod.Name)
		}
	}
}

func (s *Servers) Shutdown() {
	core.LogImportant("Shutting down GRPC")
	s.GRPC.GracefulStop()

	core.LogImportant("Shutting down HTTP")
	ctx, cancelCtx := god.NewCtxWithTimeout(5 * time.Second)
	defer cancelCtx()
	core.LogFatalIfErr(s.HTTP.Shutdown(ctx))
}
