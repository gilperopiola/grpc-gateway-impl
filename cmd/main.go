package main

import (
	"log"
	"time"

	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/cfg"
)

/* ----------------------------------- */
/*            - Welcome~! -            */
/* ----------------------------------- */
/*
/* This is the entrypoint of our app. Here we start the gRPC server and point the HTTP Gateway towards it. */

func main() {
	app := v1.NewAPI(cfg.Init())
	app.Init()

	app.Run()
	time.Sleep(1 * time.Second)
	log.Println("... ¡gRPC and HTTP OK! ...")

	app.WaitForGracefulShutdown()
}

/* ----------------------------------- */
/*              - T0D0 -               */
/* ----------------------------------- */
/* Buf file / Dockerfile / Docker-compose / Kubernetes / CI-CD / Tests /
/* Logging / Metrics / Tracing / Security / Caching / Rate limiting /
/* Postman collection / Full Swagger
/* -------------------------------------------------------------------------- */
