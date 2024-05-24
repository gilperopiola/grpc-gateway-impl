package servers

import (
	"net"
	"net/http"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var logPrefix = core.AppEmoji + " " + core.AppAlias + " | "

type Servers struct {
	GRPC *grpc.Server
	HTTP *http.Server
}

func Setup(service core.Service, toolbox core.Toolbox) core.Servers {
	var (
		grpcServerOpts   = getGRPCServerOptions(toolbox, core.TLSEnabled)
		grpcDialOpts     = getGRPCDialOptions(toolbox.GetClientCreds())
		httpServeMuxOpts = getHTTPServeMuxOptions()
		httpMiddleware   = getHTTPMiddleware()
	)

	zap.S().Info(logPrefix + " 🚀 Starting Servers")

	return &Servers{
		setupGRPCServer(service, grpcServerOpts),
		setupHTTPGateway(service, httpServeMuxOpts, httpMiddleware, grpcDialOpts),
	}
}

func setupGRPCServer(service core.Service, serverOpts core.GRPCServerOptions) *grpc.Server {
	grpcServer := grpc.NewServer(serverOpts...)
	service.RegisterGRPCServices(grpcServer)
	return grpcServer
}

func setupHTTPGateway(service core.Service, serveMuxOpts []runtime.ServeMuxOption, middlewareFn middlewareFunc, grpcDialOpts core.GRPCDialOptions) *http.Server {
	mux := runtime.NewServeMux(serveMuxOpts...)
	service.RegisterHTTPServices(mux, grpcDialOpts)
	return &http.Server{
		Addr:    core.HTTPPort,
		Handler: middlewareFn(mux),
		// TLSConfig: core.GetTLSConfig(core.GetCertPool(core.CertPath), core.GetServerName(core.ServerName)),
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (s *Servers) Run() {
	runGRPC(s.GRPC)
	runHTTP(s.HTTP)

	go func() {
		time.Sleep(time.Second) // T0D0 healtcheck??
		zap.S().Infoln(logPrefix + " 🚀 Servers OK\n")
	}()
}

func runGRPC(grpcServer *grpc.Server) {
	zap.S().Infof(logPrefix+" 🚀 GRPC Port %s", core.GRPCPort)

	lis, err := net.Listen("tcp", core.GRPCPort)
	core.LogPanicIfErr(err)

	for _, info := range grpcServer.GetServiceInfo() {
		zap.S().Infof(logPrefix+" 🐸 Service Loaded: %s", info.Metadata)
		for _, method := range info.Methods {
			zap.S().Infof(logPrefix+" \tEndpoint Loaded: %s", method.Name)
		}
	}

	go func() {
		core.LogPanicIfErr(grpcServer.Serve(lis))
	}()
}

func runHTTP(httpGateway *http.Server) {
	zap.S().Infof(logPrefix+" 🚀 HTTP Port %s", core.HTTPPort)

	go func() {
		if err := httpGateway.ListenAndServe(); err != http.ErrServerClosed {
			core.LogUnexpectedAndPanic(err)
		}
	}()
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (s *Servers) Shutdown() {
	zap.S().Info(logPrefix + " 🛑 Shutting down GRPC")
	s.GRPC.GracefulStop()

	zap.S().Info(logPrefix + " 🛑 Shutting down HTTP")
	ctx, cancel := core.NewCtxWithTimeout(5 * time.Second)
	defer cancel()
	core.LogPanicIfErr(s.HTTP.Shutdown(ctx))
}
