package main

import (
	"log"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/server"
	"github.com/gilperopiola/grpc-gateway-impl/server/config"
)

/* ----------------------------------- */
/*            - Welcome~! -            */
/* ----------------------------------- */
/*
/* This is the entrypoint of our app. Here we start the gRPC server and point the HTTP Gateway towards it. */

func main() {

	// Init new App.
	app := server.NewApp(config.New())

	// Init everything inside of the App.
	app.Init()

	// Run servers.
	app.Run()
	time.Sleep(1 * time.Second)
	log.Println("... ¡gRPC and HTTP OK! ...")

	// Wait for shutdown.
	app.WaitForGracefulShutdown()
}

/* ----------------------------------- */
/*              - T0D0 -               */
/* ----------------------------------- */
/* Buf file / Dockerfile / Docker-compose / Kubernetes / CI-CD / Tests /
/* Logging / Metrics / Tracing / Security / Caching / Rate limiting /
/* Postman collection / Full Swagger
/* -------------------------------------------------------------------------- */
