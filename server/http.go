package server

import (
	"context"
	"log"
	"net/http"
	"time"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/middleware"
	"github.com/gilperopiola/grpc-gateway-impl/server/config"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

const ()

/* ----------------------------------- */
/*           - HTTP Server -           */
/* ----------------------------------- */

// initHTTPGateway initializes the HTTP Gateway and registers the API methods there as well.
// The gateway will point towards the gRPC server's port.
// This function also adds the HTTP middleware to the server and wraps the mux with an HTTP Logger func.
func initHTTPGateway(c *config.MainConfig, middleware []runtime.ServeMuxOption, middlewareWr middleware.MuxWrapperFunc, options []grpc.DialOption) *http.Server {
	mux := runtime.NewServeMux(middleware...)

	if err := usersPB.RegisterUsersServiceHandlerFromEndpoint(context.Background(), mux, c.GRPCPort, options); err != nil {
		log.Fatalf(v1.FatalErrMsgStartingHTTP, err)
	}

	return &http.Server{Addr: c.HTTPPort, Handler: middlewareWr(mux)}
}

// runHTTPGateway runs the HTTP server on a given port.
func runHTTPGateway(server *http.Server) {
	log.Println("Running HTTP!")
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf(v1.FatalErrMsgServingHTTP, err)
		}
	}()
}

// shutdownHTTPGateway gracefully shuts down the HTTP server.
// It waits for all connections to be closed before shutting down.
func shutdownHTTPGateway(httpServer *http.Server) {
	log.Println("Shutting down HTTP server...")
	shutdownTimeout := 4 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf(v1.FatalErrMsgShuttingDownHTTP, err)
	}
}
