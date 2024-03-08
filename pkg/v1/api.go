package v1

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/cfg"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/dependencies"
	grpcV1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1/grpc"
	httpV1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1/http"
	"golang.org/x/time/rate"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

/* ----------------------------------- */
/*             - v1 API -              */
/* ----------------------------------- */

// API holds the gRPC & HTTP Servers, and all their dependencies.
// It has an embedded Service that implements all the API Handlers.
type API struct {
	Service     // API.
	*cfg.Config // API Configuration.

	// GRPCServer and HTTPGateway are the servers we are going to run.
	GRPCServer  Server
	HTTPGateway Server

	// GRPCInterceptors run before or after gRPC calls.
	// GRPCDialOptions configure the communication between the HTTP Gateway and the gRPC Server.
	GRPCInterceptors []grpc.ServerOption
	GRPCDialOptions  []grpc.DialOption

	// HTTPMiddleware are the ServeMuxOptions that run before or after HTTP calls.
	// HTTPMiddlewareWrapper are the middleware that wrap around the HTTP Gateway.
	HTTPMiddleware        []runtime.ServeMuxOption
	HTTPMiddlewareWrapper httpV1.MuxWrapperFunc

	Repository Repository // Repository to interact with the database.

	// Dependencies needed to run the API.
	// Validator, Rate Limiter, Logger, TLS Certs, etc.
	*dependencies.Dependencies
}

// NewAPI returns a new API with the given configuration.
func NewAPI(config *cfg.Config) *API {
	return &API{
		Config: config,
		Dependencies: &dependencies.Dependencies{
			TLSDependencies: &dependencies.TLSDependencies{},
		},
	}
}

// Server is an interface that abstracts the gRPC & HTTP Servers.
// Both servers have the same methods, but they are implemented differently.
type Server interface {
	Init()
	Run()
	Shutdown()
}

// Init initializes all API dependencies.
func (a *API) Init() {
	a.InitLogger()
	a.InitValidator()
	a.InitRateLimiter()
	a.InitTLSDependencies()
	a.InitRepository()
	a.InitService()
	a.InitGRPCServer()
	a.InitHTTPGateway()
}

// Run runs the gRPC & HTTP Servers.
func (a *API) Run() {
	a.GRPCServer.Run()
	a.HTTPGateway.Run()
}

// WaitForGracefulShutdown waits for a SIGINT or SIGTERM to gracefully shutdown the servers.
func (a *API) WaitForGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	a.GRPCServer.Shutdown()
	a.HTTPGateway.Shutdown()
	log.Println("Servers stopped! Bye bye~")
}

/* ----------------------------------- */
/*         - Dependency Mgmt -         */
/* ----------------------------------- */

func (a *API) InitConfig() {
	a.Config = cfg.Init()
}

func (a *API) InitRepository() {
	a.Repository = NewRepository(NewDatabase())
}

func (a *API) InitService() {
	a.Service = NewService(a.Repository)
}

func (a *API) InitGRPCServer() {
	a.GRPCInterceptors = grpcV1.AllInterceptors(a.Dependencies, a.TLSConfig.Enabled)
	a.GRPCDialOptions = grpcV1.AllDialOptions(a.ClientCreds)
	a.GRPCServer = grpcV1.NewGRPCServer(a.Config.GRPCPort, a.Service, a.GRPCInterceptors)
	a.GRPCServer.Init()
}

func (a *API) InitHTTPGateway() {
	a.HTTPMiddleware = httpV1.AllMiddleware()
	a.HTTPMiddlewareWrapper = httpV1.AllMiddlewareWrapper(a.Logger)
	a.HTTPGateway = httpV1.NewHTTPGateway(a.MainConfig, a.HTTPMiddleware, a.HTTPMiddlewareWrapper, a.GRPCDialOptions)
	a.HTTPGateway.Init()
}

func (a *API) InitLogger() {
	a.LoggerOptions = dependencies.NewLoggerOptions()
	a.Logger = dependencies.NewLogger(a.IsProd, a.LoggerOptions...)
}

func (a *API) InitValidator() {
	a.Validator = dependencies.NewValidator()
}

func (a *API) InitRateLimiter() {
	a.RateLimiter = rate.NewLimiter(rate.Limit(a.RateLimiterConfig.TokensPerSecond), a.RateLimiterConfig.MaxTokens)
}

func (a *API) InitTLSDependencies() {
	if a.TLSConfig.Enabled {
		a.ServerCert = dependencies.NewTLSCertPool(a.TLSConfig.CertPath)
		a.ServerCreds = dependencies.NewServerTransportCredentials(a.TLSConfig.CertPath, a.TLSConfig.KeyPath)
	}

	a.ClientCreds = dependencies.NewClientTransportCredentials(a.TLSConfig.Enabled, a.ServerCert)
}
