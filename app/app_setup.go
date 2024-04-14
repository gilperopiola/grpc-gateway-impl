package app

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core/interfaces"
	"github.com/gilperopiola/grpc-gateway-impl/app/modules"
	"github.com/gilperopiola/grpc-gateway-impl/app/servers"

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

/* ----------------------------------- */
/*        - Common Modules -        */
/* ----------------------------------- */

func (app App) InitInputValidator() interfaces.InputValidator {
	app.Modules.InputValidator = modules.NewInputValidator() // -> Input Validator.
	return app.Modules.InputValidator
}

func (app App) InitAuthenticator() interfaces.TokenAuthenticator {
	app.Modules.Authenticator = modules.NewJWTAuthenticator(app.JWTCfg.Secret, app.JWTCfg.SessionDays) // -> JWT Authenticator.

	return app.Modules.Authenticator
}

func (app App) InitRateLimiter() *rate.Limiter {
	app.Modules.RateLimiter = rate.NewLimiter(rate.Limit(app.RateLimiterCfg.TokensPerSecond), app.RateLimiterCfg.MaxTokens) // -> Rate Limiter.

	return app.Modules.RateLimiter
}

func (app App) InitPwdHasher() interfaces.PwdHasher {
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

func (app App) InitGRPCModule() *modules.GRPC {
	app.Modules.GRPC = &modules.GRPC{
		ServerOptions: servers.AllServerOptions(app.Modules.All, app.TLSCfg.Enabled), // -> gRPC Server Options.
		DialOptions:   servers.AllDialOptions(app.ClientCreds),                       // -> gRPC Dial Options.
	}
	return app.Modules.GRPC
}

func (app App) InitHTTPModule() *modules.HTTP {
	app.Modules.HTTP = &modules.HTTP{}
	app.Modules.HTTP.MuxOptionsMiddleware = servers.ServeMuxOpts()      // -> HTTP Middleware.
	app.Modules.HTTP.MuxWrapperMiddleware = servers.MiddlewareWrapper() // -> HTTP Mux Wrapper.
	return app.HTTP
}
