package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/servers"
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

	return App{&Core{}, &Layers{}, &Tools{}}.Setup()
}

type (
	App struct {
		*Core   // -> Servers, Config, Logger.
		*Layers // -> Service Layer and External Layer.
		*Tools  // -> JWT Auth, Reqs Validator, Rate Limiter, TLS, etc.
	}

	Core struct {
		*servers.Servers // -> GRPC and HTTP Servers.
		*core.Config     // -> Config.
		*zap.Logger      // -> Logger (also lives globally in zap.L and zap.S).
	}

	Layers struct {
		core.ServiceLayer  // -> Service, all business logic.
		core.ExternalLayer // -> Storage (DB, Cache, etc) and Clients (GRPC, HTTP, etc).
	}

	Tools struct {
		core.TokenAuthenticator // Generates & Validates JWT Tokens.
		core.RequestsValidator  // Validates GRPC requests.
		core.PwdHasher          // Hashes and ComparePwdss passwords.
		core.RateLimiter        // Limits rate of requests.
		core.TLSTool            // Configures and holds data for TLS communication.
	}
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/* - Setup App: Core, Tools & Layers - */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (app App) Setup() *App {

	// 1 -> Setup Core -> Config and Logger.
	app.Core.SetupConfigAndLogger()

	// 2 -> Setup Tools -> Interceptors, middleware, etc.
	app.Tools.Setup(app.Config)

	// 3 -> Setup Layers -> Service and External.
	//  * Service Layer = Business logic.
	//  * External Layer = Storage and Clients.
	app.Layers.Setup(&app.DBCfg, app.Tools)

	// 4 -> Setup Core -> GRPC & HTTP Servers.
	app.Core.SetupServers(app.Layers.ServiceLayer, app.Tools)

	return &app
}

// Step 1: Setup Config and Logger (on Core).
func (c *Core) SetupConfigAndLogger() {
	c.Config = core.LoadConfig()
	c.Logger = core.SetupLogger(&c.LoggerCfg)
}

// Step 2: Setup all Tools.
func (t *Tools) Setup(cfg *core.Config) {
	t.PwdHasher = tools.NewPwdHasher(cfg.PwdHasherCfg.Salt)
	t.RateLimiter = tools.NewRateLimiter(&cfg.RLimiterCfg)
	t.RequestsValidator = tools.NewProtoValidator()
	t.TokenAuthenticator = tools.NewJWTAuthenticator(cfg.JWTCfg.Secret, cfg.JWTCfg.SessionDays)
	t.TLSTool = tools.NewTLSTool(&cfg.TLSCfg)
}

// Step 3: Setup Layers.
func (l *Layers) Setup(dbCfg *core.DBCfg, toolbox core.Toolbox) {
	l.ExternalLayer = external.SetupLayer(dbCfg)
	l.ServiceLayer = service.SetupLayer(l.ExternalLayer, toolbox)
}

// Step 4: Setup GRPC & HTTP Servers (on Core).
func (c *Core) SetupServers(serviceLayer core.ServiceLayer, toolbox core.Toolbox) {
	c.Servers = servers.Setup(serviceLayer, toolbox, c.TLSCfg.Enabled)
}

/* -~-~-~-~-~ Tools & Toolbox -~-~-~-~-~- */

var _ core.Toolbox = (*Tools)(nil)

func (t *Tools) GetRequestsValidator() core.RequestsValidator { return t.RequestsValidator }
func (t *Tools) GetAuthenticator() core.TokenAuthenticator    { return t.TokenAuthenticator }
func (t *Tools) GetRateLimiter() core.RateLimiter             { return t.RateLimiter }
func (t *Tools) GetPwdHasher() core.PwdHasher                 { return t.PwdHasher }
func (t *Tools) GetTLSTool() core.TLSTool                     { return t.TLSTool }

/* -~-~-~-~-~ Run & Shutdown App -~-~-~-~-~- */

func (app *App) Run() {
	app.Servers.Run()
}

func (app *App) WaitForShutdown() { // Waits for SIGINT or SIGTERM to gracefully shutdown servers.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	zap.S().Infoln("Shutting down servers...")

	app.ExternalLayer.GetStorage().CloseDB()
	app.Servers.Shutdown()

	zap.S().Infoln("Servers stopped! Bye bye~")
	zap.L().Sync()
}
