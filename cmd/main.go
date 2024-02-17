package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/interceptors"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/middleware"
	v1Service "github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"
	"github.com/gilperopiola/grpc-gateway-impl/server"

	"github.com/bufbuild/protovalidate-go"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

/* ----------------------------------- */
/*        - gRPC Gateway Impl -        */
/* ----------------------------------- */

// Welcome~!
// This is the entrypoint of our app. Here we start the gRPC server and point the HTTP Gateway towards it.

func main() {

	// New APIComponents & Config.
	ac := APIComponents{
		Config: v1.LoadConfig(),
	}

	// ProtoValidator is used to validate the incoming gRPC & HTTP requests.
	ac.ProtoValidator = newProtoValidator()

	// GRPCInterceptors include requests logging and validation.
	ac.GRPCInterceptors = newGRPCInterceptors(ac.ProtoValidator)

	// GRPCDialOptions include the gRPC server address and the TLS credentials.
	ac.GRPCDialOptions = newGRPCDialOptions()

	// HTTPMiddleware includes an Error Handler and a Response Modifier.
	ac.HTTPMiddleware = newHTTPMiddleware()

	// Service holds the business logic of our API.
	ac.Service = newService()

	// API is our concrete implementation of the gRPC API defined in the .proto files.
	ac.API = newAPI(ac.Service)

	// GRPCServer and HTTPGateway are the servers we are going to run.
	ac.GRPCServer = newGRPCServer(ac.API, ac.GRPCInterceptors)
	ac.HTTPGateway = newHTTPGateway(ac.Config, ac.HTTPMiddleware, ac.GRPCDialOptions)

	// Run servers.
	server.RunGRPCServer(ac.GRPCServer, ac.Config.GRPCPort)
	server.RunHTTPServer(ac.HTTPGateway)
	time.Sleep(1 * time.Second)
	log.Println("... ¡gRPC and HTTP OK! ...")

	// Wait for shutdown.
	waitForGracefulShutdown(ac.GRPCServer, ac.HTTPGateway)
}

// APIComponents holds every dependency we need to initialize and run the servers. And the servers themselves.
// It's a way to keep everything together.
type APIComponents struct {
	// Core
	API         *v1.API
	Config      *v1.Config
	Service     v1Service.ServiceLayer
	GRPCServer  *grpc.Server
	HTTPGateway *http.Server

	// Non-core
	GRPCInterceptors grpc.ServerOption
	GRPCDialOptions  []grpc.DialOption
	HTTPMiddleware   []runtime.ServeMuxOption
	ProtoValidator   *protovalidate.Validator
}

/* ----------------------------------- */
/*         - Init Components -         */
/* ----------------------------------- */

func newGRPCServer(api *v1.API, interceptors grpc.ServerOption) *grpc.Server {
	return server.InitGRPCServer(api, interceptors)
}

func newHTTPGateway(config *v1.Config, middleware []runtime.ServeMuxOption, grpcDialOptions []grpc.DialOption) *http.Server {
	return server.InitHTTPGateway(config.GRPCPort, config.HTTPPort, middleware, grpcDialOptions)
}

func newAPI(service v1Service.ServiceLayer) *v1.API {
	return v1.NewAPI(service)
}

func newService() v1Service.ServiceLayer {
	return v1Service.NewService()
}

func newGRPCInterceptors(protoValidator *protovalidate.Validator) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		interceptors.NewGRPCLogger(),
		interceptors.NewGRPCValidator(protoValidator),
	)
}

func newGRPCDialOptions() []grpc.DialOption {
	return server.GetGRPCDialOptions()
}

func newHTTPMiddleware() []runtime.ServeMuxOption {
	return []runtime.ServeMuxOption{
		middleware.NewHTTPLogger(),
		middleware.NewHTTPErrorHandler(),
		middleware.NewHTTPResponseModifier(),
	}
}

func newProtoValidator() *protovalidate.Validator {
	return v1.NewProtoValidator()
}

/* ----------------------------------- */
/*              - T0D0 -               */
/* ----------------------------------- */
/* Buf file / Dockerfile / Docker-compose / Kubernetes / CI-CD / Tests /
/* Logging / Metrics / Tracing / Security / Caching / Rate limiting /
/* Postman collection / Full Swagger
/* -------------------------------------------------------------------------- */

// waitForGracefulShutdown waits for a SIGINT or SIGTERM to gracefully shutdown the servers.
func waitForGracefulShutdown(grpcServer *grpc.Server, httpServer *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // SIGINT and SIGTERM
	<-c

	server.ShutdownGRPCServer(grpcServer)
	server.ShutdownHTTPServer(httpServer)

	log.Println("Servers stopped! Bye bye~")
}
