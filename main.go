package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gilperopiola/grpc-gateway-impl/app"
)

/* ▶ ~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ */
/* -~-~-~-~- - GRPC Gateway Implementation - -~-~-~-~ */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- ▶ */

func main() {
	runApp, cleanUp := app.Setup()
	defer cleanUp()

	runApp()

	waitForSignal := make(chan os.Signal, 1)
	signal.Notify(waitForSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-waitForSignal
}

// ---             ╭───────────────╮                         ╭───────────╮
// GRPC Request -> │  GRPC Server  │ -> Interceptor Chain -> │  Service  │ -> DB Tool ~ API Calls ~ Etc ╮
//                 ╰───────────────╯                         ╰───────────╯                              |
// 		     			   ^                                                                            ↓
//		     			   |													 	          Service GRPC Response
//                         |                                                                            |
//                  Middleware Chain                            ╭───────────────╮                       |
//                     ^       | ╰--- -- -- -- -- - GRPC Out <- │  GRPC Server  │ <- Interceptor Chain <╯
//                     |       |    only http                   ╰───────────────╯
//                     |       |
//                     |       ↓
//                 ╭───────────────╮
// HTTP Request -> │  HTTP Server  │ -> HTTP Out
// ---             ╰───────────────╯
