package main

import (
	"github.com/gilperopiola/grpc-gateway-impl/app"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - Welcome~! -            */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

/* This is the entrypoint of our app */
/* It runs a gRPC Server and points an HTTP Gateway towards it */
/* Even though it's not a complex app, its architecture and overall code are extremely polished */

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

/* Update README / Put tests back in */
/* Buf file / Dockerfile / Docker-compose / Kubernetes /
/* CI-CD / Metrics / Tracing / Caching / Tests */
