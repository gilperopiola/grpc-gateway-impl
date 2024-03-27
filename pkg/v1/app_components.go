package v1

import (
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

func (a *App) InitGRPCComponents() {
	a.GRPC.ServerOptions = grpc.AllServerOptions(a.Components, a.TLSCfg.Enabled) // -> gRPC Server Options.
	a.GRPC.DialOptions = grpc.AllDialOptions(a.ClientCreds)                      // -> gRPC Dial Options.
	a.GRPC.Server = grpc.NewGRPCServer(a.GRPCPort, a.Service, a.ServerOptions)   // -> gRPC Server.

	a.GRPC.Server.Init()
}

func (a *App) InitHTTPComponents() {
	a.HTTP.Middleware = http.AllMiddleware()                                                    // -> HTTP Middleware.
	a.HTTP.MuxWrapper = http.MuxWrapper()                                                       // -> HTTP Mux Wrapper.
	a.HTTP.Gateway = http.NewHTTPGateway(&a.MainCfg, a.Middleware, a.MuxWrapper, a.DialOptions) // -> HTTP Gateway.

	a.HTTP.Gateway.Init()
}

func (a *App) InitService() {
	a.Service = service.NewService(a.Repository, a.Authenticator, a.PwdHasher) // -> Service.
}

func (a *App) InitDBAndRepository() {
	a.Database = db.NewDB(&a.DBCfg)                     // -> DB.
	a.Repository = repository.NewRepository(a.Database) // -> Repository.
}

// InitGlobalLogger sets a new *zap.Logger as global on the zap package.
func (a *App) InitGlobalLogger() {
	common.InitGlobalLogger(a.Config, common.NewLoggerOptions(a.StacktraceLevel)...) // -> Global Logger (zap).
}

/* ----------------------------------- */
/*        - Common Components -        */
/* ----------------------------------- */

func (a *App) InitInputValidator() {
	a.InputValidator = common.NewInputValidator() // -> Input Validator.
}

func (a *App) InitAuthenticator() {
	a.Authenticator = common.NewJWTAuthenticator(a.JWTCfg.Secret, a.JWTCfg.SessionDays) // -> JWT Authenticator.
}

func (a *App) InitRateLimiter() {
	a.RateLimiter = rate.NewLimiter(rate.Limit(a.RLimiterCfg.TokensPerSecond), a.RLimiterCfg.MaxTokens) // -> Rate Limiter.
}

func (a *App) InitPwdHasher() {
	a.PwdHasher = common.NewPwdHasher(a.HashSalt) // -> Password Hasher.
}

func (a *App) InitTLSComponents() {
	if a.TLSCfg.Enabled {
		a.ServerCert = common.NewTLSCertPool(a.TLSCfg.CertPath)                             // -> Server Certificate.
		a.ServerCreds = common.NewServerTransportCreds(a.TLSCfg.CertPath, a.TLSCfg.KeyPath) // -> Server Credentials.
	}
	a.ClientCreds = common.NewClientTransportCreds(a.TLSCfg.Enabled, a.ServerCert) // -> Client Credentials.
}
