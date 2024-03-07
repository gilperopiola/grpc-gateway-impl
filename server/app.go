package server

import (
	"crypto/x509"
	"log"
	"os"
	"os/signal"
	"syscall"

	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/middleware"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"
	"github.com/gilperopiola/grpc-gateway-impl/server/config"

	"github.com/bufbuild/protovalidate-go"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Server is an interface that connects the gRPC and HTTP servers.
// Both servers have the same methods, but they are implemented differently.
type Server interface {
	Init()
	Run()
	Shutdown()
}

/* ----------------------------------- */
/*         - Main Application -        */
/* ----------------------------------- */

// App holds every dependency we need to init and run the servers, and the servers themselves.
// It has an embedded *v1.API, which is our concrete implementation of the gRPC API.
type App struct {
	*v1.API

	// Cfg holds the configuration of our API.
	Cfg *config.Config

	// Service holds the business logic of our API.
	Service service.Service

	// GRPCServer and HTTPGateway are the servers we are going to run.
	GRPCServer  Server
	HTTPGateway Server

	// GRPCInterceptors run before or after gRPC calls.
	// Right now we only have a Logger and a Validator.
	//
	// GRPCDialOptions configure the communication between the HTTP Gateway and the gRPC Server.
	// Needed to establish a secure connection.
	GRPCInterceptors []grpc.ServerOption
	GRPCDialOptions  []grpc.DialOption

	// HTTPMiddleware are the ServeMuxOptions that run before or after HTTP calls.
	// Right now we only have an Error Handler and a Response Modifier.
	//
	// HTTPMiddlewareWrapper are the middleware that wrap around the HTTP server.
	// Right now we only have a Logger.
	//
	// They are divided into two different types because the ServeMuxOptions are used to configure the ServeMux,
	// and the Wrapper is used to wrap the ServeMux with middleware.
	HTTPMiddleware        []runtime.ServeMuxOption
	HTTPMiddlewareWrapper middleware.MuxWrapperFunc

	// Logger is used to log every gRPC request that comes in through the gRPC
	// It's used on an interceptor.
	//
	// LoggerOptions are the options we can pass to the Logger.
	Logger        *zap.Logger
	LoggerOptions []zap.Option

	// ProtoValidator is used to validate the incoming gRPC & HTTP requests.
	// It uses the bufbuild/protovalidate library to enforce the validation rules written in the .proto files.
	ProtoValidator *protovalidate.Validator

	// TLSServerCert is a pool of certificates to use for the
	// It is used to validate the gRPC server's certificate on the HTTP Gateway calls.
	// TLSServerCredentials and TLSClientCredentials are used to establish the
	// secure connection between the HTTP Gateway and the gRPC Server.
	TLSServerCert  *x509.CertPool
	TLSServerCreds credentials.TransportCredentials
	TLSClientCreds credentials.TransportCredentials
}

// NewApp returns a new App with the given configuration.
func NewApp(config *config.Config) *App {
	return &App{Cfg: config}
}

// Init initializes all App dependencies.
func (a *App) Init() {

	// Logger.
	a.NewLoggerOptions()
	a.NewLogger()

	// Validator.
	a.NewProtoValidator()

	// TLS Server objects are only loaded if TLS is enabled.
	if a.Cfg.TLS.Enabled {
		a.NewTLSServerCert()
		a.NewTLSServerCreds()
	}

	// TLS Client Credentials are always loaded, with an insecure option if TLS is not enabled.
	a.NewTLSClientCreds()

	// gRPC Interceptors and Dial Options.
	a.NewGRPCInterceptors()
	a.NewGRPCDialOptions()

	// HTTP Middleware and Mux Wrapper.
	a.NewHTTPMiddleware()
	a.NewHTTPMiddlewareWrapper()

	// Service and API.
	a.NewService()
	a.NewAPI()

	// gRPC and HTTP Servers.
	a.NewGRPCServer()
	a.NewHTTPGateway()
}

// Run runs the gRPC and HTTP servers.
func (a *App) Run() {
	a.GRPCServer.Run()
	a.HTTPGateway.Run()
}

// WaitForGracefulShutdown waits for a SIGINT or SIGTERM to gracefully shutdown the servers.
func (a *App) WaitForGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // SIGINT and SIGTERM
	<-c

	a.GRPCServer.Shutdown()
	a.HTTPGateway.Shutdown()

	log.Println("Servers stopped! Bye bye~")
}
