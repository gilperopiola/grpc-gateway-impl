package v1

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/cfg"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/repository/db"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"
)

/* ----------------------------------- */
/*             - App v1 -              */
/* ----------------------------------- */

// App holds our entire application.
type App struct {
	*cfg.Config         // Config holds the App configuration.
	*components.Wrapper // Components Wrapper holds all other dependencies.

	Service    service.Service       // Service Layer holds all business logic.
	Repository repository.Repository // Repository Layer manages communication with the database.
	Database   *db.DatabaseWrapper   // Database holds the DB connection.
}

// NewApp returns a new App with the given configuration and components loaded.
func NewApp(config *cfg.Config, components *components.Wrapper) *App {
	app := &App{Config: config, Wrapper: components}
	app.Load()
	return app
}

// Load initializes all App components.
func (a *App) Load() {
	a.LoadCommonComponents()
	a.LoadRepositoryAndDB()
	a.LoadService()
	a.LoadGRPC()
	a.LoadHTTP()
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

	a.Database.Close()
	a.Server.Shutdown()
	a.Gateway.Shutdown()

	log.Println("Servers stopped! Bye bye~")
}
