package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"

	"go.uber.org/zap"
)

func NewApp() *App {
	fmt.Println()
	app := App{&Core{}, &Service{}, &Toolbox{}}.Setup()
	return app
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - App (v1) -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type (
	App struct {
		*Core    // -> Servers, Config, Logger.
		*Service // -> Service Layer and External Layer.
		*Toolbox // -> JWT Auth, Reqs Validator, Rate Limiter, TLS, etc.
	}
	Core struct {
		core.Servers // -> GRPC and HTTP Servers.
		*core.Config // -> Config.
		*zap.Logger  // -> Logger (also lives globally in zap.L and zap.S).
	}
	Service struct {
		core.Service // -> All business logic.
	}
	Toolbox struct {
		core.TokenAuthenticator // -> Generates & Validates JWT Tokens.
		core.RequestsValidator  // -> Validates GRPC requests.
		core.DBTool             // -> Storage (DB, Cache, etc)
		core.APICaller          // -> Clients (GRPC, HTTP, etc).
		core.PwdHasher          // -> Hashes and ComparePwdss passwords.
		core.RateLimiter        // -> Limits rate of requests.
		core.TLSTool            // -> Configures and holds data for TLS communication.
	}
)

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
	db := app.Toolbox.GetDBTool().GetDB()
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
