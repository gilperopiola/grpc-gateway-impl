package main

import (
	"time"

	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/cfg"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components"

	"go.uber.org/zap"
)

/* ----------------------------------- */
/*            - Welcome~! -            */
/* ----------------------------------- */

/* This is the entrypoint of our app.
/* The app runs a gRPC Server and points an HTTP Gateway towards it.
/* It has a Service Layer that connects to a Repository Layer, which in turn connects to a SQL Database. */

func main() {

	// Init app.
	app := v1.NewApp(
		loadConfig(),
		loadComponentsWrapper(),
	)

	// Run app.
	app.Run()
	time.Sleep(1 * time.Second)
	zap.S().Info("Servers OK")

	// Quit app.
	app.WaitForGracefulShutdown()
}

func loadConfig() *cfg.Config {
	return cfg.Load()
}

func loadComponentsWrapper() *components.Wrapper {
	return components.NewWrapper()
}

/* ----------------------------------- */
/*              - T0D0 -               */
/* ----------------------------------- */

/* Buf file / Dockerfile / Docker-compose / Kubernetes /
/* CI-CD / Metrics / Tracing / Caching / Tests */
