package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/external"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/service"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools"

	"go.uber.org/zap"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - App (v1) -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func NewApp() *App {
	fmt.Println() // Better logs readability.

	return App{
		&Core{},
		&Layers{},
		&tools.Tools{},
	}.Setup()
}

type (
	App struct {
		*Core        // -> Config, Logger and Servers.
		*Layers      // -> Service Layer and External Layer.
		*tools.Tools // -> JWT Auth, Reqs Validator, Rate Limiter, TLS, etc.
	}
	Core struct {
		*core.Config  // -> Config.
		*zap.Logger   // -> Logger (also lives globally in zap.L and zap.S).
		*core.Servers // -> gRPC and HTTP Servers.
	}
	Layers struct {
		*service.ServiceLayer   // -> Service, all business logic.
		*external.ExternalLayer // -> Storage (DB, Cache, etc) and Clients (gRPC, HTTP, etc).
	}
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/* - Setup App: Core, Tools & Layers - */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (app App) Setup() *App {

	// 1. -> Setup Core -> Config and Logger.
	app.Core.SetupConfigAndLogger()

	// 2. -> Setup Tools -> All Tools on /tools pkg.
	app.Tools.Setup(app.Config)

	// 3. -> Setup Layers -> Service and External.
	//  * Service Layer = Business logic.
	//  * External Layer = Storage and Clients.
	app.Layers.Setup(&app.DBCfg, app.Tools)

	// 4. -> Setup Core -> gRPC & HTTP Servers.
	app.Core.SetupServers(app.Layers.ServiceLayer, app.Tools)

	return &app
}

// Step 1: Setup Config and Logger (on Core).
func (c *Core) SetupConfigAndLogger() {
	c.Config = core.LoadConfig()
	c.Logger = core.SetupLogger(&c.LoggerCfg)
}

// Step 2: Setup Tools (on /tools pkg). It's the only Setup func that doesn't live here.

// Step 3: Setup Layers.
func (l *Layers) Setup(dbCfg *core.DBCfg, tools core.Toolbox) {
	l.ExternalLayer = external.SetupLayer(dbCfg)
	l.ServiceLayer = service.SetupLayer(l.ExternalLayer, tools.GetAuthenticator(), tools.GetPwdHasher())
}

// Step 4: Setup gRPC & HTTP Servers (on Core).
func (c *Core) SetupServers(serviceLayer *service.ServiceLayer, tools core.Toolbox) {
	c.Servers = core.SetupServers(serviceLayer, tools, c.TLSCfg.Enabled)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*       - Run / Shutdown App -        */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (app *App) Run() {
	app.Servers.Run()

	go func() {
		time.Sleep(1 * time.Second) // T0D0 healtcheck??
		zap.S().Info("GRPC Gateway Implementation | App OK 🚀")
	}()
}

func (app *App) WaitForShutdown() { // Waits for a SIGINT or SIGTERM to gracefully shutdown the servers.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	zap.S().Infoln("Shutting down servers...")

	sqlDB := app.ExternalLayer.Storage.DB.GetSQL()
	if sqlDB != nil {
		sqlDB.Close()
	}

	app.Servers.Shutdown()
	zap.S().Infoln("Servers stopped! Bye bye~")
	zap.L().Sync()
}
