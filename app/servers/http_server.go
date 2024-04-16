package servers

import (
	"context"
	"net/http"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

/* ----------------------------------- */
/*          - HTTP Gateway -           */
/* ----------------------------------- */

// HTTPGateway is a wrapper around the actual HTTP Server.
type HTTPGateway struct {
	*http.Server
}

// NewHTTPGateway returns a new instance of HTTPGateway.
func NewHTTPGateway(serveOpts []runtime.ServeMuxOption, middleware func(next http.Handler) http.Handler, dialOpts []grpc.DialOption) *HTTPGateway {
	mux := runtime.NewServeMux(serveOpts...)

	if err := pbs.RegisterUsersServiceHandlerFromEndpoint(context.Background(), mux, core.GRPCPort, dialOpts); err != nil {
		zap.S().Fatalf(errs.FatalErrMsgStartingHTTP, err)
	}

	return &HTTPGateway{
		Server: &http.Server{Addr: core.HTTPPort, Handler: middleware(mux)},
	}
}

// Run runs the HTTP Gateway.
func (h *HTTPGateway) Run() {
	zap.S().Infof("Running HTTP on port %s!\n", core.HTTPPort)

	go func() {
		if err := h.ListenAndServe(); err != http.ErrServerClosed {
			zap.S().Fatalf(errs.FatalErrMsgServingHTTP, err)
		}
	}()
}

// Shutdown gracefully shuts down the HTTP Server.
// It waits for all connections to be closed before shutting down.
func (h *HTTPGateway) Shutdown() {
	zap.S().Info("Shutting down HTTP server...")
	shutdownTimeout := 4 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := h.Server.Shutdown(ctx); err != nil {
		zap.S().Fatalf(errs.FatalErrMsgShuttingDownHTTP, err)
	}
}
