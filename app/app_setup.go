package app

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/servers"
	"github.com/gilperopiola/grpc-gateway-impl/app/service"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools/db_tool/sqldb"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - Setup App -            */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (app App) Setup() *App {

	// -> Setup the Core (1/2) -> Config and Logger.
	app.Core.SetupConfigAndLogger()

	// -> Setup the Toolbox -> All Tools: From PwdHasher to TLSTool.
	app.Toolbox.Setup(app.Config)

	// -> Setup the Service -> Not much else.
	app.Service.Setup(app.Toolbox)

	// -> Setup the Core (2/2) -> GRPC & HTTP Servers.
	app.Core.SetupServers(app.Service, app.Toolbox)

	return &app
}

// Step 1: Setup Config and Logger (on Core).
func (c *Core) SetupConfigAndLogger() {
	c.Config = core.SetupConfig()
	c.Logger = core.SetupLogger(&c.LoggerCfg)
}

var _ core.Toolbox = (*Toolbox)(nil)

// Step 2: Setup all Tools in Toolbox.
func (t *Toolbox) Setup(cfg *core.Config) {
	t.PwdHasher = tools.NewPwdHasher(cfg.PwdHasherCfg.Salt)
	t.RateLimiter = tools.NewRateLimiter(&cfg.RLimiterCfg)
	t.RequestsValidator = tools.NewRequestsValidator()
	t.TokenAuthenticator = tools.NewJWTAuthenticator(cfg.JWTCfg.Secret, cfg.JWTCfg.SessionDays)
	t.TLSTool = tools.NewTLSTool(&cfg.TLSCfg)
	t.APICaller = tools.NewAPICaller()
	t.DBTool = sqldb.NewDBTool(sqldb.NewSQLDB(&cfg.DBCfg))
}

// Step 3: Setup Service.
func (s *Service) Setup(toolbox core.Toolbox) {
	s.Service = service.Setup(toolbox)
}

// Step 4: Setup GRPC & HTTP Servers (on Core).
func (c *Core) SetupServers(service core.Service, toolbox core.Toolbox) {
	c.Servers = servers.Setup(service, toolbox, c.TLSCfg.Enabled)
}
