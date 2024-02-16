package server

import (
	"context"
	"log"
	"net/http"
	"time"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/* ----------------------------------- */
/*           - HTTP Server -           */
/* ----------------------------------- */

// InitHTTPGateway initializes the HTTP Gateway and registers the API methods there as well.
// The gateway will point towards the gRPC server's port.
// This function also adds the HTTP middleware to the server.
func InitHTTPGateway(grpcPort, httpPort string, middleware []runtime.ServeMuxOption) *http.Server {
	mux := runtime.NewServeMux(middleware...)
	opts := getGRPCDialOptions()

	if err := usersPB.RegisterUsersServiceHandlerFromEndpoint(context.Background(), mux, grpcPort, opts); err != nil {
		log.Fatalf(errMsgGateway, err)
	}

	return &http.Server{Addr: httpPort, Handler: mux}
}

// RunHTTPServer runs the HTTP server on a given port.
func RunHTTPServer(server *http.Server) {
	log.Println("Running HTTP!")
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf(errMsgServeHTTP, err)
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
		log.Fatalf(errMsgShutdown, err)
	}
}

// getGRPCDialOptions returns the gRPC connection dial options.
func getGRPCDialOptions() []grpc.DialOption {
	transportCredentialsOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	return []grpc.DialOption{transportCredentialsOption}
}

const (
	errMsgServeHTTP = "Failed to serve HTTP: %v"
	errMsgGateway   = "Failed to start HTTP gateway: %v"
	errMsgShutdown  = "Failed to shutdown HTTP server: %v"

	shutdownTimeout = 4 * time.Second
)
