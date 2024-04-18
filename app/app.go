package app

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/interfaces"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/special_types"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/external"
	"github.com/gilperopiola/grpc-gateway-impl/app/layers/servers"
	"github.com/gilperopiola/grpc-gateway-impl/app/modules"

	"go.uber.org/zap"
)

func NewApp() *App {
	return App{}.Setup()
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - App (v1) -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type (
	App struct {
		Core    // -> Config and Logger.
		Modules // -> JWT Auth, Input Validator, Rate Limiter, TLS, gRPC Interceptors, HTTP Middleware, etc.
		Layers  // -> Server / Business / External Layers.
	}

	Core struct {
		*core.Config // -> Config.
		*zap.Logger  // -> Logger (also lives globally in zap.L() and zap.S()).
	}

	Modules struct {
		*modules.Passive // -> Hold data.
		*modules.Active  // -> Do things.
	}

	Layers struct {
		special_types.ServerLayer // -> gRPC and HTTP Servers.
		interfaces.BusinessLayer  // -> Service, all business logic.
		external.ExternalLayer    // -> Storage (DB, Cache, etc) and Clients (gRPC, HTTP, etc).
	}
)

func (app *App) Run() {
	servers.RunGRPCServer(app.ServerLayer.GRPCServer)
	servers.RunHTTPGateway(app.ServerLayer.HTTPServer)

	go func() {
		time.Sleep(1 * time.Second)
		zap.S().Info("Servers OK")
	}()
}

// Waits for a SIGINT or SIGTERM to gracefully shutdown the servers.
func (app *App) WaitForShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	zap.S().Infoln("Shutting down servers...")

	sqlDB := app.ExternalLayer.GetDB().GetSQL()
	if sqlDB != nil {
		sqlDB.Close()
	}

	servers.ShutdownGRPCServer(app.ServerLayer.GRPCServer)
	servers.ShutdownHTTPGateway(app.ServerLayer.HTTPServer)

	zap.S().Infoln("Servers stopped! Bye bye~")
	zap.L().Sync()
}
