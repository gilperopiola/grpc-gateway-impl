package app

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/modules"
	"github.com/gilperopiola/grpc-gateway-impl/app/servers"
	"github.com/gilperopiola/grpc-gateway-impl/app/service"
	"github.com/gilperopiola/grpc-gateway-impl/app/storage"
	"github.com/gilperopiola/grpc-gateway-impl/app/storage/db"
	"go.uber.org/zap"

	"golang.org/x/time/rate"
)

/* ----------------------------------- */
/*           - Main Modules -          */
/* ----------------------------------- */

//type (
//	setupConfigFn func() *core.Config
//	setupLoggerFn func() // logger is global
//)
//
//func (app *App) SetupCore(setupConfig setupConfigFn, setupLogger setupLoggerFn) func(*App) {
//	return func(app *App) {
//		configTrayser := treyser.NewTreyser("config", 1)
//		app.Config = &Config{setupConfig()}
//		setupLogger()
//		configTrayser.Treys()
//	}
//}

// SetupGlobalLogger sets a new *zap.Logger as global on the zap package.
func (app App) SetupGlobalLogger() *zap.Logger {
	return core.SetupLogger(app.Config.Config, core.NewLoggerOptions(app.Config.LevelStackTrace)...) // -> Global Logger (zap).
}

func (app App) InitGRPCModule() *modules.GRPC {
	app.Modules.GRPC = &modules.GRPC{}

	app.GRPC.ServerOptions = servers.AllServerOptions(app.Modules.All, app.TLSCfg.Enabled) // -> gRPC Server Options.
	app.GRPC.DialOptions = servers.AllDialOptions(app.ClientCreds)                         // -> gRPC Dial Options.
	app.GRPC.Server = servers.NewGRPCServer(app.ServiceLayer, app.ServerOptions)           // -> gRPC Server.
	app.GRPC.Server.Init()

	return app.GRPC
}

func (app App) InitHTTPModule() *modules.HTTP {
	app.Modules.HTTP = &modules.HTTP{}

	app.HTTP.MuxOptionsMiddleware = servers.ServeMuxOpts()                                                         // -> HTTP Middleware.
	app.HTTP.MuxWrapperMiddleware = servers.MiddlewareWrapper()                                                    // -> HTTP Mux Wrapper.
	app.HTTP.Gateway = servers.NewHTTPGateway(app.MuxOptionsMiddleware, app.MuxWrapperMiddleware, app.DialOptions) // -> HTTP Gateway.

	app.HTTP.Gateway.Init()

	return app.HTTP
}

func (app App) InitService() service.Service {
	app.Layers.ServiceLayer = service.NewService(app.StorageLayer, app.Authenticator, app.PwdHasher) // -> Service.

	return app.Layers.ServiceLayer
}
func (app App) InitStorage(dbLayer db.Database) storage.Storage {
	app.Layers.StorageLayer = storage.NewStorage(app.DatabaseLayer) // -> Storage.

	return app.Layers.StorageLayer
}

func (app App) InitDatabase() db.Database {
	app.Layers.StorageLayer = storage.NewStorage(app.DatabaseLayer) // -> Storage.

	return app.Layers.DatabaseLayer
}

/* ----------------------------------- */
/*        - Common Modules -        */
/* ----------------------------------- */

func (app App) InitInputValidator() modules.InputValidator {
	app.Modules.InputValidator = modules.NewInputValidator() // -> Input Validator.
	return app.Modules.InputValidator
}

func (app App) InitAuthenticator() modules.TokenAuthenticator {
	app.Modules.Authenticator = modules.NewJWTAuthenticator(app.JWTCfg.Secret, app.JWTCfg.SessionDays) // -> JWT Authenticator.

	return app.Modules.Authenticator
}

func (app App) InitRateLimiter() *rate.Limiter {
	app.Modules.RateLimiter = rate.NewLimiter(rate.Limit(app.RateLimiterCfg.TokensPerSecond), app.RateLimiterCfg.MaxTokens) // -> Rate Limiter.

	return app.Modules.RateLimiter
}

func (app App) InitPwdHasher() modules.PwdHasher {
	app.Modules.PwdHasher = modules.NewPwdHasher(app.PwdHasherCfg.Salt) // -> Password Hasher.

	return app.Modules.PwdHasher
}

func (app App) InitTLSModule() *modules.TLS {
	app.Modules.TLS = &modules.TLS{}

	if app.TLSCfg.Enabled {
		app.Modules.TLS.ServerCert = modules.NewTLSCertPool(app.TLSCfg.CertPath)                               // -> Server Certificate.
		app.Modules.TLS.ServerCreds = modules.NewServerTransportCreds(app.TLSCfg.CertPath, app.TLSCfg.KeyPath) // -> Server Credentials.
	}

	app.Modules.TLS.ClientCreds = modules.NewClientTransportCreds(app.TLSCfg.Enabled, app.ServerCert) // -> Client Credentials.

	return app.Modules.TLS
}
