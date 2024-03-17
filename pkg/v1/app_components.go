package v1

import (
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/cfg"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/common"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/grpc"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/http"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/db"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"
	"golang.org/x/time/rate"
)

/* ----------------------------------- */
/*      - App Components Loader -      */
/* ----------------------------------- */

func (a *App) LoadConfig() {
	a.Config = cfg.Load()
}

func (a *App) LoadCommonComponents() {
	a.LoadRateLimiter()
	a.LoadPwdHasher()
	a.LoadLogger()
	a.LoadTLS()
	a.LoadInputValidator()
	a.LoadAuthenticator()
}

func (a *App) LoadService() {
	a.Service = service.NewService(a.Repository, a.Authenticator, a.PwdHasher)
}

func (a *App) LoadRepositoryAndDB() {
	a.Database = db.NewDatabaseWrapper(a.DBConfig)
	a.Repository = repository.NewRepository(a.Database)
}

func (a *App) LoadGRPC() {
	a.GRPC.ServerOptions = grpc.AllServerOptions(a.Wrapper, a.TLSConfig.Enabled)
	a.GRPC.DialOptions = grpc.AllDialOptions(a.ClientCreds)
	a.GRPC.Server = grpc.NewGRPCServer(a.Config.GRPCPort, a.Service, a.ServerOptions)
	a.GRPC.Server.Init()
}

func (a *App) LoadHTTP() {
	a.HTTP.Middleware = http.AllMiddleware()
	a.HTTP.MiddlewareWrapper = http.MiddlewareWrapper(a.Logger)
	a.HTTP.Gateway = http.NewHTTPGateway(a.MainConfig, a.Middleware, a.MiddlewareWrapper, a.DialOptions)
	a.HTTP.Gateway.Init()
}

/* ----------------------------------- */
/*        - Common Components -        */
/* ----------------------------------- */

func (a *App) LoadInputValidator() {
	a.InputValidator = common.NewInputValidator()
}

func (a *App) LoadAuthenticator() {
	a.Authenticator = common.NewJWTAuthenticator(a.JWTConfig.Secret, a.JWTConfig.SessionDays)
}

func (a *App) LoadRateLimiter() {
	a.RateLimiter = rate.NewLimiter(rate.Limit(a.RateLimiterConfig.TokensPerSecond), a.RateLimiterConfig.MaxTokens)
}

func (a *App) LoadLogger() {
	a.Logger = common.NewLogger(a.IsProd, common.NewLoggerOptions()...)
}

func (a *App) LoadPwdHasher() {
	a.PwdHasher = common.NewPwdHasher(a.HashSalt)

	// Hash the DB admin password.
	a.DBConfig.AdminPassword = a.PwdHasher.Hash(a.DBConfig.AdminPassword)
}

func (a *App) LoadTLS() {
	if a.TLSConfig.Enabled {
		a.ServerCert = common.NewTLSCertPool(a.TLSConfig.CertPath)
		a.ServerCreds = common.NewServerTransportCreds(a.TLSConfig.CertPath, a.TLSConfig.KeyPath)
	}
	a.ClientCreds = common.NewClientTransportCreds(a.TLSConfig.Enabled, a.ServerCert)
}
