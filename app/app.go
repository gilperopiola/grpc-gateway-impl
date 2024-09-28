package app

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/clients"
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

// ⭐️ Holds everything together
type App struct {
	Config  *core.Config
	Servers *servers.Servers
	Service *service.Service
	Clients *clients.Clients
	Tools   *tools.Tools
}

// ╭───────────────────┬───────────────────┬────────────┬──────────────────────────────────────────────╮
// │ App Field         │ Field Type        │ Interface  │ Contains                                     │
// ├───────────────────┼───────────────────┼────────────┼──────────────────────────────────────────────┤
// │ Configuration     │ *core.Config      │     ~      │ All settings.                                │
// │ GRPC-HTTP Servers │ *servers.Servers  │     ~      │ Our GRPC and HTTP Servers.                   │
// │ Main Service      │ *service.Services │     ~      │ Endpoints and business logic.                │
// │ Clients           │ *clients.Clients  │     ~      │ DBs, APIs, Caches.                           │
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
		Config:  new(core.Config),     // The Heavens
		Servers: new(servers.Servers), // The Earth
		Service: new(service.Service), // The Factories
		Clients: new(clients.Clients), // The Travellers
		Tools:   new(tools.Tools),     // The Working Class
	}

	logs.Step(0, "Starting up")
	func() {
		app.Config = core.LoadConfig()
		logs.EnvVars()
		logs.SetupLogger(&app.Config.LoggerCfg)
	}()

	logs.Step(1, "Setup")
	func() {
		app.Tools = tools.Setup(app.Config)
		logs.SubstepOK("Tools", "🛠️ ")

		app.Clients = clients.Setup(app.Config)
		logs.SubstepOK("Clients", "🔱")

		app.Service = service.Setup(app.Clients, app.Tools)
		logs.SubstepOK("Service", "⚡")

		app.Servers = servers.Setup(app.Service, app.Tools)
		logs.SubstepOK("Servers", "📡")
	}()

	func() {
		app.Tools.AddCleanupFunc(app.Clients.CloseDB)
		app.Tools.AddCleanupFunc(app.Servers.Shutdown)
		app.Tools.AddCleanupFuncWithErr(utils.SyncLogger)
		logs.SubstepOK("Cleanup", "🧽")
	}()

	return app.Servers.Run, app.Tools.Cleanup
}

// Returning these instead of just func() for clarity's sake.
type runAppFunc func()
type cleanUpFunc func()
