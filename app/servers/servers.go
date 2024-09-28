package servers

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/logs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/utils"
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
		// GRPC Interceptors.
		grpcServerOpts = getGRPCInterceptors(tools, core.G.TLSEnabled)

		// Used by HTTP to connect to GRPC.
		grpcDialOpts = getGRPCDialOpts(tools.GetClientCreds())

		// HTTP Middleware.
		httpMiddleware   = getHTTPMiddlewareChain()
		httpServeMuxOpts = getHTTPMuxOpts()
	)

	return &Servers{
		GRPC: setupGRPC(services, grpcServerOpts),
		HTTP: setupHTTP(services, httpServeMuxOpts, httpMiddleware, grpcDialOpts...),
	}
}

func (s *Servers) Run() {
	logs.Step(2, "Run")
	go func() {
		time.Sleep(100 * time.Millisecond)
		logEndpointsPerService(s.GRPC)
		time.Sleep(1 * time.Second)
		log.Println("")
		logs.Step(3, "Enjoy")
	}()

	go s.runGRPC()
	go s.runHTTP()
}

// Stops the GRPC & HTTP Servers.
func (s *Servers) Shutdown() {
	s.shutdownGRPC()
	s.shutdownHTTP()
}

/* -~-~-~-~-~ Setup -~-~-~-~-~- */

func setupGRPC(service *service.Service, serverOpts []grpc.ServerOption) *grpc.Server {
	grpcServer := grpc.NewServer(serverOpts...)
	service.RegisterGRPCEndpoints(grpcServer)
	return grpcServer
}

func setupHTTP(service *service.Service, muxOpts []runtime.ServeMuxOption, mw middlewareFunc, dialOpts ...grpc.DialOption) *http.Server {
	mux := runtime.NewServeMux(muxOpts...)
	service.RegisterHTTPEndpoints(mux, dialOpts...)
	return &http.Server{
		Addr:    core.G.HTTPPort,
		Handler: mw(mux),
	}
}

/* -~-~-~-~-~ Run -~-~-~-~-~- */

func (s *Servers) runGRPC() {
	result, err := utils.Retry(listenGRPC, 5)
	logs.LogFatalIfErr(err)

	listener := result.(net.Listener)
	logs.LogFatalIfErr(s.GRPC.Serve(listener))
}

func (s *Servers) runHTTP() {
	result, err := utils.Retry(listenHTTP, 5)
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
	ctx, cancelCtx := god.NewCtxWithTimeout(5 * time.Second)
	defer cancelCtx()
	s.HTTP.Shutdown(ctx)
}

/* -~-~-~-~-~ Helpers -~-~-~-~-~- */

func listenGRPC() (any, error) {
	return net.Listen("tcp", core.G.GRPCPort)
}

func listenHTTP() (any, error) {
	return net.Listen("tcp", core.G.HTTPPort)
}

func logEndpointsPerService(server *grpc.Server) {
	// Loop over all services.
	for serviceName, serviceInfo := range server.GetServiceInfo() {
		log.Printf("\t 🟢 %s ▶ [%s]", serviceName, serviceInfo.Metadata)

		// Loop over all methods.
		for _, method := range serviceInfo.Methods {
			log.Printf("\t\t · %s", method.Name)
		}
	}
}
