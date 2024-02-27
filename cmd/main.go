package main

import (
	"log"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/server"
)

/* ----------------------------------- */
/*            - Welcome~! -            */
/* ----------------------------------- */
/*
/* This is the entrypoint of our app. Here we start the gRPC server and point the HTTP Gateway towards it. */

func main() {
	// Init app.
	config := server.LoadConfig()
	app := server.NewApp(config)

	// Init dependencies.
	app.InitGeneralDependencies()
	app.InitGRPCAndHTTPDependencies()
	app.InitAPI()
	app.InitServers()

	// Run servers.
	app.RunServers()
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
