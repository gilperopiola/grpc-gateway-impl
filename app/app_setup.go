package app

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/interfaces"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/special_types"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/external"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/external/storage/sqldb"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/servers"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/service"
	"github.com/gilperopiola/grpc-gateway-impl/app/modules"

	"golang.org/x/time/rate"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*      - Core, Modules & Layers -     */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (app App) Setup() *App {
	app.Core.Setup()
	app.Modules.Setup(&app)
	app.Layers.Setup(&app)
	return &app
}

func (c *Core) Setup() {
	c.Config = core.LoadConfig()
	c.Logger = core.SetupLogger(c.Config, core.NewLoggerOptions(c.Config.LevelStackTrace)...)
}

func (m *Modules) Setup(app *App) {
	m.Active = &modules.Active{}
	m.Active.PwdHasher = app.SetupPwdHasher()
	m.Active.RateLimiter = app.SetupRateLimiter()
	m.Active.InputValidator = app.SetupInputValidator()
	m.Active.Authenticator = app.SetupAuthenticator()

	m.Passive = &modules.Passive{}
	m.Passive.TLS = app.SetupTLSModule()
	m.Passive.GRPC = app.SetupGRPCModule()
	m.Passive.HTTP = app.SetupHTTPModule()
}

func (l *Layers) Setup(app *App) {
	l.ExternalLayer = app.SetupExternalLayer()
	l.BusinessLayer = app.SetupBusinessLayer()
	l.ServerLayer = app.SetupServerLayer()
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*             - Modules -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (app App) SetupInputValidator() interfaces.InputValidator {
	app.Modules.InputValidator = modules.NewInputValidator() // -> Input Validator.
	return app.Modules.InputValidator
}

func (app App) SetupAuthenticator() interfaces.TokenAuthenticator {
	app.Modules.Authenticator = modules.NewJWTAuthenticator(app.JWTCfg.Secret, app.JWTCfg.SessionDays) // -> JWT Authenticator.
	return app.Modules.Authenticator
}

func (app App) SetupRateLimiter() *rate.Limiter {
	app.Modules.RateLimiter = rate.NewLimiter(rate.Limit(app.RateLimiterCfg.TokensPerSecond), app.RateLimiterCfg.MaxTokens) // -> Rate Limiter.
	return app.Modules.RateLimiter
}

func (app App) SetupPwdHasher() interfaces.PwdHasher {
	app.Modules.PwdHasher = modules.NewPwdHasher(app.PwdHasherCfg.Salt) // -> Password Hasher.
	return app.Modules.PwdHasher
}

func (app App) SetupGRPCModule() *modules.GRPC {
	app.Modules.GRPC = &modules.GRPC{}
	app.Modules.GRPC.ServerOptions = servers.AllServerOptions(app.Modules.Active, app.TLSCfg.Enabled, app.Modules.ServerCreds) // -> gRPC Server Options.
	app.Modules.GRPC.DialOptions = servers.AllDialOptions(app.ClientCreds)                                                     // -> gRPC Dial Options.
	return app.Modules.GRPC
}

func (app App) SetupHTTPModule() *modules.HTTP {
	app.Modules.HTTP = &modules.HTTP{}
	app.Modules.HTTP.MuxOptionsMiddleware = servers.ServeMuxOpts()      // -> HTTP Middleware.
	app.Modules.HTTP.MuxWrapperMiddleware = servers.MiddlewareWrapper() // -> HTTP Mux Wrapper.
	return app.Modules.HTTP
}

func (app App) SetupTLSModule() *modules.TLS {
	app.Modules.TLS = &modules.TLS{}
	if app.TLSCfg.Enabled {
		app.Modules.TLS.ServerCert = modules.NewTLSCertPool(app.TLSCfg.CertPath)                               // -> Server Certificate.
		app.Modules.TLS.ServerCreds = modules.NewServerTransportCreds(app.TLSCfg.CertPath, app.TLSCfg.KeyPath) // -> Server Credentials.
	}
	app.Modules.TLS.ClientCreds = modules.NewClientTransportCreds(app.TLSCfg.Enabled, app.ServerCert) // -> Client Credentials.
	return app.Modules.TLS
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*              - Layers -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (app App) SetupExternalLayer() external.ExternalLayer {
	app.Layers.ExternalLayer = external.NewExternalLayer(sqldb.NewGormDB(&app.DatabaseCfg)) // -> External Layer.
	return app.Layers.ExternalLayer
}

func (app App) SetupBusinessLayer() interfaces.BusinessLayer {
	app.Layers.BusinessLayer = service.NewService(app.ExternalLayer.GetStorage(), app.Authenticator, app.PwdHasher) // -> Business Layer.
	return app.Layers.BusinessLayer
}

func (app App) SetupServerLayer() special_types.ServerLayer {
	app.Layers.ServerLayer.GRPCServer = servers.NewGRPCServer(app.BusinessLayer, app.GRPC.ServerOptions)                                           // -> gRPC Server.
	app.Layers.ServerLayer.HTTPServer = servers.NewHTTPGateway(app.HTTP.MuxOptionsMiddleware, app.HTTP.MuxWrapperMiddleware, app.GRPC.DialOptions) // -> HTTP Server.
	return app.Layers.ServerLayer
}
