package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/interfaces"
	"github.com/gilperopiola/grpc-gateway-impl/app/external"
	"github.com/gilperopiola/grpc-gateway-impl/app/external/storage/sqldb"
	"github.com/gilperopiola/grpc-gateway-impl/app/modules"
	"github.com/gilperopiola/grpc-gateway-impl/app/servers"
	"github.com/gilperopiola/grpc-gateway-impl/app/service"

	"go.uber.org/zap"
)

/* ----------------------------------- */
/*            - App (v1) -             */
/* ----------------------------------- */

type App struct {
	Config
	Layers
	Modules
}

type Config struct {
	*core.CoreCfg
}

type Modules struct {
	*modules.All
}

type Layers struct {
	interfaces.ServerLayer
	interfaces.BusinessLayer
	external.ExternalLayer
}

// New App.

func NewApp() *App {
	return App{}.Setup()
}

// Setup App.

func (app App) Setup() *App {
	app.Config.Setup()
	core.SetupLogger(app.CoreCfg, core.NewLoggerOptions(app.LevelStackTrace)...)
	app.Modules.Setup(&app)
	app.Layers.Setup(&app)
	return &app
}

func (c *Config) Setup() {
	c.CoreCfg = core.LoadConfig()
}

func (l *Layers) Setup(app *App) {
	l.ExternalLayer = external.NewExternalLayer(InitDatabase(&app.DatabaseCfg))
	l.BusinessLayer = service.NewService(l.ExternalLayer.GetStorage(), app.Authenticator, app.PwdHasher)
	l.ServerLayer.GRPCServer = app.InitGRPCServer()
	l.ServerLayer.HTTPServer = app.InitHTTPServer()
}

func (m *Modules) Setup(app *App) {
	m.All = &modules.All{}
	m.All.PwdHasher = app.InitPwdHasher()
	m.All.RateLimiter = app.InitRateLimiter()
	m.All.InputValidator = app.InitInputValidator()
	m.All.TLS = app.InitTLSModule()
	m.All.Authenticator = app.InitAuthenticator()
	m.SetupServerModules(app)
}

func (m *Modules) SetupServerModules(app *App) {
	m.All.GRPC = app.InitGRPCModule()
	m.All.HTTP = app.InitHTTPModule()
}

// SetupLogger returns a new *zap.Logger which also lives globally in the zap package.
func (app App) SetupLogger() *zap.Logger {
	return core.SetupLogger(app.CoreCfg, core.NewLoggerOptions(app.Config.LevelStackTrace)...) // -> Global Logger (zap).
}

func (app App) InitGRPCServer() *servers.GRPCServer {
	app.Layers.ServerLayer.GRPCServer = servers.NewGRPCServer(app.BusinessLayer, app.GRPC.ServerOptions) // -> gRPC Server.
	app.Layers.ServerLayer.GRPCServer.Init()
	return app.Layers.ServerLayer.GRPCServer.(*servers.GRPCServer)
}

func (app App) InitHTTPServer() *servers.HTTPGateway {
	app.Layers.ServerLayer.HTTPServer = servers.NewHTTPGateway(app.HTTP.MuxOptionsMiddleware, app.HTTP.MuxWrapperMiddleware, app.GRPC.DialOptions) // -> HTTP Gateway.
	app.Layers.ServerLayer.HTTPServer.Init()
	return app.Layers.ServerLayer.HTTPServer.(*servers.HTTPGateway)
}

func InitDatabase(dbCfg *core.DatabaseCfg) sqldb.Database {
	return sqldb.NewGormDB(dbCfg) // -> DB.
}

// Run app.

func (app *App) Run() {
	app.ServerLayer.GRPCServer.Run()
	app.ServerLayer.HTTPServer.Run()
}

// Quit app.

func (app *App) WaitForShutdown() { // Waits for a SIGINT or SIGTERM to gracefully shutdown the servers
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	zap.S().Infoln("Shutting down servers...")

	sqlDB := app.ExternalLayer.GetDB().GetSQL()
	if sqlDB != nil {
		sqlDB.Close()
	}
	app.Layers.ServerLayer.GRPCServer.Shutdown()
	app.Layers.ServerLayer.HTTPServer.Shutdown()

	zap.S().Infoln("Servers stopped! Bye bye~")
	zap.L().Sync()
}
