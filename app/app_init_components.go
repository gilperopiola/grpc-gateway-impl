package app

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/modules"
	"github.com/gilperopiola/grpc-gateway-impl/app/servers"
	"github.com/gilperopiola/grpc-gateway-impl/app/service"
	"github.com/gilperopiola/grpc-gateway-impl/app/storage"
	"github.com/gilperopiola/grpc-gateway-impl/app/storage/db"
	"github.com/gilperopiola/grpc-gateway-impl/etc/tools/treyser"

	"golang.org/x/time/rate"
)

/* ----------------------------------- */
/*         - Main Modules -         */
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
func (app *App) SetupGlobalLogger() {
	loggerTrayser := treyser.NewTreyser("logger", 1)
	core.SetupLogger(&app.Config.Config, core.NewLoggerOptions(app.LevelStackTrace)...) // -> Global Logger (zap).
	loggerTrayser.Treys()
}

func (app *App) InitGRPCModule() {
	grpcTrayser := treyser.NewTreyser("grpc", 1)
	app.GRPC.ServerOptions = servers.AllServerOptions(&app.All, app.TLSCfg.Enabled) // -> gRPC Server Options.
	app.GRPC.DialOptions = servers.AllDialOptions(app.ClientCreds)                  // -> gRPC Dial Options.
	app.GRPC.Server = servers.NewGRPCServer(app.ServiceLayer, app.ServerOptions)    // -> gRPC Server.

	app.GRPC.Server.Init()
	grpcTrayser.Treys()
}

func (app *App) InitHTTPModule() {
	httpTrayser := treyser.NewTreyser("http", 1)
	app.HTTP.MuxOptionsMiddleware = servers.ServeMuxOpts()                                                         // -> HTTP Middleware.
	app.HTTP.MuxWrapperMiddleware = servers.MiddlewareWrapper()                                                    // -> HTTP Mux Wrapper.
	app.HTTP.Gateway = servers.NewHTTPGateway(app.MuxOptionsMiddleware, app.MuxWrapperMiddleware, app.DialOptions) // -> HTTP Gateway.

	app.HTTP.Gateway.Init()
	httpTrayser.Treys()
}

func (app *App) InitService() {
	svcTrayser := treyser.NewTreyser("svc", 1)
	app.Layers.ServiceLayer = service.NewService(app.StorageLayer, app.Authenticator, app.PwdHasher) // -> Service.
	svcTrayser.Treys()
}

func (app *App) InitDBAndStorage() {
	dbTrayser := treyser.NewTreyser("db", 1)
	app.DatabaseLayer = db.NewGormDB(&app.DatabaseCfg) // -> DB.
	dbTrayser.Treys()
	repoTrayser := treyser.NewTreyser("storage", 1)
	app.StorageLayer = storage.NewStorage(app.DatabaseLayer) // -> Storage.
	repoTrayser.Treys()
}

/* ----------------------------------- */
/*        - Common Modules -        */
/* ----------------------------------- */

func (app *App) InitInputValidator() {
	inValTrayser := treyser.NewTreyser("inVal", 1)
	app.InputValidator = modules.NewInputValidator() // -> Input Validator.
	inValTrayser.Treys()
}

func (app *App) InitAuthenticator() {
	authTrayser := treyser.NewTreyser("auth", 1)
	app.Authenticator = modules.NewJWTAuthenticator(app.JWTCfg.Secret, app.JWTCfg.SessionDays) // -> JWT Authenticator.
	authTrayser.Treys()
}

func (app *App) InitRateLimiter() {
	rlTrayser := treyser.NewTreyser("rl", 1)
	app.RateLimiter = rate.NewLimiter(rate.Limit(app.RateLimiterCfg.TokensPerSecond), app.RateLimiterCfg.MaxTokens) // -> Rate Limiter.
	rlTrayser.Treys()
}

func (app *App) InitPwdHasher() {
	pwdHasherTrayser := treyser.NewTreyser("pwdHasher", 1)
	app.PwdHasher = modules.NewPwdHasher(app.PwdHasherCfg.Salt) // -> Password Hasher.
	pwdHasherTrayser.Treys()
}

func (app *App) InitTLSModule() {
	tlsTrayser := treyser.NewTreyser("tls", 1)
	if app.TLSCfg.Enabled {
		app.ServerCert = modules.NewTLSCertPool(app.TLSCfg.CertPath)                               // -> Server Certificate.
		app.ServerCreds = modules.NewServerTransportCreds(app.TLSCfg.CertPath, app.TLSCfg.KeyPath) // -> Server Credentials.
	}
	app.ClientCreds = modules.NewClientTransportCreds(app.TLSCfg.Enabled, app.ServerCert) // -> Client Credentials.
	tlsTrayser.Treys()
}
