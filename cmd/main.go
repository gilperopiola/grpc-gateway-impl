package main

import (
	"log"
	"time"

	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/cfg"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components"
)

/* ----------------------------------- */
/*            - Welcome~! -            */
/* ----------------------------------- */

/* This is the entrypoint of our app.
/* Here we start the gRPC server and point the HTTP Gateway towards it. */

func main() {

	// Init app.
	app := v1.NewApp(
		cfg.Load(),
		components.NewWrapper(),
	)

	// Run app.
	app.Run()
	time.Sleep(1 * time.Second)
	log.Println("Servers OK")

	// Quit app.
	app.WaitForGracefulShutdown()
}

/* ----------------------------------- */
/*              - T0D0 -               */
/* ----------------------------------- */
/* Buf file / Dockerfile / Docker-compose / Kubernetes / CI-CD /
/* Metrics / Tracing / Caching / Tests */
