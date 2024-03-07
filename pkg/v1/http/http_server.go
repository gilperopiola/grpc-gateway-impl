package http

import (
	"context"
	"log"
	"net/http"
	"time"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/cfg"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

/* ----------------------------------- */
/*          - HTTP Gateway -           */
/* ----------------------------------- */

type HTTPGateway struct {
	*http.Server

	cfg          *cfg.MainConfig
	middleware   []runtime.ServeMuxOption
	middlewareWr MuxWrapperFunc
	options      []grpc.DialOption
}

func NewHTTPGateway(c *cfg.MainConfig, middleware []runtime.ServeMuxOption, middlewareWr MuxWrapperFunc, options []grpc.DialOption) *HTTPGateway {
	return &HTTPGateway{
		cfg:          c,
		middleware:   middleware,
		middlewareWr: middlewareWr,
		options:      options,
	}
}

// Init initializes the HTTP Gateway and registers the API methods there as well.
// The gateway will point towards the gRPC server's port.
// This function also adds the HTTP middleware to the server and wraps the mux with an HTTP Logger func.
func (h *HTTPGateway) Init() {
	mux := runtime.NewServeMux(h.middleware...)

	if err := usersPB.RegisterUsersServiceHandlerFromEndpoint(context.Background(), mux, h.cfg.GRPCPort, h.options); err != nil {
		log.Fatalf(errs.FatalErrMsgStartingHTTP, err)
	}

	h.Server = &http.Server{Addr: h.cfg.HTTPPort, Handler: h.middlewareWr(mux)}
}

// Run runs the HTTP Gateway.
func (h *HTTPGateway) Run() {
	log.Printf("Running HTTP on port %s!\n", h.Addr)
	go func() {
		if err := h.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf(errs.FatalErrMsgServingHTTP, err)
		}
	}()
}

// Shutdown gracefully shuts down the HTTP server.
// It waits for all connections to be closed before shutting down.
func (h *HTTPGateway) Shutdown() {
	log.Println("Shutting down HTTP server...")
	shutdownTimeout := 4 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := h.Server.Shutdown(ctx); err != nil {
		log.Fatalf(errs.FatalErrMsgShuttingDownHTTP, err)
	}
}
