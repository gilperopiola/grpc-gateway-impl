package server

import (
	"crypto/x509"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/interceptors"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/middleware"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"

	"github.com/bufbuild/protovalidate-go"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// App holds every dependency we need to init and run the servers, and the servers themselves.
// It has an embedded *v1.API, which is our concrete implementation of the gRPC API.
type App struct {
	*v1.API

	// Config holds the configuration of our API.
	// Service holds the business logic of our API.
	// GRPCServer and HTTPGateway are the servers we are going to run.
	Config      *Config
	Service     service.ServiceLayer
	GRPCServer  *grpc.Server
	HTTPGateway *http.Server

	// GRPCInterceptors run before or after gRPC calls.
	// Right now we only have a Logger and a Validator.
	//
	// GRPCDialOptions configure the communication between the HTTP Gateway and the gRPC
	// Needed to establish a secure connection to the gRPC
	//
	// HTTPMiddleware run before or after HTTP calls.
	// Right now we only have an Error Handler and a Response Modifier.
	GRPCInterceptors grpc.ServerOption
	GRPCDialOptions  []grpc.DialOption
	HTTPMiddleware   []runtime.ServeMuxOption

	// Logger is used to log every gRPC request that comes in through the gRPC
	// It's used on an interceptor.
	//
	// ProtoValidator is used to validate the incoming gRPC & HTTP requests.
	// It uses the bufbuild/protovalidate library to enforce the validation rules written in the .proto files.
	//
	// TLSCertPool is a pool of certificates to use for the
	// It is used to validate the gRPC server's certificate on the HTTP Gateway calls.
	Logger         *zap.Logger
	ProtoValidator *protovalidate.Validator
	TLSCertPool    *x509.CertPool
}

// NewApp returns a new App with the given configuration.
func NewApp(config *Config) *App {
	return &App{
		Config: config,
	}
}

// InitGeneralDependencies initializes stuff.
func (a *App) InitGeneralDependencies() {
	loggerOptions := []zap.Option{
		zap.AddStacktrace(zap.DPanicLevel), // Add stack trace to panic logs.
	}
	a.Logger = v1.NewLogger(a.Config.IsProd, loggerOptions)
	a.ProtoValidator = interceptors.NewProtoValidator()
	a.TLSCertPool = loadTLSCertPool(a.Config.TLS.CertPath)
}

// InitGRPCAndHTTPDependencies initializes gRPC and HTTP stuff.
func (a *App) InitGRPCAndHTTPDependencies() {
	a.GRPCInterceptors = interceptors.GetAll(a.Logger, a.ProtoValidator)
	a.GRPCDialOptions = getAllDialOptions(a.Config.TLS.Enabled, a.TLSCertPool)
	a.HTTPMiddleware = middleware.GetAll()
}

// InitAPI initializes the API and the Service.
func (a *App) InitAPI() {
	a.Service = service.NewService()
	a.API = v1.NewAPI(a.Service)
}

// InitServers initializes the gRPC and HTTP servers.
func (a *App) InitServers() {
	a.GRPCServer = InitGRPCServer(a.API, a.Config.TLS, a.GRPCInterceptors)
	a.HTTPGateway = InitHTTPGateway(a.Config.GRPCPort, a.Config.HTTPPort, a.HTTPMiddleware, a.GRPCDialOptions, middleware.LogHTTP(a.Logger))
}

// RunServers runs the gRPC and HTTP servers.
func (a *App) RunServers() {
	runGRPCServer(a.GRPCServer, a.Config.GRPCPort)
	runHTTPServer(a.HTTPGateway)
}

// WaitForGracefulShutdown waits for a SIGINT or SIGTERM to gracefully shutdown the servers.
func (a *App) WaitForGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // SIGINT and SIGTERM
	<-c

	shutdownGRPCServer(a.GRPCServer)
	shutdownHTTPServer(a.HTTPGateway)

	log.Println("Servers stopped! Bye bye~")
}
