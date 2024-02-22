package main

import (
	"crypto/x509"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/interceptors"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/middleware"
	"github.com/gilperopiola/grpc-gateway-impl/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

/* ----------------------------------- */
/*        - gRPC Gateway Impl -        */
/* ----------------------------------- */

// Welcome~!
// This is the entrypoint of our app. Here we start the gRPC server and point the HTTP Gateway towards it.

func main() {
	// New App, loading config from .env vars or default values.
	app := App{Config: v1.LoadConfig()}

	// Load all dependencies, including the servers.
	app.loadAllDependencies()

	// Run servers.
	server.RunGRPCServer(app.GRPCServer, app.Config.GRPCPort)
	server.RunHTTPServer(app.HTTPGateway)
	time.Sleep(1 * time.Second)
	log.Println("... ¡gRPC and HTTP OK! ...")

	// Wait for shutdown.
	app.waitForGracefulShutdown()
}

// App holds every dependency we need to initialize and run the servers. And the servers themselves.
// It has an embedded *v1.API, which is our concrete implementation of the gRPC API.
type App struct {

	// Core components - API, Config, Service, GRPCServer & HTTPGateway.
	*v1.API
	Config      *v1.APIConfig
	Service     v1.ServiceLayer
	GRPCServer  *grpc.Server
	HTTPGateway *http.Server

	// GRPCInterceptors run before the gRPC server handles each request.
	// Right now we only have a Logger and a Validator.
	GRPCInterceptors grpc.ServerOption

	// GRPCDialOptions configure the communication between the HTTP Gateway and the gRPC server.
	GRPCDialOptions v1.GRPCDialOptionsI

	// HTTPMiddleware run before the HTTP Gateway handles each request.
	// Right now we only have an Error Handler and a Response Modifier.
	HTTPMiddleware v1.MiddlewareI

	// ProtoValidator is used to validate the incoming gRPC & HTTP requests.
	// It uses the bufbuild/protovalidate library to enforce the validation rules written in the .proto files.
	ProtoValidator v1.ProtoValidatorI

	// CertPool is a pool of certificates to use for the server.
	// It is used to validate the gRPC server's certificate on the HTTP Gateway calls.
	TLSCertPool *x509.CertPool
}

func (app *App) loadAllDependencies() {
	// Load TLSCertPool, ProtoValidator, GRPC Interceptors & Dial Options, HTTP Middleware.
	app.loadNonCoreComponents()
	// Service holds the business logic of our API.
	app.Service = newService()
	// API is our concrete implementation of the gRPC API defined in the .proto files.
	app.API = newAPI(app.Service)
	// GRPCServer and HTTPGateway are the servers we are going to run.
	app.GRPCServer = newGRPCServer(app.API, app.GRPCInterceptors)
	app.HTTPGateway = newHTTPGateway(app.Config, app.HTTPMiddleware, app.GRPCDialOptions)
}

// loadNonCoreComponents loads TLSCertPool, ProtoValidator, GRPC Interceptors & Dial Options, HTTP Middleware.
func (app *App) loadNonCoreComponents() {
	// TLSCertPool is a pool of certificates used to guarantee the authenticity of
	// the gRPC server's certificate on the HTTP Gateway calls.
	app.TLSCertPool = newTLSCertPool()

	// ProtoValidator is used to validate the incoming gRPC & HTTP requests.
	app.ProtoValidator = newProtoValidator()

	// GRPCInterceptors include requests logging and validation.
	app.GRPCInterceptors = newGRPCInterceptors(app.ProtoValidator)

	// GRPCDialOptions include the gRPC server address and the TLS credentials.
	app.GRPCDialOptions = newGRPCDialOptions(app.Config.TLSEnabled, app.TLSCertPool)

	// HTTPMiddleware includes an Error Handler and a Response Modifier.
	app.HTTPMiddleware = newHTTPMiddleware()
}

// waitForGracefulShutdown waits for a SIGINT or SIGTERM to gracefully shutdown the servers.
func (aw *App) waitForGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // SIGINT and SIGTERM
	<-c

	server.ShutdownGRPCServer(aw.GRPCServer)
	server.ShutdownHTTPServer(aw.HTTPGateway)

	log.Println("Servers stopped! Bye bye~")
}

/* ----------------------------------- */
/*      - Initialize Components -      */
/* ----------------------------------- */

func newGRPCInterceptors(protoValidator v1.ProtoValidatorI) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		interceptors.NewGRPCValidator(protoValidator),
	)
}

func newHTTPMiddleware() v1.MiddlewareI {
	return v1.MiddlewareI{
		middleware.NewHTTPLogger(),
		middleware.NewHTTPErrorHandler(),
		middleware.NewHTTPResponseModifier(),
	}
}

func newGRPCDialOptions(tlsEnabled bool, serverCert *x509.CertPool) v1.GRPCDialOptionsI {

	// Unless TLS is enabled, we use insecure credentials.
	transportCredentials := grpc.WithTransportCredentials(insecure.NewCredentials())
	if tlsEnabled {
		transportCredentials = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(serverCert, ""))
	}

	return []grpc.DialOption{
		transportCredentials,
		//grpc.WithUserAgent("gRPC Gateway Implementation by @gilperopiola"),
	}
}

func newGRPCServer(api *v1.API, interceptors grpc.ServerOption) *grpc.Server {
	return server.InitGRPCServer(api, interceptors)
}

func newHTTPGateway(config *v1.APIConfig, middleware v1.MiddlewareI, grpcDialOptions v1.GRPCDialOptionsI) *http.Server {
	return server.InitHTTPGateway(config.GRPCPort, config.HTTPPort, middleware, grpcDialOptions)
}

func newAPI(service v1.ServiceLayer) *v1.API {
	return v1.NewAPI(service)
}

func newService() v1.ServiceLayer {
	return v1.NewService()
}

func newProtoValidator() v1.ProtoValidatorI {
	return v1.NewProtoValidator()
}

func newTLSCertPool() *x509.CertPool {
	return server.LoadTLSCertPool()
}

/* ----------------------------------- */
/*              - T0D0 -               */
/* ----------------------------------- */
/* Buf file / Dockerfile / Docker-compose / Kubernetes / CI-CD / Tests /
/* Logging / Metrics / Tracing / Security / Caching / Rate limiting /
/* Postman collection / Full Swagger
/* -------------------------------------------------------------------------- */
