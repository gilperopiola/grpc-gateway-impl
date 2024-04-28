package servers

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Servers struct {
	GRPC *grpc.Server
	HTTP *http.Server
}

// Sets up both GRPC and HTTP servers.
func Setup(service core.Service, toolbox core.Toolbox, tlsEnabled bool) core.Servers {
	var (
		grpcServerOpts   = defaultGRPCServerOpts(toolbox, tlsEnabled)
		grpcDialOpts     = defaultGRPCDialOpts(toolbox.GetClientCreds())
		httpServeMuxOpts = defaultHTTPServeMuxOpts()
		httpMiddleware   = defaultHTTPMiddleware()
	)

	zap.S().Info("GRPC Gateway Implementation | Starting Servers 🚀")

	return &Servers{
		newGRPCServer(service, grpcServerOpts),
		newHTTPGateway(httpServeMuxOpts, httpMiddleware, grpcDialOpts),
	}
}

func (s *Servers) Run() {
	runGRPC(s.GRPC)
	runHTTP(s.HTTP)

	go func() {
		time.Sleep(time.Second) // T0D0 healtcheck??
		zap.S().Infoln("GRPC Gateway Implementation | Servers OK 🚀\n")
	}()
}

func (s *Servers) Shutdown() {
	shutdownGRPC(s.GRPC)
	shutdownHTTP(s.HTTP)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func runGRPC(grpcServer *grpc.Server) {
	zap.S().Infof("GRPC Gateway Implementation | GRPC Port %s 🚀", core.GRPCPort)

	lis, err := net.Listen("tcp", core.GRPCPort)
	core.LogPanicIfErr(err)

	go func() {
		core.LogPanicIfErr(grpcServer.Serve(lis))
	}()
}

func runHTTP(httpGateway *http.Server) {
	zap.S().Infof("GRPC Gateway Implementation | HTTP Port %s 🚀", core.HTTPPort)

	go func() {
		if err := httpGateway.ListenAndServe(); err != http.ErrServerClosed {
			core.LogUnexpectedAndPanic(err)
		}
	}()
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func shutdownGRPC(grpcServer *grpc.Server) {
	zap.S().Info("GRPC Gateway Implementation | Shutting down GRPC 🛑")
	grpcServer.GracefulStop()
}

func shutdownHTTP(httpGateway *http.Server) {
	zap.S().Info("GRPC Gateway Implementation | Shutting down HTTP 🛑")

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	core.LogPanicIfErr(httpGateway.Shutdown(ctx))
}
