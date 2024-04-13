package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/modules"
	"github.com/gilperopiola/grpc-gateway-impl/app/service"
	"github.com/gilperopiola/grpc-gateway-impl/app/storage"
	"github.com/gilperopiola/grpc-gateway-impl/app/storage/db"

	"go.uber.org/zap"
)

/* ----------------------------------- */
/*             - App v1 -              */
/* ----------------------------------- */

type App struct {
	Config
	Globals
	Modules
	Layers
}

type Config struct {
	*core.Config
}

type Globals struct {
	*zap.Logger
}

type Modules struct {
	*modules.All
}

type Layers struct {
	ServiceLayer  service.Service
	StorageLayer  storage.Storage
	DatabaseLayer db.Database
}

func NewApp() *App {
	return App{}.Setup()
}

func (app App) Setup() *App {

	app.Config.Setup()

	app.Globals.Setup(app.Config.Config)

	app.Modules.Setup1(&app)

	app.Layers.Setup(app.Authenticator, app.PwdHasher, &app.Config.DatabaseCfg)

	app.Modules.Setup2(&app)

	return &app
}

func (c *Config) Setup() {
	c.Config = core.LoadConfig()
}

func (g *Globals) Setup(cfg *core.Config) {
	g.Logger = core.SetupLogger(cfg, core.NewLoggerOptions(cfg.LevelStackTrace)...)

}

func (l *Layers) Setup(tokenAuth modules.TokenAuthenticator, pwdHasher modules.PwdHasher, dbCfg *core.DatabaseCfg) {
	l.DatabaseLayer = db.NewGormDB(dbCfg)
	l.StorageLayer = storage.NewStorage(l.DatabaseLayer)
	l.ServiceLayer = service.NewService(l.StorageLayer, tokenAuth, pwdHasher)
}

func (m *Modules) Setup1(app *App) {
	m.All = &modules.All{}
	m.All.PwdHasher = app.InitPwdHasher()
	m.All.RateLimiter = app.InitRateLimiter()
	m.All.InputValidator = app.InitInputValidator()
	m.All.TLS = app.InitTLSModule()
	m.All.Authenticator = app.InitAuthenticator()
}

func (m *Modules) Setup2(app *App) {
	m.All.GRPC = app.InitGRPCModule()
	m.All.HTTP = app.InitHTTPModule()
}

func (app *App) Run() {
	app.Server.Run()
	app.Gateway.Run()
}

func (app *App) WaitForShutdown() { // Waits for a SIGINT or SIGTERM to gracefully shutdown the servers
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	zap.S().Infoln("Shutting down servers...")

	sqlDB := app.DatabaseLayer.GetSQL()
	if sqlDB != nil {
		sqlDB.Close()
	}
	app.Server.Shutdown()
	app.Gateway.Shutdown()

	zap.S().Infoln("Servers stopped! Bye bye~")
	zap.L().Sync()
}
