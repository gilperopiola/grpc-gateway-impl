package servers

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - gRPC Server -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func NewGRPCServer(service pbs.UsersServiceServer, svOptions []grpc.ServerOption) *grpc.Server {
	grpcServer := grpc.NewServer(svOptions...)
	pbs.RegisterUsersServiceServer(grpcServer, service)
	return grpcServer
}

func RunGRPCServer(grpcServer *grpc.Server) {
	zap.S().Infof("Running gRPC on port %s!\n", core.GRPCPort)

	lis, err := net.Listen("tcp", core.GRPCPort)
	if err != nil {
		zap.S().Fatalf(errs.FatalErrMsgStartingGRPC, err)
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			zap.S().Fatalf(errs.FatalErrMsgServingGRPC, err)
		}
	}()
}

func ShutdownGRPCServer(grpcServer *grpc.Server) {
	zap.S().Info("Shutting down gRPC server...")
	grpcServer.GracefulStop()
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - HTTP Gateway -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func NewHTTPGateway(serveOpts []runtime.ServeMuxOption, middleware func(next http.Handler) http.Handler, dialOpts []grpc.DialOption) *http.Server {
	mux := runtime.NewServeMux(serveOpts...)

	if err := pbs.RegisterUsersServiceHandlerFromEndpoint(context.Background(), mux, core.GRPCPort, dialOpts); err != nil {
		zap.S().Fatalf(errs.FatalErrMsgStartingHTTP, err)
	}

	return &http.Server{
		Addr:    core.HTTPPort,
		Handler: middleware(mux),
	}
}

func RunHTTPGateway(httpGateway *http.Server) {
	zap.S().Infof("Running HTTP on port %s!\n", core.HTTPPort)

	go func() {
		if err := httpGateway.ListenAndServe(); err != http.ErrServerClosed {
			zap.S().Fatalf(errs.FatalErrMsgServingHTTP, err)
		}
	}()
}

func ShutdownHTTPGateway(httpGateway *http.Server) {
	zap.S().Info("Shutting down HTTP server...")
	shutdownTimeout := 4 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := httpGateway.Shutdown(ctx); err != nil {
		zap.S().Fatalf(errs.FatalErrMsgShuttingDownHTTP, err)
	}
}
