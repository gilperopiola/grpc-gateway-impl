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
	return App{&Core{}, &Layers{}, &tools.Tools{}}.Setup()
}

type (
	App struct {
		*Core        // -> Servers, Config and Logger.
		*Layers      // -> Business / External Layers.
		*tools.Tools // -> JWT Auth, Input Validator, Rate Limiter, TLS, gRPC Interceptors, HTTP Middleware, etc.
	}

	Core struct {
		*core.Servers // -> gRPC and HTTP Servers.
		*core.Config  // -> Config.
		*zap.Logger   // -> Logger (also lives globally in zap.L() and zap.S()).
	}

	Layers struct {
		business.Service  // -> Service, all business logic.
		external.External // -> Storage (DB, Cache, etc) and Clients (gRPC, HTTP, etc).
	}
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/* - Setup App: Core, Tools & Layers - */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (app App) Setup() *App {

	// 1. Setup Config and Logger on the Core struct.
	app.Core.SetupConfigAndLogger()

	// 2. Setup Tools struct. Tools live in their own package.
	app.Tools.Setup(app.Config)

	// 3. Setup Layers struct.
	//  * Business Layer = Service.
	//  * External Layer = DBs and such.
	app.Layers.Setup(app.Tools, &app.DatabaseCfg)

	// 4. Setup gRPC & HTTP Servers on the Core struct.
	app.Core.SetupServers(app.Tools, app.Layers.Service)

	return &app
}

func (c *Core) SetupConfigAndLogger() {
	c.Config = core.LoadConfig()
	c.Logger = core.SetupLogger(c.Config, core.SetupLoggerOptions(c.Config.LevelStackTrace)...)
}

func (l *Layers) Setup(tools core.ToolsAccessor, dbCfg *core.DatabaseCfg) {
	l.External = external.NewExternalLayer(sql.NewGormDB(dbCfg))
	l.Service = business.NewService(l.External.GetStorage(), tools.GetAuthenticator(), tools.GetPwdHasher())
}

func (c *Core) SetupServers(tools core.ToolsAccessor, businessLayer business.Service) {
	c.Servers = core.SetupServers(tools, businessLayer, c.TLSCfg.Enabled)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*       - Run & Shutdown App -        */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (app *App) Run() {
	app.Servers.Run()

	go func() {
		time.Sleep(1 * time.Second)
		zap.S().Info("Servers OK")
	}()
}

func (app *App) WaitForShutdown() { // Waits for a SIGINT or SIGTERM to gracefully shutdown the servers.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	zap.S().Infoln("Shutting down servers...")

	sql := app.External.GetDB().GetSQL()
	if sql != nil {
		sql.Close()
	}

	app.Servers.Shutdown()
	zap.S().Infoln("Servers stopped! Bye bye~")
	zap.L().Sync()
}
