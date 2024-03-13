package v1

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/cfg"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/db"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/deps"
	grpcV1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1/deps/grpc"
	httpV1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1/deps/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

/* ----------------------------------- */
/*             - v1 API -              */
/* ----------------------------------- */

// API holds the gRPC & HTTP Servers, and all their deps.
// It has an embedded Service that implements all the API Handlers.
type API struct {
	Service     // API.
	*cfg.Config // API Configuration.

	Repository Repository   // Repository to interact with the database.
	Database   *db.Database // Database connection.

	// GRPCServer and HTTPGateway are the servers we are going to run.
	GRPCServer  deps.Server
	HTTPGateway deps.Server

	// GRPCInterceptors run before or after gRPC calls.
	// GRPCDialOptions configure the communication between the HTTP Gateway and the gRPC Server.
	GRPCInterceptors []grpc.ServerOption
	GRPCDialOptions  []grpc.DialOption

	// HTTPMiddleware are the ServeMuxOptions that run before or after HTTP calls.
	// HTTPMiddlewareWrapper are the middleware that wrap around the HTTP Gateway.
	HTTPMiddleware        []runtime.ServeMuxOption
	HTTPMiddlewareWrapper httpV1.MuxWrapperFunc

	// Deps needed to run the API. Validator, Rate Limiter, Logger, TLS Certs, etc.
	*deps.Deps
}

// NewAPI returns a new API with the given configuration.
func NewAPI(config *cfg.Config) *API {
	return &API{
		Config: config,
		Deps:   deps.NewDeps(),
	}
}

// Init initializes all API deps.
func (a *API) Init() {
	a.InitLogger()
	a.InitValidator()
	a.InitAuthenticator()
	a.InitPwdHasher()
	a.InitRateLimiter()
	a.InitTLSDeps()
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

	a.Database.Close()
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
	a.DBConfig.AdminPassword = a.PwdHasher.Hash(a.DBConfig.AdminPassword) // Hash admin pwd.
	a.Database = db.NewDatabase(a.DBConfig)
	a.Repository = NewRepository(a.Database)
}

func (a *API) InitService() {
	a.Service = NewService(a.Repository, a.Authenticator, a.PwdHasher)
}

func (a *API) InitGRPCServer() {
	a.GRPCInterceptors = grpcV1.AllInterceptors(a.Deps, a.TLSConfig.Enabled)
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
	a.LoggerOptions = deps.NewLoggerOptions()
	a.Logger = deps.NewLogger(a.IsProd, a.LoggerOptions...)
}

func (a *API) InitValidator() {
	a.Validator = deps.NewValidator()
}

func (a *API) InitAuthenticator() {
	a.Authenticator = deps.NewJWTAuthenticator(a.JWTConfig.Secret, a.JWTConfig.SessionDays)
}

func (a *API) InitPwdHasher() {
	a.PwdHasher = deps.NewPwdHasher(a.JWTConfig.Secret)
}

func (a *API) InitRateLimiter() {
	a.RateLimiter = rate.NewLimiter(rate.Limit(a.RateLimiterConfig.TokensPerSecond), a.RateLimiterConfig.MaxTokens)
}

func (a *API) InitTLSDeps() {
	if a.TLSConfig.Enabled {
		a.ServerCert = deps.NewTLSCertPool(a.TLSConfig.CertPath)
		a.ServerCreds = deps.NewServerTransportCredentials(a.TLSConfig.CertPath, a.TLSConfig.KeyPath)
	}

	a.ClientCreds = deps.NewClientTransportCredentials(a.TLSConfig.Enabled, a.ServerCert)
}
