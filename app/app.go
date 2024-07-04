package app

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/servers"
	"github.com/gilperopiola/grpc-gateway-impl/app/service"
	"github.com/gilperopiola/grpc-gateway-impl/app/toolbox"

	"go.uber.org/zap"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*               - App -               */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ v1 */

// -> ⭐️ This is our core App. This holds everything.
type App struct {
	*core.Config
	*servers.Servers
	*service.Service
	*toolbox.Toolbox
}

var _ core.Servers = &servers.Servers{}
var _ core.Service = &service.Service{}
var _ core.Toolbox = &toolbox.Toolbox{}

// This will be called by main.go on init.
func NewApp() (runAppFunc, cleanUpFunc) {

	app := &App{
		Config:  &core.Config{},     // 🗺️
		Servers: &servers.Servers{}, // 🌐
		Service: &service.Service{}, // 🌟
		Toolbox: &toolbox.Toolbox{}, // 🛠️
	}

	func() {
		app.Config = core.LoadConfig()
		core.SetupLogger(&app.LoggerCfg)
	}()

	func() {
		app.Toolbox = toolbox.Setup(app.Config, app.Service.AnswerGroupInvite)
		app.Service = service.Setup(app.Toolbox)
		app.Servers = servers.Setup(app.Service, app.Toolbox)
	}()

	func() {
		app.Toolbox.AddCleanupFunc(app.CloseDB)
		app.Toolbox.AddCleanupFunc(app.Servers.Shutdown)
		app.Toolbox.AddCleanupFuncWithErr(zap.L().Sync)
	}()

	return app.Servers.Run, app.Toolbox.Cleanup
}

// NewApp returns a runAppFunc and a cleanUpFunc - so the caller can first run
// the Servers and then release gracefully all used resources when it's done.
type runAppFunc func()
type cleanUpFunc func()
