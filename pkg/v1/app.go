package v1

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/cfg"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/db"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"

	"go.uber.org/zap"
)

/* ----------------------------------- */
/*             - App v1 -              */
/* ----------------------------------- */

// App holds our entire application.
type App struct {
	*cfg.Config            // Holds the configuration.
	*components.Components // Holds all other dependencies.

	Service    service.Service       // Service Layer 		-> holds all business logic.
	Repository repository.Repository // Repository Layer 	-> manages communication with the database.
	Database   db.GormAdapter        // Database 			-> actual database connection.
}

// NewApp returns a new App with the given configuration and components.
func NewApp(config *cfg.Config, components *components.Components) *App {
	app := &App{Config: config, Components: components}
	app.Load()
	return app
}

// Load initializes all App components.
func (a *App) Load() {
	a.InitGlobalLogger()
	a.LoadCommonComponents()
	a.InitDBAndRepository()
	a.InitService()
	a.InitGRPCComponents()
	a.InitHTTPComponents()
}

// LoadCommonComponents initializes all common components.
func (a *App) LoadCommonComponents() {
	a.InitRateLimiter()
	a.InitPwdHasher()
	a.InitTLSComponents()
	a.InitInputValidator()
	a.InitAuthenticator()
}

// Run runs the gRPC & HTTP Servers.
func (a *App) Run() {
	a.Server.Run()
	a.Gateway.Run()
}

// WaitForGracefulShutdown waits for a SIGINT or SIGTERM to gracefully shutdown the servers.
func (a *App) WaitForGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	zap.S().Infoln("Shutting down servers...")

	sqlDB := a.Database.GetSQL()
	if sqlDB != nil {
		sqlDB.Close()
	}
	a.Server.Shutdown()
	a.Gateway.Shutdown()

	zap.S().Infoln("Servers stopped! Bye bye~")
	zap.L().Sync()
}
