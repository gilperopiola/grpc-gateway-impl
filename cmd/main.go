package main

import (
	"log"
	"time"

	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"
	"github.com/gilperopiola/grpc-gateway-impl/server"
)

/* ----------------------------------- */
/*            - Welcome~! -            */
/* ----------------------------------- */
/*
/* This is the entrypoint of our app. Here we start the gRPC server and point the HTTP Gateway towards it. */

func main() {
	// Init App with Config.
	app := server.App{
		Config: v1.LoadConfig(),
	}

	// Load everything.
	app.Prepare()

	// Run servers.
	server.RunGRPCServer(app.GRPCServer, app.Config.GRPCPort)
	server.RunHTTPServer(app.HTTPGateway)
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
