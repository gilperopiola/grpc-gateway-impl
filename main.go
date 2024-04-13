package main

import (
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app"
	"github.com/gilperopiola/grpc-gateway-impl/etc/tools/treyser"

	"go.uber.org/zap"
)

/* ----------------------------------- */
/*            - Welcome~! -            */
/* ----------------------------------- */

/* This is the entrypoint of our app.
/* The app runs a gRPC Server and points an HTTP Gateway towards it.
/* It has a Service Layer that connects to a Storage Layer, which in turn connects to a SQL Database. */

//var setupStructure = []func(*app.App){
//	(*app.App).SetupConfig(core.LoadConfig),
//	(*app.App).SetupGlobalLogger,
//	(*app.App).InitService,
//	(*app.App).InitDatabaseAndStorage,
//}

func main() {
	treyser := treyser.NewTreyser("main", 0)

	// Init app.
	app := app.NewApp()

	treyser.Treys()

	// Run app.
	app.Run()
	time.Sleep(1 * time.Second)
	zap.S().Info("Servers OK")

	// Exit app.
	app.WaitForShutdown()
}

/* ----------------------------------- */
/*              - T0D0 -               */
/* ----------------------------------- */

/* Buf file / Dockerfile / Docker-compose / Kubernetes /
/* CI-CD / Metrics / Tracing / Caching / Tests */
