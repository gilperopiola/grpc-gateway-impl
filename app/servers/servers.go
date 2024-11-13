package servers

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/logs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/utils"
	"github.com/gilperopiola/grpc-gateway-impl/app/service"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

/* -~-~-~-~-~ GRPC and HTTP Servers -~-~-~-~-~- */

// I just asked GPT and he said technically this project runs 1 server with 2 protocols.
// Still naming this 'Servers'.
type Servers struct {
	GRPC *grpc.Server
	HTTP *http.Server
}

func Setup(services *service.Service, tools core.Tools) *Servers {
	var (
		grpcServerOpts = getGRPCServerOpts(tools, core.G.TLSEnabled) // Interceptors
		grpcDialOpts   = getGRPCDialOpts(tools.GetClientCreds())     // Used by HTTP to connect to GRPC

		httpMiddleware   = getHTTPMiddlewareChain() // HTTP Middleware
		httpServeMuxOpts = getHTTPMuxOpts()         // HTTP Middleware
	)

	servers := Servers{
		GRPC: setupGRPC(services, grpcServerOpts),
		HTTP: setupHTTP(services, httpServeMuxOpts, httpMiddleware, grpcDialOpts...),
	}

	logs.InitModuleOK("Servers", "📡")
	return &servers
}

func (s *Servers) Run() {
	go s.runGRPC()
	go s.runHTTP()
	logGRPCRegisteredEndpoints(s.GRPC)
}

// Stops the GRPC & HTTP Servers.
func (s *Servers) Shutdown() {
	s.shutdownGRPC()
	s.shutdownHTTP()
}

/* -~-~-~-~-~ Setup -~-~-~-~-~- */

func setupGRPC(service *service.Service, serverOpts []grpc.ServerOption) *grpc.Server {
	grpcServer := grpc.NewServer(serverOpts...)
	service.RegisterInGRPC(grpcServer)
	return grpcServer
}

func setupHTTP(service *service.Service, muxOpts []runtime.ServeMuxOption, mw middlewareFunc, dialOpts ...grpc.DialOption) *http.Server {
	mux := runtime.NewServeMux(muxOpts...)
	service.RegisterInHTTP(mux, dialOpts...)
	return &http.Server{
		Addr:    core.G.HTTPPort,
		Handler: mw(mux),
	}
}

/* -~-~-~-~-~ Run -~-~-~-~-~- */

func (s *Servers) runGRPC() {
	var listenGRPC = func() (any, error) { return net.Listen("tcp", core.G.GRPCPort) }

	result, err := utils.RetryFunc(listenGRPC)
	logs.LogFatalIfErr(err)

	listener := result.(net.Listener)
	logs.LogFatalIfErr(s.GRPC.Serve(listener))
}

func (s *Servers) runHTTP() {
	var listenHTTP = func() (any, error) { return net.Listen("tcp", core.G.HTTPPort) }

	result, err := utils.RetryFunc(listenHTTP)
	logs.LogFatalIfErr(err)

	listener := result.(net.Listener)
	if err := s.HTTP.Serve(listener); err != nil && err != http.ErrServerClosed {
		logs.LogFatal(err)
	}
}

/* -~-~-~-~-~ Shutdown -~-~-~-~-~- */

// Stops the GRPC Server.
func (s *Servers) shutdownGRPC() {
	s.GRPC.GracefulStop()
}

// Stops the HTTP Server.
func (s *Servers) shutdownHTTP() {
	ctx, _ := god.NewCtxWithTimeout(5 * time.Second)
	s.HTTP.Shutdown(ctx)
}

/* -~-~-~-~-~ Helpers -~-~-~-~-~- */

func logGRPCRegisteredEndpoints(server *grpc.Server) {
	for serviceName, serviceInfo := range server.GetServiceInfo() {
		log.Printf("\t 🟢 %s ▶ [%s]", serviceName, serviceInfo.Metadata)
		for _, method := range serviceInfo.Methods {
			log.Printf("\t\t · %s", method.Name)
		}
	}
}
