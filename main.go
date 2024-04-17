package main

import (
	"github.com/gilperopiola/grpc-gateway-impl/app"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - Welcome~! -            */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

/* This is the entrypoint of our app */
/* It runs a gRPC Server and points an HTTP Gateway towards it */

func main() {

	// Init app
	app := app.NewApp()

	// Run app
	app.Run()

	// Exit app
	app.WaitForShutdown()
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*              - T0D0 -               */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Only log unexpected errors.

/* Buf file / Dockerfile / Docker-compose / Kubernetes /
/* CI-CD / Metrics / Tracing / Caching / Tests */
