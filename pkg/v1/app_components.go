package v1

import (
	"fmt"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/common"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/grpc"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/http"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/db"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"

	"golang.org/x/time/rate"
)

/* ----------------------------------- */
/*         - Main Components -         */
/* ----------------------------------- */

func (a *App) LoadService() {
	a.Service = service.NewService(a.Repository, a.Authenticator, a.PwdHasher) // -> Init Service.
}

func (a *App) LoadRepositoryAndDB() {
	a.Database = db.NewDB(a.DBConfig)                   // -> Init DB.
	a.Repository = repository.NewRepository(a.Database) // -> Init Repository.
}

func (a *App) LoadAllGRPC() {
	a.GRPC.ServerOptions = grpc.AllServerOptions(a.Wrapper, a.TLSConfig.Enabled) // -> Init gRPC Server Options.
	a.GRPC.DialOptions = grpc.AllDialOptions(a.ClientCreds)                      // -> Init gRPC Dial Options.
	a.GRPC.Server = grpc.NewGRPCServer(a.GRPCPort, a.Service, a.ServerOptions)   // -> Init gRPC Server.

	a.GRPC.Server.Init()
}

func (a *App) LoadAllHTTP() {
	a.HTTP.Middleware = http.AllMiddleware()                                                      // -> Init HTTP Middleware.
	a.HTTP.MuxWrapper = http.MuxWrapper()                                                         // -> Init HTTP Mux Wrapper.
	a.HTTP.Gateway = http.NewHTTPGateway(a.MainConfig, a.Middleware, a.MuxWrapper, a.DialOptions) // -> Init HTTP Gateway.

	a.HTTP.Gateway.Init()
}

// InitGlobalLogger sets a new *zap.Logger as global on the zap package.
func (a *App) InitGlobalLogger() {
	common.InitGlobalLogger(a.Config, common.NewLoggerOptions(a.StacktraceLevel)...)
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

func (a *App) LoadPwdHasher() {
	a.PwdHasher = common.NewPwdHasher(a.HashSalt)

	// Hash the DB admin password.
	a.DBConfig.AdminPwd = a.PwdHasher.Hash(a.DBConfig.AdminPwd)
	fmt.Println("Hashed Admin Password:", a.DBConfig.AdminPwd)
}

func (a *App) LoadTLSComponents() {
	if a.TLSConfig.Enabled {
		a.ServerCert = common.NewTLSCertPool(a.TLSConfig.CertPath)
		a.ServerCreds = common.NewServerTransportCreds(a.TLSConfig.CertPath, a.TLSConfig.KeyPath)
	}
	a.ClientCreds = common.NewClientTransportCreds(a.TLSConfig.Enabled, a.ServerCert)
}
