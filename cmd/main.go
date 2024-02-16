package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/config"
	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/interceptors"
	v1Service "github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"
	"github.com/gilperopiola/grpc-gateway-impl/server"

	"google.golang.org/grpc"
)

// Welcome~!
// This is the entrypoint of our app. Here we start the gRPC server and point the HTTP Gateway towards it.

func main() {
	// Get configuration.
	config := config.LoadConfig()

	// Init API, Interceptors, Middleware.
	api := v1.NewAPI(v1Service.NewService())
	protoValidator := interceptors.NewProtoValidator()
	interceptors := interceptors.GetInterceptorsAsServerOption(protoValidator)
	middleware := v1.GetHTTPMiddlewareAsMuxOptions()
	grpcDialOptions := server.GetGRPCDialOptions()

	// Init servers.
	grpcServer := server.InitGRPCServer(api, interceptors)
	httpGateway := server.InitHTTPGateway(config.GRPCPort, config.HTTPPort, middleware, grpcDialOptions)

	// Run servers.
	server.RunGRPCServer(grpcServer, config.GRPCPort)
	server.RunHTTPServer(httpGateway)
	time.Sleep(1 * time.Second)
	log.Println("... ¡gRPC and HTTP OK! ...")

	// Wait for shutdown.
	waitForGracefulShutdown(grpcServer, httpGateway)
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
