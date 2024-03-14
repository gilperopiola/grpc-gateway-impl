package v1

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/cfg"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/db"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/deps"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/deps/grpc"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/deps/http"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/time/rate"
	grpc_external "google.golang.org/grpc"
)

/* ----------------------------------- */
/*             - v1 API -              */
/* ----------------------------------- */

// API holds the gRPC & HTTP Servers, and all their dependencies.
// It has an embedded Service that implements all the API Handlers.
type API struct {
	service.Service // API.
	*cfg.Config     // API Configuration.

	Repository db.Repository // Repository to interact with the database.
	Database   *db.Database  // Database connection.

	// GRPCServer and HTTPGateway are the servers we are going to run.
	GRPCServer  deps.Server
	HTTPGateway deps.Server

	// GRPCInterceptors run before or after gRPC calls.
	// GRPCDialOptions configure the communication between the HTTP Gateway and the gRPC Server.
	GRPCInterceptors []grpc_external.ServerOption
	GRPCDialOptions  []grpc_external.DialOption

	// HTTPMiddleware are the ServeMuxOptions that run before or after HTTP calls.
	// HTTPMiddlewareWrapper are the middleware that wrap around the HTTP Gateway.
	HTTPMiddleware        []runtime.ServeMuxOption
	HTTPMiddlewareWrapper http.MuxWrapperFunc

	// Deps needed to run the API. Validator, Rate Limiter, Logger, TLS Certs, etc.
	*deps.Deps
}

// NewAPI returns a new API with the given configuration.
func NewAPI(config *cfg.Config) *API {
	api := &API{Config: config, Deps: deps.NewDeps()}
	api.Init()
	return api
}

// Init initializes all API components.
func (a *API) Init() {
	a.InitLogger()
	a.InitTLSDeps()
	a.InitPwdHasher()
	a.InitRateLimiter()
	a.InitValidator()
	a.InitAuthenticator()
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
	a.Config = cfg.Load()
}

func (a *API) InitRepository() {
	a.Database = db.NewDatabase(a.DBConfig)
	a.Repository = db.NewRepository(a.Database)
}

func (a *API) InitService() {
	a.Service = service.NewService(a.Repository, a.Authenticator, a.PwdHasher)
}

func (a *API) InitGRPCServer() {
	a.GRPCInterceptors = grpc.AllInterceptors(a.Deps, a.TLSConfig.Enabled)
	a.GRPCDialOptions = grpc.AllDialOptions(a.ClientCreds)
	a.GRPCServer = grpc.NewGRPCServer(a.Config.GRPCPort, a.Service, a.GRPCInterceptors)
	a.GRPCServer.Init()
}

func (a *API) InitHTTPGateway() {
	a.HTTPMiddleware = http.AllMiddleware()
	a.HTTPMiddlewareWrapper = http.AllMiddlewareWrapper(a.Logger)
	a.HTTPGateway = http.NewHTTPGateway(a.MainConfig, a.HTTPMiddleware, a.HTTPMiddlewareWrapper, a.GRPCDialOptions)
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

	// Hash the DB admin password.
	a.DBConfig.AdminPassword = a.PwdHasher.Hash(a.DBConfig.AdminPassword)
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
