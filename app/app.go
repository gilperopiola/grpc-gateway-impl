package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/servers"
	"github.com/gilperopiola/grpc-gateway-impl/app/service"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools/db_tool/sqldb"

	"go.uber.org/zap"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ */
/* -~-~-~-~-~ GRPC Gateway Implementation -~-~-~-~-~- */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ */

// -> Welcome to this lovely project :) 🌈

func NewApp() *App {
	fmt.Println()
	app := App{}.Setup()
	return app
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*               - App -               */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ v1 */

var _ core.Actions = (*Actions)(nil)

type (
	App struct {
		core.Servers // -> GRPC & HTTP.
		core.Service // -> Business logic.
		*core.Config // -> Config.
		*Actions     // -> Handy tools.
	}

	// Actions = Toolbox.
	// These below are our Tools, and with this struct we can perform any action we want.
	// Powerful.
	Actions struct {
		core.APICaller          // -> Clients (GRPC, HTTP, etc).
		core.DBTool             // -> Storage (DB, Cache, etc)
		core.FileCreator        // -> Creates folders and files.
		core.PwdHasher          // -> Hashes and compares passwords.
		core.RateLimiter        // -> Limits rate of requests.
		core.RequestsValidator  // -> Validates GRPC requests.
		core.TLSTool            // -> Holds and retrieves data for TLS communication.
		core.TokenAuthenticator // -> Generates & Validates JWT Tokens.
	}
)

/* -~-~-~-~-~ Setup App -~-~-~-~-~- */

func (app App) Setup() *App {
	app.SetupConfig()
	app.SetupLogger()
	app.SetupActions()
	app.SetupService()
	app.SetupServers()
	return &app
}

func (app *App) SetupConfig() {
	app.Config = core.SetupConfig()
}

func (app App) SetupLogger() {
	// The Logger lives globally on the zap pkg,
	// so we just initialize it here and forget about it.
	_ = core.SetupLogger(&app.LoggerCfg)
}

func (app *App) SetupActions() {
	cfg := app.Config

	app.Actions = &Actions{}
	app.Actions.APICaller = tools.NewAPICaller()
	app.Actions.DBTool = sqldb.NewDBTool(sqldb.NewSQLDB(&cfg.DBCfg))
	app.Actions.FileCreator = tools.NewFileCreator()
	app.Actions.PwdHasher = tools.NewPwdHasher(cfg.PwdHasherCfg.Salt)
	app.Actions.RateLimiter = tools.NewRateLimiter(&cfg.RLimiterCfg)
	app.Actions.RequestsValidator = tools.NewRequestsValidator()
	app.Actions.TLSTool = tools.NewTLSTool(&cfg.TLSCfg)
	app.Actions.TokenAuthenticator = tools.NewJWTAuthenticator(cfg.JWTCfg.Secret, cfg.JWTCfg.SessionDays)
}

func (app *App) SetupService() {
	app.Service = service.Setup(app.Actions)
}

func (app *App) SetupServers() {
	app.Servers = servers.Setup(app.Service, app.Actions, app.TLSCfg.Enabled)
}

/* -~-~-~-~-~ Run App -~-~-~-~-~- */

func (app *App) Run() {
	app.Servers.Run()
}

/* -~-~-~-~-~ Shutdown App -~-~-~-~-~- */

func (app *App) WaitForShutdown() { // Waits for SIGINT or SIGTERM to gracefully shutdown servers.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	zap.S().Infoln("GRPC Gateway Implementation | Shutting down servers 🛑")

	// Close SQL and Mongo connections.
	db := app.Actions.GetDB()
	if sqlDB, ok := db.(core.SQLDB); ok {
		sqlDB.Close()
	}
	if mongoDB, ok := db.(core.MongoDB); ok {
		mongoDB.Close(context.Background())
	}

	// Stop servers.
	app.Servers.Shutdown()

	zap.S().Infoln("GRPC Gateway Implementation | Servers stopped 🛑 Bye bye ~ ")
	zap.L().Sync()
}
