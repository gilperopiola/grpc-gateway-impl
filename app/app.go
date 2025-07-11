package app

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/clients"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/logs"
	"github.com/gilperopiola/grpc-gateway-impl/app/servers"
	"github.com/gilperopiola/grpc-gateway-impl/app/service"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools"
	"github.com/gilperopiola/grpc-gateway-impl/app/workers"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*               - App -               */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ v1 */

// ⭐️ Holds everything together
type App struct {
	Config  *core.Config
	Servers *servers.Servers
	Service *service.Service
	Clients *clients.Clients
	Tools   *tools.Tools
}

// ⭐️ Sets up a new App - Loads the Config, Logger,
// then the Tools, Service and Servers.
//
// Returns a func to run the Servers and another one to clean
// up used resources before exiting. Remember to defer that one.
//
// If something goes wrong, we just log it and quit.
func Setup() (runAppFunc, cleanUpFunc) {

	app := App{
		Config:  new(core.Config),     // The Heavens
		Servers: new(servers.Servers), // The Earth
		Service: new(service.Service), // The Nobles
		Clients: new(clients.Clients), // The Merchants
		Tools:   new(tools.Tools),     // The Working Class
	}

	initStep(0, func() {
		app.Config = core.LoadConfig()
		logs.SetupLogger(&app.Config.LoggerCfg)
	})

	initStep(1, func() {
		app.Tools = tools.Setup(app.Config)

		var err error
		app.Clients, err = clients.Setup(app.Config, app.Tools)
		logs.LogFatalIfErr(err)

		app.Service = service.Setup(app.Clients, app.Tools)
		app.Servers = servers.Setup(app.Service, app.Tools)
	})

	initStep(2, func() {
		app.Tools.AddCleanupFunc(func() { app.Clients.CloseDB() })
		app.Tools.AddCleanupFunc(app.Servers.Shutdown)
		app.Tools.AddCleanupFuncWithErr(logs.SyncLogger)
	})

	initStep(3, func() {})
	return app.Run, app.Cleanup
}

func initStep(step int, fn func()) {
	logs.InitStep(step)
	fn()
}

func (app *App) Run() {
	workers.RunAll(app.Service)
	//gui.Start(app.Service)
	app.Servers.Run()
}

func (app *App) Cleanup() {
	app.Tools.Cleanup()
}

// Returning these instead of just func() for clarity's sake.
type runAppFunc func()
type cleanUpFunc func()

// ╭───────────────────┬───────────────────┬────────────┬──────────────────────────────────────────────╮
// │ App Field         │ Field Type        │ Interface  │ Contains                                     │
// ├───────────────────┼───────────────────┼────────────┼──────────────────────────────────────────────┤
// │ Configuration     │ *core.Config      │     ~      │ All settings.                                │
// │ GRPC-HTTP Servers │ *servers.Servers  │     ~      │ Our GRPC and HTTP Servers.                   │
// │ Main Service      │ *service.Services │     ~      │ Endpoints and business logic.                │
// │ Clients           │ *clients.Clients  │     ~      │ DBs, APIs, Caches.                           │
// │ Tools             │ *tools.Tools      │ core.Tools │ Specific actions mainly used by the Service. │
// ╰───────────────────┴───────────────────┴────────────┴──────────────────────────────────────────────╯
// * We use a global Logger, so we don't store it anywhere. Access through the logs package.
