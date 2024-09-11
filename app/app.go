package app

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/logs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/utils"
	"github.com/gilperopiola/grpc-gateway-impl/app/servers"
	"github.com/gilperopiola/grpc-gateway-impl/app/service"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*               - App -               */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ v1 */

// ⭐️ Holds everything, but does nothing on its own.
type App struct {
	Config  *core.Config
	Servers *servers.Servers
	Service *service.Service
	Tools   *tools.Tools
}

// ╭───────────────────┬───────────────────┬────────────┬──────────────────────────────────────────────╮
// │ App Field         │ Field Type        │ Interface  │ Contains                                     │
// ├───────────────────┼───────────────────┼────────────┼──────────────────────────────────────────────┤
// │ Configuration     │ *core.Config      │     ~      │ All settings, split by module.               │
// │ GRPC-HTTP Servers │ *servers.Servers  │     ~      │ Our GRPC and HTTP Servers.                   │
// │ Main Service      │ *service.Services │     ~      │ Endpoints and business logic.                │
// │ Tools             │ *tools.Tools      │ core.Tools │ Specific actions mainly used by the Service. │
// ╰───────────────────┴───────────────────┴────────────┴──────────────────────────────────────────────╯
// * We use a global Logger, so we don't store it anywhere.

// ⭐️ Sets up a new App - Loads the Config, Logger,
// then the Tools, Service and Servers.
//
// Returns a func to run the Servers and another one to clean
// up used resources before exiting. Remember to defer that one.
//
// If something goes wrong, we just log it and quit.
func Setup() (runAppFunc, cleanUpFunc) {

	app := App{
		Config:  new(core.Config),     // The Heavens.
		Servers: new(servers.Servers), // The Earth.
		Service: new(service.Service), // The Bourgeoisie.
		Tools:   new(tools.Tools),     // The Proletariat.
	}

	func() {
		app.Config = core.LoadConfig()
		logs.SetupLogger(&app.Config.LoggerCfg)
	}()

	func() {
		app.Tools = tools.Setup(app.Config)
		app.Service = service.Setup(app.Tools)
		app.Servers = servers.Setup(app.Service, app.Tools)
	}()

	func() {
		app.Tools.AddCleanupFunc(app.Tools.CloseDB)
		app.Tools.AddCleanupFunc(app.Servers.Shutdown)
		app.Tools.AddCleanupFuncWithErr(utils.SyncLogger)
	}()

	return app.Servers.Run, app.Tools.Cleanup
}

// Returning these instead of just func() for clarity's sake.
type runAppFunc func()
type cleanUpFunc func()
