package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"
	v1Service "github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	gRPCPort = ":50051"
	httpPort = ":8080"
)

// Welcome~!
// This is the entrypoint of our app. Here we start the gRPC server and point the HTTP Gateway towards it.

func main() {

	var (
		service     = v1Service.NewService()
		grpcServer  = initGRPCServer(service)
		httpGateway = initHTTPGateway()
	)

	runGRPCServer(grpcServer)
	runHTTPServer(httpGateway)

	log.Println("... ¡OK! ...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // SIGINT and SIGTERM

	<-c

	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()

	log.Println("Shutting down HTTP server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpGateway.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown HTTP server: %v", err)
	}

	log.Println("Servers successfully stopped")
}

/* ----------------------------------- */
/*             - gRPC -                */
/* ----------------------------------- */

const (
	errMsgListenGRPC = "Failed to listen gRPC: %v"
	errMsgServeGRPC  = "Failed to serve gRPC: %v"
)

// initGRPCServer initializes the gRPC server and registers the API methods.
// The HTTP Gateway will point towards this server.
func initGRPCServer(service v1Service.ServiceLayer) *grpc.Server {
	var (
		interceptors = v1.GetInterceptors()
		grpcServer   = grpc.NewServer(interceptors)
	)
	usersPB.RegisterUsersServiceServer(grpcServer, &v1.API{Service: service})
	return grpcServer
}

// runGRPCServer runs the gRPC server on a given port.
// It listens for incoming gRPC requests and serves them.
func runGRPCServer(grpcServer *grpc.Server) {
	log.Println("Running gRPC!")

	lis, err := net.Listen("tcp", gRPCPort)
	if err != nil {
		log.Fatalf(errMsgListenGRPC, err)
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf(errMsgServeGRPC, err)
		}
	}()

	return
}

/* ----------------------------------- */
/*             - HTTP -                */
/* ----------------------------------- */

const (
	errMsgServeHTTP = "Failed to serve HTTP: %v"
	errMsgGateway   = "Failed to start HTTP gateway: %v"
)

// initHTTPGateway initializes the HTTP gateway and registers the API methods there as well.
// The gateway will point towards the gRPC server's port.
func initHTTPGateway() *http.Server {
	mux := runtime.NewServeMux(v1.GetHTTPMiddleware()...)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := usersPB.RegisterUsersServiceHandlerFromEndpoint(context.Background(), mux, gRPCPort, opts); err != nil {
		log.Fatalf(errMsgGateway, err)
	}

	return &http.Server{Addr: ":8080", Handler: mux}
}

func runHTTPServer(server *http.Server) {
	log.Println("Running HTTP!")
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Failed to serve HTTP: %v", err)
		}
	}()
}

// runHTTPGateway runs the HTTP gateway on a given port.
// It listens for incoming HTTP requests and serves them.
func runHTTPGateway(httpGateway *runtime.ServeMux) {
	log.Println("Running HTTP!")

	if err := http.ListenAndServe(httpPort, httpGateway); err != nil {
		log.Fatalf(errMsgServeHTTP, err)
	}
}

/* T0D0: Buf file / Dockerfile / Docker-compose / Kubernetes / CI-CD / Tests /
 * Logging / Metrics / Tracing / Security / Caching / Rate limiting / Postman collection
 */
