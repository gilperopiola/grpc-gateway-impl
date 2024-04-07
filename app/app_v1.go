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

// App holds our entire application.
type App struct {
	*Config
	*Layers
	*Modules
}

// Config wraps *core.Config and holds the configuration.
type Config struct {
	core.Config
}

// When requests come in, they go through:
type Layers struct {
	ServiceLayer  service.Service
	StorageLayer  storage.Storage
	DatabaseLayer db.Database
}

// Modules wraps *modules.All and holds every independent module.
type Modules struct {
	modules.All
}

// NewAppV1 returns a new App with the given configuration.
func NewAppV1() *App {
	app := &App{&Config{*core.LoadConfig()}, &Layers{}, &Modules{}}
	app.Setup()
	return app
}

// Setup initializes all Layers and Modules.
func (app *App) Setup() {
	app.SetupGlobalLogger()
	app.InitPwdHasher()
	app.InitRateLimiter()
	app.InitInputValidator()
	app.InitTLSModule()
	app.InitAuthenticator()
	app.InitDBAndStorage()
	app.InitService()
	app.InitGRPCModule()
	app.InitHTTPModule()
}

// Run runs the gRPC & HTTP Servers.
func (app *App) Run() {
	app.Server.Run()
	app.Gateway.Run()
}

// WaitForShutdown waits for a SIGINT or SIGTERM to gracefully shutdown the servers.
func (app *App) WaitForShutdown() {
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
