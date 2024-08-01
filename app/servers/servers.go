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
}

/* -~-~-~-~-~ Setup -~-~-~-~-~- */

func Setup(services *service.Service, tools core.Tools) *Servers {
	// GRPC has ServerOptions (Interceptors are there) and DialOptions (that
	// are actually used by HTTP to connect to GRPC).
	//
	// HTTP has Middleware (it's a chain of handlers that call each other) and
	// also ServeMuxOptions which configure the Multiplexer (mux for short)
	// with our custom error handler and stuff.
	var (
		grpcServerOpts   = getGRPCServerOpts(tools, core.TLSEnabled)
		grpcDialOpts     = getGRPCDialOpts(tools.GetClientCreds())
		httpMiddleware   = getHTTPMiddlewareChain()
		httpServeMuxOpts = getHTTPMuxOpts()
	)

	return &Servers{
		GRPC: setupGRPC(services, grpcServerOpts),
		HTTP: setupHTTP(services, httpServeMuxOpts, httpMiddleware, grpcDialOpts...),
	}
}

func setupGRPC(services *service.Service, serverOpts god.GRPCServerOpts) *grpc.Server {
	grpcServer := grpc.NewServer(serverOpts...)
	services.RegisterGRPCEndpoints(grpcServer)

	for _, svcInfo := range grpcServer.GetServiceInfo() {
		core.ServerLogf(" 🐸 Service Loaded: %s", svcInfo.Metadata)
		for _, svcMethod := range svcInfo.Methods {
			core.ServerLogf(" \t - Endpoint Loaded: %s", svcMethod.Name)
		}
	}

	return grpcServer
}

func setupHTTP(services *service.Service, serveMuxOpts []runtime.ServeMuxOption, mw middlewareFunc, grpcDialOpts ...grpc.DialOption) *http.Server {
	mux := runtime.NewServeMux(serveMuxOpts...)
	services.RegisterHTTPEndpoints(mux, grpcDialOpts...)

	return &http.Server{
		Addr:    core.HTTPPort,
		Handler: mw(mux),
	}
}

/* -~-~-~-~-~ Run -~-~-~-~-~- */

func (s *Servers) Run() {
	go s.runGRPC()
	go s.runHTTP()
}

func (s *Servers) runGRPC() {
	listenOnTCP := func() (any, error) { return net.Listen("tcp", core.GRPCPort) }
	result, err := core.FallbackAndRetry(listenOnTCP, func() {}, 5)
	core.LogFatalIfErr(err)

	lis := result.(net.Listener)
	core.LogFatalIfErr(s.GRPC.Serve(lis))
}

func (s *Servers) runHTTP() {
	if err := core.Retry(s.HTTP.ListenAndServe, 5); err != nil && err != http.ErrServerClosed {
		core.LogFatal(err)
	}
}

/* -~-~-~-~-~ Shutdown -~-~-~-~-~- */

// Stops the GRPC & HTTP Servers.
func (s *Servers) Shutdown() {
	s.shutdownGRPC()
	s.shutdownHTTP()
}

// Stops the GRPC Server.
func (s *Servers) shutdownGRPC() {
	s.GRPC.GracefulStop()
}

// Stops the HTTP Server.
func (s *Servers) shutdownHTTP() {
	ctx, cancelCtx := god.NewCtxWithTimeout(5 * time.Second)
	defer cancelCtx()
	s.HTTP.Shutdown(ctx)
}
