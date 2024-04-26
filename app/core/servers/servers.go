package servers

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Servers struct {
	GRPC *grpc.Server
	HTTP *http.Server
}

// Sets up both GRPC and HTTP servers.
func Setup(usersService pbs.UsersServiceServer, toolbox core.Toolbox, tlsEnabled bool) *Servers {
	var (
		grpcServerOpts       = defaultGRPCServerOpts(toolbox, tlsEnabled)
		grpcDialOpts         = defaultGRPCDialOpts(toolbox.GetClientCreds())
		httpServerOpts       = defaultHTTPMuxOpts()
		httpServerMiddleware = defaultHTTPMiddleware()
	)

	zap.S().Info("GRPC Gateway Implementation | Starting Servers 🚀")

	return &Servers{
		newGRPCServer(usersService, grpcServerOpts),
		newHTTPGateway(httpServerOpts, httpServerMiddleware, grpcDialOpts),
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (s *Servers) Run() {
	runGRPC(s.GRPC)
	runHTTP(s.HTTP)

	go func() {
		time.Sleep(time.Second) // T0D0 healtcheck??
		zap.S().Infoln("GRPC Gateway Implementation | Servers OK 🚀")
	}()
}

func runGRPC(grpcSv *grpc.Server) {
	zap.S().Infof("GRPC Gateway Implementation | GRPC Port %s 🚀", core.GRPCPort)

	lis, err := net.Listen("tcp", core.GRPCPort)
	core.LogPanicIfErr(err)

	go func() {
		core.LogPanicIfErr(grpcSv.Serve(lis))
	}()
}

func runHTTP(httpGw *http.Server) {
	zap.S().Infof("GRPC Gateway Implementation | HTTP Port %s 🚀", core.HTTPPort)

	go func() {
		if err := httpGw.ListenAndServe(); err != http.ErrServerClosed {
			core.LogUnexpectedAndPanic(err)
		}
	}()
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (s *Servers) Shutdown() {
	shutdownGRPC(s.GRPC)
	shutdownHTTP(s.HTTP)
}

func shutdownGRPC(grpcSv *grpc.Server) {
	zap.S().Info("GRPC Gateway Implementation | Shutting down GRPC 🛑")
	grpcSv.GracefulStop()
}

func shutdownHTTP(httpGw *http.Server) {
	zap.S().Info("GRPC Gateway Implementation | Shutting down HTTP 🛑")

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	core.LogPanicIfErr(httpGw.Shutdown(ctx))
}
