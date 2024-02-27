package server

import (
	"context"
	"log"
	"net/http"
	"time"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

/* ----------------------------------- */
/*           - HTTP Server -           */
/* ----------------------------------- */

// InitHTTPGateway initializes the HTTP Gateway and registers the API methods there as well.
// The gateway will point towards the gRPC server's port.
// This function also adds the HTTP middleware to the server and wraps the mux with an HTTP Logger fn.
func InitHTTPGateway(grpcPort, httpPort string, middleware v1.MiddlewareI, options v1.GRPCDialOptionsI, loggerWrapper func(next http.Handler) http.Handler) *http.Server {
	mux := runtime.NewServeMux(middleware...)

	if err := usersPB.RegisterUsersServiceHandlerFromEndpoint(context.Background(), mux, grpcPort, options); err != nil {
		log.Fatalf(msgErrStartingGateway_Fatal, err)
	}

	return &http.Server{Addr: httpPort, Handler: loggerWrapper(mux)}
}

// RunHTTPServer runs the HTTP server on a given port.
func RunHTTPServer(server *http.Server) {
	log.Println("Running HTTP!")
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf(msgErrServingHTTP_Fatal, err)
		}
	}()
}

// ShutdownHTTPServer gracefully shuts down the HTTP server.
// It waits for all connections to be closed before shutting down.
func ShutdownHTTPServer(httpServer *http.Server) {
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
