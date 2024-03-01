package server

import (
	"context"
	"log"
	"net/http"
	"time"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/middleware"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

/* ----------------------------------- */
/*           - HTTP Server -           */
/* ----------------------------------- */

// initHTTPGateway initializes the HTTP Gateway and registers the API methods there as well.
// The gateway will point towards the gRPC server's port.
// This function also adds the HTTP middleware to the server and wraps the mux with an HTTP Logger fn.
func initHTTPGateway(grpcPort, httpPort string, middleware []runtime.ServeMuxOption, options []grpc.DialOption, muxWrapper middleware.MuxWrapperFn) *http.Server {
	mux := runtime.NewServeMux(middleware...)

	if err := usersPB.RegisterUsersServiceHandlerFromEndpoint(context.Background(), mux, grpcPort, options); err != nil {
		log.Fatalf(errMsgStartingHTTP_Fatal, err)
	}

	return &http.Server{Addr: httpPort, Handler: muxWrapper(mux)}
}

// runHTTPGateway runs the HTTP server on a given port.
func runHTTPGateway(server *http.Server) {
	log.Println("Running HTTP!")
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf(errMsgServingHTTP_Fatal, err)
		}
	}()
}

// shutdownHTTPGateway gracefully shuts down the HTTP server.
// It waits for all connections to be closed before shutting down.
func shutdownHTTPGateway(httpServer *http.Server) {
	log.Println("Shutting down HTTP server...")

	ctx, cancel := context.WithTimeout(context.Background(), httpShutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf(errMsgShuttingDownHTTP_Fatal, err)
	}
}

const (
	httpShutdownTimeout = 4 * time.Second

	errMsgServingHTTP_Fatal      = "Failed to serve HTTP: %v"           // Fatal error.
	errMsgStartingHTTP_Fatal     = "Failed to start HTTP gateway: %v"   // Fatal error.
	errMsgShuttingDownHTTP_Fatal = "Failed to shutdown HTTP server: %v" // Fatal error.
)
