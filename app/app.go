package app

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/business"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/external"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/external/storage/sql"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools"

	"go.uber.org/zap"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - App (v1) -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func NewApp() *App {
	return App{
		&Core{},
		&Layers{},
		&tools.Tools{},
	}.Setup()
}

type (
	App struct {
		*Core        // -> Servers, Config and Logger.
		*Layers      // -> Business / External Layers.
		*tools.Tools // -> JWT Auth, Reqs Validator, Rate Limiter, TLS, etc.
	}

	Core struct {
		*core.Servers // -> gRPC and HTTP Servers.
		*core.Config  // -> Config.
		*zap.Logger   // -> Logger (also lives globally in zap.L() and zap.S()).
	}

	Layers struct {
		*business.ServiceLayer  // -> Service, all business logic.
		*external.ExternalLayer // -> Storage (DB, Cache, etc) and Clients (gRPC, HTTP, etc).
	}
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/* - Setup App: Core, Tools & Layers - */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (app App) Setup() *App {

	// 1. Setup Config and Logger on the Core struct.
	app.Core.SetupConfigAndLogger()

	// 2. Setup Tools struct. This goes through the tools pkg.
	app.Tools.Setup(app.Config)

	// 3. Setup Layers struct.
	//  * Service Layer = Business logic.
	//  * External Layer = Storage & Clients.
	app.Layers.Setup(app.Tools, &app.DatabaseCfg)

	// 4. Setup gRPC & HTTP Servers on the Core struct.
	app.Core.SetupServers(app.Tools, app.Layers.ServiceLayer)

	return &app
}

// Step 1: Setup Config and Logger (on Core).
func (c *Core) SetupConfigAndLogger() {
	c.Config = core.LoadConfig()
	c.Logger = core.SetupLogger(c.Config, core.SetupLoggerOptions(c.Config.LevelStackTrace)...)
}

// Step 2: Setup Tools (on /tools pkg). It's the only Setup func that doesn't live here.

// Step 3: Setup Layers.
func (l *Layers) Setup(tools core.ToolsAccessor, dbCfg *core.DatabaseCfg) {
	l.ExternalLayer = external.SetupLayer(sql.NewGormDB(dbCfg))
	l.ServiceLayer = business.NewService(l.ExternalLayer.StorageLayer, tools.GetAuthenticator(), tools.GetPwdHasher())
}

// Step 4: Setup gRPC & HTTP Servers (on Core).
func (c *Core) SetupServers(tools core.ToolsAccessor, serviceLayer *business.ServiceLayer) {
	c.Servers = core.SetupServers(tools, serviceLayer, c.TLSCfg.Enabled)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*       - Run / Shutdown App -        */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (app *App) Run() {
	app.Servers.Run()

	go func() {
		time.Sleep(1 * time.Second)
		zap.S().Info("Servers OK")
	}()
}

func (app *App) WaitForShutdown() { // Waits for a SIGINT or SIGTERM to gracefully shutdown the servers.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	zap.S().Infoln("Shutting down servers...")

	sqlDB := app.ExternalLayer.DB.GetSQL()
	if sqlDB != nil {
		sqlDB.Close()
	}

	app.Servers.Shutdown()
	zap.S().Infoln("Servers stopped! Bye bye~")
	zap.L().Sync()
}
