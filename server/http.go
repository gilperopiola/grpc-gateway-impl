package server

import (
	"context"
	"log"
	"net/http"
	"time"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

/* ----------------------------------- */
/*           - HTTP Server -           */
/* ----------------------------------- */

// InitHTTPGateway initializes the HTTP Gateway and registers the API methods there as well.
// The gateway will point towards the gRPC server's port.
// This function also adds the HTTP middleware to the server and wraps the mux with an HTTP Logger fn.
func InitHTTPGateway(grpcPort, httpPort string, middleware []runtime.ServeMuxOption, options []grpc.DialOption, loggerWrapper func(next http.Handler) http.Handler) *http.Server {
	mux := runtime.NewServeMux(middleware...)

	if err := usersPB.RegisterUsersServiceHandlerFromEndpoint(context.Background(), mux, grpcPort, options); err != nil {
		log.Fatalf(msgErrStartingGateway_Fatal, err)
	}

	return &http.Server{Addr: httpPort, Handler: loggerWrapper(mux)}
}

// runHTTPServer runs the HTTP server on a given port.
func runHTTPServer(server *http.Server) {
	log.Println("Running HTTP!")
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf(msgErrServingHTTP_Fatal, err)
		}
	}()
}

// shutdownHTTPServer gracefully shuts down the HTTP server.
// It waits for all connections to be closed before shutting down.
func shutdownHTTPServer(httpServer *http.Server) {
	log.Println("Shutting down HTTP server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf(msgErrShuttingDownHTTPServer_Fatal, err)
	}
}

const (
	shutdownTimeout = 4 * time.Second

	msgErrServingHTTP_Fatal            = "Failed to serve HTTP: %v"
	msgErrStartingGateway_Fatal        = "Failed to start HTTP gateway: %v"
	msgErrShuttingDownHTTPServer_Fatal = "Failed to shutdown HTTP server: %v"
)
