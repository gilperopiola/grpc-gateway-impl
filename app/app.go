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

// ╭───────────────────┬───────────────────┬────────────┬───────────────────────────────────────╮
// │ Field             │ Type              │ Interface  │ Contains                              │
// ├───────────────────┼───────────────────┼────────────┼───────────────────────────────────────┤
// │ Configuration     │ *core.Config      │            │ All settings, split by module.        │
// │ GRPC-HTTP Servers │ *servers.Servers  │            │ Our GRPC and HTTP Servers.            │
// │ Main Service      │ *service.Services │            │ Endpoints and business logic.         │
// │ Tools             │ *tools.Tools      │ core.Tools │ Specific actions used by our Service. │
// ╰───────────────────┴───────────────────┴────────────┴───────────────────────────────────────╯
// - We use a global Logger, so we don't store it anywhere.

// ⭐️ Our main App.
//
// Holds everything - does nothing.
type App struct {
	Config  *core.Config
	Servers *servers.Servers
	Service *service.Service
	Tools   *tools.Tools
}

// Called by main.go.
//
// Initializes a new App: Loads the Config, Logger,
// then the Tools, Service and Servers.
//
// Returns a func to run the Servers and another one to free
// used resources before exiting.
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
