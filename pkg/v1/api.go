package v1

import (
	"crypto/x509"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/cfg"
	grpcV1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1/grpc"
	httpV1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1/http"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/misc"

	"github.com/bufbuild/protovalidate-go"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

/* ----------------------------------- */
/*               - API -               */
/* ----------------------------------- */

// API holds every dependency we need to init and run the servers, and the servers themselves.
// It has an embedded Service, which powers up our API.
type API struct {
	Service

	// Cfg holds the configuration of our App.
	Cfg *cfg.Config

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
	HTTPMiddlewareWrapper httpV1.MuxWrapperFunc

	// Logger is used to log every gRPC and HTTP request that comes in.
	// It's used on an interceptor.
	//
	// LoggerOptions are the options we can pass to the Logger.
	Logger        *zap.Logger
	LoggerOptions []zap.Option

	// ProtoValidator is used to validate the incoming gRPC & HTTP requests.
	// It uses the bufbuild/protovalidate library to enforce the validation rules written in the .proto files.
	ProtoValidator *protovalidate.Validator

	// TLSServerCert is a pool of certificates to use for the Server's TLS configuration.
	//
	// TLSServerCreds and TLSClientCreds are used to establish the
	// secure connection between the HTTP Gateway and the gRPC Server.
	TLSServerCert  *x509.CertPool
	TLSServerCreds credentials.TransportCredentials
	TLSClientCreds credentials.TransportCredentials
}

// NewAPI returns a new API with the given configuration.
func NewAPI(config *cfg.Config) *API {
	return &API{Cfg: config}
}

// Server is an interface that connects the gRPC and HTTP servers.
// Both servers have the same methods, but they are implemented differently.
type Server interface {
	Init()
	Run()
	Shutdown()
}

// Init initializes all API dependencies.
func (a *API) Init() {
	a.NewLogger()
	a.NewProtoValidator()
	a.NewTLSCertAndCreds()
	a.NewService()
	a.NewGRPCServer()
	a.NewHTTPGateway()
}

// Run runs the gRPC and HTTP servers.
func (a *API) Run() {
	a.GRPCServer.Run()
	a.HTTPGateway.Run()
}

// WaitForGracefulShutdown waits for a SIGINT or SIGTERM to gracefully shutdown the servers.
func (a *API) WaitForGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // SIGINT and SIGTERM
	<-c

	a.GRPCServer.Shutdown()
	a.HTTPGateway.Shutdown()

	log.Println("Servers stopped! Bye bye~")
}

/* ----------------------------------- */
/*         - Dependency Mgmt -         */
/* ----------------------------------- */

func (a *API) NewConfig() {
	a.Cfg = cfg.New()
}

func (a *API) NewService() {
	a.Service = NewService()
}

func (a *API) NewGRPCServer() {
	a.GRPCInterceptors = grpcV1.AllInterceptors(a.Cfg, a.Logger, a.ProtoValidator, a.TLSServerCreds)
	a.GRPCDialOptions = grpcV1.AllDialOptions(a.TLSClientCreds)
	a.GRPCServer = grpcV1.NewGRPCServer(a.Cfg.GRPCPort, a.Service, a.GRPCInterceptors)
	a.GRPCServer.Init()
}

func (a *API) NewHTTPGateway() {
	a.HTTPMiddleware = httpV1.AllMiddleware()
	a.HTTPMiddlewareWrapper = httpV1.MiddlewareWrapper(a.Logger)
	a.HTTPGateway = httpV1.NewHTTPGateway(a.Cfg.MainConfig, a.HTTPMiddleware, a.HTTPMiddlewareWrapper, a.GRPCDialOptions)
	a.HTTPGateway.Init()
}

func (a *API) NewLogger() {
	a.LoggerOptions = misc.NewLoggerOptions()
	a.Logger = misc.NewLogger(a.Cfg.IsProd, a.LoggerOptions)
}

func (a *API) NewProtoValidator() {
	a.ProtoValidator = misc.NewProtoValidator()
}

func (a *API) NewTLSCertAndCreds() {
	if a.Cfg.TLS.Enabled {
		a.TLSServerCert = misc.NewTLSCertPool(a.Cfg.TLS.CertPath)
		a.TLSServerCreds = misc.NewServerTransportCredentials(a.Cfg.TLS.CertPath, a.Cfg.TLS.KeyPath)
	}
	a.TLSClientCreds = misc.NewClientTransportCredentials(a.Cfg.TLS.Enabled, a.TLSServerCert)
}
