package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/cmd/server"
	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"
	v1Service "github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"

	"google.golang.org/grpc"
)

// Welcome~!
// This is the entrypoint of our app. Here we start the gRPC server and point the HTTP Gateway towards it.

func main() {
	// Get env vars.
	grpcPort, httpPort := getEnvVars()

	// Init API, Interceptors, Middleware.
	api := v1.NewAPI(v1Service.NewService())
	interceptors := v1.GetInterceptorsAsServerOption()
	middleware := v1.GetHTTPMiddlewareAsMuxOptions()

	// Init servers.
	grpcServer := server.InitGRPCServer(api, interceptors)
	httpGateway := server.InitHTTPGateway(grpcPort, httpPort, middleware)

	// Run servers.
	server.RunGRPCServer(grpcServer, grpcPort)
	server.RunHTTPServer(httpGateway)
	time.Sleep(1 * time.Second)
	log.Println("... ¡gRPC and HTTP OK! ...")

	// Wait for shutdown.
	waitForGracefulShutdown(grpcServer, httpGateway)
}

// getEnvVars returns the values of all env vars.
func getEnvVars() (grpcPort string, httpPort string) {
	grpcPort = getEnvVar("GRPC_PORT", ":50053")
	httpPort = getEnvVar("HTTP_PORT", ":8083")
	return
}

// getEnvVar returns the value of an env var, or a fallback.
func getEnvVar(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// waitForGracefulShutdown waits for a SIGINT or SIGTERM to gracefully shutdown the servers.
func waitForGracefulShutdown(grpcServer *grpc.Server, httpServer *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // SIGINT and SIGTERM
	<-c

	server.ShutdownGRPCServer(grpcServer)
	server.ShutdownHTTPServer(httpServer)

	log.Println("Servers stopped! Bye bye~")
}

/* ----------------------------------- */
/*              - T0D0 -               */
/* ----------------------------------- */
/* Buf file / Dockerfile / Docker-compose / Kubernetes / CI-CD / Tests /
/* Logging / Metrics / Tracing / Security / Caching / Rate limiting /
/* Postman collection / Full Swagger
/* -------------------------------------------------------------------------- */
