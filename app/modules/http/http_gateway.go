package http

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

	port              string
	grpcPort          string
	middleware        []runtime.ServeMuxOption
	middlewareWrapper MuxWrapperFunc
	grpcOptions       []grpc.DialOption
}

// NewHTTPGateway returns a new instance of HTTPGateway.
func NewHTTPGateway(c *core.GeneralCfg, middleware []runtime.ServeMuxOption, muxWrapper MuxWrapperFunc, grpcOpts []grpc.DialOption) *HTTPGateway {
	return &HTTPGateway{
		port:              c.HTTPPort,
		grpcPort:          c.GRPCPort,
		middleware:        middleware,
		middlewareWrapper: muxWrapper,
		grpcOptions:       grpcOpts,
	}
}

type MuxWrapperFunc func(next http.Handler) http.Handler

// Init initializes the HTTP Gateway and registers the API endpoints. It will point towards the gRPC Server's port.
// It also adds the HTTP middleware and wraps the mux.
func (h *HTTPGateway) Init() {
	mux := runtime.NewServeMux(h.middleware...)

	if err := pbs.RegisterUsersServiceHandlerFromEndpoint(context.Background(), mux, h.grpcPort, h.grpcOptions); err != nil {
		zap.S().Fatalf(errs.FatalErrMsgStartingHTTP, err)
	}

	h.Server = &http.Server{Addr: h.port, Handler: h.middlewareWrapper(mux)}
}

// Run runs the HTTP Gateway.
func (h *HTTPGateway) Run() {
	zap.S().Infof("Running HTTP on port %s!\n", h.port)

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
