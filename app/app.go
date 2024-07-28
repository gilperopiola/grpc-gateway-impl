package app

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/servers"
	"github.com/gilperopiola/grpc-gateway-impl/app/service"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools"

	"go.uber.org/zap"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*               - App -               */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ v1 */

// -> ⭐️ This is our core App.
// -> It is divided into 4: Configuration - Servers - Service - Tools.
// ╭───────────────────┬───────────────────┬──────────────┬────────────────────────────────╮
// │ App's Module      │ Struct            │ Interface    │ Contains                       │
// ├───────────────────┼───────────────────┼──────────────┼────────────────────────────────┤
// │ Configuration     │ *core.Config      │              │ core.DBCfg, core.TLSCfg, co... │
// │ GRPC-HTTP Servers │ *servers.Servers  │              │ *grpc.Server, *http.Server     │
// │ Main Service      │ *service.Services │              │ pbs.AuthServiceServer, pbs.... │
// │ Tools             │ *tools.Tools  │ core.Tools │ core.DBTool, core.TLSTool, ... │
// ╰───────────────────┴───────────────────┴──────────────┴────────────────────────────────╯
type App struct {
	*core.Config
	*servers.Servers
	*service.Services
	*tools.Tools
}

// This will be called by main.go on init.
func NewApp() (runAppFunc, cleanUpFunc) {

	app := &App{
		Config:   &core.Config{},      // 🗺️
		Servers:  &servers.Servers{},  // 🌐
		Services: &service.Services{}, // 🌟
		Tools:    &tools.Tools{},      // 🛠️
	}

	func() {
		app.Config = core.LoadConfig()
		core.SetupLogger(&app.LoggerCfg)
	}()

	func() {
		app.Tools = tools.Setup(app.Config, app.Services.AnswerGroupInvite)
		app.Services = service.Setup(app.Tools)
		app.Servers = servers.Setup(app.Services, app.Tools)
	}()

	func() {
		app.Tools.AddCleanupFunc(app.CloseDB)
		app.Tools.AddCleanupFunc(app.Servers.Shutdown)
		app.Tools.AddCleanupFuncWithErr(zap.L().Sync)
	}()

	return app.Servers.Run, app.Tools.Cleanup
}

// NewApp returns a runAppFunc and a cleanUpFunc - so the caller can first run
// the Servers and then release gracefully all used resources when it's done.
type runAppFunc func()
type cleanUpFunc func()
